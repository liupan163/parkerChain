package blc

import "os"

func (cli *CLI) printchain() {
	if DBExists() == false {
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()
	blockchain.PrintChain()
}
