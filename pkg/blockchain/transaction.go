package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	inputs := []TxInput{}
	outputs := []TxOutput{}

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("error: not enough fund")
	}

	for key, value := range validOutputs {
		txID, err := hex.DecodeString(key)
		utils.HandleError(err)

		for _, out := range value {
			inputs = append(inputs, TxInput{
				ID:  txID,
				Out: out,
				Sig: from,
			})
		}
	}

	outputs = append(outputs, TxOutput{
		Value:  amount,
		PubKey: to,
	})

	if acc > amount {
		outputs = append(outputs, TxOutput{
			Value:  acc - amount,
			PubKey: from,
		})
	}

	tx := &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
	tx.GenerateID()

	return tx
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{
		ID:  []byte{},
		Out: -1,
		Sig: data,
	}

	txout := TxOutput{
		Value:  100,
		PubKey: to,
	}

	tx := &Transaction{
		ID:      []byte{},
		Inputs:  []TxInput{txin},
		Outputs: []TxOutput{txout},
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
