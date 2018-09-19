package blc

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
	"strconv"
	"encoding/hex"
)

// 数据库名字
const dbName = "blockchain.db"

// 表的名字
const blockTableName = "blocks"

type BlockChain struct {
	Tip []byte //最新区块的hash
	DB  *bolt.DB
}

//func (blc *BlockChain) AddBlockToBlockChain(data string, height int64, prevHash []byte) {
//	newBlock := NewBlock(data, height, prevHash)
//	blc.Blocks = append(blc.Blocks, newBlock)
//}

func (blockchain *BlockChain) Iterator() *BlockChainIterator {

	return &BlockChainIterator{blockchain.Tip, blockchain.DB}
}

func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

func (blc *BlockChain) PrintChain() {
	var block *Block
	var currentHash []byte = blc.Tip
	for {
		err := blc.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			if b != nil {
				blockByte := b.Get(currentHash)
				block = DeserializeBlock(blockByte)
				fmt.Printf("Height：%d\n", block.Height)
				fmt.Printf("PrevBlockHash：%d\n", block.PrevBlockHash)
				fmt.Printf("Data：%d\n", block.Data)
				fmt.Printf("TimeStamp：%d\n", block.TimeStamp)
				fmt.Printf("TimeStamp：%d\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
				fmt.Printf("Hash：%d\n", block.Hash)
				fmt.Printf("Nonce：%d\n", block.Nonce)
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
		currentHash = block.PrevBlockHash
	}
}

func (blockchain *BlockChain) GetBalance(address string) int64 {
	utxos := blockchain.UnUTXOs(address)
	var amount int64
	for _, out := range utxos {
		amount = amount + out.Value
	}
	return amount
}

//地址对应的区块上的交易是未花费状态的话--未花交易，添加到数组返回
func (blockchain *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][] int)

	for _, tx := range txs {

		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否能够解锁
				if in.UnLockWithAddress(address) {

					key := hex.EncodeToString(in.TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	}

	for _, tx := range txs {
	work1:
		for index, out := range tx.Vouts {
			if out.UnLockScriptPubKeyWithAddress(address) {
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						{
							if hash == txHashStr {
								var isUnSpentUTXO bool
								for _, outIndex := range indexArray {
									if index == outIndex {
										isUnSpentUTXO = true
										continue work1
									}
									if isUnSpentUTXO == false {
										utxo := &UTXO{tx.TxHash, index, out}
										unUTXOs = append(unUTXOs, utxo)
									}
								}
							} else {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						}
					}
				}
			}
		}
	}
	blockIterator := blockchain.Iterator()
	for {
		block := blockIterator.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					if in.UnLockWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}
		work:
			for index, out := range tx.Vouts {
				if out.UnLockScriptPubKeyWithAddress(address) {
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == i && hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unUTXOs
}

//转账时查找可用的UTXO
func (blc *BlockChain) FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	//1.现获取所有的UTXO
	utxos := blc.UnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)

	//2.遍历utxos
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) { //钱不够
		os.Exit(1)
	}

	return value, spendableUTXO
}

func (blc *BlockChain) AddBlockToBlackchain(data string) {
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockBytes := b.Get(blc.Tip)
			block := DeserializeBlock(blockBytes)
			newBlock := NewBlock(data, block.Height+1, block.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("1"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func CreateBlockChainWithGenesisBlock(data string) *BlockChain {

	if DBExists() {
		fmt.Println("创世区块已经有了")
		os.Exit(1)
	}
	fmt.Println("正在创建数据库")

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var genesisHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		//创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}
		if b == nil {
			genesisBlock := CreateGenesisBlock(data)
			err := b.Put([]byte("1"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			//存储最新的区块的hash
			err = b.Put([]byte("1"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.Hash
		}
		return nil
	})

	return &BlockChain{genesisHash, db}
}

func (blc *BlockChain) AddBlockToBlockchain(data string) {
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blockTableName))
		//创建新区块
		if b != nil {
			blockBytes := b.Get(blc.Tip)
			//反序列化
			block := DeserializeBlock(blockBytes)
			//加到数据库中
			newBlock := NewBlock(data, block.Height+1, block.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//更新数据库
			err = b.Put([]byte("1"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//更新数据库
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func BlockChainObject() *BlockChain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	return &BlockChain{tip, db}
}

func (blockchain *BlockChain) MineNewBlock(from []string, to []string, amount []string) {

	value, _ := strconv.Atoi(amount[0])
	tx := NewSimpleTransaction(from[0], to[0], value)

	var txs []*Transaction
	txs = append(txs, tx)

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("1"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	//建立新区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("1"), block.Hash)
			blockchain.Tip = block.Hash
		}
		return nil
	})

}
