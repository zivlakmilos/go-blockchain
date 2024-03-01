package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/zivlakmilos/go-blockchain/pkg/utils"
	"github.com/zivlakmilos/go-blockchain/pkg/wallet"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	inputs := []TxInput{}
	outputs := []TxOutput{}

	wallets, err := wallet.NewWallets()
	utils.HandleError(err)
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs := chain.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("error: not enough fund")
	}

	for key, value := range validOutputs {
		txID, err := hex.DecodeString(key)
		utils.HandleError(err)

		for _, out := range value {
			inputs = append(inputs, TxInput{
				ID:        txID,
				Out:       out,
				Signature: nil,
				PubKey:    pubKeyHash,
			})
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to))

	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	tx := &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
	tx.ID = tx.Hash()
	chain.SignTransaction(tx, w.PrivateKey)

	return tx
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{
		ID:        []byte{},
		Out:       -1,
		Signature: nil,
		PubKey:    []byte(data),
	}

	txout := NewTXOutput(100, to)

	tx := &Transaction{
		ID:      []byte{},
		Inputs:  []TxInput{txin},
		Outputs: []TxOutput{*txout},
	}
	tx.GenerateID()

	return tx
}

func (t *Transaction) GenerateID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(t)
	utils.HandleError(err)

	hash = sha256.Sum256(encoded.Bytes())
	t.ID = hash[:]
}

func (t *Transaction) IsCoinbase() bool {
	return len(t.Inputs) == 1 && len(t.Inputs[0].ID) == 0 && t.Inputs[0].Out == -1
}

func (t *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(t)
	utils.HandleError(err)

	return encoded.Bytes()
}

func (t *Transaction) Hash() []byte {
	txCopy := *t
	txCopy.ID = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (t *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if t.IsCoinbase() {
		return
	}

	for _, txIn := range t.Inputs {
		if prevTXs[hex.EncodeToString(txIn.ID)].ID == nil {
			log.Panic("error: previous transaction does not exists")
		}
	}

	txCopy := t.TrimmedCopy()

	for idx, txIn := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(txIn.ID)]
		txCopy.Inputs[idx].Signature = nil
		txCopy.Inputs[idx].PubKey = prevTX.Outputs[txIn.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[idx].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		utils.HandleError(err)
		signature := append(r.Bytes(), s.Bytes()...)

		t.Inputs[idx].Signature = signature
	}
}

func (t *Transaction) TrimmedCopy() Transaction {
	inputs := []TxInput{}
	outputs := []TxOutput{}

	for _, txIn := range t.Inputs {
		inputs = append(inputs, TxInput{
			ID:        txIn.ID,
			Out:       txIn.Out,
			Signature: nil,
			PubKey:    nil,
		})
	}

	for _, txOut := range t.Outputs {
		outputs = append(outputs, TxOutput{
			Value:      txOut.Value,
			PubKeyHash: txOut.PubKeyHash,
		})
	}

	txCopy := Transaction{
		ID:      t.ID,
		Inputs:  inputs,
		Outputs: outputs,
	}

	return txCopy
}

func (t *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if t.IsCoinbase() {
		return true
	}

	for _, txIn := range t.Inputs {
		if prevTXs[hex.EncodeToString(txIn.ID)].ID == nil {
			log.Panic("error: previous transaction does not exists")
		}
	}

	txCopy := t.TrimmedCopy()
	curve := elliptic.P256()

	for idx, txIn := range t.Inputs {
		prevTX := prevTXs[hex.EncodeToString(txIn.ID)]
		txCopy.Inputs[idx].Signature = nil
		txCopy.Inputs[idx].PubKey = prevTX.Outputs[txIn.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[idx].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(txIn.Signature)
		r.SetBytes(txIn.Signature[:(sigLen / 2)])
		s.SetBytes(txIn.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(txIn.PubKey)
		x.SetBytes(txIn.Signature[:(keyLen / 2)])
		y.SetBytes(txIn.Signature[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     &x,
			Y:     &y,
		}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false
		}
	}

	return true
}

func (t *Transaction) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("-- Transaction %x:\n", t.ID))

	for idx, txIn := range t.Inputs {
		builder.WriteString(fmt.Sprintf("  Input %d:\n", idx))
		builder.WriteString(fmt.Sprintf("    TXID:      %x:\n", txIn.ID))
		builder.WriteString(fmt.Sprintf("    Out:       %d:\n", txIn.Out))
		builder.WriteString(fmt.Sprintf("    Signature: %x:\n", txIn.Signature))
		builder.WriteString(fmt.Sprintf("    PubKey:    %x:\n", txIn.PubKey))
	}

	for idx, txOut := range t.Outputs {
		builder.WriteString(fmt.Sprintf("  Output %d:\n", idx))
		builder.WriteString(fmt.Sprintf("    Value:      %d:\n", txOut.Value))
		builder.WriteString(fmt.Sprintf("    PubKeyHash: %x:\n", txOut.PubKeyHash))
	}

	return builder.String()
}
