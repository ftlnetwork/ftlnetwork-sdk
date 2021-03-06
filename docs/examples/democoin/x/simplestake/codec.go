package simplestake

import (
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBond{}, "simplestake/BondMsg", nil)
	cdc.RegisterConcrete(MsgUnbond{}, "simplestake/UnbondMsg", nil)
}
