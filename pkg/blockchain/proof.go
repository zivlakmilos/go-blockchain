package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(255-Difficulty))

	return &ProofOfWork{
		Block:  block,
		Target: target,
	}
}

func (p *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.Block.PrevHash,
			p.Block.Data,
			utils.ToHex(int64(nonce)),
			utils.ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (p *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [sha256.Size]byte

	nonce := 0

	for ; nonce < math.MaxInt64; nonce++ {
		data := p.PrepareData(nonce)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(p.Target) < 0 {
			break
		}
		nonce++
	}

	return nonce, hash[:]
}

func (p *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := p.PrepareData(p.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(p.Target) < 0
}
