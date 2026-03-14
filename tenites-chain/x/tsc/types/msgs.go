package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgMint                = "tsc_mint"
	TypeMsgBurn                = "tsc_burn"
	TypeMsgTransfer            = "tsc_transfer"
	TypeMsgAddMintAuthority    = "tsc_add_mint_authority"
	TypeMsgRemoveMintAuthority = "tsc_remove_mint_authority"
	TypeMsgAddBurnAuthority    = "tsc_add_burn_authority"
	TypeMsgRemoveBurnAuthority = "tsc_remove_burn_authority"
)

type MsgMint struct {
	Signer       string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer"`
	AuthorityId  string `protobuf:"bytes,2,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
	WalletId     string `protobuf:"bytes,3,opt,name=wallet_id,json=walletId,proto3" json:"wallet_id"`
	Amount       string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount"`
	Purpose      string `protobuf:"bytes,5,opt,name=purpose,proto3" json:"purpose"`
	Jurisdiction string `protobuf:"bytes,6,opt,name=jurisdiction,proto3" json:"jurisdiction"`
}

type MsgMintResponse struct {
	OperationId string `protobuf:"bytes,1,opt,name=operation_id,json=operationId,proto3" json:"operation_id"`
}

func (m *MsgMint) ProtoMessage()             {}
func (m *MsgMint) Reset()                    {}
func (m *MsgMint) String() string            { return TypeMsgMint }
func (m *MsgMint) XXX_MessageName() string   { return "tenites.tsc.v1.MsgMint" }
func (m *MsgMint) Route() string             { return ModuleName }
func (m *MsgMint) Type() string              { return TypeMsgMint }
func (m *MsgMint) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}
func (m *MsgMint) ValidateBasic() error {
	if m.AuthorityId == "" {
		return ErrInvalidAuthority.Wrap("authority_id required")
	}
	if m.Amount == "" || m.Amount == "0" {
		return ErrInvalidAmount.Wrap("amount required")
	}
	if m.WalletId == "" {
		return ErrWalletNotFound.Wrap("wallet_id required")
	}
	return nil
}

func (m *MsgMintResponse) ProtoMessage()           {}
func (m *MsgMintResponse) Reset()                  {}
func (m *MsgMintResponse) String() string          { return "MsgMintResponse" }
func (m *MsgMintResponse) XXX_MessageName() string { return "tenites.tsc.v1.MsgMintResponse" }

type MsgBurn struct {
	Signer       string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer"`
	AuthorityId  string `protobuf:"bytes,2,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
	WalletId     string `protobuf:"bytes,3,opt,name=wallet_id,json=walletId,proto3" json:"wallet_id"`
	Amount       string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount"`
	Reason       string `protobuf:"bytes,5,opt,name=reason,proto3" json:"reason"`
	Jurisdiction string `protobuf:"bytes,6,opt,name=jurisdiction,proto3" json:"jurisdiction"`
}

type MsgBurnResponse struct{}

func (m *MsgBurn) ProtoMessage()             {}
func (m *MsgBurn) Reset()                    {}
func (m *MsgBurn) String() string            { return TypeMsgBurn }
func (m *MsgBurn) XXX_MessageName() string   { return "tenites.tsc.v1.MsgBurn" }
func (m *MsgBurn) Route() string             { return ModuleName }
func (m *MsgBurn) Type() string              { return TypeMsgBurn }
func (m *MsgBurn) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}
func (m *MsgBurn) ValidateBasic() error {
	if m.AuthorityId == "" {
		return ErrInvalidAuthority.Wrap("authority_id required")
	}
	if m.Amount == "" || m.Amount == "0" {
		return ErrInvalidAmount.Wrap("amount required")
	}
	return nil
}

func (m *MsgBurnResponse) ProtoMessage()           {}
func (m *MsgBurnResponse) Reset()                  {}
func (m *MsgBurnResponse) String() string          { return "MsgBurnResponse" }
func (m *MsgBurnResponse) XXX_MessageName() string { return "tenites.tsc.v1.MsgBurnResponse" }

type MsgTransfer struct {
	Signer       string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer"`
	FromWalletId string `protobuf:"bytes,2,opt,name=from_wallet_id,json=fromWalletId,proto3" json:"from_wallet_id"`
	ToWalletId   string `protobuf:"bytes,3,opt,name=to_wallet_id,json=toWalletId,proto3" json:"to_wallet_id"`
	Amount       string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount"`
	Purpose      string `protobuf:"bytes,5,opt,name=purpose,proto3" json:"purpose"`
}

type MsgTransferResponse struct{}

func (m *MsgTransfer) ProtoMessage()             {}
func (m *MsgTransfer) Reset()                    {}
func (m *MsgTransfer) String() string            { return TypeMsgTransfer }
func (m *MsgTransfer) XXX_MessageName() string   { return "tenites.tsc.v1.MsgTransfer" }
func (m *MsgTransfer) Route() string             { return ModuleName }
func (m *MsgTransfer) Type() string              { return TypeMsgTransfer }
func (m *MsgTransfer) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}
func (m *MsgTransfer) ValidateBasic() error {
	if m.FromWalletId == "" || m.ToWalletId == "" {
		return ErrWalletNotFound.Wrap("from and to wallet_id required")
	}
	if m.Amount == "" || m.Amount == "0" {
		return ErrInvalidAmount.Wrap("amount required")
	}
	return nil
}

func (m *MsgTransferResponse) ProtoMessage()           {}
func (m *MsgTransferResponse) Reset()                  {}
func (m *MsgTransferResponse) String() string          { return "MsgTransferResponse" }
func (m *MsgTransferResponse) XXX_MessageName() string { return "tenites.tsc.v1.MsgTransferResponse" }

type MsgAddMintAuthority struct {
	Signer          string   `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer"`
	TenitesId       string   `protobuf:"bytes,2,opt,name=tenites_id,json=tenitesId,proto3" json:"tenites_id"`
	Name            string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`
	DailyLimit      string   `protobuf:"bytes,4,opt,name=daily_limit,json=dailyLimit,proto3" json:"daily_limit"`
	SingleTxLimit   string   `protobuf:"bytes,5,opt,name=single_tx_limit,json=singleTxLimit,proto3" json:"single_tx_limit"`
	AllowedPurposes []string `protobuf:"bytes,6,rep,name=allowed_purposes,json=allowedPurposes,proto3" json:"allowed_purposes"`
	Jurisdictions   []string `protobuf:"bytes,7,rep,name=jurisdictions,proto3" json:"jurisdictions"`
	ExpiresAt       int64    `protobuf:"varint,8,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at"`
}

type MsgAddMintAuthorityResponse struct {
	AuthorityId string `protobuf:"bytes,1,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
}

func (m *MsgAddMintAuthority) ProtoMessage()             {}
func (m *MsgAddMintAuthority) Reset()                    {}
func (m *MsgAddMintAuthority) String() string            { return TypeMsgAddMintAuthority }
func (m *MsgAddMintAuthority) XXX_MessageName() string   { return "tenites.tsc.v1.MsgAddMintAuthority" }
func (m *MsgAddMintAuthority) Route() string             { return ModuleName }
func (m *MsgAddMintAuthority) Type() string              { return TypeMsgAddMintAuthority }
func (m *MsgAddMintAuthority) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}
func (m *MsgAddMintAuthority) ValidateBasic() error {
	if m.TenitesId == "" {
		return ErrInvalidAuthority.Wrap("tenites_id required")
	}
	if m.Name == "" {
		return ErrInvalidAuthority.Wrap("name required")
	}
	return nil
}

func (m *MsgAddMintAuthorityResponse) ProtoMessage()           {}
func (m *MsgAddMintAuthorityResponse) Reset()                  {}
func (m *MsgAddMintAuthorityResponse) String() string          { return "MsgAddMintAuthorityResponse" }
func (m *MsgAddMintAuthorityResponse) XXX_MessageName() string { return "tenites.tsc.v1.MsgAddMintAuthorityResponse" }

type MsgRemoveMintAuthority struct {
	Signer      string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer"`
	AuthorityId string `protobuf:"bytes,2,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
}

type MsgRemoveMintAuthorityResponse struct{}

func (m *MsgRemoveMintAuthority) ProtoMessage()             {}
func (m *MsgRemoveMintAuthority) Reset()                    {}
func (m *MsgRemoveMintAuthority) String() string            { return TypeMsgRemoveMintAuthority }
func (m *MsgRemoveMintAuthority) XXX_MessageName() string   { return "tenites.tsc.v1.MsgRemoveMintAuthority" }
func (m *MsgRemoveMintAuthority) Route() string             { return ModuleName }
func (m *MsgRemoveMintAuthority) Type() string              { return TypeMsgRemoveMintAuthority }
func (m *MsgRemoveMintAuthority) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}
func (m *MsgRemoveMintAuthority) ValidateBasic() error {
	if m.AuthorityId == "" {
		return ErrInvalidAuthority.Wrap("authority_id required")
	}
	return nil
}

func (m *MsgRemoveMintAuthorityResponse) ProtoMessage()           {}
func (m *MsgRemoveMintAuthorityResponse) Reset()                  {}
func (m *MsgRemoveMintAuthorityResponse) String() string          { return "MsgRemoveMintAuthorityResponse" }
func (m *MsgRemoveMintAuthorityResponse) XXX_MessageName() string { return "tenites.tsc.v1.MsgRemoveMintAuthorityResponse" }

type MsgAddBurnAuthority struct {
	Signer         string   `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer"`
	TenitesId      string   `protobuf:"bytes,2,opt,name=tenites_id,json=tenitesId,proto3" json:"tenites_id"`
	Name           string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`
	DailyLimit     string   `protobuf:"bytes,4,opt,name=daily_limit,json=dailyLimit,proto3" json:"daily_limit"`
	SingleTxLimit  string   `protobuf:"bytes,5,opt,name=single_tx_limit,json=singleTxLimit,proto3" json:"single_tx_limit"`
	AllowedReasons []string `protobuf:"bytes,6,rep,name=allowed_reasons,json=allowedReasons,proto3" json:"allowed_reasons"`
	Jurisdictions  []string `protobuf:"bytes,7,rep,name=jurisdictions,proto3" json:"jurisdictions"`
	ExpiresAt      int64    `protobuf:"varint,8,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at"`
}

type MsgAddBurnAuthorityResponse struct {
	AuthorityId string `protobuf:"bytes,1,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
}

func (m *MsgAddBurnAuthority) ProtoMessage()             {}
func (m *MsgAddBurnAuthority) Reset()                    {}
func (m *MsgAddBurnAuthority) String() string            { return TypeMsgAddBurnAuthority }
func (m *MsgAddBurnAuthority) XXX_MessageName() string   { return "tenites.tsc.v1.MsgAddBurnAuthority" }
func (m *MsgAddBurnAuthority) Route() string             { return ModuleName }
func (m *MsgAddBurnAuthority) Type() string              { return TypeMsgAddBurnAuthority }
func (m *MsgAddBurnAuthority) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}
func (m *MsgAddBurnAuthority) ValidateBasic() error {
	if m.TenitesId == "" {
		return ErrInvalidAuthority.Wrap("tenites_id required")
	}
	return nil
}

func (m *MsgAddBurnAuthorityResponse) ProtoMessage()           {}
func (m *MsgAddBurnAuthorityResponse) Reset()                  {}
func (m *MsgAddBurnAuthorityResponse) String() string          { return "MsgAddBurnAuthorityResponse" }
func (m *MsgAddBurnAuthorityResponse) XXX_MessageName() string { return "tenites.tsc.v1.MsgAddBurnAuthorityResponse" }

type MsgRemoveBurnAuthority struct {
	Signer      string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer"`
	AuthorityId string `protobuf:"bytes,2,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
}

type MsgRemoveBurnAuthorityResponse struct{}

func (m *MsgRemoveBurnAuthority) ProtoMessage()             {}
func (m *MsgRemoveBurnAuthority) Reset()                    {}
func (m *MsgRemoveBurnAuthority) String() string            { return TypeMsgRemoveBurnAuthority }
func (m *MsgRemoveBurnAuthority) XXX_MessageName() string   { return "tenites.tsc.v1.MsgRemoveBurnAuthority" }
func (m *MsgRemoveBurnAuthority) Route() string             { return ModuleName }
func (m *MsgRemoveBurnAuthority) Type() string              { return TypeMsgRemoveBurnAuthority }
func (m *MsgRemoveBurnAuthority) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Signer)
	return []sdk.AccAddress{addr}
}
func (m *MsgRemoveBurnAuthority) ValidateBasic() error {
	if m.AuthorityId == "" {
		return ErrInvalidAuthority.Wrap("authority_id required")
	}
	return nil
}

func (m *MsgRemoveBurnAuthorityResponse) ProtoMessage()           {}
func (m *MsgRemoveBurnAuthorityResponse) Reset()                  {}
func (m *MsgRemoveBurnAuthorityResponse) String() string          { return "MsgRemoveBurnAuthorityResponse" }
func (m *MsgRemoveBurnAuthorityResponse) XXX_MessageName() string { return "tenites.tsc.v1.MsgRemoveBurnAuthorityResponse" }

var ErrInvalidAuthority = ErrUnauthorized

type MsgServer interface {
	Mint(context.Context, *MsgMint) (*MsgMintResponse, error)
	Burn(context.Context, *MsgBurn) (*MsgBurnResponse, error)
	Transfer(context.Context, *MsgTransfer) (*MsgTransferResponse, error)
	AddMintAuthority(context.Context, *MsgAddMintAuthority) (*MsgAddMintAuthorityResponse, error)
	RemoveMintAuthority(context.Context, *MsgRemoveMintAuthority) (*MsgRemoveMintAuthorityResponse, error)
	AddBurnAuthority(context.Context, *MsgAddBurnAuthority) (*MsgAddBurnAuthorityResponse, error)
	RemoveBurnAuthority(context.Context, *MsgRemoveBurnAuthority) (*MsgRemoveBurnAuthorityResponse, error)
}
