package blc

func (cli *CLI) getBalance(address string) {
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	//txOutputs :=blockchain.UnUTXOs(address)
	amount := blockchain.GetBalance(address)
}
