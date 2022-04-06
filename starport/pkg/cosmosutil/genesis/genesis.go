package genesis

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	jsoniter "github.com/json-iterator/go"
	"io"
	"os"
)

const (
	genesisFilename = "genesis.json"
)

type (
	// Genesis represents a more readable version of the stargate genesis file
	Genesis struct {
		Accounts   []string
		StakeDenom string
	}
	// ChainGenesis represents the stargate genesis file
	ChainGenesis struct {
		ChainID     string `json:"chain_id"`
		GenesisTime string `json:"genesis_time"`
		AppState    struct {
			Auth struct {
				Accounts []struct {
					Address string `json:"address"`
				} `json:"accounts"`
			} `json:"auth"`
			Staking struct {
				Params struct {
					BondDenom string `json:"bond_denom"`
				} `json:"params"`
			} `json:"staking"`
		} `json:"app_state"`
	}
)

// CheckGenesisContainsAddress returns true if the address exist into the genesis file
func CheckGenesisContainsAddress(genesisPath, addr string) (bool, error) {
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	genesis, err := FromPath(genesisPath)
	if err != nil {
		return false, err
	}
	return genesis.HasAccount(addr), nil
}

// HasAccount check if account exist into the genesis account
func (g Genesis) HasAccount(address string) bool {
	for _, account := range g.Accounts {
		if account == address {
			return true
		}
	}
	return false
}

// HasAccount check if account exist into the genesis account
func (g *GenReader) HasAccount(address string) bool {
	genesis, err := g.Genesis()
	if err != nil {
		return false
	}
	for _, account := range genesis.Accounts {
		if account == address {
			return true
		}
	}
	return false
}

// StakeDenom returns the chain genesis stake denom
func (g *GenReader) StakeDenom() (string, error) {
	var genesis ChainGenesis
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return genesis.AppState.Staking.Params.BondDenom, json.NewDecoder(g).Decode(&genesis)
}

// ChainGenesis returns the chain genesis form the reader
func (g *GenReader) ChainGenesis() (genesis ChainGenesis, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(g).Decode(&genesis)
	return
}

// Genesis returns the genesis wrapper form the reader
func (g *GenReader) Genesis() (Genesis, error) {
	chainGenesis, err := g.ChainGenesis()
	if err != nil {
		return Genesis{}, err
	}
	accounts := make([]string, len(chainGenesis.AppState.Auth.Accounts))
	for i, acc := range chainGenesis.AppState.Auth.Accounts {
		accounts[i] = acc.Address
	}
	return Genesis{
		StakeDenom: chainGenesis.AppState.Staking.Params.BondDenom,
		Accounts:   accounts,
	}, nil
}

func (g *GenReader) Hash() (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, g); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (g *GenReader) String() (string, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, g); err != nil {
		return "", err
	}
	return buf.String(), nil
}