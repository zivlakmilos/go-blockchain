package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func NewBlock(data string, prevHash []byte) *Block {
	b := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
	}
	b.DeriveHash()

	return b
}

func Genesis() *Block {
	return NewBlock("Genesis", []byte{})
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func (b *Block) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Previous Hash: %x\n", b.PrevHash))
	builder.WriteString(fmt.Sprintf("Data in Block: %s\n", b.Data))
	builder.WriteString(fmt.Sprintf("Hash: %x\n", b.Hash))

	return builder.String()
}
