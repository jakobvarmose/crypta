# Crypta

A decentralized social network. Built with [IPFS](https://ipfs.io/), [Quasar](http://quasar-framework.org/) and [Vue](https://vuejs.org/).

## Getting Started

These instructions will get you up and running on your local machine. Crypta connects to the IPFS network, but it doesn't require a separate IPFS client.

### Prerequisites

First you need to install [Go](https://golang.org/) and [Node.js](https://nodejs.org/).

### Building and Running

```
go get -u -d github.com/jakobvarmose/crypta
cd $(go env GOPATH)/src/github.com/jakobvarmose/crypta
make
./crypta
```

### Development

To make development easy you can also run the GUI and API separately.

Run the API in one terminal (note the `--dev` flag):

```
cd $(go env GOPATH)/src/github.com/jakobvarmose/crypta
go build && ./crypta --dev
```

...and the GUI in another:

```
cd $(go env GOPATH)/src/github.com/jakobvarmose/crypta/web/cryptaio
npm run dev
```

## License

This project is licensed under the MIT License.
