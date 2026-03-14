package keeper

import (
        "fmt"

        errorsmod "cosmossdk.io/errors"
        "cosmossdk.io/log"
        "cosmossdk.io/math"
        storetypes "cosmossdk.io/store/types"
        abci "github.com/cometbft/cometbft/abci/types"
        "github.com/cosmos/cosmos-sdk/codec"
        sdk "github.com/cosmos/cosmos-sdk/types"
        bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

        identitykeeper "github.com/tenites/tenites-chain/x/identity/keeper"
        walletkeeper "github.com/tenites/tenites-chain/x/wallet/keeper"
        compliancekeeper "github.com/tenites/tenites-chain/x/compliance/keeper"
        "github.com/tenites/tenites-chain/x/tsc/types"
)

type ABCIQuerier interface {
        Query(req *abci.RequestQuery) (*abci.ResponseQuery, error)
}

type Keeper struct {
        cdc              codec.BinaryCodec
        storeKey         storetypes.StoreKey
        authority        string
        identityKeeper   identitykeeper.Keeper
        walletKeeper     walletkeeper.Keeper
        complianceKeeper compliancekeeper.Keeper
        bankKeeper       bankkeeper.Keeper
        abciQuerier      ABCIQuerier
}

func NewKeeper(
        cdc codec.BinaryCodec,
        storeKey storetypes.StoreKey,
        authority string,
        identityKeeper identitykeeper.Keeper,
        walletKeeper walletkeeper.Keeper,
        complianceKeeper compliancekeeper.Keeper,
        bankKeeper bankkeeper.Keeper,
) Keeper {
        return Keeper{
                cdc:              cdc,
                storeKey:         storeKey,
                authority:        authority,
                identityKeeper:   identityKeeper,
                walletKeeper:     walletKeeper,
                complianceKeeper: complianceKeeper,
                bankKeeper:       bankKeeper,
        }
}

func (k *Keeper) SetABCIQuerier(querier ABCIQuerier) {
        k.abciQuerier = querier
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
        return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthority() string {
        return k.authority
}

func (k Keeper) SetParams(ctx sdk.Context, params types.TSCParams) {
        store := ctx.KVStore(k.storeKey)
        bz := k.cdc.MustMarshal(&params)
        store.Set(types.ParamsKey, bz)
}

func (k Keeper) GetParams(ctx sdk.Context) types.TSCParams {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.ParamsKey)
        if bz == nil {
                return types.DefaultTSCParams()
        }

        var params types.TSCParams
        k.cdc.MustUnmarshal(bz, &params)
        return params
}

func (k Keeper) SetMintAuthority(ctx sdk.Context, authority types.TSCMintAuthority) {
        store := ctx.KVStore(k.storeKey)
        bz := k.cdc.MustMarshal(&authority)
        store.Set(types.MintAuthorityKey(authority.AuthorityId), bz)
        store.Set(types.TenitesIdMintAuthorityIndexKey(authority.TenitesId, authority.AuthorityId), []byte{1})
}

func (k Keeper) GetMintAuthority(ctx sdk.Context, authorityId string) (types.TSCMintAuthority, bool) {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.MintAuthorityKey(authorityId))
        if bz == nil {
                return types.TSCMintAuthority{}, false
        }

        var authority types.TSCMintAuthority
        k.cdc.MustUnmarshal(bz, &authority)
        return authority, true
}

func (k Keeper) GetAllMintAuthorities(ctx sdk.Context) []types.TSCMintAuthority {
        store := ctx.KVStore(k.storeKey)
        iterator := storetypes.KVStorePrefixIterator(store, types.MintAuthorityPrefix)
        defer iterator.Close()

        var authorities []types.TSCMintAuthority
        for ; iterator.Valid(); iterator.Next() {
                var authority types.TSCMintAuthority
                k.cdc.MustUnmarshal(iterator.Value(), &authority)
                authorities = append(authorities, authority)
        }

        return authorities
}

func (k Keeper) SetBurnAuthority(ctx sdk.Context, authority types.TSCBurnAuthority) {
        store := ctx.KVStore(k.storeKey)
        bz := k.cdc.MustMarshal(&authority)
        store.Set(types.BurnAuthorityKey(authority.AuthorityId), bz)
        store.Set(types.TenitesIdBurnAuthorityIndexKey(authority.TenitesId, authority.AuthorityId), []byte{1})
}

func (k Keeper) GetBurnAuthority(ctx sdk.Context, authorityId string) (types.TSCBurnAuthority, bool) {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.BurnAuthorityKey(authorityId))
        if bz == nil {
                return types.TSCBurnAuthority{}, false
        }

        var authority types.TSCBurnAuthority
        k.cdc.MustUnmarshal(bz, &authority)
        return authority, true
}

func (k Keeper) GetAllBurnAuthorities(ctx sdk.Context) []types.TSCBurnAuthority {
        store := ctx.KVStore(k.storeKey)
        iterator := storetypes.KVStorePrefixIterator(store, types.BurnAuthorityPrefix)
        defer iterator.Close()

        var authorities []types.TSCBurnAuthority
        for ; iterator.Valid(); iterator.Next() {
                var authority types.TSCBurnAuthority
                k.cdc.MustUnmarshal(iterator.Value(), &authority)
                authorities = append(authorities, authority)
        }

        return authorities
}

func (k Keeper) SetOperation(ctx sdk.Context, operation types.TSCOperation) {
        store := ctx.KVStore(k.storeKey)
        bz := k.cdc.MustMarshal(&operation)
        store.Set(types.OperationKey(operation.OperationId), bz)
        store.Set(types.WalletOperationIndexKey(operation.WalletId, operation.OperationId), []byte{1})
}

func (k Keeper) GetOperation(ctx sdk.Context, operationId string) (types.TSCOperation, bool) {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.OperationKey(operationId))
        if bz == nil {
                return types.TSCOperation{}, false
        }

        var operation types.TSCOperation
        k.cdc.MustUnmarshal(bz, &operation)
        return operation, true
}

func (k Keeper) SetSupplySnapshot(ctx sdk.Context, snapshot types.TSCSupplySnapshot) {
        store := ctx.KVStore(k.storeKey)
        bz := k.cdc.MustMarshal(&snapshot)
        store.Set(types.SupplySnapshotKey(snapshot.SnapshotId), bz)
}

func (k Keeper) GetSupplySnapshot(ctx sdk.Context, snapshotId string) (types.TSCSupplySnapshot, bool) {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.SupplySnapshotKey(snapshotId))
        if bz == nil {
                return types.TSCSupplySnapshot{}, false
        }

        var snapshot types.TSCSupplySnapshot
        k.cdc.MustUnmarshal(bz, &snapshot)
        return snapshot, true
}

func (k Keeper) GetTotalSupply(ctx sdk.Context) math.Int {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.TotalSupplyKey)
        if bz == nil {
                return math.ZeroInt()
        }
        
        var supply math.Int
        if err := supply.Unmarshal(bz); err != nil {
                return math.ZeroInt()
        }
        return supply
}

func (k Keeper) SetTotalSupply(ctx sdk.Context, supply math.Int) {
        store := ctx.KVStore(k.storeKey)
        bz, _ := supply.Marshal()
        store.Set(types.TotalSupplyKey, bz)
}

func (k Keeper) GetDailyMinted(ctx sdk.Context) math.Int {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.DailyMintedKey)
        if bz == nil {
                return math.ZeroInt()
        }
        
        var amount math.Int
        if err := amount.Unmarshal(bz); err != nil {
                return math.ZeroInt()
        }
        return amount
}

func (k Keeper) SetDailyMinted(ctx sdk.Context, amount math.Int) {
        store := ctx.KVStore(k.storeKey)
        bz, _ := amount.Marshal()
        store.Set(types.DailyMintedKey, bz)
}

func (k Keeper) GetDailyBurned(ctx sdk.Context) math.Int {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.DailyBurnedKey)
        if bz == nil {
                return math.ZeroInt()
        }
        
        var amount math.Int
        if err := amount.Unmarshal(bz); err != nil {
                return math.ZeroInt()
        }
        return amount
}

func (k Keeper) SetDailyBurned(ctx sdk.Context, amount math.Int) {
        store := ctx.KVStore(k.storeKey)
        bz, _ := amount.Marshal()
        store.Set(types.DailyBurnedKey, bz)
}

func (k Keeper) GetDailyResetAt(ctx sdk.Context) int64 {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.DailyResetAtKey)
        if bz == nil {
                return 0
        }
        return int64(sdk.BigEndianToUint64(bz))
}

func (k Keeper) SetDailyResetAt(ctx sdk.Context, timestamp int64) {
        store := ctx.KVStore(k.storeKey)
        bz := sdk.Uint64ToBigEndian(uint64(timestamp))
        store.Set(types.DailyResetAtKey, bz)
}

func (k Keeper) GetCounter(ctx sdk.Context) uint64 {
        store := ctx.KVStore(k.storeKey)
        bz := store.Get(types.CounterKey)
        if bz == nil {
                return 0
        }
        return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetCounter(ctx sdk.Context, counter uint64) {
        store := ctx.KVStore(k.storeKey)
        bz := sdk.Uint64ToBigEndian(counter)
        store.Set(types.CounterKey, bz)
}

func (k Keeper) checkDailyReset(ctx sdk.Context) {
        currentTime := ctx.BlockTime().Unix()
        resetAt := k.GetDailyResetAt(ctx)
        
        if currentTime >= resetAt {
                k.SetDailyMinted(ctx, math.ZeroInt())
                k.SetDailyBurned(ctx, math.ZeroInt())
                k.SetDailyResetAt(ctx, currentTime + 86400)
        }
}

func (k Keeper) AddMintAuthority(
        ctx sdk.Context,
        tenitesId string,
        name string,
        dailyLimit math.Int,
        singleTxLimit math.Int,
        allowedPurposes []string,
        jurisdictions []string,
        expiresAt int64,
        actor string,
) (types.TSCMintAuthority, error) {
        if actor != k.authority {
                return types.TSCMintAuthority{}, types.ErrUnauthorized
        }

        counter := k.GetCounter(ctx)
        authorityId := fmt.Sprintf("MINT%010d", counter+1)

        authority := types.NewTSCMintAuthority(
                authorityId,
                tenitesId,
                name,
                dailyLimit,
                singleTxLimit,
                allowedPurposes,
                jurisdictions,
                actor,
                ctx.BlockTime().Unix(),
                expiresAt,
        )

        if err := authority.Validate(); err != nil {
                return types.TSCMintAuthority{}, err
        }

        k.SetMintAuthority(ctx, authority)
        k.SetCounter(ctx, counter+1)

        k.emitAuthorityEvent(ctx, "TSC_MINT_AUTHORITY_ADDED", authorityId, tenitesId, actor, "")

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_mint_authority_added",
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("tenites_id", tenitesId),
                        sdk.NewAttribute("name", name),
                        sdk.NewAttribute("granted_by", actor),
                ),
        )

        return authority, nil
}

func (k Keeper) RemoveMintAuthority(
        ctx sdk.Context,
        authorityId string,
        actor string,
) error {
        if actor != k.authority {
                return types.ErrUnauthorized
        }

        authority, found := k.GetMintAuthority(ctx, authorityId)
        if !found {
                return types.ErrMintAuthorityNotFound
        }

        authority.Revoke(ctx.BlockTime().Unix())
        k.SetMintAuthority(ctx, authority)

        k.emitAuthorityEvent(ctx, "TSC_MINT_AUTHORITY_REVOKED", authorityId, authority.TenitesId, actor, "governance_revocation")

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_mint_authority_revoked",
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("revoked_by", actor),
                ),
        )

        return nil
}

func (k Keeper) AddBurnAuthority(
        ctx sdk.Context,
        tenitesId string,
        name string,
        dailyLimit math.Int,
        singleTxLimit math.Int,
        allowedReasons []string,
        jurisdictions []string,
        expiresAt int64,
        actor string,
) (types.TSCBurnAuthority, error) {
        if actor != k.authority {
                return types.TSCBurnAuthority{}, types.ErrUnauthorized
        }

        counter := k.GetCounter(ctx)
        authorityId := fmt.Sprintf("BURN%010d", counter+1)

        authority := types.NewTSCBurnAuthority(
                authorityId,
                tenitesId,
                name,
                dailyLimit,
                singleTxLimit,
                allowedReasons,
                jurisdictions,
                actor,
                ctx.BlockTime().Unix(),
                expiresAt,
        )

        if err := authority.Validate(); err != nil {
                return types.TSCBurnAuthority{}, err
        }

        k.SetBurnAuthority(ctx, authority)
        k.SetCounter(ctx, counter+1)

        k.emitAuthorityEvent(ctx, "TSC_BURN_AUTHORITY_ADDED", authorityId, tenitesId, actor, "")

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_burn_authority_added",
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("tenites_id", tenitesId),
                        sdk.NewAttribute("name", name),
                        sdk.NewAttribute("granted_by", actor),
                ),
        )

        return authority, nil
}

func (k Keeper) RemoveBurnAuthority(
        ctx sdk.Context,
        authorityId string,
        actor string,
) error {
        if actor != k.authority {
                return types.ErrUnauthorized
        }

        authority, found := k.GetBurnAuthority(ctx, authorityId)
        if !found {
                return types.ErrBurnAuthorityNotFound
        }

        authority.Revoke(ctx.BlockTime().Unix())
        k.SetBurnAuthority(ctx, authority)

        k.emitAuthorityEvent(ctx, "TSC_BURN_AUTHORITY_REVOKED", authorityId, authority.TenitesId, actor, "governance_revocation")

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_burn_authority_revoked",
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("revoked_by", actor),
                ),
        )

        return nil
}

func (k Keeper) MintTSC(
        ctx sdk.Context,
        authorityId string,
        walletId string,
        amount math.Int,
        purpose string,
        jurisdiction string,
        actor string,
) (types.TSCOperation, error) {
        k.checkDailyReset(ctx)
        blockTime := ctx.BlockTime().Unix()
        params := k.GetParams(ctx)

        if amount.IsNegative() || amount.IsZero() {
                return types.TSCOperation{}, types.ErrInvalidAmount
        }

        authority, found := k.GetMintAuthority(ctx, authorityId)
        if !found {
                return types.TSCOperation{}, types.ErrMintAuthorityNotFound
        }

        if authority.Address != actor && authority.TenitesId != actor && actor != k.authority {
                return types.TSCOperation{}, types.ErrUnauthorized
        }

        if err := authority.CanMint(amount, purpose, jurisdiction, blockTime); err != nil {
                return types.TSCOperation{}, err
        }

        if !params.IsJurisdictionAllowed(jurisdiction) {
                return types.TSCOperation{}, types.ErrJurisdictionNotAllowed
        }

        if amount.GT(params.SingleMintLimit) {
                return types.TSCOperation{}, types.ErrSingleMintLimitExceeded
        }

        dailyMinted := k.GetDailyMinted(ctx)
        if dailyMinted.Add(amount).GT(params.DailyMintCap) {
                return types.TSCOperation{}, types.ErrDailyMintCapExceeded
        }

        totalSupply := k.GetTotalSupply(ctx)
        if totalSupply.Add(amount).GT(params.TotalSupplyCap) {
                return types.TSCOperation{}, types.ErrTotalSupplyCapExceeded
        }

        counter := k.GetCounter(ctx)
        operationId := fmt.Sprintf("TSCOP%010d", counter+1)

        operation := types.NewTSCMintOperation(
                operationId,
                authorityId,
                walletId,
                amount,
                purpose,
                jurisdiction,
                blockTime,
        )

        err := k.walletKeeper.CreditWallet(ctx, walletId, amount, operationId, "tsc_mint", k.authority)
        if err != nil {
                operation.Fail(err.Error(), blockTime)
                k.SetOperation(ctx, operation)
                return operation, err
        }

        authority.RecordMint(amount, blockTime)
        k.SetMintAuthority(ctx, authority)

        k.SetTotalSupply(ctx, totalSupply.Add(amount))
        k.SetDailyMinted(ctx, dailyMinted.Add(amount))

        operation.Complete(fmt.Sprintf("block_%d", ctx.BlockHeight()), blockTime)
        k.SetOperation(ctx, operation)
        k.SetCounter(ctx, counter+1)

        k.emitTSCMintedEvent(ctx, operationId, authorityId, walletId, amount, purpose, jurisdiction)

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_minted",
                        sdk.NewAttribute("operation_id", operationId),
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("wallet_id", walletId),
                        sdk.NewAttribute("amount", amount.String()),
                        sdk.NewAttribute("purpose", purpose),
                        sdk.NewAttribute("jurisdiction", jurisdiction),
                        sdk.NewAttribute("new_total_supply", totalSupply.Add(amount).String()),
                ),
        )

        return operation, nil
}

func (k Keeper) BurnTSC(
        ctx sdk.Context,
        authorityId string,
        walletId string,
        amount math.Int,
        reason string,
        jurisdiction string,
        actor string,
) (types.TSCOperation, error) {
        k.checkDailyReset(ctx)
        blockTime := ctx.BlockTime().Unix()
        params := k.GetParams(ctx)

        if amount.IsNegative() || amount.IsZero() {
                return types.TSCOperation{}, types.ErrInvalidAmount
        }

        authority, found := k.GetBurnAuthority(ctx, authorityId)
        if !found {
                return types.TSCOperation{}, types.ErrBurnAuthorityNotFound
        }

        if authority.Address != actor && authority.TenitesId != actor && actor != k.authority {
                return types.TSCOperation{}, types.ErrUnauthorized
        }

        if err := authority.CanBurn(amount, reason, jurisdiction, blockTime); err != nil {
                return types.TSCOperation{}, err
        }

        if !params.IsJurisdictionAllowed(jurisdiction) {
                return types.TSCOperation{}, types.ErrJurisdictionNotAllowed
        }

        if amount.GT(params.SingleBurnLimit) {
                return types.TSCOperation{}, types.ErrSingleBurnLimitExceeded
        }

        dailyBurned := k.GetDailyBurned(ctx)
        if dailyBurned.Add(amount).GT(params.DailyBurnCap) {
                return types.TSCOperation{}, types.ErrDailyBurnCapExceeded
        }

        totalSupply := k.GetTotalSupply(ctx)
        newSupply := totalSupply.Sub(amount)
        if newSupply.IsNegative() {
                return types.TSCOperation{}, errorsmod.Wrapf(
                        types.ErrInsufficientSupply,
                        "burn amount %s exceeds total supply %s",
                        amount.String(),
                        totalSupply.String(),
                )
        }

        counter := k.GetCounter(ctx)
        operationId := fmt.Sprintf("TSCOP%010d", counter+1)

        operation := types.NewTSCBurnOperation(
                operationId,
                authorityId,
                walletId,
                amount,
                reason,
                jurisdiction,
                blockTime,
        )

        err := k.walletKeeper.DebitWallet(ctx, walletId, amount, operationId, "tsc_burn", k.authority)
        if err != nil {
                operation.Fail(err.Error(), blockTime)
                k.SetOperation(ctx, operation)
                return operation, types.ErrInsufficientBalance
        }

        authority.RecordBurn(amount, blockTime)
        k.SetBurnAuthority(ctx, authority)

        k.SetTotalSupply(ctx, newSupply)
        k.SetDailyBurned(ctx, dailyBurned.Add(amount))

        operation.Complete(fmt.Sprintf("block_%d", ctx.BlockHeight()), blockTime)
        k.SetOperation(ctx, operation)
        k.SetCounter(ctx, counter+1)

        k.emitTSCBurnedEvent(ctx, operationId, authorityId, walletId, amount, reason, jurisdiction)

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_burned",
                        sdk.NewAttribute("operation_id", operationId),
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("wallet_id", walletId),
                        sdk.NewAttribute("amount", amount.String()),
                        sdk.NewAttribute("reason", reason),
                        sdk.NewAttribute("jurisdiction", jurisdiction),
                        sdk.NewAttribute("new_total_supply", newSupply.String()),
                ),
        )

        return operation, nil
}

func (k Keeper) TransferTSC(
        ctx sdk.Context,
        fromWalletId string,
        toWalletId string,
        amount math.Int,
        purpose string,
        actor string,
) (types.TSCOperation, error) {
        blockTime := ctx.BlockTime().Unix()
        params := k.GetParams(ctx)

        if amount.IsNegative() || amount.IsZero() {
                return types.TSCOperation{}, types.ErrInvalidAmount
        }

        fromWallet, found := k.walletKeeper.GetWallet(ctx, fromWalletId)
        if !found {
                return types.TSCOperation{}, types.ErrWalletNotFound
        }

        toWallet, found := k.walletKeeper.GetWallet(ctx, toWalletId)
        if !found {
                return types.TSCOperation{}, types.ErrWalletNotFound
        }

        if fromWallet.OwnerTenitesId != actor && actor != k.authority {
                return types.TSCOperation{}, types.ErrUnauthorized
        }

        if params.ComplianceRequired && amount.GTE(params.TravelRuleThreshold) {
                origIdentity, found := k.identityKeeper.GetIdentity(ctx, fromWallet.OwnerTenitesId)
                if !found {
                        return types.TSCOperation{}, types.ErrIdentityNotFound
                }
                if int32(origIdentity.KycTier) < params.MinKycTierForTransfer {
                        return types.TSCOperation{}, types.ErrKycTierInsufficient
                }

                benefIdentity, found := k.identityKeeper.GetIdentity(ctx, toWallet.OwnerTenitesId)
                if !found {
                        return types.TSCOperation{}, types.ErrIdentityNotFound
                }
                if int32(benefIdentity.KycTier) < params.MinKycTierForTransfer {
                        return types.TSCOperation{}, types.ErrKycTierInsufficient
                }
        }

        counter := k.GetCounter(ctx)
        operationId := fmt.Sprintf("TSCOP%010d", counter+1)

        operation := types.NewTSCTransferOperation(
                operationId,
                fromWalletId,
                toWalletId,
                amount,
                purpose,
                fromWallet.Jurisdiction,
                blockTime,
        )

        err := k.walletKeeper.DebitWallet(ctx, fromWalletId, amount, operationId, "tsc_transfer", k.authority)
        if err != nil {
                operation.Fail(err.Error(), blockTime)
                k.SetOperation(ctx, operation)
                return operation, types.ErrInsufficientBalance
        }

        err = k.walletKeeper.CreditWallet(ctx, toWalletId, amount, operationId, "tsc_transfer", k.authority)
        if err != nil {
                k.walletKeeper.CreditWallet(ctx, fromWalletId, amount, operationId, "tsc_transfer_rollback", k.authority)
                operation.Fail(err.Error(), blockTime)
                k.SetOperation(ctx, operation)
                return operation, err
        }

        operation.Complete(fmt.Sprintf("block_%d", ctx.BlockHeight()), blockTime)
        k.SetOperation(ctx, operation)
        k.SetCounter(ctx, counter+1)

        k.emitTSCTransferredEvent(ctx, operationId, fromWalletId, toWalletId, amount, purpose)

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_transferred",
                        sdk.NewAttribute("operation_id", operationId),
                        sdk.NewAttribute("from_wallet_id", fromWalletId),
                        sdk.NewAttribute("to_wallet_id", toWalletId),
                        sdk.NewAttribute("amount", amount.String()),
                        sdk.NewAttribute("purpose", purpose),
                ),
        )

        return operation, nil
}

func (k Keeper) UpdateParams(
        ctx sdk.Context,
        newParams types.TSCParams,
        actor string,
) error {
        if actor != k.authority {
                return types.ErrUnauthorized
        }

        if err := newParams.Validate(); err != nil {
                return types.ErrInvalidParams
        }

        oldParams := k.GetParams(ctx)
        newParams.UpdatedAt = ctx.BlockTime().Unix()
        k.SetParams(ctx, newParams)

        k.emitParamsUpdatedEvent(ctx, oldParams, newParams, actor)

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_params_updated",
                        sdk.NewAttribute("updated_by", actor),
                        sdk.NewAttribute("daily_mint_cap", newParams.DailyMintCap.String()),
                        sdk.NewAttribute("daily_burn_cap", newParams.DailyBurnCap.String()),
                        sdk.NewAttribute("total_supply_cap", newParams.TotalSupplyCap.String()),
                ),
        )

        return nil
}

func (k Keeper) CreateDailySnapshot(ctx sdk.Context) types.TSCSupplySnapshot {
        blockTime := ctx.BlockTime()
        snapshotDate := blockTime.Format("2006-01-02")
        snapshotId := fmt.Sprintf("SNAP_%s", snapshotDate)

        totalSupply := k.GetTotalSupply(ctx)
        dailyMinted := k.GetDailyMinted(ctx)
        dailyBurned := k.GetDailyBurned(ctx)

        snapshot := types.NewTSCSupplySnapshot(
                snapshotId,
                totalSupply,
                totalSupply,
                math.ZeroInt(),
                dailyMinted,
                dailyBurned,
                snapshotDate,
                blockTime.Unix(),
        )

        k.SetSupplySnapshot(ctx, snapshot)

        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_supply_snapshot",
                        sdk.NewAttribute("snapshot_id", snapshotId),
                        sdk.NewAttribute("total_supply", totalSupply.String()),
                        sdk.NewAttribute("daily_minted", dailyMinted.String()),
                        sdk.NewAttribute("daily_burned", dailyBurned.String()),
                        sdk.NewAttribute("snapshot_date", snapshotDate),
                ),
        )

        return snapshot
}

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
        k.SetParams(ctx, data.Params)
        k.SetTotalSupply(ctx, data.TotalSupply)
        k.SetCounter(ctx, data.Counter)
        k.SetDailyResetAt(ctx, ctx.BlockTime().Unix() + 86400)
        k.SetDailyMinted(ctx, math.ZeroInt())
        k.SetDailyBurned(ctx, math.ZeroInt())

        for _, authority := range data.MintAuthorities {
                k.SetMintAuthority(ctx, authority)
        }

        for _, authority := range data.BurnAuthorities {
                k.SetBurnAuthority(ctx, authority)
        }

        for _, operation := range data.Operations {
                k.SetOperation(ctx, operation)
        }

        for _, snapshot := range data.Snapshots {
                k.SetSupplySnapshot(ctx, snapshot)
        }
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
        return &types.GenesisState{
                Params:          k.GetParams(ctx),
                MintAuthorities: k.GetAllMintAuthorities(ctx),
                BurnAuthorities: k.GetAllBurnAuthorities(ctx),
                Operations:      []types.TSCOperation{},
                Snapshots:       []types.TSCSupplySnapshot{},
                TotalSupply:     k.GetTotalSupply(ctx),
                Counter:         k.GetCounter(ctx),
        }
}

func (k Keeper) emitTSCMintedEvent(
        ctx sdk.Context,
        operationId string,
        authorityId string,
        walletId string,
        amount math.Int,
        purpose string,
        jurisdiction string,
) {
        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_audit",
                        sdk.NewAttribute("event_type", "TSC_MINTED"),
                        sdk.NewAttribute("operation_id", operationId),
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("wallet_id", walletId),
                        sdk.NewAttribute("amount", amount.String()),
                        sdk.NewAttribute("purpose", purpose),
                        sdk.NewAttribute("jurisdiction", jurisdiction),
                        sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
                        sdk.NewAttribute("timestamp", fmt.Sprintf("%d", ctx.BlockTime().Unix())),
                ),
        )
}

func (k Keeper) emitTSCBurnedEvent(
        ctx sdk.Context,
        operationId string,
        authorityId string,
        walletId string,
        amount math.Int,
        reason string,
        jurisdiction string,
) {
        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_audit",
                        sdk.NewAttribute("event_type", "TSC_BURNED"),
                        sdk.NewAttribute("operation_id", operationId),
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("wallet_id", walletId),
                        sdk.NewAttribute("amount", amount.String()),
                        sdk.NewAttribute("reason", reason),
                        sdk.NewAttribute("jurisdiction", jurisdiction),
                        sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
                        sdk.NewAttribute("timestamp", fmt.Sprintf("%d", ctx.BlockTime().Unix())),
                ),
        )
}

func (k Keeper) emitTSCTransferredEvent(
        ctx sdk.Context,
        operationId string,
        fromWalletId string,
        toWalletId string,
        amount math.Int,
        purpose string,
) {
        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_audit",
                        sdk.NewAttribute("event_type", "TSC_TRANSFERRED"),
                        sdk.NewAttribute("operation_id", operationId),
                        sdk.NewAttribute("from_wallet_id", fromWalletId),
                        sdk.NewAttribute("to_wallet_id", toWalletId),
                        sdk.NewAttribute("amount", amount.String()),
                        sdk.NewAttribute("purpose", purpose),
                        sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
                        sdk.NewAttribute("timestamp", fmt.Sprintf("%d", ctx.BlockTime().Unix())),
                ),
        )
}

func (k Keeper) emitAuthorityEvent(
        ctx sdk.Context,
        eventType string,
        authorityId string,
        tenitesId string,
        actor string,
        reason string,
) {
        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_audit",
                        sdk.NewAttribute("event_type", eventType),
                        sdk.NewAttribute("authority_id", authorityId),
                        sdk.NewAttribute("tenites_id", tenitesId),
                        sdk.NewAttribute("actor", actor),
                        sdk.NewAttribute("reason", reason),
                        sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
                        sdk.NewAttribute("timestamp", fmt.Sprintf("%d", ctx.BlockTime().Unix())),
                ),
        )
}

func (k Keeper) emitParamsUpdatedEvent(
        ctx sdk.Context,
        oldParams types.TSCParams,
        newParams types.TSCParams,
        actor string,
) {
        ctx.EventManager().EmitEvent(
                sdk.NewEvent(
                        "tsc_audit",
                        sdk.NewAttribute("event_type", "TSC_PARAMS_UPDATED"),
                        sdk.NewAttribute("actor", actor),
                        sdk.NewAttribute("old_daily_mint_cap", oldParams.DailyMintCap.String()),
                        sdk.NewAttribute("new_daily_mint_cap", newParams.DailyMintCap.String()),
                        sdk.NewAttribute("old_daily_burn_cap", oldParams.DailyBurnCap.String()),
                        sdk.NewAttribute("new_daily_burn_cap", newParams.DailyBurnCap.String()),
                        sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
                        sdk.NewAttribute("timestamp", fmt.Sprintf("%d", ctx.BlockTime().Unix())),
                ),
        )
}

func SupplyInvariant(k Keeper) sdk.Invariant {
        return func(ctx sdk.Context) (string, bool) {
                totalSupply := k.GetTotalSupply(ctx)
                params := k.GetParams(ctx)

                if totalSupply.IsNegative() {
                        return sdk.FormatInvariant(
                                types.ModuleName, "negative-supply",
                                fmt.Sprintf("total supply is negative: %s", totalSupply.String()),
                        ), true
                }

                if totalSupply.GT(params.TotalSupplyCap) {
                        return sdk.FormatInvariant(
                                types.ModuleName, "supply-exceeds-cap",
                                fmt.Sprintf("total supply %s exceeds cap %s", totalSupply.String(), params.TotalSupplyCap.String()),
                        ), true
                }

                dailyMinted := k.GetDailyMinted(ctx)
                if dailyMinted.GT(params.DailyMintCap) {
                        return sdk.FormatInvariant(
                                types.ModuleName, "daily-mint-exceeded",
                                fmt.Sprintf("daily minted %s exceeds cap %s", dailyMinted.String(), params.DailyMintCap.String()),
                        ), true
                }

                dailyBurned := k.GetDailyBurned(ctx)
                if dailyBurned.GT(params.DailyBurnCap) {
                        return sdk.FormatInvariant(
                                types.ModuleName, "daily-burn-exceeded",
                                fmt.Sprintf("daily burned %s exceeds cap %s", dailyBurned.String(), params.DailyBurnCap.String()),
                        ), true
                }

                return sdk.FormatInvariant(
                        types.ModuleName, "supply-invariants",
                        "all supply invariants hold",
                ), false
        }
}
