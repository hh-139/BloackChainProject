package assignment02IBC_master

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

const checkString = "000" //difficulty level is 3 i.e. first three characters of hash are zero.

type Block struct {
	Transaction []Transaction
	Number      int
	PrevBlock   *Block
	PrevHash    string
	Hash        string
	Nonce       int
	Votes       map[string]string
	NextMiner   string
}

func CalculateHash(trans []Transaction, prev *Block) (int, string, string) {
	var hh string
	if prev == nil {
		hh = "0"
		i := 63
		for i > 0 {
			hh = hh + "0"
			i -= 1
		}
	} else {
		hh = prev.Hash
	}

	var nonce int
	nonce = 0
	var tranString string
	for _, tran := range trans {
		tranString += tran.Receiver + tran.Sender + strconv.Itoa(tran.Amount)
	}
	s := hh + tranString + strconv.Itoa(nonce)
	h := sha256.New()
	h.Write([]byte(s))
	sha256Hash := hex.EncodeToString(h.Sum(nil))
	for sha256Hash[:len(checkString)] != checkString {
		s = s[:len(s)-len(strconv.Itoa(nonce))]
		nonce += 1
		s = s + strconv.Itoa(nonce)
		h.Write([]byte(s))
		sha256Hash = hex.EncodeToString(h.Sum(nil))
	}
	return nonce, hh, sha256Hash
}

func InsertBlock(transaction []Transaction,nextMiner string,votes map[string]string ,chainHead *Block) *Block {
	nonce, prev, curr := CalculateHash(transaction, chainHead)
	if chainHead != nil {
		chainHead.Number += 1
	}
	//tran []Transaction.Transaction
	//tran = append(tran, transaction)
	BlockChain := Block{
		Transaction: transaction,
		PrevBlock:   chainHead,
		PrevHash:    prev,
		Hash:        curr,
		Nonce:       nonce,
		Votes:       votes,
		NextMiner:   nextMiner,
	}
	return &BlockChain
}

func CalculateAmount(user string, chain Block) int {
	var Amount int
	for chain.PrevBlock != nil {
		for _, tran := range chain.Transaction {
			if tran.Receiver == user {
				Amount += tran.Amount
			} else if tran.Sender == user {
				Amount -= tran.Amount
			}
		}
		chain = *chain.PrevBlock
	}
	for _, tran := range chain.Transaction {
		if tran.Receiver == user {
			Amount += tran.Amount
		} else if tran.Sender == user {
			Amount -= tran.Amount
		}
	}
	return Amount
}

func isEmpty(chain *Block) bool {
	if chain.PrevBlock == nil && chain.Transaction[chain.Number].IsEmpty() {
		return true
	}
	return false
}

func ListBlocks(chain *Block) {
	if isEmpty(chain) {
		println("Blockchain is empty.")
	} else {
		for chain.PrevBlock != nil {
			fmt.Println("Transaction: ", chain.Transaction, "Nonce: ", chain.Nonce,
				"\nPrevious Hash: ", chain.PrevHash, "\nCurrent Hash: ", chain.Hash,"\nNext Miner: ",chain.NextMiner)
			chain = chain.PrevBlock
		}
		fmt.Println("Transaction: ", chain.Transaction, "Nonce: ", chain.Nonce,
			"\nPrevious Hash: ", chain.PrevHash, "\nCurrent Hash: ", chain.Hash,"\nNext Miner: ",chain.NextMiner)
	}
}

func VerifyChain(chainHead *Block) bool {
	for chainHead.PrevBlock != nil {
		if chainHead.PrevHash != chainHead.PrevBlock.Hash {
			println("Hashes doesn't match.")
			println("Hash of previous Block: ", chainHead.PrevBlock.Hash)
			println("Previous hash stored on current Block: ", chainHead.PrevHash)
			return false
		}
		chainHead = chainHead.PrevBlock
	}
	return true
}
