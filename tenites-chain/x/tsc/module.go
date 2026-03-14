package tsc

  import (
        "encoding/json"

        abci "github.com/cometbft/cometbft/abci/types"
        "github.com/cosmos/cosmos-sdk/client"
        "github.com/cosmos/cosmos-sdk/codec"
        codectypes "github.com/cosmos/cosmos-sdk/codec/types"
        sdk "github.com/cosmos/cosmos-sdk/types"
        "github.com/cosmos/cosmos-sdk/types/module"
        "github.com/grpc-ecosystem/grpc-gateway/runtime"
        "github.com/spf13/cobra"

        "github.com/tenites/tenites-chain/x/tsc/keeper"
        "github.com/tenites/tenites-chain/x/tsc/types"
        "github.com/tenites/tenites-chain/x/tsc/client/cli"
)

  var (
        _ module.AppModule      = AppModule{}
        _ module.AppModuleBasic = AppModuleBasic{}
  )

  type AppModuleBasic struct{}

  func (AppModuleBasic) Name() string { return types.ModuleName }

  func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

  func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
        types.RegisterInterfaces(registry)
  }

  func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
        return json.RawMessage(`{}`)
  }

  func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
        return nil
  }

  func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

  func (AppModuleBasic) GetTxCmd() *cobra.Command    { return cli.GetTxCmd() }
  func (AppModuleBasic) GetQueryCmd() *cobra.Command { return cli.GetQueryCmd() }

  type AppModule struct {
        AppModuleBasic
        keeper keeper.Keeper
  }

  func NewAppModule(k keeper.Keeper) AppModule {
        return AppModule{keeper: k}
  }

  func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
        ir.RegisterRoute(types.ModuleName, "supply-invariants", keeper.SupplyInvariant(am.keeper))
  }

  func (am AppModule) RegisterServices(cfg module.Configurator) {
        types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
        types.RegisterQueryServer(cfg.QueryServer(), keeper.NewGRPCQueryServer(am.keeper))
  }

  func (am AppModule) ConsensusVersion() uint64 { return 1 }

  func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
        return []abci.ValidatorUpdate{}
  }

  func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
        return json.RawMessage(`{}`)
  }
  
func (AppModule) IsAppModule() {}
func (AppModule) IsOnePerModuleType() {}
