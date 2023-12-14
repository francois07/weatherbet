package blockchain

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"internal/merkle"
	"log"
	"strings"
	"time"
)

type ScriptSig struct {
	Signature string
	PubKey    crypto.PublicKey
}

type TransactionInput struct {
	TXID          string
	VOUT          int
	ScriptSig     ScriptSig
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
	PrevBlockHash string
	MerkleRoot    string
	Time          time.Time
	Difficulty    int
	Nonce         int
}

type Block struct {
	Header       BlockHeader
	Transactions []Transaction
}

type BlockChain struct {
	GenesisBock Block
	Chain       []Block
}

func serialize[T any](target T) string {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(target)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func NewTransaction(inputs []TransactionInput, outputs []TransactionOutput) Transaction {
	serializedInputs := serialize(inputs)
	serializedOutputs := serialize(outputs)
	transactionHash := sha256.Sum256([]byte(serializedInputs + serializedOutputs))

	return Transaction{
		TXID:    fmt.Sprint(transactionHash),
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func NewBlock(prevBlock Block, transactions []Transaction) Block {
	var ids []string

	for _, t := range transactions {
		ids = append(ids, t.TXID)
	}

	return Block{
		Header: BlockHeader{
			PrevBlockHash: prevBlock.Hash(),
			MerkleRoot:    merkle.MerkleRoot(ids),
			Time:          time.Now(),
			Difficulty:    2,
			Nonce:         0,
		},
		Transactions: transactions,
	}
}

func (b *Block) Hash() string {
	serializedBlock := serialize(b)
	blockHash := sha256.Sum256([]byte(serializedBlock))

	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine() {
	for !strings.HasPrefix(b.Hash(), strings.Repeat("0", b.Header.Difficulty)) {
		b.Header.Nonce += 1
	}
}
