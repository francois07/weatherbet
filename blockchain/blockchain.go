package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"internal/merkle"
	"log"
	"script"
	"strings"
	"time"
)

type TransactionInput struct {
	TXID          string
	VOUT          int
	ScriptArgs    map[string]string
	ScriptSigSize int
}

type TransactionOutput struct {
	Value  int
	Script string
}

type Transaction struct {
	TXID     string
	Inputs   []TransactionInput
	Outputs  []TransactionOutput
	LockTime time.Duration
}

type BlockHeader struct {
	PrevBlockHash [32]byte
	MerkleRoot    [32]byte
	Time          string
	Difficulty    int
	Nonce         int
	Height        int
}

type Block struct {
	Header       BlockHeader
	Transactions []Transaction
}

type BlockChain struct {
	GenesisBlock Block
	Chain        []Block
	UTXO         map[string][]TransactionOutput
}

func serialize[T any](target T) string {
	return fmt.Sprint(target)
}

func NewTransaction(inputs []TransactionInput, outputs []TransactionOutput) Transaction {
	serializedInputs := serialize(inputs)
	serializedOutputs := serialize(outputs)
	transactionHash := sha256.Sum256([]byte(serializedInputs + serializedOutputs))

	return Transaction{
		TXID:    hex.EncodeToString(transactionHash[:]),
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func NewBlock(prevBlock Block, transactions []Transaction) Block {
	var ids [][32]byte

	for _, t := range transactions {
		hexHash, err := hex.DecodeString(t.TXID)
		if err != nil {
			log.Fatal(err)
		}

		ids = append(ids, [32]byte(hexHash))
	}

	return Block{
		Header: BlockHeader{
			PrevBlockHash: prevBlock.Hash(),
			MerkleRoot:    merkle.MerkleRoot(ids),
			Time:          time.Now().String(),
			Difficulty:    2,
			Nonce:         0,
			Height:        prevBlock.Header.Height + 1,
		},
		Transactions: transactions,
	}
}

func (c *BlockChain) IsUnspent(TXID string, idx int) bool {
	val, ok := c.UTXO[TXID]
	return ok && (idx < len(val))
}

func (c *BlockChain) Unlock(input TransactionInput) (int, bool) {
	if !c.IsUnspent(input.TXID, input.VOUT) {
		return 0, false
	}
	utxo := c.UTXO[input.TXID][input.VOUT]

	return utxo.Value, script.EvalScript(utxo.Script, input.ScriptArgs)
}

func (c *BlockChain) Spend(input TransactionInput) bool {
	_, ok := c.Unlock(input)

	if ok {
		if input.VOUT < len(c.UTXO[input.TXID])-1 {
			c.UTXO[input.TXID] = append(c.UTXO[input.TXID][:input.VOUT], c.UTXO[input.TXID][input.VOUT+1:]...)
		} else {
			c.UTXO[input.TXID] = c.UTXO[input.TXID][:input.VOUT]
		}
	}

	return ok
}

func (c *BlockChain) IsValidTransaction(t Transaction) bool {
	balance := 0
	totalSpent := 0

	for _, input := range t.Inputs {
		val, ok := c.Unlock(input)
		if !ok {
			return false
		}
		balance += val
	}

	for _, output := range t.Outputs {
		totalSpent += output.Value
	}

	return totalSpent <= balance
}

func (b *Block) Hash() [32]byte {
	serializedBlock := serialize(b)
	blockHash := sha256.Sum256([]byte(serializedBlock))

	return blockHash
}

func (b *Block) IsValid() bool {
	hash := b.Hash()
	strHash := hex.EncodeToString(hash[:])
	return strings.HasPrefix(strHash, strings.Repeat("0", b.Header.Difficulty))
}

func (b *Block) Mine(stopMining chan bool, callback func(Block)) {
	for {
		select {
		case <-stopMining:
			fmt.Println("Mining stopped")
			return
		default:
			b.Header.Nonce += 1

			if b.IsValid() {
				fmt.Printf("Block mined successfully with Nonce %d and hash %x\n", b.Header.Nonce, b.Hash())
				fmt.Println(b)
				callback(*b)
				return
			}
		}
	}
}

func (c *BlockChain) AddBlock(b Block) bool {
	if !b.IsValid() {
		return false
	}

	transactions := b.Transactions
	for _, transaction := range transactions {
		if !c.IsValidTransaction(transaction) {
			return false
		}
	}

	for _, transaction := range transactions {
		for _, input := range transaction.Inputs {
			c.Spend(input)
		}
		c.UTXO[transaction.TXID] = transaction.Outputs
	}

	c.Chain = append(c.Chain, b)

	return true
}

func NewChain(genesis Block) BlockChain {
	chain := BlockChain{
		GenesisBlock: genesis,
		Chain:        []Block{genesis},
		UTXO:         make(map[string][]TransactionOutput),
	}

	for _, transaction := range genesis.Transactions {
		chain.UTXO[transaction.TXID] = transaction.Outputs
	}

	return chain
}

func (c *BlockChain) IsValid() bool {
	for _, block := range c.Chain {
		if !block.IsValid() {
			return false
		}
	}

	return true
}
