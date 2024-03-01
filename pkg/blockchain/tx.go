package blockchain

import (
	"bytes"

	"github.com/zivlakmilos/go-blockchain/pkg/utils"
	"github.com/zivlakmilos/go-blockchain/pkg/wallet"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func NewTXOutput(value int, address string) *TxOutput {
	o := &TxOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	o.Lock([]byte(address))

	return o
}

func (i *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(i.PubKey)
	return bytes.Equal(pubKeyHash, lockingHash)
}

func (o *TxOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	o.PubKeyHash = pubKeyHash
}

func (o *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(o.PubKeyHash, pubKeyHash)
}
