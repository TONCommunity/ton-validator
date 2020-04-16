package fift

import (
	"github.com/mercuryoio/ton-validator/utils"
	tonlib "github.com/mercuryoio/tonlib-go/v2"
	"log"
	"strconv"
	"strings"
)

//Config config
type Config struct {
	FiftBin                 *string
	FiftPath                *string
	WalletFif               *string
	RecoverFif              *string
	ValidatorElectReqFif    *string
	ValidatorElectSignedFif *string
	ValidatorWalletFile     *string
	Verbose                 *bool
}

//NewClient init new connection to the database
func NewClient(config *Config) *Config {
	return config
}

//FiftValidatorElectReq validator elect req
func (c *Config) FiftValidatorElectReq(walletAddr string, electionTimestamp int64, maxFactor, adnlKey string) (string, error) {

	args := []string{"-s", *c.ValidatorElectReqFif, walletAddr, strconv.Itoa(int(electionTimestamp)), maxFactor, adnlKey}

	output, err := utils.CmdExec(*c.FiftBin, *c.Verbose, args...)
	if err != nil {
		log.Println("fiftValidatorElectReq err:", err)
		return string(output), err
	}
	i := strings.Index(output, "Creating")
	output = output[i:]

	var strArr []string
	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		strArr = append(strArr, line)
	}
	output = strings.TrimSpace(strArr[1])
	return string(output), nil

}

//FiftValidatorElectSigned validator elect signed
func (c *Config) FiftValidatorElectSigned(walletAddr string, electionTimestamp int64, maxFactor, adnlKey, pubKey, signature string) (string, error) {
	args := []string{"-s", *c.ValidatorElectSignedFif, walletAddr, strconv.Itoa(int(electionTimestamp)), maxFactor, adnlKey, pubKey, signature}

	output, err := utils.CmdExec(*c.FiftBin, *c.Verbose, args...)
	if err != nil {
		log.Println("fiftValidatorElectSigned err:", err)
		return string(output), err
	}
	return string(output), err
}

//FiftWalletQuery wallet query
func (c *Config) FiftWalletQuery(walletFile, destAddr string, seqno int64, amount int, bocFile string) (string, error) {
	var args []string
	if bocFile == "" {
		args = []string{"-s", *c.WalletFif, walletFile, destAddr, strconv.Itoa(int(seqno)), strconv.Itoa(amount) + "."}
	} else {
		args = []string{"-s", *c.WalletFif, walletFile, destAddr, strconv.Itoa(int(seqno)), strconv.Itoa(amount) + ".", "-B", bocFile}
	}
	output, err := utils.CmdExec(*c.FiftBin, *c.Verbose, args...)
	if err != nil {
		log.Println("fiftWalletQuery err:", err)
		log.Println(string(output))
		return "", err
	}
	i := strings.Index(output, "(Saved to file ")
	output = output[i:]
	output = strings.TrimPrefix(output, "(Saved to file ")
	output = strings.Replace(output, ")", "", 1)
	output = strings.TrimSpace(output)
	return output, nil
}

//FiftGenRecoverQueryFile recover query file
func (c *Config) FiftGenRecoverQueryFile() (string, error) {
	args := []string{"-s", *c.RecoverFif}
	output, err := utils.CmdExec(*c.FiftBin, *c.Verbose, args...)
	if err != nil {
		return "", err
	}

	i := strings.Index(output, "Saved to file ")
	output = output[i:]
	output = strings.TrimPrefix(output, "Saved to file ")
	output = strings.TrimSpace(output)
	return output, nil
}

//RecoverStake recover stake
func (c *Config) RecoverStake(cln *tonlib.Client, walletFile, destAddress string, seqno int64) error {
	recoverQueryFile, err := c.FiftGenRecoverQueryFile()
	if err != nil {
		return err
	}
	walletQueryFile, err := c.FiftWalletQuery(walletFile, destAddress, seqno, 1, recoverQueryFile)
	if err != nil {
		return err
	}
	cln.TonlibSendFile(walletQueryFile)
	return nil
}
