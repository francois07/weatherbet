package blockchain_test

import (
	"blockchain"
	"testing"
)

func TestBlockAddFailure(t *testing.T) {
	genesis := blockchain.Block{
		Transactions: []blockchain.Transaction{
			blockchain.NewTransaction(
				[]blockchain.TransactionInput{},
				[]blockchain.TransactionOutput{
					{
						Value:  200,
						Script: "test --- test OPDup test1 OPEqualVerify",
					},
				},
			),
		},
	}
	chain := blockchain.BlockChain{
		GenesisBlock: genesis,
	}
	chain.Chain = append(chain.Chain, genesis)
	chain.UTXO = make(map[string][]blockchain.TransactionOutput)
	chain.UTXO[genesis.Transactions[0].TXID] = genesis.Transactions[0].Outputs

	transaction := blockchain.NewTransaction(
		[]blockchain.TransactionInput{
			{
				TXID:       genesis.Transactions[0].TXID,
				VOUT:       0,
				ScriptArgs: map[string]string{"test": "test"},
			},
		},
		[]blockchain.TransactionOutput{},
	)

	ok := chain.AddBlock(blockchain.NewBlock(genesis, []blockchain.Transaction{transaction}))

	if ok {
		t.Fatalf("Got %v, expected false", ok)
	}
	if len(chain.Chain) != 1 {
		t.Fatalf("len(chain.Chain) == %v, expected 1", len(chain.Chain))
	}
}

func TestBlockAddSuccess(t *testing.T) {
	genesis := blockchain.Block{
		Transactions: []blockchain.Transaction{
			blockchain.NewTransaction(
				[]blockchain.TransactionInput{},
				[]blockchain.TransactionOutput{
					{
						Value:  200,
						Script: "test --- test OPDup test1 OPEqualVerify",
					},
				},
			),
		},
	}
	chain := blockchain.BlockChain{
		GenesisBlock: genesis,
	}
	chain.Chain = append(chain.Chain, genesis)
	chain.UTXO = make(map[string][]blockchain.TransactionOutput)
	chain.UTXO[genesis.Transactions[0].TXID] = genesis.Transactions[0].Outputs

	transaction := blockchain.NewTransaction(
		[]blockchain.TransactionInput{
			{
				TXID:       genesis.Transactions[0].TXID,
				VOUT:       0,
				ScriptArgs: map[string]string{"test": "test1"},
			},
		},
		[]blockchain.TransactionOutput{},
	)

	ok := chain.AddBlock(blockchain.NewBlock(genesis, []blockchain.Transaction{transaction}))

	if !ok {
		t.Fatalf("Got %v, expected true", ok)
	}
	if len(chain.Chain) != 2 {
		t.Fatalf("len(chain.Chain) == %v, expected 2", len(chain.Chain))
	}

	var length int
	for _, outputs := range chain.UTXO {
		length += len(outputs)
	}
	if length > 0 {
		t.Fatalf("len(chain.UTXO) == %v, expected 0", length)
	}
}
