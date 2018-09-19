package blc

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)

//int64转化成字节数组
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

//将json字符串转化成数组
func JSONToArray(jsonString string) []string {
	//json到[]string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
}

// 字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}