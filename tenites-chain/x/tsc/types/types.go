package types

import (
        "fmt"

        "cosmossdk.io/math"
)

type AuthorityStatus string

const (
        AuthorityStatusActive    AuthorityStatus = "active"
        AuthorityStatusSuspended AuthorityStatus = "suspended"
        AuthorityStatusRevoked   AuthorityStatus = "revoked"
)

type WalletTypeAllowed string

const (
        WalletTypeTreasury   WalletTypeAllowed = "treasury"
        WalletTypeCustodian  WalletTypeAllowed = "custodian"
        WalletTypeLiquidity  WalletTypeAllowed = "liquidity"
        WalletTypeSettlement WalletTypeAllowed = "settlement"
)

type MintPurpose string

const (
        MintPurposeIssuance      MintPurpose = "issuance"
        MintPurposeLiquidity     MintPurpose = "liquidity"
        MintPurposeSettlement    MintPurpose = "settlement"
        MintPurposeReserveBackup MintPurpose = "reserve_backup"
)

type BurnReason string

const (
        BurnReasonRedemption  BurnReason = "redemption"
        BurnReasonRecall      BurnReason = "recall"
        BurnReasonCompliance  BurnReason = "compliance"
        BurnReasonExpiry      BurnReason = "expiry"
        BurnReasonRebalancing BurnReason = "rebalancing"
)

type TSCParams struct {
        Denom                 string   `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom"`
        Decimals              int32    `protobuf:"varint,2,opt,name=decimals,proto3" json:"decimals"`
        TotalSupplyCap        math.Int `protobuf:"bytes,3,opt,name=total_supply_cap,json=totalSupplyCap,proto3,customtype=cosmossdk.io/math.Int" json:"total_supply_cap"`
        DailyMintCap          math.Int `protobuf:"bytes,4,opt,name=daily_mint_cap,json=dailyMintCap,proto3,customtype=cosmossdk.io/math.Int" json:"daily_mint_cap"`
        DailyBurnCap          math.Int `protobuf:"bytes,5,opt,name=daily_burn_cap,json=dailyBurnCap,proto3,customtype=cosmossdk.io/math.Int" json:"daily_burn_cap"`
        SingleMintLimit       math.Int `protobuf:"bytes,6,opt,name=single_mint_limit,json=singleMintLimit,proto3,customtype=cosmossdk.io/math.Int" json:"single_mint_limit"`
        SingleBurnLimit       math.Int `protobuf:"bytes,7,opt,name=single_burn_limit,json=singleBurnLimit,proto3,customtype=cosmossdk.io/math.Int" json:"single_burn_limit"`
        AllowedWalletTypes    []string `protobuf:"bytes,8,rep,name=allowed_wallet_types,json=allowedWalletTypes,proto3" json:"allowed_wallet_types"`
        AllowedJurisdictions  []string `protobuf:"bytes,9,rep,name=allowed_jurisdictions,json=allowedJurisdictions,proto3" json:"allowed_jurisdictions"`
        MinKycTierForMint     int32    `protobuf:"varint,10,opt,name=min_kyc_tier_for_mint,json=minKycTierForMint,proto3" json:"min_kyc_tier_for_mint"`
        MinKycTierForBurn     int32    `protobuf:"varint,11,opt,name=min_kyc_tier_for_burn,json=minKycTierForBurn,proto3" json:"min_kyc_tier_for_burn"`
        MinKycTierForTransfer int32    `protobuf:"varint,12,opt,name=min_kyc_tier_for_transfer,json=minKycTierForTransfer,proto3" json:"min_kyc_tier_for_transfer"`
        ComplianceRequired    bool     `protobuf:"varint,13,opt,name=compliance_required,json=complianceRequired,proto3" json:"compliance_required"`
        TravelRuleThreshold   math.Int `protobuf:"bytes,14,opt,name=travel_rule_threshold,json=travelRuleThreshold,proto3,customtype=cosmossdk.io/math.Int" json:"travel_rule_threshold"`
        UpdatedAt             int64    `protobuf:"varint,15,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
}

func DefaultTSCParams() TSCParams {
        return TSCParams{
                Denom:                 "utsc",
                Decimals:              6,
                TotalSupplyCap:        math.NewInt(10000000000000000),
                DailyMintCap:          math.NewInt(1000000000000),
                DailyBurnCap:          math.NewInt(1000000000000),
                SingleMintLimit:       math.NewInt(100000000000),
                SingleBurnLimit:       math.NewInt(100000000000),
                AllowedWalletTypes:    []string{"treasury", "custodian", "liquidity", "settlement"},
                AllowedJurisdictions:  []string{"NG", "GH", "KE", "CI", "ZA", "TZ", "UG", "RW", "CM", "ET", "EG", "MA", "TN"},
                MinKycTierForMint:     3,
                MinKycTierForBurn:     3,
                MinKycTierForTransfer: 2,
                ComplianceRequired:    true,
                TravelRuleThreshold:   math.NewInt(100000000),
                UpdatedAt:             0,
        }
}

func (p *TSCParams) Validate() error {
        if p.Denom == "" {
                return fmt.Errorf("denom cannot be empty")
        }
        if p.Decimals < 0 || p.Decimals > 18 {
                return fmt.Errorf("decimals must be between 0 and 18")
        }
        if p.TotalSupplyCap.IsNegative() {
                return fmt.Errorf("total supply cap cannot be negative")
        }
        if p.DailyMintCap.IsNegative() {
                return fmt.Errorf("daily mint cap cannot be negative")
        }
        if p.DailyBurnCap.IsNegative() {
                return fmt.Errorf("daily burn cap cannot be negative")
        }
        if p.SingleMintLimit.GT(p.DailyMintCap) {
                return fmt.Errorf("single mint limit cannot exceed daily mint cap")
        }
        if p.SingleBurnLimit.GT(p.DailyBurnCap) {
                return fmt.Errorf("single burn limit cannot exceed daily burn cap")
        }
        if len(p.AllowedWalletTypes) == 0 {
                return fmt.Errorf("at least one wallet type must be allowed")
        }
        if len(p.AllowedJurisdictions) == 0 {
                return fmt.Errorf("at least one jurisdiction must be allowed")
        }
        return nil
}

func (p *TSCParams) IsWalletTypeAllowed(walletType string) bool {
        for _, allowed := range p.AllowedWalletTypes {
                if allowed == walletType {
                        return true
                }
        }
        return false
}

func (p *TSCParams) IsJurisdictionAllowed(jurisdiction string) bool {
        for _, allowed := range p.AllowedJurisdictions {
                if allowed == jurisdiction {
                        return true
                }
        }
        return false
}

func (p *TSCParams) Reset() {
}

func (p *TSCParams) ProtoMessage() {}
func (p *TSCParams) String() string { return "TSCParams" }

type TSCMintAuthority struct {
        AuthorityId     string          `protobuf:"bytes,1,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
        TenitesId       string          `protobuf:"bytes,2,opt,name=tenites_id,json=tenitesId,proto3" json:"tenites_id"`
        Name            string          `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`
        Status          AuthorityStatus `protobuf:"bytes,4,opt,name=status,proto3,casttype=AuthorityStatus" json:"status"`
        DailyLimit      math.Int        `protobuf:"bytes,5,opt,name=daily_limit,json=dailyLimit,proto3,customtype=cosmossdk.io/math.Int" json:"daily_limit"`
        SingleTxLimit   math.Int        `protobuf:"bytes,6,opt,name=single_tx_limit,json=singleTxLimit,proto3,customtype=cosmossdk.io/math.Int" json:"single_tx_limit"`
        DailyUsed       math.Int        `protobuf:"bytes,7,opt,name=daily_used,json=dailyUsed,proto3,customtype=cosmossdk.io/math.Int" json:"daily_used"`
        DailyResetAt    int64           `protobuf:"varint,8,opt,name=daily_reset_at,json=dailyResetAt,proto3" json:"daily_reset_at"`
        TotalMinted     math.Int        `protobuf:"bytes,9,opt,name=total_minted,json=totalMinted,proto3,customtype=cosmossdk.io/math.Int" json:"total_minted"`
        AllowedPurposes []string        `protobuf:"bytes,10,rep,name=allowed_purposes,json=allowedPurposes,proto3" json:"allowed_purposes"`
        Jurisdictions   []string        `protobuf:"bytes,11,rep,name=jurisdictions,proto3" json:"jurisdictions"`
        GrantedBy       string          `protobuf:"bytes,12,opt,name=granted_by,json=grantedBy,proto3" json:"granted_by"`
        GrantedAt       int64           `protobuf:"varint,13,opt,name=granted_at,json=grantedAt,proto3" json:"granted_at"`
        ExpiresAt       int64           `protobuf:"varint,14,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at"`
        UpdatedAt       int64           `protobuf:"varint,15,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
        Address         string          `protobuf:"bytes,16,opt,name=address,proto3" json:"address"`
}

func NewTSCMintAuthority(
        authorityId string,
        tenitesId string,
        name string,
        dailyLimit math.Int,
        singleTxLimit math.Int,
        allowedPurposes []string,
        jurisdictions []string,
        grantedBy string,
        blockTime int64,
        expiresAt int64,
) TSCMintAuthority {
        return TSCMintAuthority{
                AuthorityId:     authorityId,
                TenitesId:       tenitesId,
                Name:            name,
                Status:          AuthorityStatusActive,
                DailyLimit:      dailyLimit,
                SingleTxLimit:   singleTxLimit,
                DailyUsed:       math.ZeroInt(),
                DailyResetAt:    blockTime + 86400,
                TotalMinted:     math.ZeroInt(),
                AllowedPurposes: allowedPurposes,
                Jurisdictions:   jurisdictions,
                GrantedBy:       grantedBy,
                GrantedAt:       blockTime,
                ExpiresAt:       expiresAt,
                UpdatedAt:       blockTime,
                Address:         grantedBy,
        }
}

func (a *TSCMintAuthority) Validate() error {
        if a.AuthorityId == "" {
                return fmt.Errorf("authority_id cannot be empty")
        }
        if a.TenitesId == "" {
                return fmt.Errorf("tenites_id cannot be empty")
        }
        if a.Name == "" {
                return fmt.Errorf("name cannot be empty")
        }
        if a.DailyLimit.IsNegative() {
                return fmt.Errorf("daily_limit cannot be negative")
        }
        if a.SingleTxLimit.IsNegative() {
                return fmt.Errorf("single_tx_limit cannot be negative")
        }
        if a.SingleTxLimit.GT(a.DailyLimit) {
                return fmt.Errorf("single_tx_limit cannot exceed daily_limit")
        }
        if len(a.AllowedPurposes) == 0 {
                return fmt.Errorf("at least one purpose must be allowed")
        }
        return nil
}

func (a *TSCMintAuthority) IsActive() bool {
        return a.Status == AuthorityStatusActive
}

func (a *TSCMintAuthority) CanMint(amount math.Int, purpose string, jurisdiction string, currentTime int64) error {
        if !a.IsActive() {
                return fmt.Errorf("AUTHORITY_NOT_ACTIVE: authority is %s", a.Status)
        }

        if a.ExpiresAt > 0 && currentTime >= a.ExpiresAt {
                return fmt.Errorf("AUTHORITY_EXPIRED: expired at %d", a.ExpiresAt)
        }

        if amount.GT(a.SingleTxLimit) {
                return fmt.Errorf("SINGLE_TX_LIMIT_EXCEEDED: limit %s, requested %s", a.SingleTxLimit.String(), amount.String())
        }

        purposeAllowed := false
        for _, p := range a.AllowedPurposes {
                if p == purpose {
                        purposeAllowed = true
                        break
                }
        }
        if !purposeAllowed {
                return fmt.Errorf("PURPOSE_NOT_ALLOWED: %s", purpose)
        }

        jurisdictionAllowed := false
        for _, j := range a.Jurisdictions {
                if j == jurisdiction || j == "*" {
                        jurisdictionAllowed = true
                        break
                }
        }
        if !jurisdictionAllowed {
                return fmt.Errorf("JURISDICTION_NOT_ALLOWED: %s", jurisdiction)
        }

        if currentTime < a.DailyResetAt {
                newUsed := a.DailyUsed.Add(amount)
                if newUsed.GT(a.DailyLimit) {
                        return fmt.Errorf("DAILY_LIMIT_EXCEEDED: limit %s, used %s, requested %s",
                                a.DailyLimit.String(), a.DailyUsed.String(), amount.String())
                }
        }

        return nil
}

func (a *TSCMintAuthority) RecordMint(amount math.Int, currentTime int64) {
        if currentTime >= a.DailyResetAt {
                a.DailyUsed = amount
                a.DailyResetAt = currentTime + 86400
        } else {
                a.DailyUsed = a.DailyUsed.Add(amount)
        }
        a.TotalMinted = a.TotalMinted.Add(amount)
        a.UpdatedAt = currentTime
}

func (a *TSCMintAuthority) Suspend(currentTime int64) {
        a.Status = AuthorityStatusSuspended
        a.UpdatedAt = currentTime
}

func (a *TSCMintAuthority) Revoke(currentTime int64) {
        a.Status = AuthorityStatusRevoked
        a.UpdatedAt = currentTime
}

func (a *TSCMintAuthority) Reactivate(currentTime int64) {
        a.Status = AuthorityStatusActive
        a.UpdatedAt = currentTime
}

func (a *TSCMintAuthority) Reset() {
}

func (a *TSCMintAuthority) ProtoMessage() {}
func (a *TSCMintAuthority) String() string { return "TSCMintAuthority" }

type TSCBurnAuthority struct {
        AuthorityId    string          `protobuf:"bytes,1,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
        TenitesId      string          `protobuf:"bytes,2,opt,name=tenites_id,json=tenitesId,proto3" json:"tenites_id"`
        Name           string          `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`
        Status         AuthorityStatus `protobuf:"bytes,4,opt,name=status,proto3,casttype=AuthorityStatus" json:"status"`
        DailyLimit     math.Int        `protobuf:"bytes,5,opt,name=daily_limit,json=dailyLimit,proto3,customtype=cosmossdk.io/math.Int" json:"daily_limit"`
        SingleTxLimit  math.Int        `protobuf:"bytes,6,opt,name=single_tx_limit,json=singleTxLimit,proto3,customtype=cosmossdk.io/math.Int" json:"single_tx_limit"`
        DailyUsed      math.Int        `protobuf:"bytes,7,opt,name=daily_used,json=dailyUsed,proto3,customtype=cosmossdk.io/math.Int" json:"daily_used"`
        DailyResetAt   int64           `protobuf:"varint,8,opt,name=daily_reset_at,json=dailyResetAt,proto3" json:"daily_reset_at"`
        TotalBurned    math.Int        `protobuf:"bytes,9,opt,name=total_burned,json=totalBurned,proto3,customtype=cosmossdk.io/math.Int" json:"total_burned"`
        AllowedReasons []string        `protobuf:"bytes,10,rep,name=allowed_reasons,json=allowedReasons,proto3" json:"allowed_reasons"`
        Jurisdictions  []string        `protobuf:"bytes,11,rep,name=jurisdictions,proto3" json:"jurisdictions"`
        GrantedBy      string          `protobuf:"bytes,12,opt,name=granted_by,json=grantedBy,proto3" json:"granted_by"`
        GrantedAt      int64           `protobuf:"varint,13,opt,name=granted_at,json=grantedAt,proto3" json:"granted_at"`
        ExpiresAt      int64           `protobuf:"varint,14,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at"`
        UpdatedAt      int64           `protobuf:"varint,15,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
        Address        string          `protobuf:"bytes,16,opt,name=address,proto3" json:"address"`
}

func NewTSCBurnAuthority(
        authorityId string,
        tenitesId string,
        name string,
        dailyLimit math.Int,
        singleTxLimit math.Int,
        allowedReasons []string,
        jurisdictions []string,
        grantedBy string,
        blockTime int64,
        expiresAt int64,
) TSCBurnAuthority {
        return TSCBurnAuthority{
                AuthorityId:    authorityId,
                TenitesId:      tenitesId,
                Name:           name,
                Status:         AuthorityStatusActive,
                DailyLimit:     dailyLimit,
                SingleTxLimit:  singleTxLimit,
                DailyUsed:      math.ZeroInt(),
                DailyResetAt:   blockTime + 86400,
                TotalBurned:    math.ZeroInt(),
                AllowedReasons: allowedReasons,
                Jurisdictions:  jurisdictions,
                GrantedBy:      grantedBy,
                GrantedAt:      blockTime,
                ExpiresAt:      expiresAt,
                UpdatedAt:      blockTime,
                Address:        grantedBy,
        }
}

func (a *TSCBurnAuthority) Validate() error {
        if a.AuthorityId == "" {
                return fmt.Errorf("authority_id cannot be empty")
        }
        if a.TenitesId == "" {
                return fmt.Errorf("tenites_id cannot be empty")
        }
        if a.Name == "" {
                return fmt.Errorf("name cannot be empty")
        }
        if a.DailyLimit.IsNegative() {
                return fmt.Errorf("daily_limit cannot be negative")
        }
        if a.SingleTxLimit.IsNegative() {
                return fmt.Errorf("single_tx_limit cannot be negative")
        }
        if a.SingleTxLimit.GT(a.DailyLimit) {
                return fmt.Errorf("single_tx_limit cannot exceed daily_limit")
        }
        if len(a.AllowedReasons) == 0 {
                return fmt.Errorf("at least one reason must be allowed")
        }
        return nil
}

func (a *TSCBurnAuthority) IsActive() bool {
        return a.Status == AuthorityStatusActive
}

func (a *TSCBurnAuthority) CanBurn(amount math.Int, reason string, jurisdiction string, currentTime int64) error {
        if !a.IsActive() {
                return fmt.Errorf("AUTHORITY_NOT_ACTIVE: authority is %s", a.Status)
        }

        if a.ExpiresAt > 0 && currentTime >= a.ExpiresAt {
                return fmt.Errorf("AUTHORITY_EXPIRED: expired at %d", a.ExpiresAt)
        }

        if amount.GT(a.SingleTxLimit) {
                return fmt.Errorf("SINGLE_TX_LIMIT_EXCEEDED: limit %s, requested %s", a.SingleTxLimit.String(), amount.String())
        }

        reasonAllowed := false
        for _, r := range a.AllowedReasons {
                if r == reason {
                        reasonAllowed = true
                        break
                }
        }
        if !reasonAllowed {
                return fmt.Errorf("REASON_NOT_ALLOWED: %s", reason)
        }

        jurisdictionAllowed := false
        for _, j := range a.Jurisdictions {
                if j == jurisdiction || j == "*" {
                        jurisdictionAllowed = true
                        break
                }
        }
        if !jurisdictionAllowed {
                return fmt.Errorf("JURISDICTION_NOT_ALLOWED: %s", jurisdiction)
        }

        if currentTime < a.DailyResetAt {
                newUsed := a.DailyUsed.Add(amount)
                if newUsed.GT(a.DailyLimit) {
                        return fmt.Errorf("DAILY_LIMIT_EXCEEDED: limit %s, used %s, requested %s",
                                a.DailyLimit.String(), a.DailyUsed.String(), amount.String())
                }
        }

        return nil
}

func (a *TSCBurnAuthority) RecordBurn(amount math.Int, currentTime int64) {
        if currentTime >= a.DailyResetAt {
                a.DailyUsed = amount
                a.DailyResetAt = currentTime + 86400
        } else {
                a.DailyUsed = a.DailyUsed.Add(amount)
        }
        a.TotalBurned = a.TotalBurned.Add(amount)
        a.UpdatedAt = currentTime
}

func (a *TSCBurnAuthority) Suspend(currentTime int64) {
        a.Status = AuthorityStatusSuspended
        a.UpdatedAt = currentTime
}

func (a *TSCBurnAuthority) Revoke(currentTime int64) {
        a.Status = AuthorityStatusRevoked
        a.UpdatedAt = currentTime
}

func (a *TSCBurnAuthority) Reactivate(currentTime int64) {
        a.Status = AuthorityStatusActive
        a.UpdatedAt = currentTime
}

func (a *TSCBurnAuthority) Reset() {
}

func (a *TSCBurnAuthority) ProtoMessage() {}
func (a *TSCBurnAuthority) String() string { return "TSCBurnAuthority" }

type TSCSupplySnapshot struct {
        SnapshotId        string   `protobuf:"bytes,1,opt,name=snapshot_id,json=snapshotId,proto3" json:"snapshot_id"`
        TotalSupply       math.Int `protobuf:"bytes,2,opt,name=total_supply,json=totalSupply,proto3,customtype=cosmossdk.io/math.Int" json:"total_supply"`
        CirculatingSupply math.Int `protobuf:"bytes,3,opt,name=circulating_supply,json=circulatingSupply,proto3,customtype=cosmossdk.io/math.Int" json:"circulating_supply"`
        LockedSupply      math.Int `protobuf:"bytes,4,opt,name=locked_supply,json=lockedSupply,proto3,customtype=cosmossdk.io/math.Int" json:"locked_supply"`
        DailyMinted       math.Int `protobuf:"bytes,5,opt,name=daily_minted,json=dailyMinted,proto3,customtype=cosmossdk.io/math.Int" json:"daily_minted"`
        DailyBurned       math.Int `protobuf:"bytes,6,opt,name=daily_burned,json=dailyBurned,proto3,customtype=cosmossdk.io/math.Int" json:"daily_burned"`
        NetDailyChange    math.Int `protobuf:"bytes,7,opt,name=net_daily_change,json=netDailyChange,proto3,customtype=cosmossdk.io/math.Int" json:"net_daily_change"`
        SnapshotDate      string   `protobuf:"bytes,8,opt,name=snapshot_date,json=snapshotDate,proto3" json:"snapshot_date"`
        CreatedAt         int64    `protobuf:"varint,9,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
}

func NewTSCSupplySnapshot(
        snapshotId string,
        totalSupply math.Int,
        circulatingSupply math.Int,
        lockedSupply math.Int,
        dailyMinted math.Int,
        dailyBurned math.Int,
        snapshotDate string,
        blockTime int64,
) TSCSupplySnapshot {
        return TSCSupplySnapshot{
                SnapshotId:        snapshotId,
                TotalSupply:       totalSupply,
                CirculatingSupply: circulatingSupply,
                LockedSupply:      lockedSupply,
                DailyMinted:       dailyMinted,
                DailyBurned:       dailyBurned,
                NetDailyChange:    dailyMinted.Sub(dailyBurned),
                SnapshotDate:      snapshotDate,
                CreatedAt:         blockTime,
        }
}

func (s *TSCSupplySnapshot) Reset() {
}

func (s *TSCSupplySnapshot) ProtoMessage() {}
func (s *TSCSupplySnapshot) String() string { return "TSCSupplySnapshot" }

type TSCOperation struct {
        OperationId   string   `protobuf:"bytes,1,opt,name=operation_id,json=operationId,proto3" json:"operation_id"`
        OperationType string   `protobuf:"bytes,2,opt,name=operation_type,json=operationType,proto3" json:"operation_type"`
        AuthorityId   string   `protobuf:"bytes,3,opt,name=authority_id,json=authorityId,proto3" json:"authority_id"`
        WalletId      string   `protobuf:"bytes,4,opt,name=wallet_id,json=walletId,proto3" json:"wallet_id"`
        Amount        math.Int `protobuf:"bytes,5,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount"`
        Purpose       string   `protobuf:"bytes,6,opt,name=purpose,proto3" json:"purpose"`
        Jurisdiction  string   `protobuf:"bytes,7,opt,name=jurisdiction,proto3" json:"jurisdiction"`
        ComplianceRef string   `protobuf:"bytes,8,opt,name=compliance_ref,json=complianceRef,proto3" json:"compliance_ref"`
        TxHash        string   `protobuf:"bytes,9,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash"`
        Status        string   `protobuf:"bytes,10,opt,name=status,proto3" json:"status"`
        FailureReason string   `protobuf:"bytes,11,opt,name=failure_reason,json=failureReason,proto3" json:"failure_reason"`
        CreatedAt     int64    `protobuf:"varint,12,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
        CompletedAt   int64    `protobuf:"varint,13,opt,name=completed_at,json=completedAt,proto3" json:"completed_at"`
}

func NewTSCMintOperation(
        operationId string,
        authorityId string,
        walletId string,
        amount math.Int,
        purpose string,
        jurisdiction string,
        blockTime int64,
) TSCOperation {
        return TSCOperation{
                OperationId:   operationId,
                OperationType: "mint",
                AuthorityId:   authorityId,
                WalletId:      walletId,
                Amount:        amount,
                Purpose:       purpose,
                Jurisdiction:  jurisdiction,
                Status:        "pending",
                CreatedAt:     blockTime,
        }
}

func NewTSCBurnOperation(
        operationId string,
        authorityId string,
        walletId string,
        amount math.Int,
        reason string,
        jurisdiction string,
        blockTime int64,
) TSCOperation {
        return TSCOperation{
                OperationId:   operationId,
                OperationType: "burn",
                AuthorityId:   authorityId,
                WalletId:      walletId,
                Amount:        amount,
                Purpose:       reason,
                Jurisdiction:  jurisdiction,
                Status:        "pending",
                CreatedAt:     blockTime,
        }
}

func NewTSCTransferOperation(
        operationId string,
        fromWalletId string,
        toWalletId string,
        amount math.Int,
        purpose string,
        jurisdiction string,
        blockTime int64,
) TSCOperation {
        return TSCOperation{
                OperationId:   operationId,
                OperationType: "transfer",
                WalletId:      fromWalletId + "->" + toWalletId,
                Amount:        amount,
                Purpose:       purpose,
                Jurisdiction:  jurisdiction,
                Status:        "pending",
                CreatedAt:     blockTime,
        }
}

func (o *TSCOperation) Complete(txHash string, blockTime int64) {
        o.Status = "completed"
        o.TxHash = txHash
        o.CompletedAt = blockTime
}

func (o *TSCOperation) Fail(reason string, blockTime int64) {
        o.Status = "failed"
        o.FailureReason = reason
        o.CompletedAt = blockTime
}

func (o *TSCOperation) Reset() {
}

func (o *TSCOperation) ProtoMessage() {}
func (o *TSCOperation) String() string { return "TSCOperation" }

type GenesisState struct {
        Params          TSCParams           `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
        MintAuthorities []TSCMintAuthority  `protobuf:"bytes,2,rep,name=mint_authorities,json=mintAuthorities,proto3" json:"mint_authorities"`
        BurnAuthorities []TSCBurnAuthority  `protobuf:"bytes,3,rep,name=burn_authorities,json=burnAuthorities,proto3" json:"burn_authorities"`
        Operations      []TSCOperation      `protobuf:"bytes,4,rep,name=operations,proto3" json:"operations"`
        Snapshots       []TSCSupplySnapshot `protobuf:"bytes,5,rep,name=snapshots,proto3" json:"snapshots"`
        TotalSupply     math.Int            `protobuf:"bytes,6,opt,name=total_supply,json=totalSupply,proto3,customtype=cosmossdk.io/math.Int" json:"total_supply"`
        Counter         uint64              `protobuf:"varint,7,opt,name=counter,proto3" json:"counter"`
}

func DefaultGenesisState() GenesisState {
        return GenesisState{
                Params:          DefaultTSCParams(),
                MintAuthorities: []TSCMintAuthority{},
                BurnAuthorities: []TSCBurnAuthority{},
                Operations:      []TSCOperation{},
                Snapshots:       []TSCSupplySnapshot{},
                TotalSupply:     math.ZeroInt(),
                Counter:         0,
        }
}

func (g *GenesisState) Reset() {
}

func (g *GenesisState) ProtoMessage() {}
func (g *GenesisState) String() string { return "GenesisState" }
