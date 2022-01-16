package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (pow *ProofOfWork) NewNonce(nonce int) ([]byte, error) {
	nonceHex, err := ToHex(int64(nonce))
	if err != nil {
		log.Println("[Err][PoW::NewNonce] - While hex The Nonce integer", err)
		return nil, err
	}

	difficultyHex, err := ToHex(int64(Difficulty))
	if err != nil {
		log.Println("[Err][PoW::NewNonce] - While hex the difficulty integer", err)
		return nil, err
	}

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			[]byte(fmt.Sprintf("%d", pow.Block.Timestamp)),
			nonceHex,
			difficultyHex,
		},
		[]byte{},
	)

	return data, nil
}

func (pow *ProofOfWork) Exec() (int, []byte, error) {
	// log.Println(fmt.Sprintf("[PoW::Exec] - Start PoW for the Block: %d", pow.Block.Id))
	var intHash big.Int
	var hash [32]byte

	resolved := false
	nonce := 0

	for nonce < math.MaxInt64 {
		data, err := pow.NewNonce(nonce)
		if err != nil {
			break
		}

		hash = sha256.Sum256(data)

		// log.Println(fmt.Sprintf("[PoW::Exec] - Hash: %x", hash))
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			resolved = true
			break
		} else {
			nonce++
		}

	}

	if !resolved {
		log.Println("[Err][PoW::Exec] - Cant mine this block!")
		return 0, nil, errors.New("cant mine this block")
	}

	// log.Println(fmt.Sprintf("[PoW::Exec] - Finished PoW for the Block: %d", pow.Block.Id))
	return nonce, hash[:], nil
}

func (pow *ProofOfWork) Validate() (bool, error) {
	var intHash big.Int

	data, err := pow.NewNonce(pow.Block.Nonce)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1, nil
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func ToHex(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
