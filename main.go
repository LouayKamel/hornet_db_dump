package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gohornet/hornet/pkg/model/hornet"
	"github.com/gohornet/hornet/pkg/model/milestone"
	"github.com/iotaledger/iota.go/transaction"
	"github.com/iotaledger/iota.go/trinary"
)

type Dump struct {
	TxHash            trinary.Trytes
	TrunkHash         trinary.Trytes
	BranchHash        trinary.Trytes
	BundleHash        trinary.Trytes
	confirmationIndex milestone.Index
	IsSolid           bool
	IsConfirmed       bool
	IsConflicting     uint8
	IsHead            bool
	IsTail            bool
	IsValue           bool
	Trytes            trinary.Trytes
}

var (
	db                       *bolt.DB
	err                      error
	transactionStorageBucket = []byte{1}
	metadataStorageBucket    = []byte{2}
	totalCount, successCount = 0, 0
)

func main() {
	cfgDbPath := flag.String("dbPath", "", "directory that contains the tangle.db")
	cfgOutputFile := flag.String("output", "output.txt", "output file to store the dump")
	flag.Parse()
	dbPath := *cfgDbPath
	outputFile := *cfgOutputFile

	if dbPath == "" {
		log.Fatal("dbPath is required")
	}
	db, err = bolt.Open(dbPath+"/tangle.db", 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Panic(err)
	}

	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf(fmt.Errorf("cannot open file: %w", err).Error())
	}
	defer f.Close()

	baremetal(f)
}

func baremetal(f *os.File) {
	_ = forEachInBucket(transactionStorageBucket, func(k, v []byte) error {
		totalCount++
		tx := hornet.NewTransaction(k)
		_, err = tx.UnmarshalObjectStorageValue(v)
		if err != nil {
			log.Printf("error unmarshaling tx: %s", err.Error())
			return nil
		}

		bytes, err := findInBucket(metadataStorageBucket, k)
		if err != nil {
			log.Println("error finding metadata: ", err.Error())
			return nil
		}
		if len(bytes) == 0 {
			log.Printf("metada for %s not found\n", tx.GetTxHash().Trytes())
			return nil
		}
		metadataBytes := bytes[0]
		txMetadata := hornet.NewTransactionMetadata(k)
		_, err = txMetadata.UnmarshalObjectStorageValue(metadataBytes)
		if err != nil {
			log.Printf("error unmarshaling metadta: %s", err.Error())
			return nil
		}

		trytes, _ := transaction.TransactionToTrytes(tx.Tx)
		_, confirmationIndex := txMetadata.GetConfirmed()
		var isConflicting uint8
		if txMetadata.IsConflicting() {
			isConflicting = 1
		} else {
			isConflicting = 0
		}

		dump := Dump{
			TxHash:            tx.GetTxHash().Trytes(),
			confirmationIndex: confirmationIndex,
			IsConflicting:     isConflicting,
			Trytes:            trytes,
		}

		line := fmt.Sprintf("%v,%v,%b,%d\n", dump.TxHash, dump.Trytes, isConflicting, confirmationIndex)

		if _, err := f.WriteString(line); err != nil {
			log.Println(fmt.Errorf("err writing to file: %w", err).Error())
			return nil
		}

		log.Println(tx.GetTxHash().Trytes(), " ...done")
		successCount++
		return nil
	})

	log.Println("Total txs: ", totalCount)
	log.Println("Success: ", successCount)
}

func forEachInBucket(bucketName []byte, fn func(_k, _v []byte) error) error {
	return db.View(func(tx *bolt.Tx) error {
		var bk *bolt.Bucket

		//find the bucket
		_ = tx.ForEach(func(nm []byte, b *bolt.Bucket) error {
			if bytes.Equal(nm, bucketName) {
				bk = b
				// stop finding
				return errors.New("")
			}
			return nil
		})
		if bk == nil {
			return fmt.Errorf("bucket not found")
		}

		if bk.Tx().DB() == nil {
			return fmt.Errorf("tx closed")
		}
		c := bk.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := fn(k, v); err != nil {
				return err
			}
		}

		return nil
	})
}

// prefix is a key prefix. e.g if prefix = "123", then you are looking for pairs whos keys are prefixed with "123"
func findInBucket(bucketName []byte, prefix []byte) (res [][]byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		var bk *bolt.Bucket
		//find the bucket
		_ = tx.ForEach(func(nm []byte, b *bolt.Bucket) error {
			if bytes.Equal(nm, bucketName) {
				bk = b
				// stop finding
				return errors.New("")
			}
			return nil
		})
		if bk == nil {
			return fmt.Errorf("bucket not found")
		}

		if bk.Tx().DB() == nil {
			err = fmt.Errorf("tx closed")
			return nil
		}
		c := bk.Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			res = append(res, v)
		}
		return nil

	})
	return
}
