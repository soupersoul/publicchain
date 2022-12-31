package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/soupersoul/publicchain/BLC"
)

const (
	version         = byte(0x00)
	addressChecksum = 4
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PubKey     []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	return &Wallet{
		PrivateKey: private,
		PubKey:     public,
	}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

type UTXO struct {
	TxHash  string
	VoutNum uint64
	Amount  uint64
}

func (w *Wallet) Pay(toAddr []byte, money uint64) {
	blockChain := BLC.BlockChainIns
	var utxos []UTXO
	for txo, amount := range blockChain.UTXOS[string(w.PubKey)] {
		txs := strings.Split(txo, "-")
		voutNum, _ := strconv.Atoi(txs[1])
		utxos = append(utxos, UTXO{TxHash: txs[0], VoutNum: uint64(voutNum), Amount: amount})
	}
	sort.Slice(utxos, func(d, e int) bool {
		return utxos[d].Amount < utxos[e].Amount
	})
	var total uint64
	var txInputs []*BLC.TxInput
	for _, utxo := range utxos {
		total += utxo.Amount
		if total >= money {
			break
		}
		txInputs = append(txInputs, &BLC.TxInput{
			TxHash:  []byte(utxo.TxHash),
			VoutNum: utxo.VoutNum,
			PubKey:  w.PubKey,
		})
	}
	if total < money {
		panic("insufficient funds")
	}
	var txOutputs []*BLC.TxOutput
	output := &BLC.TxOutput{Value: money, OutPubKey: string(toAddr)}
	txOutputs = append(txOutputs, output)
	if total > money {
		outputLeft := &BLC.TxOutput{Value: total - money, OutPubKey: string(toAddr)}
		txOutputs = append(txOutputs, outputLeft)
	}

	trans := &BLC.Transaction{Vins: txInputs, Vouts: txOutputs, PubKey: []byte(w.PubKey)}
	hashBytes, _ := BLC.TransanctionHash(trans)
	trans.TxHash = hashBytes
	r, s, _ := ecdsa.Sign(rand.Reader, &w.PrivateKey, trans.TxHash)
	signature := append(r.Bytes(), s.Bytes()...)
	trans.ScriptSig = string(signature)
	publishTransaction(trans)
}

func publishTransaction(tx *BLC.Transaction) {

}

/*
func (w *Wallet) GetAddress() []byte {
	// raw pubKey -> sha256 -> hash160   => 20 Byte
	// version + hash160
	// sha256 twice -> last 4 byte (checksum)
	// version + hash160 + checksum (25byte) -> base58
	pubKeyHash := HashPubKey(w.PubKey)
	payload := append([]byte{version}, pubKeyHash...)
	checkSum := checksum(payload)
	fullPayload := append(payload, checkSum...)
	encoding := base58.FlickrEncoding
	address, err := encoding.Encode(fullPayload)
	if err != nil {
		panic(err)
	}
	return address
}

func HashPubKey(pubKey []byte) []byte {
	pubSha256 := sha256.Sum256(pubKey)
	hasher := ripemd160.New()
	hasher.Write(pubSha256[:])
	bytes := hasher.Sum(nil)
	return bytes
}

func ValidateAddress(address string) bool {
	encoding := base58.FlickrEncoding
	pubKeyHash, err := encoding.Decode([]byte(address))
	if err != nil {
		panic(err)
	}
	actChecksum := pubKeyHash[len(pubKeyHash)-addressChecksum:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksum]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Equal(actChecksum, targetChecksum)
}

func checksum(data []byte) []byte {
	bytes1 := sha256.Sum256(data)
	bytes2 := sha256.Sum256(bytes1[:])
	return bytes2[:addressChecksum]
}
*/
