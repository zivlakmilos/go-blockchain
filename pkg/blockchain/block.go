package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func NewBlock(data string, prevHash []byte) *Block {
	b := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
	}

	p := NewProofOfWork(b)
	nonce, hash := p.Run()

	b.Hash = hash
	b.Nonce = nonce

	return b
}

func Genesis() *Block {
	return NewBlock("Genesis", []byte{})
}

func (b *Block) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Previous Hash: %x\n", b.PrevHash))
	builder.WriteString(fmt.Sprintf("Data in Block: %s\n", b.Data))
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
