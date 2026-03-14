package keeper

import (
        "context"
        "fmt"

        "cosmossdk.io/math"
        sdk "github.com/cosmos/cosmos-sdk/types"
        "github.com/tenites/tenites-chain/x/tsc/types"
)

type MsgServer struct {
        Keeper
}

var _ types.MsgServer = MsgServer{}

func NewMsgServer(k Keeper) MsgServer {
        return MsgServer{Keeper: k}
}

func (ms MsgServer) Mint(
        goCtx context.Context,
        msg *types.MsgMint,
) (*types.MsgMintResponse, error) {
        if msg == nil {
                return nil, fmt.Errorf("nil message")
        }
        ctx := sdk.UnwrapSDKContext(goCtx)

        amount, ok := math.NewIntFromString(msg.Amount)
        if !ok {
                return nil, fmt.Errorf("invalid amount: %s", msg.Amount)
        }

        op, err := ms.Keeper.MintTSC(
                ctx,
                msg.AuthorityId,
                msg.WalletId,
                amount,
                msg.Purpose,
                msg.Jurisdiction,
                msg.Signer,
        )
        if err != nil {
                return nil, err
        }

        return &types.MsgMintResponse{
                OperationId: op.OperationId,
        }, nil
}

func (ms MsgServer) Burn(
        goCtx context.Context,
        msg *types.MsgBurn,
) (*types.MsgBurnResponse, error) {
        if msg == nil {
                return nil, fmt.Errorf("nil message")
        }
        ctx := sdk.UnwrapSDKContext(goCtx)

        amount, ok := math.NewIntFromString(msg.Amount)
        if !ok {
                return nil, fmt.Errorf("invalid amount: %s", msg.Amount)
        }

        _, err := ms.Keeper.BurnTSC(
                ctx,
                msg.AuthorityId,
                msg.WalletId,
                amount,
                msg.Reason,
                msg.Jurisdiction,
                msg.Signer,
        )
        if err != nil {
                return nil, err
        }

        return &types.MsgBurnResponse{}, nil
}

func (ms MsgServer) Transfer(
        goCtx context.Context,
        msg *types.MsgTransfer,
) (*types.MsgTransferResponse, error) {
        if msg == nil {
                return nil, fmt.Errorf("nil message")
        }
        ctx := sdk.UnwrapSDKContext(goCtx)

        amount, ok := math.NewIntFromString(msg.Amount)
        if !ok {
                return nil, fmt.Errorf("invalid amount: %s", msg.Amount)
        }

        _, err := ms.Keeper.TransferTSC(
                ctx,
                msg.FromWalletId,
                msg.ToWalletId,
                amount,
                msg.Purpose,
                msg.Signer,
        )
        if err != nil {
                return nil, err
        }

        return &types.MsgTransferResponse{}, nil
}

func (ms MsgServer) AddMintAuthority(
        goCtx context.Context,
        msg *types.MsgAddMintAuthority,
) (*types.MsgAddMintAuthorityResponse, error) {
        if msg == nil {
                return nil, fmt.Errorf("nil message")
        }
        ctx := sdk.UnwrapSDKContext(goCtx)

        dailyLimit, ok := math.NewIntFromString(msg.DailyLimit)
        if !ok {
                return nil, fmt.Errorf("invalid daily_limit: %s", msg.DailyLimit)
        }

        singleTxLimit, ok := math.NewIntFromString(msg.SingleTxLimit)
        if !ok {
                return nil, fmt.Errorf("invalid single_tx_limit: %s", msg.SingleTxLimit)
        }

        authority, err := ms.Keeper.AddMintAuthority(
                ctx,
                msg.TenitesId,
                msg.Name,
                dailyLimit,
                singleTxLimit,
                msg.AllowedPurposes,
                msg.Jurisdictions,
                msg.ExpiresAt,
                msg.Signer,
        )
        if err != nil {
                return nil, err
        }

        return &types.MsgAddMintAuthorityResponse{
                AuthorityId: authority.AuthorityId,
        }, nil
}

func (ms MsgServer) RemoveMintAuthority(
        goCtx context.Context,
        msg *types.MsgRemoveMintAuthority,
) (*types.MsgRemoveMintAuthorityResponse, error) {
        if msg == nil {
                return nil, fmt.Errorf("nil message")
        }
        ctx := sdk.UnwrapSDKContext(goCtx)

        err := ms.Keeper.RemoveMintAuthority(ctx, msg.AuthorityId, msg.Signer)
        if err != nil {
                return nil, err
        }

        return &types.MsgRemoveMintAuthorityResponse{}, nil
}

func (ms MsgServer) AddBurnAuthority(
        goCtx context.Context,
        msg *types.MsgAddBurnAuthority,
) (*types.MsgAddBurnAuthorityResponse, error) {
        if msg == nil {
                return nil, fmt.Errorf("nil message")
        }
        ctx := sdk.UnwrapSDKContext(goCtx)

        dailyLimit, ok := math.NewIntFromString(msg.DailyLimit)
        if !ok {
                return nil, fmt.Errorf("invalid daily_limit: %s", msg.DailyLimit)
        }

        singleTxLimit, ok := math.NewIntFromString(msg.SingleTxLimit)
        if !ok {
                return nil, fmt.Errorf("invalid single_tx_limit: %s", msg.SingleTxLimit)
        }

        authority, err := ms.Keeper.AddBurnAuthority(
                ctx,
                msg.TenitesId,
                msg.Name,
                dailyLimit,
                singleTxLimit,
                msg.AllowedReasons,
                msg.Jurisdictions,
                msg.ExpiresAt,
                msg.Signer,
        )
        if err != nil {
                return nil, err
        }

        return &types.MsgAddBurnAuthorityResponse{
                AuthorityId: authority.AuthorityId,
        }, nil
}

func (ms MsgServer) RemoveBurnAuthority(
        goCtx context.Context,
        msg *types.MsgRemoveBurnAuthority,
) (*types.MsgRemoveBurnAuthorityResponse, error) {
        if msg == nil {
                return nil, fmt.Errorf("nil message")
        }
        ctx := sdk.UnwrapSDKContext(goCtx)

        err := ms.Keeper.RemoveBurnAuthority(ctx, msg.AuthorityId, msg.Signer)
        if err != nil {
                return nil, err
        }

        return &types.MsgRemoveBurnAuthorityResponse{}, nil
}
