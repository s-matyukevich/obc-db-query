package main

import (
	"flag"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/openblockchain/obc-peer/protos"
	"github.com/tecbot/gorocksdb"
)

func main() {
	dbmap := map[string]func([]byte) (interface{}, error){
		"blockchainCF": func(b []byte) (interface{}, error) {
			res := &protos.Block{}
			err := proto.Unmarshal(b, res)
			return res, err
		},
	}
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
		fmt.Pritf(err.Error())
		return
	}

	opt := gorocksdb.NewDefaultReadOptions()
	defer opt.Destroy()
	iterator := db.NewIteratorCF(opt, cfHandler)

	res := make(map[string]interface{}, 0)
	for ; iterator.IsValid(); iterator.IsValid() {
		val := iterator.Value()
		key := iterator.Key()
		block := &protos.Block{}
		err := proto.Unmarshal(val, block)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		res = append(res, block)
	}
	pretty.Printf("%# v\n", res)
}
