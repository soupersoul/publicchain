package BLC

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	dbName                  = "blocktrain.db"
	bucket                  = "blocks"
	BlockchainReward uint64 = 10
)

type Blockchain struct {
	lastBlockHash []byte
	db            *bolt.DB
	utxos         map[string]map[string]uint64
}

func (b *Blockchain) Close() {
	b.db.Close()
}

func (bc *Blockchain) AddBlock(trans []*Transaction, minerAddress string) {
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Fatal("failed to get bucket")
		}
		var height uint64
		if len(bc.lastBlockHash) > 0 {
			lastBlockData := b.Get(bc.lastBlockHash)
			lastBlock, _ := DeserializeBlock(lastBlockData)
			height = lastBlock.Height + 1
		}

		newBlock := NewBlock(trans, height, bc.lastBlockHash)
		newBlockData, _ := newBlock.Serialize()
		if err := b.Put(newBlock.Hash, newBlockData); err != nil {
			log.Fatal(err)
		}
		bc.lastBlockHash = newBlock.Hash
		return nil
	})

	for _, tx := range trans {
		if tx.IsCoinbase() {
			out := tx.Vouts[0]
			bc.utxos[out.outPubKey] = map[string]uint64{fmt.Sprintf("%s-%v", tx.TxHash, 0): out.Value}
			continue
		}
		for _, txIn := range tx.Vins {
			delete(bc.utxos[string(tx.PubKey)], fmt.Sprintf("%s-%v", tx.TxHash, txIn.VoutNum))
		}
		for i, txOut := range tx.Vouts {
			v, ok := bc.utxos[txOut.outPubKey]
			if !ok {
				bc.utxos[txOut.outPubKey] = make(map[string]uint64)
				v = bc.utxos[txOut.outPubKey]
			}
			v[fmt.Sprintf("%s-%v", tx.TxHash, i)] = txOut.Value
		}
	}
}
