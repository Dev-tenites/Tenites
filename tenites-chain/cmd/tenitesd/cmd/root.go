package cmd

import (
        "fmt"
        "io"
        "os"
        "time"

        "cosmossdk.io/log"
        dbm "github.com/cosmos/cosmos-db"
        cmtcfg "github.com/cometbft/cometbft/config"
        "github.com/cosmos/cosmos-sdk/client"
        "github.com/cosmos/cosmos-sdk/client/config"
        "github.com/cosmos/cosmos-sdk/client/debug"
        "github.com/cosmos/cosmos-sdk/client/flags"
        "github.com/cosmos/cosmos-sdk/client/keys"
        "github.com/cosmos/cosmos-sdk/client/pruning"
        "github.com/cosmos/cosmos-sdk/client/rpc"
        "github.com/cosmos/cosmos-sdk/client/snapshot"
        "github.com/cosmos/cosmos-sdk/server"
        serverconfig "github.com/cosmos/cosmos-sdk/server/config"
        servertypes "github.com/cosmos/cosmos-sdk/server/types"
        sdk "github.com/cosmos/cosmos-sdk/types"
        "github.com/cosmos/cosmos-sdk/types/tx/signing"
        authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
        "github.com/cosmos/cosmos-sdk/x/auth/tx"
        txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
        authcodec "github.com/cosmos/cosmos-sdk/codec/address"
        bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
        tsccli "github.com/tenites/tenites-chain/x/tsc/client/cli"
        ibccli "github.com/cosmos/ibc-go/v8/modules/core/client/cli"
        transfercli "github.com/cosmos/ibc-go/v8/modules/apps/transfer/client/cli"
        "github.com/cosmos/cosmos-sdk/x/crisis"
        genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
        "github.com/spf13/cobra"
        "github.com/spf13/viper"

        "github.com/tenites/tenites-chain/app"
)

func NewRootCmd() *cobra.Command {
        cfg := sdk.GetConfig()
        cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
        cfg.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
        cfg.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
        cfg.Seal()

        encodingConfig := app.EncodingConfig()

        initClientCtx := client.Context{}.
                WithCodec(encodingConfig.Codec).
                WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
                WithTxConfig(encodingConfig.TxConfig).
                WithLegacyAmino(encodingConfig.Amino).
                WithInput(os.Stdin).
                WithAccountRetriever(encodingConfig.AccountRetriever).
                WithHomeDir(app.DefaultNodeHome).
                WithViper("")

        rootCmd := &cobra.Command{
                Use:   "tenitesd",
                Short: "Tenites Chain - Production-grade cross-border payment settlement blockchain",
                Long: `Tenites Chain is a Cosmos SDK sovereign Layer 1 blockchain designed for 
cross-border payment settlement in African corridors. It features jurisdiction-aware 
compliance, trust scoring, regulatory oversight, and atomic settlement.`,
                PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
                        initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
                        if err != nil {
                                return err
                        }

                        initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
                        if err != nil {
                                return err
                        }

                        if !initClientCtx.Offline {
                                txConfigOpts := tx.ConfigOptions{
                                        EnabledSignModes:           append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL),
                                        TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
                                }
                                txConfig, err := tx.NewTxConfigWithOptions(
                                        initClientCtx.Codec,
                                        txConfigOpts,
                                )
                                if err != nil {
                                        return err
                                }

                                initClientCtx = initClientCtx.WithTxConfig(txConfig)
                        }

                        if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
                                return err
                        }

                        customAppTemplate, customAppConfig := initAppConfig()
                        customCMTConfig := initCometBFTConfig()

                        return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
                },
        }

        initRootCmd(rootCmd, encodingConfig)

        return rootCmd
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig app.EncodingConfigType) {
        rootCmd.AddCommand(
                genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
                debug.Cmd(),
                pruning.Cmd(newApp, app.DefaultNodeHome),
                snapshot.Cmd(newApp),
        )

        server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

        rootCmd.AddCommand(
                rpc.WaitTxCmd(),
                genutilcli.GenesisCoreCommand(encodingConfig.TxConfig, app.ModuleBasics, app.DefaultNodeHome),
                genutilcli.Commands(encodingConfig.TxConfig, app.ModuleBasics, app.DefaultNodeHome),
                queryCommand(),
                txCommand(),
                keys.Commands(),
        )
}

func addModuleInitFlags(startCmd *cobra.Command) {
        crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
        cmd := &cobra.Command{
                Use:                        "query",
                Aliases:                    []string{"q"},
                Short:                      "Querying subcommands",
                DisableFlagParsing:         false,
                SuggestionsMinimumDistance: 2,
                RunE:                       client.ValidateCmd,
        }

        cmd.AddCommand(
                rpc.QueryEventForTxCmd(),
                server.QueryBlockCmd(),
                authcmd.QueryTxsByEventsCmd(),
                server.QueryBlocksCmd(),
                authcmd.QueryTxCmd(),
                server.QueryBlockResultsCmd(),
                tsccli.GetQueryCmd(),
                ibccli.GetQueryCmd(),
                transfercli.GetQueryCmd(),
        )

        return cmd
}

func txCommand() *cobra.Command {
        cmd := &cobra.Command{
                Use:                        "tx",
                Short:                      "Transactions subcommands",
                DisableFlagParsing:         false,
                SuggestionsMinimumDistance: 2,
                RunE:                       client.ValidateCmd,
        }

        bankCmd := &cobra.Command{
                Use:   "bank",
                Short: "Bank transaction subcommands",
                RunE:  client.ValidateCmd,
        }
        bankCmd.AddCommand(bankcli.NewSendTxCmd(authcodec.NewBech32Codec(app.Bech32PrefixAccAddr)))

        cmd.AddCommand(
                authcmd.GetSignCommand(),
                authcmd.GetSignBatchCommand(),
                authcmd.GetMultiSignCommand(),
                authcmd.GetMultiSignBatchCmd(),
                authcmd.GetValidateSignaturesCommand(),
                authcmd.GetBroadcastCommand(),
                authcmd.GetEncodeCommand(),
                authcmd.GetDecodeCommand(),
                authcmd.GetSimulateCmd(),
                bankCmd,
                tsccli.GetTxCmd(),
                ibccli.GetTxCmd(),
                transfercli.NewTransferTxCmd(),
        )

        return cmd
}

func newApp(
        logger log.Logger,
        db dbm.DB,
        traceStore io.Writer,
        appOpts servertypes.AppOptions,
) servertypes.Application {
        baseappOptions := server.DefaultBaseappOptions(appOpts)

        return app.NewTenitesApp(
                logger,
                db,
                traceStore,
                true,
                appOpts,
                baseappOptions...,
        )
}

func appExport(
        logger log.Logger,
        db dbm.DB,
        traceStore io.Writer,
        height int64,
        forZeroHeight bool,
        jailAllowedAddrs []string,
        appOpts servertypes.AppOptions,
        modulesToExport []string,
) (servertypes.ExportedApp, error) {
        var tenitesApp *app.TenitesApp

        homePath, ok := appOpts.Get(flags.FlagHome).(string)
        if !ok || homePath == "" {
                return servertypes.ExportedApp{}, fmt.Errorf("application home is not set")
        }

        viperAppOpts, ok := appOpts.(*viper.Viper)
        if !ok {
                return servertypes.ExportedApp{}, fmt.Errorf("appOpts is not viper.Viper")
        }

        viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
        appOpts = viperAppOpts

        tenitesApp = app.NewTenitesApp(
                logger,
                db,
                traceStore,
                height == -1,
                appOpts,
        )

        if height != -1 {
                if err := tenitesApp.LoadHeight(height); err != nil {
                        return servertypes.ExportedApp{}, err
                }
        }

        return tenitesApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

func initAppConfig() (string, interface{}) {
        type CustomAppConfig struct {
                serverconfig.Config `mapstructure:",squash"`
        }

        srvCfg := serverconfig.DefaultConfig()
        srvCfg.StateSync.SnapshotInterval = 1000
        srvCfg.StateSync.SnapshotKeepRecent = 2
        srvCfg.MinGasPrices = "0.0001utsc"
        srvCfg.API.Enable = true
        srvCfg.API.Swagger = true
        srvCfg.API.EnableUnsafeCORS = false
        srvCfg.GRPC.Enable = true
        srvCfg.GRPCWeb.Enable = true

        customAppConfig := CustomAppConfig{
                Config: *srvCfg,
        }

        customAppTemplate := serverconfig.DefaultConfigTemplate

        return customAppTemplate, customAppConfig
}

func initCometBFTConfig() *cmtcfg.Config {
        cfg := cmtcfg.DefaultConfig()
        cfg.Consensus.TimeoutPropose = time.Second * 3
        cfg.Consensus.TimeoutProposeDelta = time.Millisecond * 500
        cfg.Consensus.TimeoutPrevote = time.Second * 1
        cfg.Consensus.TimeoutPrevoteDelta = time.Millisecond * 500
        cfg.Consensus.TimeoutPrecommit = time.Second * 1
        cfg.Consensus.TimeoutPrecommitDelta = time.Millisecond * 500
        cfg.Consensus.TimeoutCommit = time.Second * 5
        cfg.P2P.MaxNumInboundPeers = 100
        cfg.P2P.MaxNumOutboundPeers = 40
        cfg.Mempool.Size = 10000
        cfg.Mempool.MaxTxsBytes = 1073741824
        return cfg
}
