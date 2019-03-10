# go-blockchain

[![Build Status](https://travis-ci.org/danitello/go-blockchain.svg?branch=master)](https://travis-ci.org/danitello/go-blockchain)

Blockchain implementation secured by proof of work

```bash
go get
go run main.go create-wallet # returns ADDR1
go run main.go create-wallet # returns ADDR2
go run main.go init-chain -address <ADDR1> # receives coinbase
go run main.go balance -address <ADDR1>
go run main.go balance -address <ADDR2>
go run main.go send -from <ADDR1> -to <ADDR2> -amount <A_NUMBER>
go run main.go balance -address <ADDR1>
go run main.go balance -address <ADDR2>
go run main.go print-chain
```
This will likely change as more functionality is added.

## Objective
This project is primarily to get a deeper understanding of - and experience with - the Go language, blockchain/cryptography, and surrounding concepts.

## Reference
References implementations made by [Tensor Programming](http://tensor-programming.com/) (based on [Bitcoin](https://bitcoin.org) spec) and [Ethereum](https://www.ethereum.org/)