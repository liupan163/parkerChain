package blc

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const targetBit = 20

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (proof *ProofOfWork) IsValid() bool {
	//1. proofOfWork.Block.Hash
	//2.proof.Target
	var hashInt big.Int
	hashInt.SetBytes(proof.Block.Hash)
	if proof.Target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

func (proof *ProofOfWork) Run() ([]byte, int64) {
	nonce := 0
	var hashInt big.Int
	var hash [32]byte
	for {
		dataBytes := proof.prepareData(nonce)
		hash = sha256.Sum256(dataBytes)

		hashInt.SetBytes(hash[:])

		//难度小鱼
		if proof.Target.Cmp(&hashInt) == 1 {
			fmt.Println("hashInt", hashInt)
			break
		}
		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}

func NewProofOfWork(block *Block) *ProofOfWork {

	//1.big.Int对象 1
	// 2
	//0000 0001
	// 8 - 2 = 6
	// 0100 0000  64
	// 0010 0000
	// 0000 0000 0000 0001 0000 0000 0000 0000 0000 0000 .... 0000

	//1. 创建一个初始值为1的target
	target := big.NewInt(1)
	//2. 左移256 - targetBit
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.HashTransactions(),
		IntToHex(pow.Block.TimeStamp),
		IntToHex(int64(targetBit)),
		IntToHex(int64(nonce)),
		IntToHex(int64(pow.Block.Height)),
	}, []byte{})
	return data
}
