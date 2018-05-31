/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package server

import (
	log "github.com/cihub/seelog"
	"github.com/mysterium/node/cmd"
	"github.com/mysterium/node/communication"
	"github.com/mysterium/node/identity"
	"github.com/mysterium/node/ip"
	"github.com/mysterium/node/location"
	"github.com/mysterium/node/nat"
	"github.com/mysterium/node/openvpn"
	"github.com/mysterium/node/openvpn/discovery"
	"github.com/mysterium/node/openvpn/middlewares/state"
	"github.com/mysterium/node/openvpn/tls"
	"github.com/mysterium/node/server"
	dto_discovery "github.com/mysterium/node/service_discovery/dto"
	"github.com/mysterium/node/session"
	"github.com/mysterium/node/version"
	"sync"
	"time"
)

// Command represent entrypoint for Mysterium server with top level components
type Command struct {
	identityLoader   func() (identity.Identity, error)
	createSigner     identity.SignerFactory
	ipResolver       ip.Resolver
	mysteriumClient  server.Client
	natService       nat.NATService
	locationDetector location.Detector

	dialogWaiterFactory func(identity identity.Identity) communication.DialogWaiter
	dialogWaiter        communication.DialogWaiter

	sessionManagerFactory func(primitives *tls.Primitives, serverIP string) session.Manager

	vpnServerFactory func(sessionManager session.Manager, primitives *tls.Primitives, openvpnStateCallback state.Callback) *openvpn.Server

	vpnServer      *openvpn.Server
	checkOpenvpn   func() error
	protocol       string
	WaitUnregister *sync.WaitGroup
}

// Start starts server - does not block
func (cmd *Command) Start() (err error) {
	printLicense()
	log.Info("[Server version]", version.AsString())
	err = cmd.checkOpenvpn()
	if err != nil {
		return err
	}

	providerID, err := cmd.identityLoader()
	if err != nil {
		return err
	}

	cmd.dialogWaiter = cmd.dialogWaiterFactory(providerID)
	providerContact, err := cmd.dialogWaiter.Start()

	// if for some reason we will need truly external IP, use GetPublicIP()
	vpnServerIP, err := cmd.ipResolver.GetOutboundIP()
	if err != nil {
		return err
	}

	cmd.natService.Add(nat.RuleForwarding{
		SourceAddress: "10.8.0.0/24",
		TargetIP:      vpnServerIP,
	})
	if err = cmd.natService.Start(); err != nil {
		return err
	}

	location, err := cmd.locationDetector.DetectLocation()
	if err != nil {
		return err
	}
	log.Info("Country detected: ", location.Country)
	serviceLocation := dto_discovery.Location{Country: location.Country}

	proposal := discovery.NewServiceProposalWithLocation(providerID, providerContact, serviceLocation, cmd.protocol)

	primitives, err := tls.NewTLSPrimitives(serviceLocation, providerID)
	if err != nil {
		return err
	}

	sessionManager := cmd.sessionManagerFactory(primitives, vpnServerIP)

	dialogHandler := session.NewDialogHandler(proposal.ID, sessionManager)
	if err := cmd.dialogWaiter.ServeDialogs(dialogHandler); err != nil {
		return err
	}

	stopDiscoveryAnnouncement := make(chan int)
	vpnStateCallback := func(state openvpn.State) {
		switch state {
		case openvpn.ConnectedState:
			log.Info("Open vpn service started")
		case openvpn.ExitingState:
			log.Info("Open vpn service exiting")
			close(stopDiscoveryAnnouncement)
		}
	}
	cmd.vpnServer = cmd.vpnServerFactory(sessionManager, primitives, vpnStateCallback)
	if err := cmd.vpnServer.Start(); err != nil {
		return err
	}

	signer := cmd.createSigner(providerID)

	cmd.WaitUnregister.Add(1)
	go cmd.discoveryAnnouncementLoop(proposal, cmd.mysteriumClient, signer, stopDiscoveryAnnouncement)

	return nil
}

func printLicense() {
	license := cmd.GetStartupLicense(
		"run program with '-warranty' option",
		"run program with '-conditions' option",
	)
	log.Info("\n" + license)
}

// Wait blocks until server is stopped
func (cmd *Command) Wait() error {
	cmd.WaitUnregister.Wait()

	return cmd.vpnServer.Wait()
}

// Kill stops server
func (cmd *Command) Kill() error {
	cmd.vpnServer.Stop()

	err := cmd.dialogWaiter.Stop()
	if err != nil {
		return err
	}

	err = cmd.natService.Stop()

	return err
}

func (cmd *Command) discoveryAnnouncementLoop(proposal dto_discovery.ServiceProposal, mysteriumClient server.Client, signer identity.Signer, stopPinger <-chan int) {
	for {
		err := mysteriumClient.RegisterProposal(proposal, signer)
		if err != nil {
			log.Errorf("Failed to register proposal: %v, retrying after 1 min.", err)
			time.Sleep(1 * time.Minute)
		} else {
			break
		}
	}
	cmd.pingProposalLoop(proposal, mysteriumClient, signer, stopPinger)

}

func (cmd *Command) pingProposalLoop(proposal dto_discovery.ServiceProposal, mysteriumClient server.Client, signer identity.Signer, stopPinger <-chan int) {
	for {
		select {
		case <-time.After(1 * time.Minute):
			err := mysteriumClient.PingProposal(proposal, signer)
			if err != nil {
				log.Error("Failed to ping proposal: ", err)
				// do not stop server on missing ping to discovery. More on this in MYST-362 and MYST-370
			}
		case <-stopPinger:
			log.Info("Stopping discovery announcement")
			err := mysteriumClient.UnregisterProposal(proposal, signer)
			if err != nil {
				log.Error("Failed to unregister proposal: ", err)
			}
			time.Sleep(200 * time.Millisecond) // sleep for prints to be printed out
			cmd.WaitUnregister.Done()
			return
		}
	}
}
