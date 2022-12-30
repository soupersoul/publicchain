package BLC

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"github.com/boltdb/bolt"
)

type Miner struct {
	PubKey      string
	unpackTrans []*Transaction
	blockChain  *Blockchain
}

var Difficulty uint = 8

type ProofOfWork struct {
	Block *Block
	// Target *big.Int
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-Difficulty)
	var nounce uint64 = 0
	var hash [32]byte
	for {
		blockBytes := pow.Block.ToBytes(nounce)
		hash = sha256.Sum256(blockBytes)
		fmt.Printf("\r%X", hash)
		hashInt := (&big.Int{}).SetBytes(hash[:])
		if hashInt.Cmp(target) < 0 {
			break
		}
		nounce++
	}
	return hash[:], nounce
}

func (m *Miner) receiveTransaction(tx *Transaction) {
	m.unpackTrans = append(m.unpackTrans, tx)
}

func (m *Miner) mineBlock() {
	var trans []*Transaction
	if len(m.unpackTrans) >= BlockSize {
		trans = m.unpackTrans[0:BlockSize]
	} else {
		trans = m.unpackTrans[:len(m.unpackTrans)-1]
	}
	m.blockChain.AddBlock(trans, m.PubKey)
	m.unpackTrans = m.unpackTrans[len(trans):]
}

func (m *Miner) NewBlockchain() *Blockchain {
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
	m.blockChain = blockChain
	return blockChain
}
