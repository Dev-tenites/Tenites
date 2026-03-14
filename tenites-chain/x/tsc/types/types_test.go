package types

import (
        "testing"

        "cosmossdk.io/math"
        sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestTSCParams_Validate(t *testing.T) {
        tests := []struct {
                name    string
                params  TSCParams
                wantErr bool
        }{
                {
                        name:    "default params are valid",
                        params:  DefaultTSCParams(),
                        wantErr: false,
                },
                {
                        name: "empty denom",
                        params: TSCParams{
                                Denom:                "",
                                Decimals:             6,
                                TotalSupplyCap:       math.NewInt(1000000),
                                DailyMintCap:         math.NewInt(100000),
                                DailyBurnCap:         math.NewInt(100000),
                                SingleMintLimit:      math.NewInt(10000),
                                SingleBurnLimit:      math.NewInt(10000),
                                AllowedWalletTypes:   []string{"treasury"},
                                AllowedJurisdictions: []string{"NG"},
                        },
                        wantErr: true,
                },
                {
                        name: "decimals out of range",
                        params: TSCParams{
                                Denom:                "utsc",
                                Decimals:             20,
                                TotalSupplyCap:       math.NewInt(1000000),
                                DailyMintCap:         math.NewInt(100000),
                                DailyBurnCap:         math.NewInt(100000),
                                SingleMintLimit:      math.NewInt(10000),
                                SingleBurnLimit:      math.NewInt(10000),
                                AllowedWalletTypes:   []string{"treasury"},
                                AllowedJurisdictions: []string{"NG"},
                        },
                        wantErr: true,
                },
                {
                        name: "single limit exceeds daily",
                        params: TSCParams{
                                Denom:                "utsc",
                                Decimals:             6,
                                TotalSupplyCap:       math.NewInt(1000000),
                                DailyMintCap:         math.NewInt(100000),
                                DailyBurnCap:         math.NewInt(100000),
                                SingleMintLimit:      math.NewInt(200000),
                                SingleBurnLimit:      math.NewInt(10000),
                                AllowedWalletTypes:   []string{"treasury"},
                                AllowedJurisdictions: []string{"NG"},
                        },
                        wantErr: true,
                },
                {
                        name: "no wallet types",
                        params: TSCParams{
                                Denom:                "utsc",
                                Decimals:             6,
                                TotalSupplyCap:       math.NewInt(1000000),
                                DailyMintCap:         math.NewInt(100000),
                                DailyBurnCap:         math.NewInt(100000),
                                SingleMintLimit:      math.NewInt(10000),
                                SingleBurnLimit:      math.NewInt(10000),
                                AllowedWalletTypes:   []string{},
                                AllowedJurisdictions: []string{"NG"},
                        },
                        wantErr: true,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := tt.params.Validate()
                        if (err != nil) != tt.wantErr {
                                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
                        }
                })
        }
}

func TestTSCMintAuthority_CanMint(t *testing.T) {
        blockTime := int64(1000000)
        
        authority := NewTSCMintAuthority(
                "MINT001",
                "TEN001",
                "Test Minter",
                math.NewInt(1000000),
                math.NewInt(100000),
                []string{"issuance", "liquidity"},
                []string{"NG", "GH"},
                "governance",
                blockTime,
                blockTime + 86400*365,
        )

        tests := []struct {
                name         string
                amount       math.Int
                purpose      string
                jurisdiction string
                currentTime  int64
                wantErr      bool
                errContains  string
        }{
                {
                        name:         "valid mint",
                        amount:       math.NewInt(50000),
                        purpose:      "issuance",
                        jurisdiction: "NG",
                        currentTime:  blockTime + 100,
                        wantErr:      false,
                },
                {
                        name:         "exceeds single tx limit",
                        amount:       math.NewInt(200000),
                        purpose:      "issuance",
                        jurisdiction: "NG",
                        currentTime:  blockTime + 100,
                        wantErr:      true,
                        errContains:  "SINGLE_TX_LIMIT_EXCEEDED",
                },
                {
                        name:         "purpose not allowed",
                        amount:       math.NewInt(50000),
                        purpose:      "settlement",
                        jurisdiction: "NG",
                        currentTime:  blockTime + 100,
                        wantErr:      true,
                        errContains:  "PURPOSE_NOT_ALLOWED",
                },
                {
                        name:         "jurisdiction not allowed",
                        amount:       math.NewInt(50000),
                        purpose:      "issuance",
                        jurisdiction: "KE",
                        currentTime:  blockTime + 100,
                        wantErr:      true,
                        errContains:  "JURISDICTION_NOT_ALLOWED",
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := authority.CanMint(tt.amount, tt.purpose, tt.jurisdiction, tt.currentTime)
                        if (err != nil) != tt.wantErr {
                                t.Errorf("CanMint() error = %v, wantErr %v", err, tt.wantErr)
                        }
                        if tt.wantErr && tt.errContains != "" && err != nil {
                                if !contains(err.Error(), tt.errContains) {
                                        t.Errorf("CanMint() error = %v, should contain %s", err, tt.errContains)
                                }
                        }
                })
        }
}

func TestTSCMintAuthority_DailyLimit(t *testing.T) {
        blockTime := int64(1000000)
        
        authority := NewTSCMintAuthority(
                "MINT001",
                "TEN001",
                "Test Minter",
                math.NewInt(100000),
                math.NewInt(50000),
                []string{"issuance"},
                []string{"*"},
                "governance",
                blockTime,
                0,
        )

        authority.RecordMint(math.NewInt(60000), blockTime+100)

        err := authority.CanMint(math.NewInt(50000), "issuance", "NG", blockTime+200)
        if err == nil {
                t.Error("Expected daily limit exceeded error")
        }
        if !contains(err.Error(), "DAILY_LIMIT_EXCEEDED") {
                t.Errorf("Expected DAILY_LIMIT_EXCEEDED, got: %v", err)
        }

        err = authority.CanMint(math.NewInt(30000), "issuance", "NG", blockTime+200)
        if err != nil {
                t.Errorf("Expected no error for amount within limit, got: %v", err)
        }
}

func TestTSCMintAuthority_DailyReset(t *testing.T) {
        blockTime := int64(1000000)
        
        authority := NewTSCMintAuthority(
                "MINT001",
                "TEN001",
                "Test Minter",
                math.NewInt(100000),
                math.NewInt(100000),
                []string{"issuance"},
                []string{"*"},
                "governance",
                blockTime,
                0,
        )

        authority.RecordMint(math.NewInt(90000), blockTime+100)

        err := authority.CanMint(math.NewInt(50000), "issuance", "NG", blockTime+200)
        if err == nil {
                t.Error("Expected daily limit exceeded error before reset")
        }

        err = authority.CanMint(math.NewInt(50000), "issuance", "NG", blockTime+86400+100)
        if err != nil {
                t.Errorf("Expected no error after daily reset, got: %v", err)
        }
}

func TestTSCMintAuthority_Expiration(t *testing.T) {
        blockTime := int64(1000000)
        expiresAt := blockTime + 3600
        
        authority := NewTSCMintAuthority(
                "MINT001",
                "TEN001",
                "Test Minter",
                math.NewInt(100000),
                math.NewInt(100000),
                []string{"issuance"},
                []string{"*"},
                "governance",
                blockTime,
                expiresAt,
        )

        err := authority.CanMint(math.NewInt(1000), "issuance", "NG", blockTime+100)
        if err != nil {
                t.Errorf("Expected no error before expiration, got: %v", err)
        }

        err = authority.CanMint(math.NewInt(1000), "issuance", "NG", expiresAt+100)
        if err == nil {
                t.Error("Expected authority expired error")
        }
        if !contains(err.Error(), "AUTHORITY_EXPIRED") {
                t.Errorf("Expected AUTHORITY_EXPIRED, got: %v", err)
        }
}

func TestTSCBurnAuthority_CanBurn(t *testing.T) {
        blockTime := int64(1000000)
        
        authority := NewTSCBurnAuthority(
                "BURN001",
                "TEN001",
                "Test Burner",
                math.NewInt(1000000),
                math.NewInt(100000),
                []string{"redemption", "recall"},
                []string{"NG", "GH"},
                "governance",
                blockTime,
                0,
        )

        tests := []struct {
                name         string
                amount       math.Int
                reason       string
                jurisdiction string
                currentTime  int64
                wantErr      bool
                errContains  string
        }{
                {
                        name:         "valid burn",
                        amount:       math.NewInt(50000),
                        reason:       "redemption",
                        jurisdiction: "NG",
                        currentTime:  blockTime + 100,
                        wantErr:      false,
                },
                {
                        name:         "exceeds single tx limit",
                        amount:       math.NewInt(200000),
                        reason:       "redemption",
                        jurisdiction: "NG",
                        currentTime:  blockTime + 100,
                        wantErr:      true,
                        errContains:  "SINGLE_TX_LIMIT_EXCEEDED",
                },
                {
                        name:         "reason not allowed",
                        amount:       math.NewInt(50000),
                        reason:       "expiry",
                        jurisdiction: "NG",
                        currentTime:  blockTime + 100,
                        wantErr:      true,
                        errContains:  "REASON_NOT_ALLOWED",
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := authority.CanBurn(tt.amount, tt.reason, tt.jurisdiction, tt.currentTime)
                        if (err != nil) != tt.wantErr {
                                t.Errorf("CanBurn() error = %v, wantErr %v", err, tt.wantErr)
                        }
                })
        }
}

func TestTSCMintAuthority_StatusTransitions(t *testing.T) {
        blockTime := int64(1000000)
        
        authority := NewTSCMintAuthority(
                "MINT001",
                "TEN001",
                "Test Minter",
                math.NewInt(100000),
                math.NewInt(100000),
                []string{"issuance"},
                []string{"*"},
                "governance",
                blockTime,
                0,
        )

        if !authority.IsActive() {
                t.Error("New authority should be active")
        }

        authority.Suspend(blockTime + 100)
        if authority.IsActive() {
                t.Error("Suspended authority should not be active")
        }
        if authority.Status != AuthorityStatusSuspended {
                t.Errorf("Expected status suspended, got %s", authority.Status)
        }

        authority.Reactivate(blockTime + 200)
        if !authority.IsActive() {
                t.Error("Reactivated authority should be active")
        }

        authority.Revoke(blockTime + 300)
        if authority.IsActive() {
                t.Error("Revoked authority should not be active")
        }
        if authority.Status != AuthorityStatusRevoked {
                t.Errorf("Expected status revoked, got %s", authority.Status)
        }
}

func TestTSCSupplySnapshot(t *testing.T) {
        blockTime := int64(1000000)
        
        snapshot := NewTSCSupplySnapshot(
                "SNAP_2024-01-01",
                math.NewInt(10000000),
                math.NewInt(8000000),
                math.NewInt(2000000),
                math.NewInt(500000),
                math.NewInt(200000),
                "2024-01-01",
                blockTime,
        )

        expectedNet := math.NewInt(300000)
        if !snapshot.NetDailyChange.Equal(expectedNet) {
                t.Errorf("Expected net daily change %s, got %s", expectedNet.String(), snapshot.NetDailyChange.String())
        }
}

func TestTSCOperation(t *testing.T) {
        blockTime := int64(1000000)
        
        mintOp := NewTSCMintOperation(
                "OP001",
                "MINT001",
                "WALLET001",
                math.NewInt(1000000),
                "issuance",
                "NG",
                blockTime,
        )

        if mintOp.OperationType != "mint" {
                t.Errorf("Expected operation type mint, got %s", mintOp.OperationType)
        }
        if mintOp.Status != "pending" {
                t.Errorf("Expected status pending, got %s", mintOp.Status)
        }

        mintOp.Complete("txhash123", blockTime+100)
        if mintOp.Status != "completed" {
                t.Errorf("Expected status completed, got %s", mintOp.Status)
        }
        if mintOp.TxHash != "txhash123" {
                t.Errorf("Expected txhash txhash123, got %s", mintOp.TxHash)
        }

        burnOp := NewTSCBurnOperation(
                "OP002",
                "BURN001",
                "WALLET001",
                math.NewInt(500000),
                "redemption",
                "NG",
                blockTime,
        )

        burnOp.Fail("insufficient balance", blockTime+100)
        if burnOp.Status != "failed" {
                t.Errorf("Expected status failed, got %s", burnOp.Status)
        }
        if burnOp.FailureReason != "insufficient balance" {
                t.Errorf("Expected failure reason 'insufficient balance', got %s", burnOp.FailureReason)
        }
}

func TestMsgMint_ValidateBasic(t *testing.T) {
        tests := []struct {
                name    string
                msg     MsgMint
                wantErr bool
        }{
                {
                        name: "valid mint msg",
                        msg: MsgMint{
                                Signer:       "tenites1abc123",
                                AuthorityId:  "MINT001",
                                WalletId:     "WALLET001",
                                Amount:       "1000000",
                                Purpose:      "issuance",
                                Jurisdiction: "NG",
                        },
                        wantErr: false,
                },
                {
                        name: "missing authority_id",
                        msg: MsgMint{
                                Signer:   "tenites1abc123",
                                WalletId: "WALLET001",
                                Amount:   "1000000",
                        },
                        wantErr: true,
                },
                {
                        name: "missing amount",
                        msg: MsgMint{
                                Signer:      "tenites1abc123",
                                AuthorityId: "MINT001",
                                WalletId:    "WALLET001",
                        },
                        wantErr: true,
                },
                {
                        name: "zero amount",
                        msg: MsgMint{
                                Signer:      "tenites1abc123",
                                AuthorityId: "MINT001",
                                WalletId:    "WALLET001",
                                Amount:      "0",
                        },
                        wantErr: true,
                },
                {
                        name: "missing wallet_id",
                        msg: MsgMint{
                                Signer:      "tenites1abc123",
                                AuthorityId: "MINT001",
                                Amount:      "1000000",
                        },
                        wantErr: true,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := tt.msg.ValidateBasic()
                        if (err != nil) != tt.wantErr {
                                t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
                        }
                })
        }
}

func TestMsgBurn_ValidateBasic(t *testing.T) {
        tests := []struct {
                name    string
                msg     MsgBurn
                wantErr bool
        }{
                {
                        name: "valid burn msg",
                        msg: MsgBurn{
                                Signer:       "tenites1abc123",
                                AuthorityId:  "BURN001",
                                WalletId:     "WALLET001",
                                Amount:       "1000000",
                                Reason:       "redemption",
                                Jurisdiction: "NG",
                        },
                        wantErr: false,
                },
                {
                        name: "missing authority_id",
                        msg: MsgBurn{
                                Signer: "tenites1abc123",
                                Amount: "1000000",
                        },
                        wantErr: true,
                },
                {
                        name: "missing amount",
                        msg: MsgBurn{
                                Signer:      "tenites1abc123",
                                AuthorityId: "BURN001",
                        },
                        wantErr: true,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := tt.msg.ValidateBasic()
                        if (err != nil) != tt.wantErr {
                                t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
                        }
                })
        }
}

func TestMsgTransfer_ValidateBasic(t *testing.T) {
        tests := []struct {
                name    string
                msg     MsgTransfer
                wantErr bool
        }{
                {
                        name: "valid transfer",
                        msg: MsgTransfer{
                                Signer:       "tenites1abc123",
                                FromWalletId: "WALLET001",
                                ToWalletId:   "WALLET002",
                                Amount:       "1000000",
                                Purpose:      "settlement",
                        },
                        wantErr: false,
                },
                {
                        name: "missing from_wallet_id",
                        msg: MsgTransfer{
                                Signer:     "tenites1abc123",
                                ToWalletId: "WALLET002",
                                Amount:     "1000000",
                        },
                        wantErr: true,
                },
                {
                        name: "missing to_wallet_id",
                        msg: MsgTransfer{
                                Signer:       "tenites1abc123",
                                FromWalletId: "WALLET001",
                                Amount:       "1000000",
                        },
                        wantErr: true,
                },
                {
                        name: "zero amount",
                        msg: MsgTransfer{
                                Signer:       "tenites1abc123",
                                FromWalletId: "WALLET001",
                                ToWalletId:   "WALLET002",
                                Amount:       "0",
                        },
                        wantErr: true,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := tt.msg.ValidateBasic()
                        if (err != nil) != tt.wantErr {
                                t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
                        }
                })
        }
}

func TestMsgAddMintAuthority_ValidateBasic(t *testing.T) {
        tests := []struct {
                name    string
                msg     MsgAddMintAuthority
                wantErr bool
        }{
                {
                        name: "valid add mint authority",
                        msg: MsgAddMintAuthority{
                                Signer:    "tenites1abc123",
                                TenitesId: "TEN001",
                                Name:      "Treasury Minter",
                        },
                        wantErr: false,
                },
                {
                        name: "missing tenites_id",
                        msg: MsgAddMintAuthority{
                                Signer: "tenites1abc123",
                                Name:   "Treasury Minter",
                        },
                        wantErr: true,
                },
                {
                        name: "missing name",
                        msg: MsgAddMintAuthority{
                                Signer:    "tenites1abc123",
                                TenitesId: "TEN001",
                        },
                        wantErr: true,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := tt.msg.ValidateBasic()
                        if (err != nil) != tt.wantErr {
                                t.Errorf("ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
                        }
                })
        }
}

func TestMsgRemoveMintAuthority_ValidateBasic(t *testing.T) {
        msg := MsgRemoveMintAuthority{Signer: "tenites1abc", AuthorityId: "MINT001"}
        if err := msg.ValidateBasic(); err != nil {
                t.Errorf("unexpected error: %v", err)
        }
        msg2 := MsgRemoveMintAuthority{Signer: "tenites1abc", AuthorityId: ""}
        if err := msg2.ValidateBasic(); err == nil {
                t.Error("expected error for empty authority_id")
        }
}

func TestMsgAddBurnAuthority_ValidateBasic(t *testing.T) {
        msg := MsgAddBurnAuthority{Signer: "tenites1abc", TenitesId: "TEN001", Name: "Burner"}
        if err := msg.ValidateBasic(); err != nil {
                t.Errorf("unexpected error: %v", err)
        }
        msg2 := MsgAddBurnAuthority{Signer: "tenites1abc", TenitesId: "", Name: "Burner"}
        if err := msg2.ValidateBasic(); err == nil {
                t.Error("expected error for empty tenites_id")
        }
}

func TestMsgRemoveBurnAuthority_ValidateBasic(t *testing.T) {
        msg := MsgRemoveBurnAuthority{Signer: "tenites1abc", AuthorityId: "BURN001"}
        if err := msg.ValidateBasic(); err != nil {
                t.Errorf("unexpected error: %v", err)
        }
        msg2 := MsgRemoveBurnAuthority{Signer: "tenites1abc", AuthorityId: ""}
        if err := msg2.ValidateBasic(); err == nil {
                t.Error("expected error for empty authority_id")
        }
}

func TestMsg_GetSigners_DualField(t *testing.T) {
        sdk.GetConfig().SetBech32PrefixForAccount("tenites", "tenitespub")

        addrBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
                0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}
        validAddr := sdk.AccAddress(addrBytes).String()

        t.Run("MsgMint signer is bech32 address, authority_id is domain field", func(t *testing.T) {
                msg := &MsgMint{
                        Signer:      validAddr,
                        AuthorityId: "MINT0000000001",
                        WalletId:    "WALLET001",
                        Amount:      "1000000",
                }
                signers := msg.GetSigners()
                if len(signers) != 1 {
                        t.Fatalf("expected 1 signer, got %d", len(signers))
                }
                if signers[0].String() != validAddr {
                        t.Errorf("signer mismatch: expected %s, got %s", validAddr, signers[0].String())
                }
                if msg.AuthorityId != "MINT0000000001" {
                        t.Errorf("authority_id should remain as domain ID, got %s", msg.AuthorityId)
                }
        })

        t.Run("MsgBurn signer is bech32 address, authority_id is domain field", func(t *testing.T) {
                msg := &MsgBurn{
                        Signer:      validAddr,
                        AuthorityId: "BURN0000000001",
                        Amount:      "500000",
                }
                signers := msg.GetSigners()
                if len(signers) != 1 {
                        t.Fatalf("expected 1 signer, got %d", len(signers))
                }
                if signers[0].String() != validAddr {
                        t.Errorf("signer mismatch: expected %s, got %s", validAddr, signers[0].String())
                }
        })

        t.Run("MsgTransfer signer is bech32 address", func(t *testing.T) {
                msg := &MsgTransfer{
                        Signer:       validAddr,
                        FromWalletId: "WALLET001",
                        ToWalletId:   "WALLET002",
                        Amount:       "100000",
                }
                signers := msg.GetSigners()
                if len(signers) != 1 {
                        t.Fatalf("expected 1 signer, got %d", len(signers))
                }
                if signers[0].String() != validAddr {
                        t.Errorf("signer mismatch: expected %s, got %s", validAddr, signers[0].String())
                }
        })

        t.Run("MsgAddMintAuthority signer is bech32 address", func(t *testing.T) {
                msg := &MsgAddMintAuthority{
                        Signer:    validAddr,
                        TenitesId: "TEN001",
                        Name:      "Treasury Minter",
                }
                signers := msg.GetSigners()
                if len(signers) != 1 {
                        t.Fatalf("expected 1 signer, got %d", len(signers))
                }
                if signers[0].String() != validAddr {
                        t.Errorf("signer mismatch: expected %s, got %s", validAddr, signers[0].String())
                }
        })

        t.Run("MsgRemoveMintAuthority signer is bech32 address", func(t *testing.T) {
                msg := &MsgRemoveMintAuthority{
                        Signer:      validAddr,
                        AuthorityId: "MINT001",
                }
                signers := msg.GetSigners()
                if len(signers) != 1 {
                        t.Fatalf("expected 1 signer, got %d", len(signers))
                }
        })

        t.Run("MsgAddBurnAuthority signer is bech32 address", func(t *testing.T) {
                msg := &MsgAddBurnAuthority{
                        Signer:    validAddr,
                        TenitesId: "TEN001",
                        Name:      "Burner",
                }
                signers := msg.GetSigners()
                if len(signers) != 1 {
                        t.Fatalf("expected 1 signer, got %d", len(signers))
                }
        })

        t.Run("MsgRemoveBurnAuthority signer is bech32 address", func(t *testing.T) {
                msg := &MsgRemoveBurnAuthority{
                        Signer:      validAddr,
                        AuthorityId: "BURN001",
                }
                signers := msg.GetSigners()
                if len(signers) != 1 {
                        t.Fatalf("expected 1 signer, got %d", len(signers))
                }
        })
}

func TestMsg_XXX_MessageName(t *testing.T) {
        tests := []struct {
                name     string
                msg      interface{ XXX_MessageName() string }
                expected string
        }{
                {"MsgMint", &MsgMint{}, "tenites.tsc.v1.MsgMint"},
                {"MsgBurn", &MsgBurn{}, "tenites.tsc.v1.MsgBurn"},
                {"MsgTransfer", &MsgTransfer{}, "tenites.tsc.v1.MsgTransfer"},
                {"MsgAddMintAuthority", &MsgAddMintAuthority{}, "tenites.tsc.v1.MsgAddMintAuthority"},
                {"MsgRemoveMintAuthority", &MsgRemoveMintAuthority{}, "tenites.tsc.v1.MsgRemoveMintAuthority"},
                {"MsgAddBurnAuthority", &MsgAddBurnAuthority{}, "tenites.tsc.v1.MsgAddBurnAuthority"},
                {"MsgRemoveBurnAuthority", &MsgRemoveBurnAuthority{}, "tenites.tsc.v1.MsgRemoveBurnAuthority"},
                {"MsgMintResponse", &MsgMintResponse{}, "tenites.tsc.v1.MsgMintResponse"},
                {"MsgBurnResponse", &MsgBurnResponse{}, "tenites.tsc.v1.MsgBurnResponse"},
                {"MsgTransferResponse", &MsgTransferResponse{}, "tenites.tsc.v1.MsgTransferResponse"},
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        if got := tt.msg.XXX_MessageName(); got != tt.expected {
                                t.Errorf("XXX_MessageName() = %s, want %s", got, tt.expected)
                        }
                })
        }
}

func contains(s, substr string) bool {
        return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
        for i := 0; i <= len(s)-len(substr); i++ {
                if s[i:i+len(substr)] == substr {
                        return true
                }
        }
        return false
}
