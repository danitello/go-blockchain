package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danitello/go-blockchain/common/errutil"
)

const walletFile = "./tmp/wallets.dat"

// Wallets keeps track of all current Wallet structs
type Wallets struct {
	Wallets map[string]*Wallet
}

// InitWallets makes a new Wallets struct and loads it with previous Wallets data if possible
func InitWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

// CreateWallet makes a new wallet and adds it to the Wallets
func (ws *Wallets) CreateWallet() string {
	wallet := InitWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

// GetAddresses retrieves all of the address from the Wallets
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet retrieves a specific wallet by address
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// LoadFromFile loads Wallets data from disk
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	data, err := ioutil.ReadFile(walletFile)
	errutil.Handle(err)

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&wallets)
	errutil.Handle(err)

	ws.Wallets = wallets.Wallets

	return nil
}

// SaveToFile writes the Wallets data to disk
func (ws *Wallets) SaveToFile() {
	var data bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(ws)
	errutil.Handle(err)

	err = ioutil.WriteFile(walletFile, data.Bytes(), 0644)
	errutil.Handle(err)
}
