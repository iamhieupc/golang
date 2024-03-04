package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

type BlockChain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return &block
}

// add new block to blockchain
func (bc *BlockChain) AddBlock(data string) {
	length := len(bc.blocks)
	prevBlock := bc.blocks[length-1]
	block := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, block)
}

// the first block
func NewGenesisBlock() *Block {
	return NewBlock("hieupc blockchain initial", []byte{})
}

func NewBlockchain() *BlockChain {
	fmt.Println("test: ", BlockChain{[]*Block{NewGenesisBlock()}})
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

// func GetBlockChain() {}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
