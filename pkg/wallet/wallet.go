package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"

	"github.com/zivlakmilos/go-blockchain/pkg/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksupLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type PrivateKey struct {
	D          *big.Int
	PublicKeyX *big.Int
	PublicKeyY *big.Int
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}

	return &wallet
}

func (w *Wallet) Address() []byte {
	publichHash := publicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, publichHash...)
	checksum := checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := utils.Base58Encode(fullHash)

	return address
}

func (w *Wallet) GobEncode() ([]byte, error) {
	privateKey := &PrivateKey{
		D:          w.PrivateKey.D,
		PublicKeyX: w.PrivateKey.X,
		PublicKeyY: w.PrivateKey.Y,
	}

	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(privateKey)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(w.PublicKey)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (w *Wallet) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	var privateKey PrivateKey

	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&privateKey)
	if err != nil {
		return err
	}

	w.PrivateKey = ecdsa.PrivateKey{
		D: privateKey.D,
		PublicKey: ecdsa.PublicKey{
			X: privateKey.PublicKeyX,
			Y: privateKey.PublicKeyY,
		},
	}

	w.PublicKey = make([]byte, buf.Len())
	_, err = buf.Read(w.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	public := append(private.X.Bytes(), private.Y.Bytes()...)
	return *private, public
}

func publicKeyHash(publicKey []byte) []byte {
	publicHash := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(publicHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksupLength]
}
