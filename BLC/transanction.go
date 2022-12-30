package BLC

import (
	"bytes"
	"encoding/gob"
	"math/big"
)

type Transaction struct {
	TxHash []byte
	Vins   []*TxInput
	Vouts  []*TxOutput
	PubKey []byte
}

type TxInput struct {
	TxHash    []byte
	VoutNum   uint64
	ScriptSig string
}

type TxOutput struct {
	Value     uint64
	outPubKey string
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
	txInput := &TxInput{[]byte{}, 0, "Genesis Data"}
	txOutput := &TxOutput{10, address}
	trans := &Transaction{Vins: []*TxInput{txInput}, Vouts: []*TxOutput{txOutput}, PubKey: []byte(address)}
	hashBytes, _ := TransanctionHash(trans)
	trans.TxHash = hashBytes
	return trans
}

func NewTransanction(from string, to string, amount uint64) {
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	var total uint64 = 100
	input := &TxInput{[]byte("hash"), total, from}
	txInputs = append(txInputs, input)

	output := &TxOutput{amount, to}
	outputLeft := &TxOutput{total - amount, from}
	txOutputs = append(txOutputs, output)
	txOutputs = append(txOutputs, outputLeft)
	trans := &Transaction{Vins: txInputs, Vouts: txOutputs, PubKey: []byte(from)}
	hashBytes, _ := TransanctionHash(trans)
	trans.TxHash = hashBytes
}

func (tx *Transaction) IsCoinbase() bool {
	var hashInt big.Int
	hashInt.SetBytes(tx.Vins[0].TxHash)
	return hashInt.Cmp(big.NewInt(0)) == 1
}
