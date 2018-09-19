package blc

type TXInput struct {
	TxHash    []byte //交易hash
	Vout      int    //存储TXOutput在里面的索引
	ScriptSig string //用户名
}

func (txInput *TXInput) UnLockWithAddress(address string) bool {
	return txInput.ScriptSig == address
}
