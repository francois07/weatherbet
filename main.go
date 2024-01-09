package main

import (
	"blockchain"
	"client"
	"os"
)

func main() {
	genesis := blockchain.Block{
		Transactions: []blockchain.Transaction{
			blockchain.NewTransaction(
				[]blockchain.TransactionInput{},
				[]blockchain.TransactionOutput{
					{
						Value:  200,
						Script: "sign pubKey --- pubKey OPDup OPHash 3e4c25fe2d8751520c0b444d3e43a955feb782f10b25c68acebfe8c29dc63c91 OPEqualVerify sign OPDup pubKey OPDup OPHash pubKey OPDup OPCheckSig",
					},
				},
			),
		},
	}
	chain := blockchain.NewChain(genesis)
	var peers []string

	if len(os.Args) > 1 {
		peers = append(peers, os.Args[1:]...)
	}

	blockClient := client.NewClient(&chain, peers)
	blockClient.Start()
}
