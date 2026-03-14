package app

import (
        "fmt"

        errorsmod "cosmossdk.io/errors"
        storetypes "cosmossdk.io/store/types"
        txsigning "cosmossdk.io/x/tx/signing"
        sdk "github.com/cosmos/cosmos-sdk/types"
        sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
        "github.com/cosmos/cosmos-sdk/types/tx/signing"
        "github.com/cosmos/cosmos-sdk/x/auth/ante"
        authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
        bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

        compliancekeeper "github.com/tenites/tenites-chain/x/compliance/keeper"
        identitykeeper "github.com/tenites/tenites-chain/x/identity/keeper"
)

type HandlerOptions struct {
        AccountKeeper   ante.AccountKeeper
        BankKeeper      bankkeeper.Keeper
        SignModeHandler *txsigning.HandlerMap
        SigGasConsumer  func(meter storetypes.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
        ComplianceKeeper compliancekeeper.Keeper
        IdentityKeeper   identitykeeper.Keeper
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
        if options.AccountKeeper == nil {
                return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
        }

        _ = options.BankKeeper

        if options.SignModeHandler == nil {
                return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
        }

        anteDecorators := []sdk.AnteDecorator{
                ante.NewSetUpContextDecorator(),
                ante.NewExtensionOptionsDecorator(nil),
                ante.NewValidateBasicDecorator(),
                ante.NewTxTimeoutHeightDecorator(),
                ante.NewValidateMemoDecorator(options.AccountKeeper),
                ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
                ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, nil, nil),
                ante.NewSetPubKeyDecorator(options.AccountKeeper),
                ante.NewValidateSigCountDecorator(options.AccountKeeper),
                ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
                ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
                ante.NewIncrementSequenceDecorator(options.AccountKeeper),
                NewComplianceDecorator(options.ComplianceKeeper, options.IdentityKeeper),
                NewWriteBarrierDecorator(),
        }

        return sdk.ChainAnteDecorators(anteDecorators...), nil
}

type ComplianceDecorator struct {
        complianceKeeper compliancekeeper.Keeper
        identityKeeper   identitykeeper.Keeper
}

func NewComplianceDecorator(ck compliancekeeper.Keeper, ik identitykeeper.Keeper) ComplianceDecorator {
        return ComplianceDecorator{
                complianceKeeper: ck,
                identityKeeper:   ik,
        }
}

func (cd ComplianceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
        msgs := tx.GetMsgs()

        for _, msg := range msgs {
                if err := cd.validateCompliance(ctx, msg); err != nil {
                        ctx.EventManager().EmitEvent(
                                sdk.NewEvent(
                                        "compliance_rejection",
                                        sdk.NewAttribute("reason_code", err.Error()),
                                        sdk.NewAttribute("msg_type", sdk.MsgTypeURL(msg)),
                                ),
                        )
                        return ctx, err
                }
        }

        return next(ctx, tx, simulate)
}

func (cd ComplianceDecorator) validateCompliance(ctx sdk.Context, msg sdk.Msg) error {
        signers, _, err := cd.getSigners(msg)
        if err != nil {
                return nil
        }

        for _, signer := range signers {
                identity, found := cd.identityKeeper.GetIdentityByOwner(ctx, signer.String())
                if !found {
                        continue
                }

                if identity.Status != "active" {
                        return errorsmod.Wrapf(
                                sdkerrors.ErrUnauthorized,
                                "COMPLIANCE_IDENTITY_NOT_ACTIVE: identity %s has status %s",
                                identity.TenitesId,
                                identity.Status,
                        )
                }
        }

        if err := cd.checkCorridorCompliance(ctx, msg); err != nil {
                return err
        }

        return nil
}

func (cd ComplianceDecorator) checkCorridorCompliance(ctx sdk.Context, msg sdk.Msg) error {
        type corridorMsg interface {
                GetSourceJurisdiction() string
                GetDestJurisdiction() string
        }

        cm, ok := msg.(corridorMsg)
        if !ok {
                return nil
        }

        source := cm.GetSourceJurisdiction()
        dest := cm.GetDestJurisdiction()
        if source == "" || dest == "" {
                return nil
        }

        corridorAllowed := cd.complianceKeeper.IsCorridorAllowed(ctx, source, dest)
        if !corridorAllowed {
                return errorsmod.Wrapf(
                        sdkerrors.ErrUnauthorized,
                        "COMPLIANCE_CORRIDOR_BLOCKED: corridor %s→%s is not permitted",
                        source,
                        dest,
                )
        }

        return nil
}

func (cd ComplianceDecorator) getSigners(msg sdk.Msg) ([]sdk.AccAddress, int, error) {
        signerIface, ok := msg.(interface{ GetSigners() []sdk.AccAddress })
        if !ok {
                return nil, 0, fmt.Errorf("message does not implement GetSigners")
        }
        addrs := signerIface.GetSigners()
        return addrs, len(addrs), nil
}

type WriteBarrierDecorator struct{}

func NewWriteBarrierDecorator() WriteBarrierDecorator {
        return WriteBarrierDecorator{}
}

func (wbd WriteBarrierDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
        if ctx.IsCheckTx() && !simulate {
                for _, msg := range tx.GetMsgs() {
                        msgTypeURL := sdk.MsgTypeURL(msg)
                        if isWriteMessage(msgTypeURL) {
                                ctx.EventManager().EmitEvent(
                                        sdk.NewEvent(
                                                "write_barrier_check",
                                                sdk.NewAttribute("msg_type", msgTypeURL),
                                        ),
                                )
                        }
                }
        }
        return next(ctx, tx, simulate)
}

func isWriteMessage(msgTypeURL string) bool {
        switch msgTypeURL {
        case "/tenites.tsc.v1.MsgMintTSC",
                "/tenites.tsc.v1.MsgBurnTSC",
                "/tenites.tsc.v1.MsgTransferTSC",
                "/tenites.obligation.v1.MsgCreateObligation",
                "/tenites.settlement.v1.MsgSettleObligation":
                return true
        }
        return false
}
