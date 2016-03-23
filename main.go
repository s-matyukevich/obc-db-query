package main

import (
	"flag"
	"fmt"
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/openblockchain/obc-peer/protos"
	"github.com/tecbot/gorocksdb"
)

func main() {
	var dbPath string
	flag.StringVar(&dbPath, "db", "", "Database path")
	flag.Parse()

	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	opts.SetCreateIfMissing(false)

	db, h, err := gorocksdb.OpenDbColumnFamilies(opts, dbPath,
		[]string{"blockchainCF"},
		[]*gorocksdb.Options{opts})

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	opt := gorocksdb.NewDefaultReadOptions()
	defer opt.Destroy()
	iterator := db.NewIteratorCF(opt, h[0])

	res := make(map[string]interface{}, 0)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		val := iterator.Value()
		keyData := string(key.Data())
		if keyData == "blockCount"{
			res[keyData] = binary.BigEndian.Uint64(val.Data())
		} else {
			block := &protos.Block{}
			err := proto.Unmarshal(val.Data(), block)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
			res[keyData] = block
		}
	}
	pretty.Printf("%# v\n", res)
}
