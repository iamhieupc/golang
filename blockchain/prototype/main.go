package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

// the difficult of algo
const targetBits = 16
const maxNonce = math.MaxInt64

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

type BlockChain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)                  // 0x0000.....1
	target.Lsh(target, uint(256-targetBits)) // 0x10000000000000000000000000000  => 29bytes
	pow := &ProofOfWork{b, target}

	return pow
}

func IntToHex(number int64) []byte {
	hexString := fmt.Sprintf("%x", number)
	return []byte(hexString) // 0x...
}

func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,       // 0x..
			pow.block.Data,                // [12 231 21 21 21 21]
			IntToHex(pow.block.Timestamp), // 0x...
			IntToHex(int64(targetBits)),   // 0x...
			IntToHex(int64(nonce)),        // 0x...
		},
		[]byte{}, // []
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		// fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")

	return nonce, hash[:]

}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValis := hashInt.Cmp(pow.target) == -1

	return isValis
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return &block
}

// add new block to blockchain
func (bc *BlockChain) AddBlock(data string) {
	// get len of blockchain
	length := len(bc.blocks)

	// get prev block
	prevBlock := bc.blocks[length-1]

	block := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, block)
}

// the first block
func NewGenesisBlock() *Block {
	return NewBlock("hieupc blockchain initial", []byte{})
}

func NewBlockchain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

// func GetBlockChain() {}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to 0xA566f5q34q3674637abc")
	bc.AddBlock("Send 2 more BTC to 0xA566f5q34q3674637abc")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
