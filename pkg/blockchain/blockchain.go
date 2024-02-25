package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

const (
	dbPath      = "/tmp/blocks"
	dbFile      = "/tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func NewBlockChain(address string) *BlockChain {
	if dbExists() {
		fmt.Printf("Blockchain already exists")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	utils.HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			cbtx := CoinbaseTx(address, genesisData)
			genesis := Genesis(cbtx)
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

func OpenBlockChain(address string) *BlockChain {
	if !dbExists() {
		fmt.Printf("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	utils.HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
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

func (c *BlockChain) AddBlock(txs []*Transaction) {
	block := NewBlock(txs, c.LastHash)

	err := c.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		utils.HandleError(err)
		err = txn.Set([]byte("lh"), block.Hash)
		return err
	})
	utils.HandleError(err)

	c.LastHash = block.Hash
}

func (c *BlockChain) FindUnspentTransactions(address string) []Transaction {
	unspentTxs := []Transaction{}

	spentTXOs := map[string][]int{}

	iter := c.Iterator()
	for iter.Next() {
		block := iter.Value()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			for outIdx, txo := range tx.Outputs {
				if !txo.CanBeUnlocked(address) {
					continue
				}

				if spentTXOs[txID] != nil {
					found := false
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							found = true
						}
					}
					if found {
						continue
					}
				}

				unspentTxs = append(unspentTxs, *tx)
			}

			if tx.IsCoinbase() {
				continue
			}

			for _, txi := range tx.Inputs {
				if txi.CanUnlock(address) {
					txInID := hex.EncodeToString(txi.ID)
					spentTXOs[txInID] = append(spentTXOs[txInID], txi.Out)
				}
			}
		}
	}

	return unspentTxs
}

func (c *BlockChain) FindUTXO(address string) []TxOutput {
	UTXOs := []TxOutput{}
	unspentTransactions := c.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (c *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := map[string][]int{}
	accumulated := 0

	unspentTxs := c.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {
		if accumulated >= amount {
			break
		}

		txID := hex.EncodeToString(tx.ID)

		for idx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], idx)

				if accumulated >= amount {
					break
				}
			}
		}
	}

	return accumulated, unspentOuts
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

func (c *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		CurrentHash:  c.LastHash,
		CurrentBlock: nil,
		Database:     c.Database,
	}
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

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
