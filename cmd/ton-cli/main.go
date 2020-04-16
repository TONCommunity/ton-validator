package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/mercuryoio/ton-validator/database"
	"github.com/peterbourgon/ff/ffcli"
)

func main() {
	var (
		rootFlagSet   = flag.NewFlagSet("ton-cli", flag.ExitOnError)
		nodeFlagSet   = flag.NewFlagSet("ton-cli node", flag.ExitOnError)
		nodeEnabled   = nodeFlagSet.Int("enabled", 2, "\t\"Filter nodes: 0 - disabled, 1 - enabled, 2 - all\"")
		walletID      = nodeFlagSet.Int("wallet", 1, "\t\"Filter by wallet ID\"")
		walletFlagSet = flag.NewFlagSet("ton-cli wallet", flag.ExitOnError)
		walletEnabled = walletFlagSet.Int("enabled", 2, "\t\"Filter wallets: 0 - disabled, 1 - enabled, 2 - all\"")
		stakeFlagSet  = flag.NewFlagSet("ton-cli stake", flag.ExitOnError)
	)

	s, err := database.NewClient("./ton.db")
	if err != nil {
		fmt.Println("Failed to connect to db:", err)
		os.Exit(1)
	}

	addWallet := &ffcli.Command{
		Name:       "add",
		ShortUsage: "add <wallet_address> <wallet_file_path>",
		ShortHelp:  "Add wallet.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 2 {
				return fmt.Errorf("Add wallet requires exactly 2 arguments, but you provided %d", n)
			}
			walletID, err := s.AddWallet(args[1], args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Added wallet: %s %s with ID: %d\n", args[0], args[1], walletID)
			return nil
		},
	}

	delWallet := &ffcli.Command{
		Name:       "del",
		ShortUsage: "del [<id> ...]",
		ShortHelp:  "Delete wallet by ID.",
		Exec: func(_ context.Context, args []string) error {
			for _, id := range args {
				i, err := strconv.Atoi(id)
				if err != nil {
					return fmt.Errorf("Failed to remove wallet with ID: %d : %v", i, err)
				}
				err = s.DelWallet(i)
				if err != nil {
					return err
				}
				fmt.Println("Removed wallet with id: ", i)
			}
			return nil
		},
	}

	listWallets := &ffcli.Command{
		Name:       "list",
		ShortUsage: "list",
		ShortHelp:  "List wallets",
		FlagSet:    walletFlagSet,
		Exec: func(_ context.Context, args []string) error {
			wallets, err := s.GetWallets(*walletEnabled)
			if err != nil {
				return err
			}
			for _, wallet := range wallets {
				fmt.Println("ID:", wallet.ID, "\tAddress:", wallet.Addr, "\tWallet File:", wallet.FilePath, "\tEnabled:", wallet.Enabled)
			}
			return nil
		},
	}

	wallet := &ffcli.Command{
		Name:        "wallet",
		ShortUsage:  "wallet [<arg> ...]",
		ShortHelp:   "Wallet management.",
		Subcommands: []*ffcli.Command{addWallet, listWallets, delWallet},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	addNode := &ffcli.Command{
		Name:       "add",
		ShortUsage: "add <node_host:port> <client_cert> <server.pub> <wallet_id>",
		ShortHelp:  "Add node.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 4 {
				return fmt.Errorf("Add node requires exactly 4 arguments, but you provided %d", n)
			}
			walletID, _ := strconv.Atoi(args[3])
			nodeID, err := s.AddNode(args[0], args[1], args[2], walletID)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Added node: %s %s %s to wallet %d with id: %d\n", args[0], args[1], args[2], walletID, nodeID)
			return nil
		},
	}

	listNodes := &ffcli.Command{
		Name:       "list",
		ShortUsage: "list",
		ShortHelp:  "List nodes.",
		FlagSet:    nodeFlagSet,
		Exec: func(_ context.Context, args []string) error {
			nodes, err := s.GetNodes(*walletID, *nodeEnabled)
			if err != nil {
				return err
			}
			for _, node := range nodes {
				fmt.Println("ID:", node.ID, "\tAddress:", node.HostPort, "\tClient cert:", node.ClientCert, "\tserver.pub:", node.ServerPub, "\tEnabled:", node.Enabled)
			}
			return nil
		},
	}

	delNode := &ffcli.Command{
		Name:       "del",
		ShortUsage: "del [<id> ...]",
		ShortHelp:  "Delete node from db by ID.",
		Exec: func(_ context.Context, args []string) error {
			for _, id := range args {
				i, err := strconv.Atoi(id)
				if err != nil {
					fmt.Println(err)
				}
				err = s.DelNode(i)
				if err != nil {
					return err
				}
				fmt.Println("Removed node with id: ", i)
			}
			return nil
		},
	}

	node := &ffcli.Command{
		Name:        "node",
		ShortUsage:  "node [<arg> ...]",
		ShortHelp:   "Node management",
		Subcommands: []*ffcli.Command{addNode, listNodes, delNode},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	listStakes := &ffcli.Command{
		Name:       "list",
		ShortUsage: "list",
		ShortHelp:  "List stakes.",
		FlagSet:    stakeFlagSet,
		Exec: func(_ context.Context, args []string) error {
			return nil
		},
	}

	stake := &ffcli.Command{
		Name:        "stake",
		ShortUsage:  "stake [<arg> ...]",
		ShortHelp:   "Participate in election with stake.",
		Subcommands: []*ffcli.Command{listStakes},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	activeElection := &ffcli.Command{
		Name:       "active",
		ShortUsage: "active",
		ShortHelp:  "Get active election id.",
		FlagSet:    stakeFlagSet,
		Exec: func(_ context.Context, args []string) error {
			fmt.Println("To Be Done")
			return nil
		},
	}

	listElections := &ffcli.Command{
		Name:       "list",
		ShortUsage: "list",
		ShortHelp:  "List elections.",
		FlagSet:    stakeFlagSet,
		Exec: func(_ context.Context, args []string) error {

			return nil
		},
	}

	election := &ffcli.Command{
		Name:        "election",
		ShortUsage:  "election [<arg> ...]",
		ShortHelp:   "Show elections information.",
		Subcommands: []*ffcli.Command{listElections, activeElection},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "ton-cli [flags] <subcommand>",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{wallet, node, stake, election},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
