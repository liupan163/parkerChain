package blc

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	//1高度
	Height int64
	//2上个hash
	PrevBlockHash []byte
	//3交易数据
	Txs []*Transaction
	//4时间戳
	TimeStamp int64
	//5hash
	Hash []byte
	//6、Nonce
	Nonce int64
}

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Txs {
		txHashes = append(txHashes, txHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (block *Block) SetHash() {
	//
	heightBytes := IntToHex(block.Height)
	//时间戳转[]byte
	timeString := strconv.FormatInt(block.TimeStamp, 2)

	timeBytes := []byte(timeString)
	//拼接所有数据
	blockBytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timeBytes, block.Hash}, []byte{})
	//生产hash
	hash := sha256.Sum256(blockBytes)

	block.Hash = hash[:]
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//1、new Block
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {

	block := &Block{height, prevBlockHash, []byte(data), time.Now().Unix(), nil, 0}
	//block.SetHash()

	pow := NewProofOfWork(block)

	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//2、单独写一个方法，生成创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0})
}
