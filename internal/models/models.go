package models 

import (

) 

type AddressRequest struct {
	WalletAddress string `json:"walletAddress"`
	ChainName     string `json:"chainName"`
	CryptoName    string `json:"cryptoName"`
	Tag           string `json:"tag"`
}

type Address struct {
	ID            uint64 `json:"id"`
	WalletAddress string `json:"wallet_address"`
	ChainName     string `json:"chain_name"`
	CryptoName    string `json:"crypto_name"`
	Tag           string `json:"tag"`
	Balance       int64  `json:"balance"`
}

type TagUpdateRequest struct {
	ID  uint64 `json:"id"`
	Tag string `json:"tag"`
}