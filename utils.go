package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBlockChainAuthAndClient gets blocal chain auth and client
func GetBlockChainAuthAndClient(RPCAddr string, PrivateKey string, chainID int64, value int64) (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(RPCAddr)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		return nil, nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, nil, err
	}

	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, err
	}

	var auth *bind.TransactOpts
	if chainID > 0 {
		auth, _ = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	} else {
		auth = bind.NewKeyedTransactor(privateKey)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(value) // when deploy contract. this should not set
	auth.GasLimit = uint64(5000000) // in units
	auth.GasPrice = price

	// Deploy a new awesome contract for the binding demo
	return client, auth, nil
}
