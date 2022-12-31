package BLC

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Miner struct {
	PubKey      string
	unpackTrans []*Transaction
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
	for _, tx := range trans {
		if tx.Verify() {
			panic("invalid transaction")
		}
	}
	BlockChainIns.AddBlock(trans, m.PubKey)
	m.unpackTrans = m.unpackTrans[len(trans):]
}

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Hash  []byte
	Left  *MerkleNode
	Right *MerkleNode
}

func (m *Miner) MakeMerkleTree([]*Transaction) *MerkleTree {
	return nil
}
