package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"math/big"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (pow *ProofOfWork) NewNonce(nonce int) (error, []byte) {
	err, nonceHex := ToHex(int64(nonce))
	if err != nil {
		log.Println("[Err][PoW::NewNonce] - While hex The Nonce integer", err)
		return err, nil
	}

	err, difficultyHex := ToHex(int64(Difficulty))
	if err != nil {
		log.Println("[Err][PoW::NewNonce] - While hex the difficulty integer", err)
		return err, nil
	}

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			nonceHex,
			difficultyHex,
		},
		[]byte{},
	)

	return nil, data
}

func (pow *ProofOfWork) Exec() (error, int, []byte) {
	// log.Println(fmt.Sprintf("[PoW::Exec] - Start PoW for the Block: %d", pow.Block.Id))
	var intHash big.Int
	var hash [32]byte

	resolved := false
	nonce := 0

	for nonce < math.MaxInt64 {
		err, data := pow.NewNonce(nonce)
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
		return errors.New("Cant mine this block!"), 0, nil
	}

	// log.Println(fmt.Sprintf("[PoW::Exec] - Finished PoW for the Block: %d", pow.Block.Id))
	return nil, nonce, hash[:]
}

func (pow *ProofOfWork) Validate() (error, bool) {
	var intHash big.Int

	err, data := pow.NewNonce(pow.Block.Nonce)
	if err != nil {
		return err, false
	}

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return nil, intHash.Cmp(pow.Target) == -1
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
