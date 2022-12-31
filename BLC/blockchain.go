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

var BlockChainIns *Blockchain

type Blockchain struct {
	lastBlockHash []byte
	db            *bolt.DB
	UTXOS         map[string]map[string]uint64
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
		coinbase := CoinbaseTransaction(minerAddress)
		trans = append(trans, coinbase)
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
			bc.UTXOS[out.OutPubKey] = map[string]uint64{fmt.Sprintf("%s-%v", tx.TxHash, 0): out.Value}
			continue
		}
		for _, txIn := range tx.Vins {
			delete(bc.UTXOS[string(tx.PubKey)], fmt.Sprintf("%s-%v", tx.TxHash, txIn.VoutNum))
		}
		for i, txOut := range tx.Vouts {
			v, ok := bc.UTXOS[txOut.OutPubKey]
			if !ok {
				bc.UTXOS[txOut.OutPubKey] = make(map[string]uint64)
				v = bc.UTXOS[txOut.OutPubKey]
			}
			v[fmt.Sprintf("%s-%v", tx.TxHash, i)] = txOut.Value
		}
	}
}

func (bc *Blockchain) validateTransaction(tx *Transaction) bool {
	return true
}

func InitBlockchain(m Miner) {
	coinbase := CoinbaseTransaction("Genesis address")
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			_, err = tx.CreateBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
	blockChain := &Blockchain{
		lastBlockHash: []byte{},
		db:            db,
	}
	blockChain.AddBlock([]*Transaction{coinbase}, m.PubKey)
	BlockChainIns = blockChain
}
