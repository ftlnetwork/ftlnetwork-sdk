package auth

import (
	"fmt"

	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	CollectedFees sdk.Coins `json:"collected_fees"` // collected fees
	Params        Params    `json:"params"`
}

// Create a new genesis state
func NewGenesisState(collectedFees sdk.Coins, params Params) GenesisState {
	return GenesisState{
		CollectedFees: collectedFees,
		Params:        params,
	}
}

// Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(sdk.Coins{}, DefaultParams())
}

// Init store state from genesis data
func InitGenesis(ctx sdk.Context, ak AccountKeeper, fck FeeCollectionKeeper, data GenesisState) {
	ak.SetParams(ctx, data.Params)
	fck.setCollectedFees(ctx, data.CollectedFees)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak AccountKeeper, fck FeeCollectionKeeper) GenesisState {
	collectedFees := fck.GetCollectedFees(ctx)
	params := ak.GetParams(ctx)

	return NewGenesisState(collectedFees, params)
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.Params.TxSigLimit == 0 {
		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
	}
	if data.Params.SigVerifyCostED25519 == 0 {
		return fmt.Errorf("invalid ED25519 signature verification cost: %d", data.Params.SigVerifyCostED25519)
	}
	if data.Params.SigVerifyCostSecp256k1 == 0 {
		return fmt.Errorf("invalid SECK256k1 signature verification cost: %d", data.Params.SigVerifyCostSecp256k1)
	}
	if data.Params.MaxMemoCharacters == 0 {
		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
	}
	if data.Params.MemoCostPerByte == 0 {
		return fmt.Errorf("invalid memo cost per byte: %d", data.Params.MemoCostPerByte)
	}

	return nil
}
