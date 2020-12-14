package assignment02IBC_master

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

type Node struct {
	Address    string
	BlockChain *Block
	Nodes      map[string]bool
	Votes      map[string]string
	VoteFor    string
}

func (node *Node) SendBlockChain(address string) {
	fmt.Println("sending blockchain to " + address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode("ReceiveBlockChain," + node.Address)

	err = encoder.Encode(node.BlockChain)
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("sending blockchain completed to " + address)
}
func (node *Node) ReceiveBlockChain( /*conn net.Conn*/ decoder *gob.Decoder) {
	fmt.Println("receiving blockchain")
	var blockChain *Block
	//decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&blockChain)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("decode successful in client receive blockchain")
	if VerifyChain(blockChain) && node.VerifyVote(blockChain) {
		node.BlockChain = blockChain
		node.SendVote()
		fmt.Println("blockchain received")
		ListBlocks(node.BlockChain)
		node.ListAccounts()
		return
	}
	fmt.Println("received blockchain not verified")
	fmt.Println("receiving blockchain")
}
func (node *Node) ListAccounts() {
	fmt.Println("displaying all nodes with their balances")
	for n := range node.Nodes {
		fmt.Println(n+": ", CalculateAmount(n, *node.BlockChain))
	}
}
func (node *Node) MineBlock(transaction Transaction) {
	fmt.Println("mining blockchain")
	if (transaction.Sender != "early joiners reward") && (CalculateAmount(transaction.Sender, *node.BlockChain) < transaction.Amount) {
		fmt.Println("not enough balance")
		return
	}
	if transaction.IsEmpty() {
		fmt.Println("transaction is empty")
		return
	}
	t2 := Transaction{
		Amount:   1,
		Sender:   "mining reward",
		Receiver: node.Address,
	}

	var trns []Transaction
	trns = append(trns, transaction)
	trns = append(trns, t2)
	node.BlockChain = InsertBlock(trns, node.DetermineNextMiner(), node.Votes, node.BlockChain)
	fmt.Println("mining blockchain completed")
}
func (node *Node) VerifyBlock() {

}
func (node *Node) InitiateTransaction(to string, amount int) {
	fmt.Println("initiating transaction")
	if CalculateAmount(node.Address, *node.BlockChain) < amount {
		fmt.Println("not enough amount")
		return
	}
	trans := Transaction{
		Amount:   amount,
		Sender:   node.Address,
		Receiver: to,
	}
	conn, err := net.Dial("tcp", node.BlockChain.NextMiner)
	if err != nil {
		fmt.Println(err)
		fmt.Println("transaction failed")
		return
	}
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode("ReceiveTransaction," + node.Address)
	if err != nil {
		fmt.Println(err)
		fmt.Println("transaction failed")
		return
	}
	err = encoder.Encode(trans)
	if err != nil {
		fmt.Println(err)
		fmt.Println("transaction failed")
		return
	}
	_ = conn.Close()
	fmt.Println("initiating transaction completed")
}
func (node *Node) ReceiveTransaction(decoder *gob.Decoder) {
	fmt.Println("receiving transaction")
	var transaction Transaction
	err := decoder.Decode(&transaction)
	if err != nil {
		fmt.Println(err)
		fmt.Println("receiving transaction failed")
		return
	}
	fmt.Println("receiving transaction completed")
	node.MineBlock(transaction)
	node.FloodBlockChain()
}
func (node *Node) SendVote() {
	fmt.Println("sending vote")
	conn, err := net.Dial("tcp", node.BlockChain.NextMiner)
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(conn)

	err = encoder.Encode("ReceiveVote," + node.Address)
	err = encoder.Encode(node.VoteFor)
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Println("sending vote failed")
		return
	}
	fmt.Println("sending vote completed")
}
func (node *Node) ReceiveVote(decoder *gob.Decoder, voter string) {
	fmt.Println("receiving vote")
	var castedVote string
	err := decoder.Decode(&castedVote)
	if err != nil {
		fmt.Println(err)
		fmt.Println("receiving vote failed")
		return
	}
	node.Votes[voter] = castedVote
	fmt.Println("receiving vote completed")
}

func (node *Node) VerifyVote(blockChain *Block) bool {
	fmt.Println("verifying vote")
	return blockChain.Votes[node.Address] == node.VoteFor
}
func (node *Node) ChangeMyVote(to string) {
	node.VoteFor = to
	node.SendVote()
}
func (node *Node) DetermineNextMiner() string {
	fmt.Println("determining next miner")
	keys := make([]string, len(node.Nodes))

	i := 0
	for k := range node.Nodes {
		keys[i] = k
		i++
	}

	if len(node.Nodes) < 5 {
		nn := rand.Intn(len(keys))
		return keys[nn]
	} else {

		nMiner := keys[rand.Intn(len(keys))]
		maxCount := 0
		voteCount := make(map[string]int)
		for key := range node.Votes {
			voteCount[node.Votes[key]] += 1
			if voteCount[node.Votes[key]] > maxCount {
				maxCount = voteCount[node.Votes[key]]
				nMiner = node.Votes[key]
			}
		}
		node.Votes = make(map[string]string)
		return nMiner
	}
}
func (node *Node) SendNodes(address string) {
	fmt.Println("sending nodes")
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := gob.NewEncoder(conn)

	err = encoder.Encode("ReceiveNodes," + node.Address)
	err = encoder.Encode(node.Nodes)
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("sending nodes")
}
func (node *Node) ReceiveNodes( /*conn net.Conn*/ decoder *gob.Decoder) {
	fmt.Println("receiving nodes")

	var nodes map[string]bool
	//fmt.Println("decode started in receive nodes")
	//decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&nodes)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("decode successful in client receive nodes")
	for n := range nodes {
		if n == node.Address {
			continue
		}
		//randomly initializing nodes own vote when nothing decided
		if (rand.Intn(999) == 12) && (node.VoteFor == "") {
			node.VoteFor = n
		}
		node.AddNode(n)
	}
	fmt.Println("receiving nodes completed")
}
func (node *Node) AddNode(address string) {
	if node.Nodes == nil {
		node.Nodes = make(map[string]bool)
	}
	node.Nodes[address] = true
	//fmt.Println("no of nodes inside function: ")
	//fmt.Println(len(node.Nodes))

}
func (node *Node) FloodNodes() {
	fmt.Println("flooding nodes")
	for currentNode := range node.Nodes {
		node.SendNodes(currentNode)
	}
	fmt.Println("flooding nodes completed")
}
func (node *Node) FloodBlockChain() {
	fmt.Println("flooding blockchain")
	for currentNode := range node.Nodes {
		node.SendBlockChain(currentNode)
	}
	fmt.Println("flooding blockchain completed")
}
func (node *Node) HandleConnections(conn net.Conn) {
	decoder := gob.NewDecoder(conn)
	//fmt.Println("decode started in client handle connection")
	var message string
	err := decoder.Decode(&message)
	if err != nil {
		fmt.Println("error from client handle connection")
		fmt.Println(err)
		return
	}
	messageContents := strings.Split(message, ",")
	action := messageContents[0]
	sender := messageContents[1]
	//fmt.Println("decode successful in client handle connection. message: " + message)
	switch action {
	case "ReceiveNodes":
		node.ReceiveNodes(decoder)
	case "ReceiveBlockChain":
		node.ReceiveBlockChain(decoder)
	case "ReceiveTransaction":
		node.ReceiveTransaction(decoder)
	case "ReceiveVote":
		node.ReceiveVote(decoder, sender)
	case "AddNode":
		node.AddNode(sender)
		node.SendBlockChain(sender)
		node.FloodNodes()
	}

}

func (node *Node) ListenConnections() {
	ln, err := net.Listen("tcp", node.Address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error from client listen connection")
			log.Println(err)
			return
		}
		//fmt.Println("connection accepted successfully")
		go node.HandleConnections(conn)

	}

}

func (node *Node) HandleConnectionsSatoshi(conn net.Conn, channel chan bool) {
	decoder := gob.NewDecoder(conn)

	var message string
	err := decoder.Decode(&message)
	if err != nil {
		fmt.Println("error from client handle connection")
		fmt.Println(err)
		return
	}
	messageContents := strings.Split(message, ",")
	action := messageContents[0]
	if action != "AddNode" {
		fmt.Println("Incorrect Acton")
		return
	}
	address := messageContents[1]
	fmt.Println("connected to: " + address)
	node.AddNode(address)
	node.MineBlock(Transaction{
		Amount:   10,
		Sender:   "early joiners reward",
		Receiver: address,
	})
	//fmt.Println("no of nodes: ")
	//fmt.Println(len(node.Nodes))
	if len(node.Nodes) < 4 {
		channel <- false
	} else {
		channel <- true
	}
	return
}

func (node *Node) ListenConnectionsSatoshi() {

	ln, err := net.Listen("tcp", node.Address)
	if err != nil {
		log.Fatal(err)
	}
	channel := make(chan bool)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go node.HandleConnectionsSatoshi(conn, channel)
		if <-channel {
			break
		}
	}
	_ = ln.Close()

}
