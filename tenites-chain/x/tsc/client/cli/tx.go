package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/tenites/tenites-chain/x/tsc/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Transaction commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdMintTSC(),
		GetCmdBurnTSC(),
		GetCmdTransferTSC(),
		GetCmdAddMintAuthority(),
		GetCmdRemoveMintAuthority(),
		GetCmdAddBurnAuthority(),
		GetCmdRemoveBurnAuthority(),
	)

	return cmd
}

func GetCmdMintTSC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [authority-id] [wallet-id] [amount] [purpose] [jurisdiction]",
		Short: "Mint TSC to a wallet using mint authority",
		Long: `Mint TSC tokens to a specified wallet.

The caller must have a valid mint authority that permits the specified purpose and jurisdiction.

Example:
  $ teniteschaind tx tsc mint MINT0000000001 WALLET001 1000000000 issuance NG --from mykey`,
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgMint{
				Signer:       clientCtx.GetFromAddress().String(),
				AuthorityId:  args[0],
				WalletId:     args[1],
				Amount:       args[2],
				Purpose:      args[3],
				Jurisdiction: args[4],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdBurnTSC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [authority-id] [wallet-id] [amount] [reason] [jurisdiction]",
		Short: "Burn TSC from a wallet using burn authority",
		Long: `Burn TSC tokens from a specified wallet.

The caller must have a valid burn authority that permits the specified reason and jurisdiction.

Example:
  $ teniteschaind tx tsc burn BURN0000000001 WALLET001 1000000000 redemption NG --from mykey`,
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgBurn{
				Signer:       clientCtx.GetFromAddress().String(),
				AuthorityId:  args[0],
				WalletId:     args[1],
				Amount:       args[2],
				Reason:       args[3],
				Jurisdiction: args[4],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdTransferTSC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [from-wallet-id] [to-wallet-id] [amount] [purpose]",
		Short: "Transfer TSC between wallets",
		Long: `Transfer TSC tokens between wallets with compliance checks.

For amounts above the travel rule threshold, both parties must meet KYC requirements.

Example:
  $ teniteschaind tx tsc transfer WALLET001 WALLET002 1000000000 settlement --from mykey`,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgTransfer{
				Signer:       clientCtx.GetFromAddress().String(),
				FromWalletId: args[0],
				ToWalletId:   args[1],
				Amount:       args[2],
				Purpose:      args[3],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdAddMintAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-mint-authority [tenites-id] [name] [daily-limit] [single-tx-limit]",
		Short: "Add a new mint authority (governance only)",
		Long: `Add a new TSC mint authority.

This command can only be executed by governance. The authority grants the ability to mint TSC.

Example:
  $ teniteschaind tx tsc add-mint-authority TEN001 "Treasury Minter" 1000000000000 100000000000 --from governance`,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgAddMintAuthority{
				Signer:          clientCtx.GetFromAddress().String(),
				TenitesId:       args[0],
				Name:            args[1],
				DailyLimit:      args[2],
				SingleTxLimit:   args[3],
				AllowedPurposes: []string{"issuance", "liquidity", "settlement", "reserve_backup"},
				Jurisdictions:   []string{"*"},
				ExpiresAt:       0,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRemoveMintAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-mint-authority [authority-id]",
		Short: "Revoke a mint authority (governance only)",
		Long: `Revoke an existing TSC mint authority.

This command can only be executed by governance. The authority will be marked as revoked.

Example:
  $ teniteschaind tx tsc remove-mint-authority MINT0000000001 --from governance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgRemoveMintAuthority{
				Signer:      clientCtx.GetFromAddress().String(),
				AuthorityId: args[0],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdAddBurnAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-burn-authority [tenites-id] [name] [daily-limit] [single-tx-limit]",
		Short: "Add a new burn authority (governance only)",
		Long: `Add a new TSC burn authority.

This command can only be executed by governance. The authority grants the ability to burn TSC.

Example:
  $ teniteschaind tx tsc add-burn-authority TEN001 "Treasury Burner" 1000000000000 100000000000 --from governance`,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgAddBurnAuthority{
				Signer:         clientCtx.GetFromAddress().String(),
				TenitesId:      args[0],
				Name:           args[1],
				DailyLimit:     args[2],
				SingleTxLimit:  args[3],
				AllowedReasons: []string{"redemption", "recall", "compliance", "expiry", "rebalancing"},
				Jurisdictions:  []string{"*"},
				ExpiresAt:      0,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRemoveBurnAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-burn-authority [authority-id]",
		Short: "Revoke a burn authority (governance only)",
		Long: `Revoke an existing TSC burn authority.

This command can only be executed by governance. The authority will be marked as revoked.

Example:
  $ teniteschaind tx tsc remove-burn-authority BURN0000000001 --from governance`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgRemoveBurnAuthority{
				Signer:      clientCtx.GetFromAddress().String(),
				AuthorityId: args[0],
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
