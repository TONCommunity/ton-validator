package main

import (
	"flag"
	"os"

	"github.com/peterbourgon/ff"
)

var (
	fiftBin                 string
	fiftPath                string
	liteClient              string
	liteclientConfig        string
	dbFile                  string
	tonlibConfig            string
	validatorConsole        string
	validatorWalletAddr     string
	maxFactor               string
	validatorHost           string
	stakeAmount             int
	walletFif               string
	recoverFif              string
	validatorElectReqFif    string
	validatorElectSignedFif string
	validatorClientCert     string
	validatorServerPub      string
	validatorWalletFile     string
	verbose                 bool
	verboseTonlib           int
)

// GetConfig Gets the conf in the config file
func GetConfig() error {
	fs := flag.NewFlagSet("ton-validator", flag.ExitOnError)
	fs.StringVar(&fiftBin, "fift-bin", "fift", "path to fift binary")
	fs.StringVar(&fiftPath, "fift-path", "crypto/fift/lib/", "path to fift lib")
	fs.StringVar(&liteClient, "lite-client", "lite-client", "path to lite-client binary")
	fs.StringVar(&liteclientConfig, "lite-client-config", "ton-lite-client-test1.config.json", "path to lite-client config")
	fs.StringVar(&tonlibConfig, "tonlib-config", "tonlib.config.json", "tonlib config")
	fs.StringVar(&dbFile, "db-file", "./ton.db", "path to db file")
	fs.StringVar(&validatorConsole, "validator-console", "validator-engine-console", "path to validator-engine-console binary")
	fs.StringVar(&maxFactor, "max-factor", "2.7", "max factor")
	fs.IntVar(&stakeAmount, "stake-amount", 20000, "stake amount")
	fs.StringVar(&walletFif, "wallet-fif", "crypto/smartcont/wallet.fif", "path to wallet.fif file")
	fs.StringVar(&recoverFif, "recover-fif", "crypto/smartcont/recover-stake.fif", "path to recover-stake.fif file")
	fs.StringVar(&validatorElectReqFif, "validator-elect-req-fif", "crypto/smartcont/validator-elect-req.fif", "path to validator-elect-req.fif file")
	fs.StringVar(&validatorElectSignedFif, "validator-elect-signed-fif", "crypto/smartcont/validator-elect-signed.fif", "path to validator-elect-signed.fif file")
	fs.BoolVar(&verbose, "verbose", false, "tool verbosity")
	fs.IntVar(&verboseTonlib, "verbose-tonlib", 0, "tonlib versbosity")
	_ = fs.String("config", "", "config file (optional)")

	err := ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
		ff.WithEnvVarPrefix("TON"),
	)
	return err
}
