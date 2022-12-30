package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"strconv"
)

const (
	BlockSize = 10  // transactions number
)

type Block struct {
	Height       uint64
	Timestamp    int64
	PreBlockHash []byte
	Trans        []*Transaction
	Hash         []byte
	Nounce       uint64
}

func (b *Block) TransHash() []byte {
	var txHashes [][]byte
	for _, tx := range b.Trans {
		txHashes = append(txHashes, tx.TxHash)
	}
	hash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return hash[:]
}

func (b *Block) GenHash() {
	heightBytes := uint64ToBytes(b.Height)
	timeString := strconv.FormatInt(b.Timestamp, 2)
	timeBytes := []byte(timeString)
	blockBytes := bytes.Join([][]byte{heightBytes, b.PreBlockHash, timeBytes, b.TransHash()}, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	return &block, err
}

func (b *Block) ToBytes(nounce uint64) []byte {
	heightBytes := uint64ToBytes(b.Height)
	nounceBytes := uint64ToBytes(nounce)
	timeString := strconv.FormatInt(b.Timestamp, 2)
	timeBytes := []byte(timeString)
	blockBytes := bytes.Join([][]byte{heightBytes, b.PreBlockHash, timeBytes, b.TransHash(), nounceBytes}, []byte{})
	return blockBytes
}

func NewBlock(trans []*Transaction, height uint64, preBlockHash []byte) *Block {
	block := &Block{
		Trans:        trans,
		Height:       height,
		PreBlockHash: preBlockHash,
	}
	hash, nounce := (&ProofOfWork{block}).Run()
	block.Hash = hash
	block.Nounce = nounce
	return block
}

func uint64ToBytes(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func CreateGenesisBlock(trans []*Transaction) *Block {
	block := NewBlock(trans, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	return block
}
