package main

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

func (pow *ProofOfWork) InitNonce(nonce int) []byte {
	_, nonceHex := ToHex(int64(nonce))
	_, difficultyHex := ToHex(int64(Difficulty))

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			nonceHex,
			difficultyHex,
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Exec() (error, int, []byte) {
	log.Println(fmt.Sprintf("[PoW::Exec] - Start PoW for the Block: %d", pow.Block.Id))
	var intHash big.Int
	var hash [32]byte

	resolved := false
	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitNonce(nonce)
		hash = sha256.Sum256(data)

		log.Println(fmt.Sprintf("[PoW::Exec] - Hash: %x", hash))
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			resolved = true
			break
		} else {
			nonce++
		}

	}

	if !resolved {
		return errors.New("Cant mine this block!"), 0, nil
	}

	log.Println(fmt.Sprintf("[PoW::Exec] - Finished PoW for the Block: %d", pow.Block.Id))
	return nil, nonce, hash[:]
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func ToHex(num int64) (error, []byte) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return err, nil
	}
	return nil, buff.Bytes()
}
