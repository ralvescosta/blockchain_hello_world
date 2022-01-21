package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey //ecdsa = elliptical curve digital signature algorithm
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil, err
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub, nil
}

func PublicKeyHash(publicKey []byte) ([]byte, error) {
	hashedPublicKey := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	if err != nil {
		return nil, err
	}
	publicRipeMd := hasher.Sum(nil)

	return publicRipeMd, nil
}

func Checksum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func (w *Wallet) Address() ([]byte, error) {
	pubHash, err := PublicKeyHash(w.PublicKey)
	if err != nil {
		return nil, err
	}

	versionedHash := append([]byte{version}, pubHash...)

	checksum := Checksum(versionedHash)

	finalHash := append(versionedHash, checksum...)

	address := base58Encode(finalHash)

	return address, nil
}

func NewWallet() (*Wallet, error) {
	privateKey, publicKey, err := NewKeyPair()
	if err != nil {
		return nil, err
	}

	wallet := Wallet{privateKey, publicKey}

	return &wallet, nil
}
