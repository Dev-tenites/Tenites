package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgMint{},
		&MsgBurn{},
		&MsgTransfer{},
		&MsgAddMintAuthority{},
		&MsgRemoveMintAuthority{},
		&MsgAddBurnAuthority{},
		&MsgRemoveBurnAuthority{},
	)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgMint{}, "tsc/MsgMint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "tsc/MsgBurn", nil)
	cdc.RegisterConcrete(&MsgTransfer{}, "tsc/MsgTransfer", nil)
	cdc.RegisterConcrete(&MsgAddMintAuthority{}, "tsc/MsgAddMintAuthority", nil)
	cdc.RegisterConcrete(&MsgRemoveMintAuthority{}, "tsc/MsgRemoveMintAuthority", nil)
	cdc.RegisterConcrete(&MsgAddBurnAuthority{}, "tsc/MsgAddBurnAuthority", nil)
	cdc.RegisterConcrete(&MsgRemoveBurnAuthority{}, "tsc/MsgRemoveBurnAuthority", nil)
}
