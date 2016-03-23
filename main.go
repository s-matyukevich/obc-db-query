package main

import (
	"encoding/binary"
	"flag"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/openblockchain/obc-peer/protos"
	"github.com/tecbot/gorocksdb"
)

func main() {
	var dbPath string
	flag.StringVar(&dbPath, "d", "", "Database path")
	flag.Parse()

	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	opts.SetCreateIfMissing(false)

	db, h, err := gorocksdb.OpenDbColumnFamilies(opts, dbPath,
		[]string{"default", "blockchainCF", "indexesCF", "stateDeltaCF", "stateCF"},
		[]*gorocksdb.Options{opts, opts, opts, opts, opts})

	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}

	opt := gorocksdb.NewDefaultReadOptions()
	defer opt.Destroy()
	iterator := db.NewIteratorCF(opt, h[1])

	res := make(map[string]interface{}, 0)
	iterator.SeekToFirst()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		val := iterator.Value()
		keyData := string(append([]byte(nil), key.Data()...))
		if keyData == "blockCount" {
			data := append([]byte(nil), val.Data()...)
			res[keyData] = binary.BigEndian.Uint64(data)
		} else {
			block := &protos.Block{}
			data := append([]byte(nil), val.Data()...)
			err := proto.Unmarshal(data, block)
			if err != nil {
				fmt.Printf(err.Error() + "\n")
				return
			}
			res[keyData] = block
		}
	}
	pretty.Printf("%# v\n", res)
}
