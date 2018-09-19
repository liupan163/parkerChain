package blc

func (cli *CLI)createGenesisBlockChain(data string)  {
	blockchain := CreateBlockChainWithGenesisBlock(data)
	defer blockchain.DB.Close()
}
