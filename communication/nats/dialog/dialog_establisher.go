package dialog

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mysterium/node/communication"
	"github.com/mysterium/node/communication/nats"
	"github.com/mysterium/node/communication/nats/discovery"
	"github.com/mysterium/node/identity"
	dto_discovery "github.com/mysterium/node/service_discovery/dto"
)

// NewDialogEstablisher constructs new DialogEstablisher which works thru NATS connection.
func NewDialogEstablisher(myID identity.Identity, signer identity.Signer) *dialogEstablisher {

	return &dialogEstablisher{
		myID:     myID,
		mySigner: signer,
		peerAddressFactory: func(contact dto_discovery.Contact) (*discovery.AddressNATS, error) {
			address, err := discovery.NewAddressForContact(contact)
			if err == nil {
				err = address.Connect()
			}

			return address, err
		},
	}
}

const establisherLogPrefix = "[NATS.DialogEstablisher] "

type dialogEstablisher struct {
	myID               identity.Identity
	mySigner           identity.Signer
	peerAddressFactory func(contact dto_discovery.Contact) (*discovery.AddressNATS, error)
}

func (establisher *dialogEstablisher) EstablishDialog(
	peerID identity.Identity,
	peerContact dto_discovery.Contact,
) (communication.Dialog, error) {

	log.Info(establisherLogPrefix, fmt.Sprintf("Connecting to: %#v", peerContact))
	peerAddress, err := establisher.peerAddressFactory(peerContact)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to: %#v. %s", peerContact, err)
	}

	peerCodec := establisher.newCodecForPeer(peerID)

	peerSender := establisher.newSenderToPeer(peerAddress, peerCodec)
	err = establisher.negotiateDialog(peerSender)
	if err != nil {
		return nil, err
	}

	dialog := establisher.newDialogToPeer(peerID, peerAddress, peerCodec)
	log.Info(establisherLogPrefix, fmt.Sprintf("Dialog established with: %#v", peerContact))

	return dialog, nil
}

func (establisher *dialogEstablisher) negotiateDialog(sender communication.Sender) error {
	response, err := sender.Request(&dialogCreateProducer{
		&dialogCreateRequest{
			PeerID: establisher.myID.Address,
		},
	})
	if err != nil {
		return fmt.Errorf("dialog creation error. %s", err)
	}
	if response.(*dialogCreateResponse).Reason != 200 {
		return fmt.Errorf("dialog creation rejected. %#v", response)
	}

	return nil
}

func (establisher *dialogEstablisher) newCodecForPeer(peerID identity.Identity) *codecSecured {

	return NewCodecSecured(
		communication.NewCodecJSON(),
		establisher.mySigner,
		identity.NewVerifierIdentity(peerID),
	)
}

func (establisher *dialogEstablisher) newSenderToPeer(
	peerAddress *discovery.AddressNATS,
	peerCodec *codecSecured,
) communication.Sender {

	return nats.NewSender(
		peerAddress.GetConnection(),
		peerCodec,
		peerAddress.GetTopic(),
	)
}

func (establisher *dialogEstablisher) newDialogToPeer(
	peerID identity.Identity,
	peerAddress *discovery.AddressNATS,
	peerCodec *codecSecured,
) *dialog {

	subTopic := peerAddress.GetTopic() + "." + establisher.myID.Address
	return &dialog{
		peerID:   peerID,
		Sender:   nats.NewSender(peerAddress.GetConnection(), peerCodec, subTopic),
		Receiver: nats.NewReceiver(peerAddress.GetConnection(), peerCodec, subTopic),
	}
}
