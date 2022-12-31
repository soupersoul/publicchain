package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"math/big"
)

type Transaction struct {
	TxHash    []byte
	Vins      []*TxInput
	Vouts     []*TxOutput
	PubKey    []byte
	ScriptSig string
}

type TxInput struct {
	TxHash  []byte
	VoutNum uint64
	PubKey  []byte
}

type TxOutput struct {
	Value     uint64
	OutPubKey string
}

func TransanctionHash(tx *Transaction) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(tx); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func CoinbaseTransaction(address string) *Transaction {
	txInput := &TxInput{TxHash: []byte{}, VoutNum: 0, PubKey: []byte{}}
	txOutput := &TxOutput{10, address}
	trans := &Transaction{Vins: []*TxInput{txInput}, Vouts: []*TxOutput{txOutput}, PubKey: []byte(address)}
	hashBytes, _ := TransanctionHash(trans)
	trans.TxHash = hashBytes
	return trans
}

/*
func NewTransanction(from string, to string, amount uint64) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	var total uint64 = 100
	input := &TxInput{TxHash: []byte("hash"), VoutNum: total, ScriptSig: from}
	txInputs = append(txInputs, input)

	output := &TxOutput{amount, to}
	outputLeft := &TxOutput{total - amount, from}
	txOutputs = append(txOutputs, output)
	txOutputs = append(txOutputs, outputLeft)
	trans := &Transaction{Vins: txInputs, Vouts: txOutputs, PubKey: []byte(from)}
	hashBytes, _ := TransanctionHash(trans)
	trans.TxHash = hashBytes
	return trans
}
*/

func (tx *Transaction) IsCoinbase() bool {
	var hashInt big.Int
	hashInt.SetBytes(tx.Vins[0].TxHash)
	return hashInt.Cmp(big.NewInt(0)) == 1
}

func (tx *Transaction) Verify() bool {
	r, s, x, y := big.Int{}, big.Int{}, big.Int{}, big.Int{}
	sigLen := len(tx.ScriptSig)
	r.SetBytes([]byte(tx.ScriptSig)[:(sigLen / 2)])
	s.SetBytes([]byte(tx.ScriptSig)[(sigLen / 2):])
	keyLen := len(tx.PubKey)
	x.SetBytes(tx.PubKey[:(keyLen / 2)])
	y.SetBytes(tx.PubKey[(keyLen / 2):])
	rawPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
	// TODO: check tx hash
	return ecdsa.Verify(&rawPubKey, tx.TxHash, &r, &s)
}
