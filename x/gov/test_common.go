package gov

import (
	"bytes"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/bank"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/mock"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/stake"
	stakeTypes "github.com/ftlnetwork/ftlnetwork-sdk/x/stake/types"
)

// initialize the mock application for this module
func getMockApp(t *testing.T, numGenAccs int) (*mock.App, Keeper, stake.Keeper, []sdk.AccAddress, []crypto.PubKey, []crypto.PrivKey) {
	mapp := mock.NewApp()

	stake.RegisterCodec(mapp.Cdc)
	RegisterCodec(mapp.Cdc)

	keyStake := sdk.NewKVStoreKey(stake.StoreKey)
	tkeyStake := sdk.NewTransientStoreKey(stake.TStoreKey)
	keyGov := sdk.NewKVStoreKey(StoreKey)

	pk := mapp.ParamsKeeper
	ck := bank.NewBaseKeeper(mapp.AccountKeeper)
	sk := stake.NewKeeper(mapp.Cdc, keyStake, tkeyStake, ck, pk.Subspace(stake.DefaultParamspace), stake.DefaultCodespace)
	keeper := NewKeeper(mapp.Cdc, keyGov, pk, pk.Subspace("testgov"), ck, sk, DefaultCodespace)

	mapp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mapp.QueryRouter().AddRoute(QuerierRoute, NewQuerier(keeper))

	mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper, sk))

	require.NoError(t, mapp.CompleteSetup(keyStake, tkeyStake, keyGov))

	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(numGenAccs, sdk.Coins{sdk.NewInt64Coin(stakeTypes.DefaultBondDenom, 42)})

	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, sk, addrs, pubKeys, privKeys
}

// gov and stake endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		tags := EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// gov and stake initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakeKeeper stake.Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakeGenesis := stake.DefaultGenesisState()
		stakeGenesis.Pool.LooseTokens = sdk.NewDec(100000)

		validators, err := stake.InitGenesis(ctx, stakeKeeper, stakeGenesis)
		if err != nil {
			panic(err)
		}
		InitGenesis(ctx, keeper, DefaultGenesisState())
		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

// TODO: Remove once address interface has been implemented (ref: #2186)
func SortValAddresses(addrs []sdk.ValAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}

	SortByteArrays(byteAddrs)

	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// Sorts Addresses
func SortAddresses(addrs []sdk.AccAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}
	SortByteArrays(byteAddrs)
	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// Public
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}
