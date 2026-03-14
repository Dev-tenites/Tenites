package types

const (
        ModuleName = "tsc"
        StoreKey   = ModuleName
        RouterKey  = ModuleName
)

var (
        ParamsKey              = []byte{0x01}
        MintAuthorityPrefix    = []byte{0x02}
        BurnAuthorityPrefix    = []byte{0x03}
        OperationPrefix        = []byte{0x04}
        SupplySnapshotPrefix   = []byte{0x05}
        TotalSupplyKey         = []byte{0x06}
        DailyMintedKey         = []byte{0x07}
        DailyBurnedKey         = []byte{0x08}
        DailyResetAtKey        = []byte{0x09}
        CounterKey             = []byte{0x0A}
        
        TenitesIdMintAuthorityIndexPrefix = []byte{0x10}
        TenitesIdBurnAuthorityIndexPrefix = []byte{0x11}
        WalletOperationIndexPrefix        = []byte{0x12}

        CommitProofPrefix                 = []byte{0x20}
        UsedNoncePrefix                   = []byte{0x21}
        PolicyRequirementPrefix           = []byte{0x22}
        ObligationPrefix                  = []byte{0x23}
        ObligationStatusIndexPrefix       = []byte{0x24}
        InvariantViolationPrefix          = []byte{0x25}
        UsedRequestIdPrefix               = []byte{0x26}
        TrustedCustodianPrefix            = []byte{0x30}
        CustodianNoncePrefix              = []byte{0x31}
)

func MintAuthorityKey(authorityId string) []byte {
        return append(MintAuthorityPrefix, []byte(authorityId)...)
}

func BurnAuthorityKey(authorityId string) []byte {
        return append(BurnAuthorityPrefix, []byte(authorityId)...)
}

func OperationKey(operationId string) []byte {
        return append(OperationPrefix, []byte(operationId)...)
}

func SupplySnapshotKey(snapshotId string) []byte {
        return append(SupplySnapshotPrefix, []byte(snapshotId)...)
}

func TenitesIdMintAuthorityIndexKey(tenitesId string, authorityId string) []byte {
        return append(append(TenitesIdMintAuthorityIndexPrefix, []byte(tenitesId+"/")...), []byte(authorityId)...)
}

func TenitesIdBurnAuthorityIndexKey(tenitesId string, authorityId string) []byte {
        return append(append(TenitesIdBurnAuthorityIndexPrefix, []byte(tenitesId+"/")...), []byte(authorityId)...)
}

func WalletOperationIndexKey(walletId string, operationId string) []byte {
        return append(append(WalletOperationIndexPrefix, []byte(walletId+"/")...), []byte(operationId)...)
}

func CommitProofKey(proofId string) []byte {
        return append(CommitProofPrefix, []byte(proofId)...)
}

func UsedNonceKey(nonce string) []byte {
        return append(UsedNoncePrefix, []byte(nonce)...)
}

func UsedRequestIdKey(requestId string) []byte {
        return append(UsedRequestIdPrefix, []byte(requestId)...)
}

func TrustedCustodianKey(custodianId string) []byte {
        return append(TrustedCustodianPrefix, []byte(custodianId)...)
}

func CustodianNonceKey(custodianId string) []byte {
        return append(CustodianNoncePrefix, []byte(custodianId)...)
}

func PolicyRequirementKey(policySetId string) []byte {
        return append(PolicyRequirementPrefix, []byte(policySetId)...)
}

func ObligationKey(obligationId string) []byte {
        return append(ObligationPrefix, []byte(obligationId)...)
}

func ObligationStatusIndexKey(status string, obligationId string) []byte {
        return append(append(ObligationStatusIndexPrefix, []byte(status+"/")...), []byte(obligationId)...)
}

func InvariantViolationKey(violationId string) []byte {
        return append(InvariantViolationPrefix, []byte(violationId)...)
}
