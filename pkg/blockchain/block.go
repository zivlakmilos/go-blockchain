package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := &Block{
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
	}

	p := NewProofOfWork(b)
	nonce, hash := p.Run()

	b.Hash = hash
	b.Nonce = nonce

	return b
}

func Genesis(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var hash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	hash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return hash[:]
}

func (b *Block) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Previous Hash: %x\n", b.PrevHash))
	builder.WriteString(fmt.Sprintf("Hash: %x\n", b.Hash))

	p := NewProofOfWork(b)
	builder.WriteString(fmt.Sprintf("PoW: %v\n", p.Validate()))

	return builder.String()
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)

	utils.HandleError(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	utils.HandleError(err)

	return &block
}
