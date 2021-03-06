package oracle

import (
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
)

// GetInfoKey returns the key for OracleInfo
func GetInfoKey(p Payload, cdc *codec.Codec) []byte {
	bz := cdc.MustMarshalBinaryLengthPrefixed(p)
	return append([]byte{0x00}, bz...)
}

// GetSignPrefix returns the prefix for signs
func GetSignPrefix(p Payload, cdc *codec.Codec) []byte {
	bz := cdc.MustMarshalBinaryLengthPrefixed(p)
	return append([]byte{0x01}, bz...)
}

// GetSignKey returns the key for sign
func GetSignKey(p Payload, signer sdk.AccAddress, cdc *codec.Codec) []byte {
	return append(GetSignPrefix(p, cdc), signer...)
}
