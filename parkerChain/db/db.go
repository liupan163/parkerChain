package db

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"github.com/scryinfo/parkerChain/blc"
)

func main() {
	//建库db
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//建表，更新数据库
	//err = db.Update(func(tx *bolt.Tx) error {
	//	//建新表
	//	b, err := tx.CreateBucket([]byte("BlockBucket"))
	//	//获取表单
	//	//b:=tx.Bucket([]byte("BlockBucket"))
	//	if err != nil {
	//		return fmt.Errorf("create bucket: %s", err)
	//	}
	//	//往表里面添加数据
	//	if b != nil {
	//		err := b.Put([]byte("1"), block.Serialize())
	//		if err != nil {
	//			log.Panic("数据存储失败")
	//		}
	//	}
	//	return nil
	//})
	//建表失败
	//if err != nil {
	//	log.Panic(err)
	//}

	// 查看数据
	err = db.View(func(tx *bolt.Tx) error {
		// 获取BlockBucket表对象
		b := tx.Bucket([]byte("BlockBucket"))
		// 往表里面存储数据
		if b != nil {
			blockData := b.Get([]byte("l"))
			fmt.Printf("%s\n", blockData)
			deSerializeBlock := blc.Deserialize(blockData)
			fmt.Printf("%s\n", deSerializeBlock)
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})
	//更新失败
	if err != nil {
		log.Panic(err)
	}

}
