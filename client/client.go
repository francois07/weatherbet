package client

import (
	"blockchain"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Router          *gin.Engine
	Scheduler       *cron.Cron
	BlockChain      *blockchain.BlockChain
	TransactionPool []blockchain.Transaction
	Peers           []string
}

func NewClient(chain *blockchain.BlockChain, peers []string) Client {
	client := Client{
		Router:     gin.Default(),
		BlockChain: chain,
		Scheduler:  cron.New(),
		Peers:      peers,
	}

	client.Router.GET("/api", client.getBlockChain)
	client.Router.GET("/api/transactions", client.getTransactions)
	client.Router.GET("/api/utxo", client.getUTXO)
	client.Router.GET("/api/peers", client.getPeers)
	client.Router.POST("/api/transactions", client.postTransaction)
	client.Router.POST("/api", client.postBlock)

	client.Scheduler.AddFunc("@every 1m", client.MineCandidateBlock)

	return client
}

func (client *Client) Start() {
	fmt.Printf("Starting client with %d peers\n", len(client.Peers))

	client.syncPeerList()

	for _, peer := range client.Peers {
		client.SyncBlockchain(peer)
	}

	client.Scheduler.Start()
	client.Router.Run()
}

func (client *Client) MineCandidateBlock() {
	transactions := client.TransactionPool
	if len(transactions) < 1 {
		fmt.Println("No transactions in transaction pool, skipping mining")
		return
	} else {
		fmt.Printf("%d transactions in transaction pool, starting mining...\n", len(transactions))
	}

	client.TransactionPool = nil
	prevBlock := client.BlockChain.Chain[len(client.BlockChain.Chain)-1]

	candidateBlock := blockchain.NewBlock(prevBlock, transactions)
	stopMining := make(chan bool)

	go candidateBlock.Mine(stopMining, func(b blockchain.Block) { client.AddBlockAndPropagate(b) })
}

func (client *Client) getBlockChain(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, client.BlockChain)
}

func (client *Client) postTransaction(c *gin.Context) {
	var newTransaction struct {
		Inputs   []blockchain.TransactionInput
		Outputs  []blockchain.TransactionOutput
		LockTime time.Duration
	}

	if err := c.BindJSON(&newTransaction); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction fields"})
		return
	}

	transaction := blockchain.NewTransaction(newTransaction.Inputs, newTransaction.Outputs)
	if !client.BlockChain.IsValidTransaction(transaction) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction"})
		return
	}

	client.TransactionPool = append(client.TransactionPool, transaction)
}

func (client *Client) getTransactions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, client.TransactionPool)
}

func (client *Client) getUTXO(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, client.BlockChain.UTXO)
}

func (client *Client) getPeers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, client.Peers)
}

func (client *Client) syncPeerList() {
	var newPeers []string

	for _, peer := range client.Peers {
		resp, err := http.Get(fmt.Sprintf("%s/api/peers", peer))
		if err != nil {
			continue
		}

		respBody, errParse := io.ReadAll(resp.Body)
		if errParse != nil {
			continue
		}

		var body []string
		json.Unmarshal(respBody, &body)

		newPeers = append(newPeers, body...)
	}

	fmt.Printf("Added %d new peers\n", len(newPeers))
}

func (client *Client) SyncBlockchain(peer string) {
	resp, err := http.Get(fmt.Sprintf("%s/api", peer))
	if err != nil {
		fmt.Printf("Error fetching blockchain from %s: %v\n", peer, err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body from %s: %v\n", peer, err)
		return
	}

	var bc blockchain.BlockChain
	err = json.Unmarshal(respBody, &bc)
	if err != nil {
		fmt.Printf("Error unmarshalling blockchain from %s: %v\n", peer, err)
		return
	}

	if !bc.IsValid() {
		fmt.Printf("Invalid blockchain received from %s\n", peer)
		return
	}

	if len(bc.Chain) > len(client.BlockChain.Chain) {
		*client.BlockChain = bc
		fmt.Printf("Synced chain of length %d with %s\n", len(bc.Chain), peer)
	}
}

func (client *Client) postBlock(c *gin.Context) {
	var block blockchain.Block

	if err := c.BindJSON(&block); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid block fields"})
		return
	}
	fmt.Println(block, block.Hash())

	if block.Header.Height < len(client.BlockChain.Chain) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Block height is too low"})
		return
	}

	if !block.IsValid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid block"})
		return
	}

	for idx, transaction := range block.Transactions {
		if !client.BlockChain.IsValidTransaction(transaction) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid transaction with index %d", idx)})
			return
		}
	}

	client.AddBlockAndPropagate(block)
}

func (client *Client) AddBlockAndPropagate(block blockchain.Block) {
	reqBody, err := json.Marshal(block)
	if err != nil {
		return
	}

	ok := client.BlockChain.AddBlock(block)
	if ok {
		fmt.Printf("Added block with hash %x\n", block.Hash())
	}

	for _, peer := range client.Peers {
		resp, err := http.Post(fmt.Sprintf("%s/api", peer), "application/json", bytes.NewBuffer(reqBody))
		if err == nil && resp.StatusCode == 200 {
			fmt.Printf("Propagated block to %s\n", peer)
		} else {
			fmt.Printf("%s refused block\n", peer)
		}
	}
}
