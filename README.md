# Mysterium Node - decentralized VPN built on blockchain

[![Go Report Card](https://goreportcard.com/badge/github.com/MysteriumNetwork/node)](https://goreportcard.com/report/github.com/MysteriumNetwork/node)
[![Build Status](https://travis-ci.org/MysteriumNetwork/node.svg?branch=master)](https://travis-ci.org/MysteriumNetwork/node)

Cross-platform software to run a node in Mysterium Network. It contains Mysterium server (node), 
client API (tequila API) and client-cli (console client) for Mysterium Network.
 
Currently node supports OpenVPN as its underlying VPN transport. 

## Getting Started

- Homepage https://mysterium.network
- [Whitepaper](https://mysterium.network/whitepaper.pdf)
- [Latest](https://github.com/MysteriumNetwork/node/releases/latest) release
- [Installation guide](./INSTALL.md)
- [Node wiki](https://github.com/MysteriumNetwork/node/wiki/) 

### Prerequisites

To run a node as docker container You will need [docker](https://www.docker.com/). 
You should be able to run a node on any OS that supports docker. 
Tested on these OSes so far: _Dabian 9_, _Ubuntu 16.04_ and _Centos 7_. 

You can check latest docker node versions here: https://hub.docker.com/r/mysteriumnetwork/mysterium-node/


### Installation

Go to [docker](https://www.docker.com/) on how to get a recent docker version for Your OS.

```bash
sudo apt-get install docker-ce
sudo docker run --cap-add NET_ADMIN --net host --publish "1194:1194" --name mysterium-node -d mysteriumnetwork/mysterium-node:latest
```
### Running
```bash
sudo docker start mysterium-node
sudo docker stop mysterium-node
```
### Debugging
```bash
sudo docker logs -f mysterium-node
```
More detailed installation options described [here](./INSTALL.md).
For possible issues while running a node refer to our [FAQ](https://github.com/MysteriumNetwork/node/wiki/Node-operation) section.

## Built With

* [go](https://golang.org/) - The Go Programming Language
* [travis](https://travis-ci.org/) - Travis continuous integration tool
* [docker](https://www.docker.com/what-docker) - Containerize applications
* [openvpn](https://openvpn.net) - Solid VPN solution

## Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors
* **Valdas Petrulis** - *Lead developer, go evangelist* - [Waldz](https://github.com/Waldz)
* **Tadas Valiukas** - *Senior developer, experienced bug maker* - [tadovas](https://github.com/tadovas)
* **Donatas Kučinskas** - *Senior developer, clean code savvy* - [donce](https://github.com/donce)
* **Antanas Masevičius** - *Network engineer / developer, net guru* - [zolia](https://github.com/zolia)
* **Paulius Mozuras** - *Software developer, python lover* - [interro](https://github.com/interro)
* **Ignas Bernotas** - *Full stack developer, lifestyle expert* - [ignasbernotas](https://github.com/ignasbernotas)
* **Andrej Novikov** - *Senior developer, JS enthusiast, UX jazzman* - [shroomist](https://github.com/shroomist)

See also full list of [contributors](https://github.com/MysteriumNetwork/node/contributors) who participated in this project.

## License

This project is licensed under the terms of the GNU General Public License v3.0 (see [details](./LICENSE)).
