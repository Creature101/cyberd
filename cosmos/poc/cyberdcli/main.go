package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	cyberdcmd "github.com/cybercongress/cyberd/cosmos/poc/commands"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

// cyberdcli is the entry point for this binary
var (
	cyberdcli = &cobra.Command{
		Use:   "cyberdcli",
		Short: "Cyberd node client",
	}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.

	rpc.AddCommands(cyberdcli) // Node management commands
	cyberdcli.AddCommand(client.LineBreak)
	tx.AddCommands(cyberdcli, cdc) // Txs info commands
	cyberdcli.AddCommand(client.LineBreak)
	cyberdcli.AddCommand(rpc.BlockCommand()) // Block info command
	cyberdcli.AddCommand(client.LineBreak)

	cyberdcli.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, app.GetAccountDecoder(cdc)),
			cyberdcmd.GetLinksCmd("link", cdc),
		)...)

	cyberdcli.AddCommand(
		client.PostCommands(
			cyberdcmd.LinkTxCmd(cdc),
			bankcmd.SendTxCmd(cdc),
		)...)

	cyberdcli.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc), // Commands to start local rpc proxy to node
		keys.Commands(), // Commands to generate and handle keys
		client.LineBreak,
		version.VersionCmd,
	)

	executor := cli.PrepareMainCmd(cyberdcli, "CBD", os.ExpandEnv("$HOME/.cyberdcli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}