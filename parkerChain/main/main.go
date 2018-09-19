package main

import (
	"github.com/scryinfo/parkerChain/blc"
	"fmt"
)

func main() {
	//block := blc.NewBlock("Genenis Block", 1, []byte{0, 0, 0})
	//genesisBlock := blc.CreateGenesisBlock("Genenis Block")
	//blockChain := blc.CreateBlockChainWithGenesisBlock()
	//fmt.Println(blockChain)
	//fmt.Println(blockChain.Blocks)
	//fmt.Println(blockChain.Blocks[0])
	//
	//blockChain.AddBlockToBlockChain("send 100 2liuapn",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	//blockChain.AddBlockToBlockChain("send 100 2parker",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	//blockChain.AddBlockToBlockChain("send 100",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)
	//blockChain.AddBlockToBlockChain("send 100",blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,blockChain.Blocks[len(blockChain.Blocks)-1].Hash)

	block := blc.NewBlock("TEST", 1, []byte{0, 0, 0, 0})
	fmt.Println("block.hash", block.Hash)
	fmt.Println("block.hash", block.Nonce)
	//proofOfWork := blc.NewProofOfWork(block)

	bytes := block.Serialize()
	block = blc.DeserializeBlock(bytes)

	blockChain := blc.CreateBlockChainWithGenesisBlock()
	defer blockChain.DB.Close()

	//新区块
	blockChain.AddBlockToBlackchain("Send S.T to me")
	blockChain.AddBlockToBlackchain("Send S.T to her")
}
