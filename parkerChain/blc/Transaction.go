package blc

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

type Transaction struct {
	TxHash []byte // 交易hash
	Vins   []*TXInput
	Vouts  []*TXOutput
}

func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

func NewCoinbaseTransaction(address string) *Transaction {
	//代表消费
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TXOutput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	txCoinbase.HashTransaction() //设置hash
}

func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash
}

func NewSimpleTransaction(from string, to string, amount int, blockchain *BlockChain, txs []*Transaction) *Transaction {

	money, spendableUTXODic := blockchain.FindSpendableUTXOS(from, amount, txs)

	var txInputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic{
		txHashBytes,_ :=hex.DecodeString(txHash)
		for _,index :=range indexArray{
			txInput := &TXOutput{bytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}


	//转账
	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput = &TXOutput{int64(money) - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	tx.HashTransaction()
	return tx
}
