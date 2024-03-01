package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"os"
)

const walletFile = "/tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	w := &Wallets{
		Wallets: map[string]*Wallet{},
	}

	err := w.LoadFile()

	return w, err
}

func (w *Wallets) GetAllAddresses() []string {
	addresses := []string{}

	for address := range w.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (w *Wallets) GetWallet(address string) Wallet {
	return *w.Wallets[address]
}

func (w *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := string(wallet.Address())

	w.Wallets[address] = wallet

	return address
}

func (w *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (w *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	w.Wallets = wallets.Wallets

	return nil
}
