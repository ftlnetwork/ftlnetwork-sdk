package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	"github.com/ftlnetwork/ftlnetwork-sdk/docs/examples/democoin/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/docs/examples/democoin/x/cool"
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/auth"
)

func setGenesis(bapp *DemocoinApp, trend string, accs ...auth.BaseAccount) error {
	genaccs := make([]*types.GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = types.NewGenesisAccount(&types.AppAccount{acc, "foobart"})
	}

	genesisState := types.GenesisState{
		Accounts:    genaccs,
		CoolGenesis: cool.Genesis{trend},
	}

	stateBytes, err := codec.MarshalJSONIndent(bapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.ValidatorUpdate{}
	bapp.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	bapp.Commit()

	return nil
}

func TestGenesis(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	bapp := NewDemocoinApp(logger, db)

	// Construct some genesis bytes to reflect democoin/types/AppAccount
	pk := ed25519.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pk.Address())
	coins, err := sdk.ParseCoins("77foocoin,99barcoin")
	require.Nil(t, err)
	baseAcc := auth.BaseAccount{
		Address: addr,
		Coins:   coins,
	}
	acc := &types.AppAccount{baseAcc, "foobart"}

	err = setGenesis(bapp, "ice-cold", baseAcc)
	require.Nil(t, err)
	// A checkTx context
	ctx := bapp.BaseApp.NewContext(true, abci.Header{})
	res1 := bapp.accountKeeper.GetAccount(ctx, baseAcc.Address)
	require.Equal(t, acc, res1)

	// reload app and ensure the account is still there
	bapp = NewDemocoinApp(logger, db)
	bapp.InitChain(abci.RequestInitChain{AppStateBytes: []byte("{}")})
	ctx = bapp.BaseApp.NewContext(true, abci.Header{})
	res1 = bapp.accountKeeper.GetAccount(ctx, baseAcc.Address)
	require.Equal(t, acc, res1)
}
