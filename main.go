package main

import (
	"context"
	"log"
	"github.com/mengguang/newchain/ethclient"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])
	return hash
}
func main() {
	conn,err := ethclient.Dial("https://rpc_address_here")
	if err != nil {
		log.Fatal(err)
	}
	blk, err := conn.BlockByNumber(context.Background(),nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(blk.Number())
	log.Println(blk.Hash().Hex())
	log.Printf("%x\n",blk.Extra()[32:])
	hash := sigHash(blk.Header()).Bytes()
	//hash := blk.Hash().Bytes()
	sig := blk.Extra()[32:]
	pubkey, err := crypto.SigToPub(hash , sig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pubkey)
	publicKey := crypto.FromECDSAPub(pubkey)
	if err != nil {
		log.Fatal(err)
	}
	valid := crypto.VerifySignature(publicKey,hash,sig[:64])
	log.Printf("valid: %v\n",valid)
	log.Println(crypto.PubkeyToAddress(*pubkey).Hex())
}
