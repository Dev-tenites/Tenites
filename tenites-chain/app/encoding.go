package app

import (
        "fmt"

        gogoproto "github.com/cosmos/gogoproto/proto"
        "google.golang.org/protobuf/proto"
        "google.golang.org/protobuf/reflect/protoreflect"

        "github.com/cosmos/cosmos-sdk/client"
        "github.com/cosmos/cosmos-sdk/codec"
        "github.com/cosmos/cosmos-sdk/codec/types"
        "github.com/cosmos/cosmos-sdk/std"
        "github.com/cosmos/cosmos-sdk/x/auth/tx"
        authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
        authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
        "cosmossdk.io/x/tx/signing"
)

type EncodingConfigType struct {
        InterfaceRegistry types.InterfaceRegistry
        Codec             codec.Codec
        TxConfig          client.TxConfig
        Amino             *codec.LegacyAmino
        AccountRetriever  client.AccountRetriever
}

func tscSignerGetter(fieldName string) signing.GetSignersFunc {
        return func(msg proto.Message) ([][]byte, error) {
                ref := msg.ProtoReflect()
                fd := ref.Descriptor().Fields().ByName(protoreflect.Name(fieldName))
                if fd == nil {
                        return nil, fmt.Errorf("field %s not found in %s", fieldName, ref.Descriptor().FullName())
                }
                signer := ref.Get(fd).String()
                if signer == "" {
                        return nil, fmt.Errorf("empty signer field %s in %s", fieldName, ref.Descriptor().FullName())
                }
                addrCodec := authcodec.NewBech32Codec(Bech32PrefixAccAddr)
                addrBytes, err := addrCodec.StringToBytes(signer)
                if err != nil {
                        return nil, fmt.Errorf("invalid signer address %s: %w", signer, err)
                }
                return [][]byte{addrBytes}, nil
        }
}

func MakeEncodingConfig() EncodingConfigType {
        amino := codec.NewLegacyAmino()

        signingOpts := signing.Options{
                AddressCodec:          authcodec.NewBech32Codec(Bech32PrefixAccAddr),
                ValidatorAddressCodec: authcodec.NewBech32Codec(Bech32PrefixValAddr),
        }

        tscMsgTypes := []protoreflect.FullName{
                "tenites.tsc.v1.MsgMint",
                "tenites.tsc.v1.MsgBurn",
                "tenites.tsc.v1.MsgTransfer",
                "tenites.tsc.v1.MsgAddMintAuthority",
                "tenites.tsc.v1.MsgRemoveMintAuthority",
                "tenites.tsc.v1.MsgAddBurnAuthority",
                "tenites.tsc.v1.MsgRemoveBurnAuthority",
        }
        for _, mt := range tscMsgTypes {
                signingOpts.DefineCustomGetSigners(mt, tscSignerGetter("signer"))
        }

        interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
                ProtoFiles:     gogoproto.HybridResolver,
                SigningOptions: signingOpts,
        })
        if err != nil {
                panic(err)
        }
        cdc := codec.NewProtoCodec(interfaceRegistry)
        txCfg := tx.NewTxConfig(cdc, tx.DefaultSignModes)

        std.RegisterLegacyAminoCodec(amino)
        std.RegisterInterfaces(interfaceRegistry)
        ModuleBasics.RegisterLegacyAminoCodec(amino)
        ModuleBasics.RegisterInterfaces(interfaceRegistry)

        return EncodingConfigType{
                InterfaceRegistry: interfaceRegistry,
                Codec:             cdc,
                TxConfig:          txCfg,
                Amino:             amino,
                AccountRetriever:  authtypes.AccountRetriever{},
        }
}

func EncodingConfig() EncodingConfigType {
        return MakeEncodingConfig()
}
