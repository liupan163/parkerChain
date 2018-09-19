package blc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

const version = byte(0x00)
const addressChecksumLen = 4

func NewWallet() *Wallet {
	privatekey, publicKey := newKeyPaire()
	return &Wallet{privatekey, publicKey}
}

func (w *Wallet) GetAddress() []byte {
	ripemd160Hash := w.Ripemd160Hash(w.PublicKey)
	version_ripemd160Hash := append([]byte{version}, ripemd160Hash...)
	checkSumBytes := CheckSum(version_ripemd160Hash)
	bytes := append(version_ripemd160Hash, checkSumBytes...)
	return Bas58Encode(bytes)

}

func CheckSum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}

func (w *Wallet) Ripemd160Hash(publicKey []byte) []byte {
	//256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	//160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

func newKeyPaire() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//...切片打散被传入
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
