package blc

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {
	BC *BlockChain
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchainwithgenesis -data -- 交易数据.")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")

}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//abandon
func (cli *CLI) addBlock(txs []*Transaction) {
	if DBExists() == false {
		fmt.Println("db not exist")
		os.Exit(1)
	}
	blockchain := BlockChainObject()
	defer blockchain.DB.Close()

	blockchain.AddBlockToBlackchain(txs)
}

func (cli *CLI) Run() {
	isValidArgs()

	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址......")
	flagTo := sendBlockCmd.String("to", "", "转账源地址......")
	flagAmout := sendBlockCmd.String("amount", "", "转账金额")

	// flagAddBlockData := addBlockCmd.String("data", "http://parkerChain", "TX data......")
	flagCreateBlockChainWithData := createBlockChainCmd.String("data", "http://parkerChain", "TX data......")
	getbalanceWithAddress := getbalanceCmd.String("address", "", "要查询某一个账号的余额......")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		//case "addblock":
		//	err := addBlockCmd.Parse(os.Args[2:])
		//	if err != nil {
		//		log.Panic(err)
		//	}
	case "parintchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmout == "" {
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*flagFrom)
		fmt.Println(*flagTo)
		fmt.Println(*flagAmout)

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmout)
		cli.send(from, to, amount)
	}

	//if addBlockCmd.Parsed() {
	//	if *flagAddBlockData == "" {
	//		printUsage()
	//		os.Exit(1)
	//	}
	//
	//	//fmt.Println(*flagAddBlockData)
	//	cli.addBlock(*flagAddBlockData) //changed
	//}
	if printChainCmd.Parsed() {
		cli.printchain()
	}

	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockChainWithData == "" {
			fmt.Println("交易数据不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(*flagCreateBlockChainWithData)
	}

	if getbalanceCmd.Parsed() {
		if *getbalanceWithAddress == "" {
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceWithAddress)
	}

}