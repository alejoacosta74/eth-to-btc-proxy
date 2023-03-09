
# `qproxy`: An eth to qtum signature conversion proxy

`qproxy` is a proxy server that converts a signed Ethereum transaction to a Qtum (Bitcoin) signed transaction and sends it to the Qtum node for broadcasting.

Simply put, `qproxy` enables a ethereum user or app to transact on the Qtum blockchain seamlessly, without having to worry about the underlying utxo model.

## How it works
The proxy runs a JSON RPC server that implements the method `eth_sendRawTransaction` (as defined in the Ethereum JSON RPC API)

Whenever a Ethereum raw transaction is received, the tx is decoded and the signature is verified.

Next, a new utxo based transaction is assembled with the same values from the original Ethereum transaction that are relevant to the Qtum blockchain, such as `amount`, `data`, `to`, `gas`, etc.

Finally the new utxo based transaction is signed and broadcasted to the Qtum network, and the tx hash is returned to the original caller as part of the `eth_senRawTransaction` result.

## Requirements

Two basic conditions must be met in order for the whole workflow to run succesfully:
1. The private key of the sender/signer must be available to the proxy server. This can be achieved by sending the private key via the `personal_importRawKey` [JSON RPC method](https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-personal#personal_importrawkey)

2. On the target blockchain (i.e. Qtum), there should exist a UTXO set that is spendable by the owner of the private key, with enough balance to cover the `amount` to be sent as well as to cover the gas fee.

## Usage

### Run a qtum node
A qtum node can be easily deployed with docker running the following command:

```
docker run -d -p 3889:3889 -p 3888:3888 -v /path/to/qtum/data:/root/.qtum qtumproject/qtum:latest
```

### Run the proxy

```
make build
qtumproxy --qtumrpc=127.0.0.1:3889 --user=qtum --pass=qtum -l info
```

### Share private key with proxy server

```
curl -X -d '{"jsonrpc":"2.0","method": "personal_importRawKey", "params": [string, string],"id":1}' -H 'content-type: application/json;' http://127.0.0.1:8080/rpc
```

### Send a raw transaction

```
curl -X -d '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":[" "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"],"id":1}' -H 'content-type: application/json;' http://127.0.0.1:8080/rpc
```

## Features
- A reverse-proxy is available at endpoint `/proxy` that can be used to send requests to an ethereum node (like Ganache) and log both JSON RPC request and response (usefull for debugging and testing)
- Command line configuration can be passed as flags, environment vars or within a `.config.yml` file
- JSON RPC service implementation based on *go-ethereum* rpc module.
- QTUM rpc client implementation based on *btcd* bitcoin rpc client
- Support for different log levels (info, trace, debug)

## Run tests

- Unit tests
  
   Unit tests can be run with: `make unit-test`

- Integration tests

   Integration tests are run against a `qtumd` node running on a Docker container that is started and stopped automatically by the test suite.

   A test ethereum tx is built using an [ethereum cli](https://github.com/alejoacosta74/ethcli), which also signs and sends the tx to `qproxy`.
   The tx is then decoded and verified by `qproxy`, and a new qtum tx is built and signed using the same private key.

   Run integration tests with: `make integration-test` 

## TODO

1. Add SSL support
2. Implement bitcoin wallet to store private keys persistently
3. ~~Implement ethereum signature verification~~ :white_check_mark:
4. Implement all possible ethereum interaction use cases:

   a. ~~Send eth~~ :white_check_mark:

   b. Create contract

   c. Call contract

5. ~~Add support for EIP 1155 signature~~ :white_check_mark:
6. Add persistent mapping between ethereum tx hash and qtum tx hash
7. Implement automated integration tests for all use cases using Qtum regtest network