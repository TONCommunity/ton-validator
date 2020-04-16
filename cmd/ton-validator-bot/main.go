package main

import (
	"fmt"
	"github.com/mercuryoio/ton-validator/database"
	"github.com/mercuryoio/ton-validator/utils"
	"github.com/mercuryoio/ton-validator/wrappers/fift"
	"github.com/mercuryoio/ton-validator/wrappers/liteclient"
	"github.com/mercuryoio/ton-validator/wrappers/validator"
	tonlib "github.com/mercuryoio/tonlib-go/v2"
	"log"
	"os"
	"time"
)

//GetTonlibClient Get tonlib client
func GetTonlibClient() *tonlib.Client {
	options, err := tonlib.ParseConfigFile(tonlibConfig)
	if err != nil {
		log.Fatal("failed parse config error. ", err)
	}

	req := tonlib.TonInitRequest{
		Type:    "init",
		Options: *options,
	}

	cln, err := tonlib.NewClient(&req, tonlib.Config{}, 60, verbose, int32(verboseTonlib))
	if err != nil {
		log.Fatalln("Init client error", err)
	}

	return cln
}

func main() {
	GetConfig()
	cln := GetTonlibClient()
	s, err := database.NewClient(dbFile)
	if err != nil {
		fmt.Println("Failed to connect to db:", err)
		os.Exit(1)
	}

	fiftConfig := fift.Config{
		FiftBin:                 &fiftBin,
		FiftPath:                &fiftPath,
		WalletFif:               &walletFif,
		RecoverFif:              &recoverFif,
		ValidatorElectReqFif:    &validatorElectReqFif,
		ValidatorElectSignedFif: &validatorElectSignedFif,
		ValidatorWalletFile:     &validatorWalletFile,
		Verbose:                 &verbose,
	}
	f := fift.NewClient(&fiftConfig)

	liteConfig := liteclient.Config{
		LiteClient:       &liteClient,
		LiteclientConfig: &liteclientConfig,
		Verbose:          &verbose,
	}
	lc := liteclient.NewClient(&liteConfig)

	validatorConfig := validator.Config{
		ValidatorConsole: &validatorConsole,
		Verbose:          &verbose,
	}
	vc := validator.NewClient(&validatorConfig)
	err = s.SyncWalletsBalance(cln)
	if err != nil {
		log.Println(err)
	}

	currentElectorAddress, err := lc.GetCurrentElectorAddress()
	if err != nil {
		log.Println("Current elector Address failed", err)
		os.Exit(1)
	}
	fmt.Println("Current elector Address:", currentElectorAddress)

	periods, _ := lc.GetElectionConfig()
	fmt.Println("Network configuration:")
	fmt.Printf("\tvalidators_elected_for: %d\telections_start_before: %d\telections_end_before: %d\tstake_held_for: %d\t\n", periods.ValidatorsElectedFor, periods.ElectionsStartBefore, periods.ElectionsEndBefore, periods.StakeHeldFor)

	stakeConfig, err := lc.GetStakeConfig()
	fmt.Println("Network stake config:")
	fmt.Printf("\tmin_stake: %s\tmax_stake: %s\tmin_total_stake: %s\tmax_stake_factor: %d (%d)\t\n", utils.FormatGrams(stakeConfig.MinStake), utils.FormatGrams(stakeConfig.MaxStake), utils.FormatGrams(stakeConfig.MinTotalStake), stakeConfig.MaxStakeFactor, (stakeConfig.MaxStakeFactor / 65536))

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err = cln.UpdateTonConnection()
		if err != nil {
			log.Fatalln("UpdateTonConnection", err)
		}

		activeElectionID, err := cln.GetActiveElectionID(currentElectorAddress)

		if err != nil {
			log.Fatalln("GetActiveElectionID failed", err)
		}

		if activeElectionID != 0 {

			log.Println("Active election ID:", activeElectionID)
			/* minEffStake, maxEffStake, err := cln.GetElectionStackes(currentElectorAddress, stakeConfig.minStake, stakeConfig.maxStake, stakeConfig.maxStakeFactor, 13)

			fmt.Println("minimal effective stake:", formatGrams(minEffStake), "maximal effective stake", formatGrams(maxEffStake)) */
			election, err := s.GetElection(activeElectionID)
			if err != nil {
				log.Println("Failed to get election from db:", err)
			}
			if election.ElectionID != 0 {
				log.Println("Election from db", election.ElectionID)
			} else {
				election = database.Election{
					ElectionID:      activeElectionID,
					StartAt:         activeElectionID - periods.ElectionsStartBefore,
					CloseAt:         activeElectionID - periods.ElectionsEndBefore,
					NextElectionsAt: activeElectionID + periods.ElectionsStartBefore,
				}
				_, err := s.AddElection(election)
				if err != nil {
					log.Fatalln("Failed to add election to db:", err)
				}
			}
		}

		wallets, err := s.GetWallets(1)
		if err != nil {
			log.Fatalln(err)
		}
		if len(wallets) == 0 {
			log.Fatalln("No wallets found")
		}
		for _, wallet := range wallets {
			err = cln.UpdateTonConnection()
			if err != nil {
				log.Fatalln("UpdateTonConnection", err)
			}
			AccountState, err := cln.GetAccountState(*tonlib.NewAccountAddress(wallet.Addr))
			if err != nil {
				log.Println("getAccountState failed", err)
				os.Exit(1)
			}

			log.Println("Wallet", wallet.Addr, "balance", utils.FormatGrams(wallet.Balance))

			if wallet.Balance < int64(AccountState.Balance) {
				log.Printf("Balance changed: %s (+%s)", utils.FormatGrams(int64(AccountState.Balance)), utils.FormatGrams(int64(AccountState.Balance)-wallet.Balance))
				s.UpdateWalletBalance(wallet.ID, int64(AccountState.Balance))
			} else if wallet.Balance > int64(AccountState.Balance) {
				log.Printf("Balance changed: %s (-%s)", utils.FormatGrams(int64(AccountState.Balance)), utils.FormatGrams(wallet.Balance-int64(AccountState.Balance)))
				s.UpdateWalletBalance(wallet.ID, int64(AccountState.Balance))
			}
			err = cln.UpdateTonConnection()
			if err != nil {
				log.Fatalln("UpdateTonConnection", err)
			}
			unpackedAddress, err := cln.UnpackAccountAddress(wallet.Addr)

			if err != nil {
				log.Println("UnpackAccountAddress failed", err)
				break
			}
			err = cln.UpdateTonConnection()
			if err != nil {
				log.Fatalln("UpdateTonConnection", err)
			}
			reward, _ := cln.CheckReward(utils.PubKeyToHex(unpackedAddress.Addr), currentElectorAddress)
			if reward != 0 {
				log.Printf("Sending request to recover %12s GRAMs\n", utils.FormatGrams(reward))
				err = cln.UpdateTonConnection()
				if err != nil {
					log.Fatalln("UpdateTonConnection", err)
				}
				recoverSeqno, err := cln.GetWalletSeqno(wallet.Addr)
				if err != nil {
					log.Fatalln("GetWalletSeqno failed:", err)
				}
				f.RecoverStake(cln, wallet.FilePath, currentElectorAddress, recoverSeqno)
			}

			if wallet.Balance < stakeConfig.MinStake {
				log.Println("Account balance is too low, can't stake, skipping")
				continue
			}

			if activeElectionID == 0 {
				continue
			}

			nodes, err := s.GetNodes(wallet.ID, 1)
			if err != nil {
				log.Fatalln(err)
			}
			if len(nodes) == 0 {
				log.Println("No nodes found for wallet", wallet.Addr)
				continue
			}
			for _, node := range nodes {
				log.Println(node.HostPort)
				if !vc.CheckNodeSync(node) {
					fmt.Println("Validator node is out of sync")
				}

				participates := s.GetParticipates(node.ID, activeElectionID)
				if len(participates) == 0 {
					log.Println("Not participating")
				} else {
					for i, p := range participates {
						log.Println("Participate #", i, p)
					}
					//continue
				}

				validatorKey, err := s.GetKey("key", node.ID, activeElectionID)
				if err != nil {
					log.Println("Failed to get key from", node.HostPort, err)
				}
				if validatorKey.Key == "" {
					validatorKey, err = vc.ValidatorCreateNewKey(node, activeElectionID)
					if err != nil {
						log.Fatalln("validatorCreateNewKey failed", err)
					}
					log.Println("Created new key:", validatorKey)
					keyID, err := s.AddKey(validatorKey)
					if err != nil {
						log.Println("failed to save key to db", keyID, err)
					}
					if vc.ValidatorAddPermKey(node, validatorKey.Key, activeElectionID) {
						log.Println("Added permKey", validatorKey.Key, activeElectionID)
					}

					if vc.ValidatorAddTempKey(node, validatorKey.Key, validatorKey.Key, activeElectionID+int64(periods.ValidatorsElectedFor)+10000) {
						log.Println("Added tempKey", validatorKey.Key, activeElectionID)
					}
				}

				pubKey, err := s.GetKey("pubkey", node.ID, activeElectionID)
				if err != nil {
					log.Println("Failed to get pubkey from", node.HostPort, err)
				}
				if pubKey.Key == "" {
					pubKey, err = vc.ValidatorGetPublicKey(node, validatorKey.Key, activeElectionID)
					if err != nil {
						log.Fatalln("validatorGetPublicKey failed", node.HostPort, err)
					}
					log.Println(pubKey)
					pubKeyID, err := s.AddKey(pubKey)
					if err != nil {
						log.Println("failed to save key to db", pubKeyID, pubKey.Key, err)
					}
				}
				err = cln.UpdateTonConnection()
				if err != nil {
					log.Fatalln("UpdateTonConnection", err)
				}
				amount, err := cln.CheckParticipatesIn(utils.PubKeyToHex(pubKey.Key), currentElectorAddress)
				if err != nil {
					log.Fatalln("CheckParticipatesIn failed:", err)
				}
				if amount > 0 {
					log.Println("Already participating as", utils.PubKeyToHex(pubKey.Key), "with", utils.FormatGrams(amount), "stake")
					continue
				}

				validatorAdnlKey, err := s.GetKey("adnlkey", node.ID, activeElectionID)
				if err != nil {
					log.Println("Failed to get adnlkey from", node.HostPort, err)
				}
				if validatorAdnlKey.Key == "" {

					validatorAdnlKey, err = vc.ValidatorCreateNewKey(node, activeElectionID)
					if err != nil {
						log.Fatalln("adnl validatorCreateNewKey failed", err)
					}
					log.Println(validatorAdnlKey.Key)
					validatorAdnlKey.Type = "adnlkey"
					keyID, err := s.AddKey(validatorAdnlKey)
					if err != nil {
						log.Println("failed to save key to db", keyID, err)
					}
					if vc.ValidatorAddAdnl(node, validatorAdnlKey.Key, 0) {
						log.Println("Added ADNL for key hash:", validatorAdnlKey.Key)
					}

					if vc.ValidatorAddValidatorAddr(node, validatorKey.Key, validatorAdnlKey.Key, activeElectionID+70000) {
						log.Println("Added validator addres for key hash:", validatorKey.Key, validatorAdnlKey.Key)
					}
				}

				log.Println("loaded adnlkey from db", validatorAdnlKey.Key)

				fiftElectReq, _ := f.FiftValidatorElectReq(wallet.Addr, activeElectionID, maxFactor, validatorAdnlKey.Key)
				log.Println("fiftelectreq", fiftElectReq)

				signature, err := vc.ValidatorSign(node, validatorKey.Key, fiftElectReq)
				if err != nil {
					log.Fatalln("validatorSign failed", err)
				}

				f.FiftValidatorElectSigned(wallet.Addr, activeElectionID, maxFactor, validatorAdnlKey.Key, pubKey.Key, signature)
				err = cln.UpdateTonConnection()
				if err != nil {
					log.Fatalln("UpdateTonConnection", err)
				}
				seqno, err := cln.GetWalletSeqno(wallet.Addr)
				if err != nil {
					log.Fatalln("GetWalletSeqno failed:", err)
				}
				walletQueryFile, err := f.FiftWalletQuery(wallet.FilePath, currentElectorAddress, seqno, stakeAmount, "validator-query.boc")
				if err != nil {
					log.Println(err)
				}
				err = cln.UpdateTonConnection()
				if err != nil {
					log.Fatalln("UpdateTonConnection", err)
				}
				cln.TonlibSendFile(walletQueryFile)
				participate := database.Participate{
					NodeID:      node.ID,
					ElectionID:  activeElectionID,
					StakeAmount: int64(stakeAmount),
					MaxFactor:   maxFactor,
				}
				_, err = s.AddParticipate(participate)
				if err != nil {
					log.Println("Failed to add participate record to DB:", err)
				}
				time.Sleep(10 * time.Second)
			}

		}
	}

}
