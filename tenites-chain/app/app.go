package app

import (
        "encoding/json"
        "io"
        "os"
        "path/filepath"

        dbm "github.com/cosmos/cosmos-db"
        "github.com/spf13/cast"

        "cosmossdk.io/log"
        storetypes "cosmossdk.io/store/types"

        abci "github.com/cometbft/cometbft/abci/types"
        tmos "github.com/cometbft/cometbft/libs/os"

        autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
        reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"

        "github.com/cosmos/cosmos-sdk/baseapp"
        "github.com/cosmos/cosmos-sdk/client"
        cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
        nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
        "github.com/cosmos/cosmos-sdk/codec"
        authcodec "github.com/cosmos/cosmos-sdk/codec/address"
        "github.com/cosmos/cosmos-sdk/codec/types"
        "github.com/cosmos/cosmos-sdk/runtime"
        runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
        "github.com/cosmos/cosmos-sdk/server"
        "github.com/cosmos/cosmos-sdk/server/api"
        "github.com/cosmos/cosmos-sdk/server/config"
        servertypes "github.com/cosmos/cosmos-sdk/server/types"
        sdk "github.com/cosmos/cosmos-sdk/types"
        "github.com/cosmos/cosmos-sdk/types/module"
        "github.com/cosmos/cosmos-sdk/version"
        "github.com/cosmos/cosmos-sdk/x/auth"
        authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
        authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
        authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
        "github.com/cosmos/cosmos-sdk/x/auth/vesting"
        vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
        "github.com/cosmos/cosmos-sdk/x/bank"
        bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
        banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
        "github.com/cosmos/ibc-go/modules/capability"
        capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
        capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
        "github.com/cosmos/cosmos-sdk/x/consensus"
        consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
        consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
        "github.com/cosmos/cosmos-sdk/x/crisis"
        crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
        crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
        distr "github.com/cosmos/cosmos-sdk/x/distribution"
        distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
        distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
        "github.com/cosmos/cosmos-sdk/x/genutil"
        genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
        "github.com/cosmos/cosmos-sdk/x/gov"
        govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
        govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
        govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
        govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
        "github.com/cosmos/cosmos-sdk/x/mint"
        mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
        minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
        "github.com/cosmos/cosmos-sdk/x/params"
        paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
        paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
        paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
        paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
        "github.com/cosmos/cosmos-sdk/x/slashing"
        slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
        slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
        "github.com/cosmos/cosmos-sdk/x/staking"
        stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
        stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
        "cosmossdk.io/x/upgrade"
        upgradekeeper "cosmossdk.io/x/upgrade/keeper"
        upgradetypes "cosmossdk.io/x/upgrade/types"

        ibc "github.com/cosmos/ibc-go/v8/modules/core"
        ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
        ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
        ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
        "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
        transferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
        transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
        porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"

        aakeeper "github.com/tenites/tenites-chain/x/aa/keeper"
        aatypes "github.com/tenites/tenites-chain/x/aa/types"
        auditkeeper "github.com/tenites/tenites-chain/x/audit/keeper"
        audittypes "github.com/tenites/tenites-chain/x/audit/types"
        compliancekeeper "github.com/tenites/tenites-chain/x/compliance/keeper"
        compliancetypes "github.com/tenites/tenites-chain/x/compliance/types"
        disputekeeper "github.com/tenites/tenites-chain/x/dispute/keeper"
        disputetypes "github.com/tenites/tenites-chain/x/dispute/types"
        constitutionkeeper "github.com/tenites/tenites-chain/x/governance/keeper"
        constitutiontypes "github.com/tenites/tenites-chain/x/governance/types"
        identitykeeper "github.com/tenites/tenites-chain/x/identity/keeper"
        identitytypes "github.com/tenites/tenites-chain/x/identity/types"
        intentkeeper "github.com/tenites/tenites-chain/x/intent/keeper"
        intenttypes "github.com/tenites/tenites-chain/x/intent/types"
        obligationkeeper "github.com/tenites/tenites-chain/x/obligation/keeper"
        obligationtypes "github.com/tenites/tenites-chain/x/obligation/types"
        policykeeper "github.com/tenites/tenites-chain/x/policy/keeper"
        policytypes "github.com/tenites/tenites-chain/x/policy/types"
        regulatorkeeper "github.com/tenites/tenites-chain/x/regulator/keeper"
        regulatortypes "github.com/tenites/tenites-chain/x/regulator/types"
        settlementkeeper "github.com/tenites/tenites-chain/x/settlement/keeper"
        settlementtypes "github.com/tenites/tenites-chain/x/settlement/types"
        trustkeeper "github.com/tenites/tenites-chain/x/trust/keeper"
        trusttypes "github.com/tenites/tenites-chain/x/trust/types"
        tsckeeper "github.com/tenites/tenites-chain/x/tsc/keeper"
        tsctypes "github.com/tenites/tenites-chain/x/tsc/types"
        valsetkeeper "github.com/tenites/tenites-chain/x/valset/keeper"
        valsettypes "github.com/tenites/tenites-chain/x/valset/types"
        walletkeeper "github.com/tenites/tenites-chain/x/wallet/keeper"
        wallettypes "github.com/tenites/tenites-chain/x/wallet/types"

        tenaa "github.com/tenites/tenites-chain/x/aa"
        tenaudit "github.com/tenites/tenites-chain/x/audit"
        tencompliance "github.com/tenites/tenites-chain/x/compliance"
        tendispute "github.com/tenites/tenites-chain/x/dispute"
        tenconstitution "github.com/tenites/tenites-chain/x/governance"
        tenidentity "github.com/tenites/tenites-chain/x/identity"
        tenintent "github.com/tenites/tenites-chain/x/intent"
        tenobligation "github.com/tenites/tenites-chain/x/obligation"
        tenpolicy "github.com/tenites/tenites-chain/x/policy"
        tenregulator "github.com/tenites/tenites-chain/x/regulator"
        tensettlement "github.com/tenites/tenites-chain/x/settlement"
        tentrust "github.com/tenites/tenites-chain/x/trust"
        tentsc "github.com/tenites/tenites-chain/x/tsc"
        tenvalset "github.com/tenites/tenites-chain/x/valset"
        tenwallet "github.com/tenites/tenites-chain/x/wallet"
)

const Name = "tenites"

var (
        DefaultNodeHome string

        Bech32PrefixAccAddr  = "tenites"
        Bech32PrefixAccPub   = "tenitespub"
        Bech32PrefixValAddr  = "tenitesvaloper"
        Bech32PrefixValPub   = "tenitesvaloperpub"
        Bech32PrefixConsAddr = "tenitesvalcons"
        Bech32PrefixConsPub  = "tenitesvalconspub"
)

func init() {
        userHomeDir, err := os.UserHomeDir()
        if err != nil {
                panic(err)
        }

        DefaultNodeHome = filepath.Join(userHomeDir, ".tenitesd")
}

var (
        ModuleBasics = module.NewBasicManager(
                auth.AppModuleBasic{},
                genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
                bank.AppModuleBasic{},
                capability.AppModuleBasic{},
                staking.AppModuleBasic{},
                mint.AppModuleBasic{},
                distr.AppModuleBasic{},
                gov.NewAppModuleBasic(
                        []govclient.ProposalHandler{
                                paramsclient.ProposalHandler,
                        },
                ),
                params.AppModuleBasic{},
                crisis.AppModuleBasic{},
                slashing.AppModuleBasic{},
                upgrade.AppModuleBasic{},
                consensus.AppModuleBasic{},
                vesting.AppModuleBasic{},
                ibc.AppModuleBasic{},
                transfer.AppModuleBasic{},
                ibctm.AppModuleBasic{},
                tenidentity.AppModuleBasic{},
                tenwallet.AppModuleBasic{},
                tentrust.AppModuleBasic{},
                tencompliance.AppModuleBasic{},
                tenobligation.AppModuleBasic{},
                tensettlement.AppModuleBasic{},
                tendispute.AppModuleBasic{},
                tenregulator.AppModuleBasic{},
                tenconstitution.AppModuleBasic{},
                tenaa.AppModuleBasic{},
                tenintent.AppModuleBasic{},
                tentsc.AppModuleBasic{},
                tenpolicy.AppModuleBasic{},
                tenaudit.AppModuleBasic{},
                tenvalset.AppModuleBasic{},
        )

        maccPerms = map[string][]string{
                authtypes.FeeCollectorName:     nil,
                distrtypes.ModuleName:          nil,
                minttypes.ModuleName:           {authtypes.Minter},
                stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
                stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
                govtypes.ModuleName:            {authtypes.Burner},
                ibcexported.ModuleName:         nil,
                transfertypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
                wallettypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
                settlementtypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
                tsctypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
        }
)

var _ servertypes.Application = (*TenitesApp)(nil)

type TenitesApp struct {
        *baseapp.BaseApp

        legacyAmino       *codec.LegacyAmino
        appCodec          codec.Codec
        interfaceRegistry types.InterfaceRegistry
        txConfig          client.TxConfig

        keys    map[string]*storetypes.KVStoreKey
        tkeys   map[string]*storetypes.TransientStoreKey
        memKeys map[string]*storetypes.MemoryStoreKey

        AccountKeeper    authkeeper.AccountKeeper
        BankKeeper       bankkeeper.Keeper
        CapabilityKeeper *capabilitykeeper.Keeper
        StakingKeeper    *stakingkeeper.Keeper
        IBCKeeper        *ibckeeper.Keeper
        TransferKeeper   transferkeeper.Keeper
        ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
        ScopedTransferKeeper capabilitykeeper.ScopedKeeper
        SlashingKeeper   slashingkeeper.Keeper
        MintKeeper       mintkeeper.Keeper
        DistrKeeper      distrkeeper.Keeper
        GovKeeper        *govkeeper.Keeper
        CrisisKeeper     *crisiskeeper.Keeper
        UpgradeKeeper    *upgradekeeper.Keeper
        ParamsKeeper     paramskeeper.Keeper
        ConsensusKeeper  consensuskeeper.Keeper

        IdentityKeeper     identitykeeper.Keeper
        WalletKeeper       walletkeeper.Keeper
        TrustKeeper        trustkeeper.Keeper
        ComplianceKeeper   compliancekeeper.Keeper
        ObligationKeeper   obligationkeeper.Keeper
        SettlementKeeper   settlementkeeper.Keeper
        DisputeKeeper      disputekeeper.Keeper
        RegulatorKeeper    regulatorkeeper.Keeper
        ConstitutionKeeper constitutionkeeper.Keeper
        TscKeeper          tsckeeper.Keeper
        AaKeeper           aakeeper.Keeper
        AuditKeeper        auditkeeper.Keeper
        IntentKeeper       intentkeeper.Keeper
        PolicyKeeper       policykeeper.Keeper
        ValsetKeeper       *valsetkeeper.Keeper

        ModuleManager      *module.Manager
        BasicModuleManager module.BasicManager
        configurator       module.Configurator
}

func NewTenitesApp(
        logger log.Logger,
        db dbm.DB,
        traceStore io.Writer,
        loadLatest bool,
        appOpts servertypes.AppOptions,
        baseAppOptions ...func(*baseapp.BaseApp),
) *TenitesApp {
        encodingConfig := MakeEncodingConfig()
        appCodec := encodingConfig.Codec
        legacyAmino := encodingConfig.Amino
        interfaceRegistry := encodingConfig.InterfaceRegistry
        txConfig := encodingConfig.TxConfig

        bApp := baseapp.NewBaseApp(Name, logger, db, txConfig.TxDecoder(), baseAppOptions...)
        bApp.SetCommitMultiStoreTracer(traceStore)
        bApp.SetVersion(version.Version)
        bApp.SetInterfaceRegistry(interfaceRegistry)
        bApp.SetTxEncoder(txConfig.TxEncoder())

        keys := storetypes.NewKVStoreKeys(
                authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
                minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
                govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
                capabilitytypes.StoreKey, crisistypes.StoreKey, consensustypes.StoreKey,
                ibcexported.StoreKey, transfertypes.StoreKey,
                identitytypes.StoreKey, wallettypes.StoreKey, trusttypes.StoreKey,
                compliancetypes.StoreKey, obligationtypes.StoreKey, settlementtypes.StoreKey,
                disputetypes.StoreKey, regulatortypes.StoreKey, constitutiontypes.StoreKey,
                tsctypes.StoreKey, aatypes.StoreKey, audittypes.StoreKey,
                intenttypes.StoreKey, policytypes.StoreKey, valsettypes.StoreKey,
        )
        tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
        memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

        app := &TenitesApp{
                BaseApp:           bApp,
                legacyAmino:       legacyAmino,
                appCodec:          appCodec,
                interfaceRegistry: interfaceRegistry,
                txConfig:          txConfig,
                keys:              keys,
                tkeys:             tkeys,
                memKeys:           memKeys,
        }

        app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

        app.ConsensusKeeper = consensuskeeper.NewKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[consensustypes.StoreKey]),
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                runtime.EventService{},
        )
        bApp.SetParamStore(app.ConsensusKeeper.ParamsStore)

        app.CapabilityKeeper = capabilitykeeper.NewKeeper(
                appCodec,
                keys[capabilitytypes.StoreKey],
                memKeys[capabilitytypes.MemStoreKey],
        )

        app.ScopedIBCKeeper = app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
        app.ScopedTransferKeeper = app.CapabilityKeeper.ScopeToModule(transfertypes.ModuleName)
        app.CapabilityKeeper.Seal()

        app.AccountKeeper = authkeeper.NewAccountKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[authtypes.StoreKey]),
                authtypes.ProtoBaseAccount,
                maccPerms,
                authcodec.NewBech32Codec(Bech32PrefixAccAddr),
                Bech32PrefixAccAddr,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.BankKeeper = bankkeeper.NewBaseKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[banktypes.StoreKey]),
                app.AccountKeeper,
                BlockedAddresses(),
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                logger,
        )

        app.StakingKeeper = stakingkeeper.NewKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
                app.AccountKeeper,
                app.BankKeeper,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                authcodec.NewBech32Codec(Bech32PrefixValAddr),
                authcodec.NewBech32Codec(Bech32PrefixConsAddr),
        )

        app.MintKeeper = mintkeeper.NewKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[minttypes.StoreKey]),
                app.StakingKeeper,
                app.AccountKeeper,
                app.BankKeeper,
                authtypes.FeeCollectorName,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.DistrKeeper = distrkeeper.NewKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
                app.AccountKeeper,
                app.BankKeeper,
                app.StakingKeeper,
                authtypes.FeeCollectorName,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.SlashingKeeper = slashingkeeper.NewKeeper(
                appCodec,
                legacyAmino,
                runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
                app.StakingKeeper,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.StakingKeeper.SetHooks(
                stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
        )

        app.CrisisKeeper = crisiskeeper.NewKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
                cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
                app.BankKeeper,
                authtypes.FeeCollectorName,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.AccountKeeper.AddressCodec(),
        )

        app.UpgradeKeeper = upgradekeeper.NewKeeper(
                make(map[int64]bool),
                runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
                appCodec,
                DefaultNodeHome,
                app.BaseApp,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.IBCKeeper = ibckeeper.NewKeeper(
                appCodec,
                keys[ibcexported.StoreKey],
                app.GetSubspace(ibcexported.ModuleName),
                app.StakingKeeper,
                app.UpgradeKeeper,
                app.ScopedIBCKeeper,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.TransferKeeper = transferkeeper.NewKeeper(
                appCodec,
                keys[transfertypes.StoreKey],
                app.GetSubspace(transfertypes.ModuleName),
                app.IBCKeeper.ChannelKeeper,
                app.IBCKeeper.ChannelKeeper,
                app.IBCKeeper.PortKeeper,
                app.AccountKeeper,
                app.BankKeeper,
                app.ScopedTransferKeeper,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)
        ibcRouter := porttypes.NewRouter()
        ibcRouter.AddRoute(transfertypes.ModuleName, transferIBCModule)
        app.IBCKeeper.SetRouter(ibcRouter)

        app.IdentityKeeper = identitykeeper.NewKeeper(
                appCodec,
                keys[identitytypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.TrustKeeper = trustkeeper.NewKeeper(
                appCodec,
                keys[trusttypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.ComplianceKeeper = compliancekeeper.NewKeeper(
                appCodec,
                keys[compliancetypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.IdentityKeeper,
        )

        app.WalletKeeper = walletkeeper.NewKeeper(
                appCodec,
                keys[wallettypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.IdentityKeeper,
                app.ComplianceKeeper,
                app.BankKeeper,
        )

        app.ObligationKeeper = obligationkeeper.NewKeeper(
                appCodec,
                keys[obligationtypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.IdentityKeeper,
                app.WalletKeeper,
                app.ComplianceKeeper,
                app.TrustKeeper,
        )

        app.SettlementKeeper = settlementkeeper.NewKeeper(
                appCodec,
                keys[settlementtypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.IdentityKeeper,
                app.WalletKeeper,
                app.ObligationKeeper,
                app.ComplianceKeeper,
                app.TrustKeeper,
                app.BankKeeper,
        )

        app.DisputeKeeper = disputekeeper.NewKeeper(
                appCodec,
                keys[disputetypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.IdentityKeeper,
                app.ObligationKeeper,
                app.SettlementKeeper,
                app.TrustKeeper,
        )

        app.RegulatorKeeper = regulatorkeeper.NewKeeper(
                appCodec,
                keys[regulatortypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.IdentityKeeper,
                app.WalletKeeper,
        )

        app.ConstitutionKeeper = constitutionkeeper.NewKeeper(
                appCodec,
                keys[constitutiontypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.TscKeeper = tsckeeper.NewKeeper(
                appCodec,
                keys[tsctypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
                app.IdentityKeeper,
                app.WalletKeeper,
                app.ComplianceKeeper,
                app.BankKeeper,
        )

        app.AaKeeper = aakeeper.NewKeeper(
                appCodec,
                keys[aatypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.AuditKeeper = auditkeeper.NewKeeper(
                appCodec,
                keys[audittypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.IntentKeeper = intentkeeper.NewKeeper(
                appCodec,
                keys[intenttypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.PolicyKeeper = policykeeper.NewKeeper(
                appCodec,
                keys[policytypes.StoreKey],
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.ValsetKeeper = valsetkeeper.NewKeeper(
                keys[valsettypes.StoreKey],
                &app.AuditKeeper,
        )

        govConfig := govtypes.DefaultConfig()
        app.GovKeeper = govkeeper.NewKeeper(
                appCodec,
                runtime.NewKVStoreService(keys[govtypes.StoreKey]),
                app.AccountKeeper,
                app.BankKeeper,
                app.StakingKeeper,
                app.DistrKeeper,
                app.MsgServiceRouter(),
                govConfig,
                authtypes.NewModuleAddress(govtypes.ModuleName).String(),
        )

        app.GovKeeper.SetLegacyRouter(govv1beta1.NewRouter().
                AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
                AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)),
        )

        app.setAnteHandler()

        app.ModuleManager = module.NewManager(
                genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app, txConfig),
                auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
                vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
                bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
                capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
                crisis.NewAppModule(app.CrisisKeeper, false, app.GetSubspace(crisistypes.ModuleName)),
                gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
                mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
                slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
                distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
                staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
                upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
                params.NewAppModule(app.ParamsKeeper),
                consensus.NewAppModule(appCodec, app.ConsensusKeeper),
                ibc.NewAppModule(app.IBCKeeper),
                transfer.NewAppModule(app.TransferKeeper),
                tenidentity.NewAppModule(app.IdentityKeeper),
                tenwallet.NewAppModule(app.WalletKeeper),
                tentrust.NewAppModule(app.TrustKeeper),
                tencompliance.NewAppModule(app.ComplianceKeeper),
                tenobligation.NewAppModule(app.ObligationKeeper),
                tensettlement.NewAppModule(app.SettlementKeeper),
                tendispute.NewAppModule(app.DisputeKeeper),
                tenregulator.NewAppModule(app.RegulatorKeeper),
                tenconstitution.NewAppModule(app.ConstitutionKeeper),
                tentsc.NewAppModule(app.TscKeeper),
                tenaa.NewAppModule(app.AaKeeper),
                tenaudit.NewAppModule(app.AuditKeeper),
                tenintent.NewAppModule(app.IntentKeeper),
                tenpolicy.NewAppModule(app.PolicyKeeper),
                tenvalset.NewAppModule(app.ValsetKeeper),
        )

        app.ModuleManager.SetOrderBeginBlockers(
                upgradetypes.ModuleName,
                capabilitytypes.ModuleName,
                ibcexported.ModuleName,
                transfertypes.ModuleName,
                minttypes.ModuleName,
                distrtypes.ModuleName,
                slashingtypes.ModuleName,
                stakingtypes.ModuleName,
                genutiltypes.ModuleName,
                authtypes.ModuleName,
                banktypes.ModuleName,
                govtypes.ModuleName,
                crisistypes.ModuleName,
                paramstypes.ModuleName,
                consensustypes.ModuleName,
                vestingtypes.ModuleName,
                identitytypes.ModuleName,
                wallettypes.ModuleName,
                trusttypes.ModuleName,
                compliancetypes.ModuleName,
                obligationtypes.ModuleName,
                settlementtypes.ModuleName,
                disputetypes.ModuleName,
                regulatortypes.ModuleName,
                constitutiontypes.ModuleName,
                tsctypes.ModuleName,
                aatypes.ModuleName,
                audittypes.ModuleName,
                intenttypes.ModuleName,
                policytypes.ModuleName,
                valsettypes.ModuleName,
        )

        app.ModuleManager.SetOrderEndBlockers(
                crisistypes.ModuleName,
                govtypes.ModuleName,
                stakingtypes.ModuleName,
                ibcexported.ModuleName,
                transfertypes.ModuleName,
                capabilitytypes.ModuleName,
                authtypes.ModuleName,
                banktypes.ModuleName,
                distrtypes.ModuleName,
                slashingtypes.ModuleName,
                minttypes.ModuleName,
                genutiltypes.ModuleName,
                paramstypes.ModuleName,
                upgradetypes.ModuleName,
                consensustypes.ModuleName,
                vestingtypes.ModuleName,
                identitytypes.ModuleName,
                wallettypes.ModuleName,
                trusttypes.ModuleName,
                compliancetypes.ModuleName,
                obligationtypes.ModuleName,
                settlementtypes.ModuleName,
                disputetypes.ModuleName,
                regulatortypes.ModuleName,
                constitutiontypes.ModuleName,
                tsctypes.ModuleName,
                aatypes.ModuleName,
                audittypes.ModuleName,
                intenttypes.ModuleName,
                policytypes.ModuleName,
                valsettypes.ModuleName,
        )

        app.ModuleManager.SetOrderInitGenesis(
                capabilitytypes.ModuleName,
                ibcexported.ModuleName,
                transfertypes.ModuleName,
                authtypes.ModuleName,
                banktypes.ModuleName,
                distrtypes.ModuleName,
                stakingtypes.ModuleName,
                slashingtypes.ModuleName,
                govtypes.ModuleName,
                minttypes.ModuleName,
                crisistypes.ModuleName,
                genutiltypes.ModuleName,
                paramstypes.ModuleName,
                upgradetypes.ModuleName,
                consensustypes.ModuleName,
                vestingtypes.ModuleName,
                identitytypes.ModuleName,
                wallettypes.ModuleName,
                trusttypes.ModuleName,
                compliancetypes.ModuleName,
                obligationtypes.ModuleName,
                settlementtypes.ModuleName,
                disputetypes.ModuleName,
                regulatortypes.ModuleName,
                constitutiontypes.ModuleName,
                tsctypes.ModuleName,
                aatypes.ModuleName,
                audittypes.ModuleName,
                intenttypes.ModuleName,
                policytypes.ModuleName,
                valsettypes.ModuleName,
        )

        app.ModuleManager.RegisterInvariants(app.CrisisKeeper)
        app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
        app.ModuleManager.RegisterServices(app.configurator)

        autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))

        reflectionSvc, err := runtimeservices.NewReflectionService()
        if err != nil {
                panic(err)
        }
        reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

        app.SetInitChainer(app.InitChainer)
        app.SetBeginBlocker(app.BeginBlocker)
        app.SetEndBlocker(app.EndBlocker)

        app.MountKVStores(keys)
        app.MountTransientStores(tkeys)
        app.MountMemoryStores(memKeys)

        if loadLatest {
                if err := app.LoadLatestVersion(); err != nil {
                        tmos.Exit(err.Error())
                }
        }

        return app
}

func (app *TenitesApp) setAnteHandler() {
        anteHandler, err := NewAnteHandler(
                HandlerOptions{
                        AccountKeeper:    app.AccountKeeper,
                        BankKeeper:       app.BankKeeper,
                        SignModeHandler:  app.txConfig.SignModeHandler(),
                        SigGasConsumer:   nil,
                        ComplianceKeeper: app.ComplianceKeeper,
                        IdentityKeeper:   app.IdentityKeeper,
                },
        )
        if err != nil {
                panic(err)
        }

        app.SetAnteHandler(anteHandler)
}

func (app *TenitesApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
        var genesisState map[string]json.RawMessage
        if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
                panic(err)
        }
        return app.ModuleManager.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *TenitesApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
        return app.ModuleManager.BeginBlock(ctx)
}

func (app *TenitesApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
        return app.ModuleManager.EndBlock(ctx)
}

func (app *TenitesApp) Name() string { return app.BaseApp.Name() }

func (app *TenitesApp) LegacyAmino() *codec.LegacyAmino {
        return app.legacyAmino
}

func (app *TenitesApp) AppCodec() codec.Codec {
        return app.appCodec
}

func (app *TenitesApp) InterfaceRegistry() types.InterfaceRegistry {
        return app.interfaceRegistry
}

func (app *TenitesApp) TxConfig() client.TxConfig {
        return app.txConfig
}

func (app *TenitesApp) GetSubspace(moduleName string) paramstypes.Subspace {
        subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
        return subspace
}

func (app *TenitesApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
        clientCtx := apiSvr.ClientCtx
        nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
        ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}

func (app *TenitesApp) RegisterTxService(clientCtx client.Context) {
        authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *TenitesApp) RegisterTendermintService(clientCtx client.Context) {
        cmtservice.RegisterTendermintService(
                clientCtx,
                app.BaseApp.GRPCQueryRouter(),
                app.interfaceRegistry,
                app.Query,
        )
}

func (app *TenitesApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
        nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

func (app *TenitesApp) LoadHeight(height int64) error {
        return app.LoadVersion(height)
}

func (app *TenitesApp) ExportAppStateAndValidators(
        forZeroHeight bool,
        jailAllowedAddrs []string,
        modulesToExport []string,
) (servertypes.ExportedApp, error) {
        ctx := app.NewContext(true)

        height := app.LastBlockHeight()
        if forZeroHeight {
                height = 0
                app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
        }

        genState, err := app.ModuleManager.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
        if err != nil {
                return servertypes.ExportedApp{}, err
        }

        appState, err := json.MarshalIndent(genState, "", "  ")
        if err != nil {
                return servertypes.ExportedApp{}, err
        }

        validators, err := staking.WriteValidators(ctx, app.StakingKeeper)
        if err != nil {
                return servertypes.ExportedApp{}, err
        }

        return servertypes.ExportedApp{
                AppState:        appState,
                Validators:      validators,
                Height:          height,
                ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
        }, nil
}

func (app *TenitesApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) {
}

func BlockedAddresses() map[string]bool {
        modAccAddrs := make(map[string]bool)
        for acc := range maccPerms {
                modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
        }
        return modAccAddrs
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
        paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

        paramsKeeper.Subspace(authtypes.ModuleName)
        paramsKeeper.Subspace(banktypes.ModuleName)
        paramsKeeper.Subspace(stakingtypes.ModuleName)
        paramsKeeper.Subspace(minttypes.ModuleName)
        paramsKeeper.Subspace(distrtypes.ModuleName)
        paramsKeeper.Subspace(slashingtypes.ModuleName)
        paramsKeeper.Subspace(govtypes.ModuleName)
        paramsKeeper.Subspace(crisistypes.ModuleName)
        paramsKeeper.Subspace(ibcexported.ModuleName)
        paramsKeeper.Subspace(transfertypes.ModuleName).WithKeyTable(transfertypes.ParamKeyTable())

        return paramsKeeper
}
