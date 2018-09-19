package blc

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (blockChainIterator *BlockChainIterator) Next() *Block {
	var block *Block
	err := blockChainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			currentBlockByte := b.Get(blockChainIterator.CurrentHash)
			block := DeserializeBlock(currentBlockByte)
			blockChainIterator.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}