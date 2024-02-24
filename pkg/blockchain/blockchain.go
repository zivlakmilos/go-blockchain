package blockchain

import (
	"log"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

const dbPath = "/tmp/blocks"

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func NewBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	utils.HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			log.Printf("No existing blockchain found")
			genesis := Genesis()
			log.Printf("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			utils.HandleError(err)

			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		}

		item, err := txn.Get([]byte("lh"))
		utils.HandleError(err)
		lastHash, err = item.ValueCopy(lastHash)
		return err
	})
	utils.HandleError(err)

	return &BlockChain{
		LastHash: lastHash,
		Database: db,
	}
}

func (c *BlockChain) AddBlock(data string) {
	block := NewBlock(data, c.LastHash)

	err := c.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		utils.HandleError(err)
		err = txn.Set([]byte("lh"), block.Hash)
		return err
	})
	utils.HandleError(err)

	c.LastHash = block.Hash
}

func (c *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		CurrentHash:  c.LastHash,
		CurrentBlock: nil,
		Database:     c.Database,
	}
}

func (c *BlockChain) String() string {
	var builder strings.Builder

	iter := c.Iterator()
	for iter.Next() {
		block := iter.Value()
		builder.WriteString(block.String())
		builder.WriteString("\n")
	}

	return builder.String()
}

type BlockChainIterator struct {
	CurrentBlock *Block
	Database     *badger.DB
	CurrentHash  []byte
}

func (i *BlockChainIterator) Next() bool {
	var block *Block

	if i.CurrentHash == nil || len(i.CurrentHash) == 0 {
		return false
	}

	err := i.Database.View(func(txn *badger.Txn) error {
		var encodedBlock []byte

		item, err := txn.Get(i.CurrentHash)
		utils.HandleError(err)
		encodedBlock, err = item.ValueCopy(encodedBlock)
		block = Deserialize(encodedBlock)

		return err
	})
	utils.HandleError(err)

	i.CurrentHash = block.PrevHash
	i.CurrentBlock = block

	return true
}

func (i BlockChainIterator) Value() *Block {
	return i.CurrentBlock
}
