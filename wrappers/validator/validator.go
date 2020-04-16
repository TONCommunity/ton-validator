package validator

import (
	"fmt"
	"github.com/mercuryoio/ton-validator/database"
	"github.com/mercuryoio/ton-validator/utils"
	"log"
	"strconv"
	"strings"
)

//Config config
type Config struct {
	ValidatorConsole *string
	Verbose          *bool
}

//NewClient set config
func NewClient(config *Config) *Config {
	return config
}

//ValidatorStats stats
type ValidatorStats struct {
	unixtime                        int64
	masterchainblock                string
	masterchainblocktime            int64
	gcmasterchainblock              string
	keymasterchainblock             string
	knownkeymasterchainblock        string
	rotatemasterchainblock          string
	stateserializermasterchainseqno int64
	shardclientmasterchainseqno     int64
}

//ValidatorAddPermKey add perm key
func (c *Config) ValidatorAddPermKey(node database.Node, keyHash string, electionDate int64) bool {
	expireAt := electionDate + 70000
	/* args := fmt.Sprintf("-c addpermkey %s %d %d", keyHash, electionDate, expireAt)
	output, err := validatorEngineConsoleReq(args) */
	valCmd := fmt.Sprintf("-c addpermkey %s %d %d", keyHash, electionDate, expireAt)
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}

	output, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Println(err)
		return false
	}
	if !strings.Contains(output, "success") {
		return false
	}

	return true
}

//ValidatorAddTempKey add temp key
func (c *Config) ValidatorAddTempKey(node database.Node, permKeyHash string, keyHash string, expireAt int64) bool {
	//args := fmt.Sprintf("-c addtempkey %s %s %d", permKeyHash, keyHash, expireAt)
	//output, err := validatorEngineConsoleReq(args)
	valCmd := fmt.Sprintf("-c addtempkey %s %s %d", permKeyHash, keyHash, expireAt)
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}

	_, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Fatalln("validatorEngineConsoleReq err:", err)
		return false
	}
	return true
}

//ValidatorAddAdnl add adnl
func (c *Config) ValidatorAddAdnl(node database.Node, keyHash string, category int) bool {
	//	args := fmt.Sprintf("-c addadnl %s %d", keyHash, category)
	//	output, err := validatorEngineConsoleReq(args)
	valCmd := fmt.Sprintf("-c addadnl %s %d", keyHash, category)
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}

	_, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Fatalln("validatorEngineConsoleReq err:", err)
		return false
	}
	return true
}

//ValidatorAddValidatorAddr add validator addr
func (c *Config) ValidatorAddValidatorAddr(node database.Node, permKeyHash string, keyHash string, expireAt int64) bool {
	valCmd := fmt.Sprintf("-c addvalidatoraddr %s %s %d", permKeyHash, keyHash, expireAt)
	//output, err := validatorEngineConsoleReq(args)
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}
	_, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Fatalln("validatorEngineConsoleReq err:", err)
		return false
	}
	return true
}

//ValidatorSign sign message
func (c *Config) ValidatorSign(node database.Node, keyHash string, data string) (string, error) {
	valCmd := fmt.Sprintf("-c sign %s %s", keyHash, data)
	//output, err := validatorEngineConsoleReq(args)
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}
	output, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Fatalln("validatorEngineConsoleReq err:", err)
		return string(output), err
	}
	i := strings.Index(output, "got signature")
	output = output[i:]
	output = strings.TrimPrefix(output, "got signature ")
	output = strings.TrimSpace(output)
	return string(output), nil
}

//ValidatorCreateNewKey create new key
func (c *Config) ValidatorCreateNewKey(node database.Node, electionID int64) (database.Key, error) {
	valCmd := "-c newkey"
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}

	output, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Fatalln("validatorEngineConsoleReq err:", err)
		return database.Key{}, err
	}
	i := strings.Index(output, "new key")
	output = output[i:]
	output = strings.TrimPrefix(output, "new key ")
	output = strings.TrimSpace(output)
	key := database.Key{
		Key:        output,
		ElectionID: electionID,
		NodeID:     node.ID,
		Type:       "key",
	}
	return key, err
}

//ValidatorGetPublicKey get public key
func (c *Config) ValidatorGetPublicKey(node database.Node, signingKey string, electionID int64) (database.Key, error) {
	valCmd := fmt.Sprintf("-c exportpub %s", signingKey)
	//output, err := validatorEngineConsoleReq(args)
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}
	output, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Println("getValidatorPublicKey err:", err)
		return database.Key{}, err
	}
	i := strings.Index(output, "got public key:")
	output = output[i:]
	output = strings.TrimPrefix(output, "got public key: ")
	output = strings.TrimSpace(output)
	key := database.Key{
		Key:        output,
		ElectionID: electionID,
		NodeID:     node.ID,
		Type:       "pubkey",
	}
	return key, err
}

//ValGetStats getstats
func (c *Config) ValGetStats(node database.Node) ValidatorStats {
	valCmd := "-c getstats"
	args := []string{"-k", node.ClientCert, "-p", node.ServerPub, "-a", node.HostPort, "-v 0", valCmd}
	output, err := utils.CmdExec(*c.ValidatorConsole, *c.Verbose, args...)
	if err != nil {
		log.Println("Validator getstats err:", err)
		//return string(output), err
	}
	i := strings.Index(output, "unixtime")
	output = output[i:]
	lines := strings.Split(output, "\n")
	var values []string
	for _, s := range lines {
		output = s
		output = strings.TrimPrefix(output, "unixtime")
		output = strings.TrimPrefix(output, "masterchainblocktime")
		output = strings.TrimPrefix(output, "masterchainblock")
		output = strings.TrimPrefix(output, "gcmasterchainblock")
		output = strings.TrimPrefix(output, "keymasterchainblock")
		output = strings.TrimPrefix(output, "knownkeymasterchainblock")
		output = strings.TrimPrefix(output, "rotatemasterchainblock")
		output = strings.TrimPrefix(output, "stateserializermasterchainseqno")
		output = strings.TrimPrefix(output, "shardclientmasterchainseqno")
		output = strings.TrimSpace(output)
		values = append(values, output)

	}
	unixtime, _ := strconv.ParseInt(values[0], 10, 64)
	masterchainblocktime, _ := strconv.ParseInt(values[2], 10, 64)
	stateserializermasterchainseqno, _ := strconv.ParseInt(values[7], 10, 64)
	shardclientmasterchainseqno, _ := strconv.ParseInt(values[8], 10, 64)
	stats := ValidatorStats{
		unixtime:                        unixtime,
		masterchainblock:                values[1],
		masterchainblocktime:            masterchainblocktime,
		gcmasterchainblock:              values[3],
		keymasterchainblock:             values[4],
		knownkeymasterchainblock:        values[5],
		rotatemasterchainblock:          values[6],
		stateserializermasterchainseqno: stateserializermasterchainseqno,
		shardclientmasterchainseqno:     shardclientmasterchainseqno,
	}
	return stats
}

//CheckNodeSync check if node in sync
func (c *Config) CheckNodeSync(node database.Node) bool {
	stats := c.ValGetStats(node)
	diff := stats.unixtime - stats.masterchainblocktime

	if diff > 25 {
		return false
	}
	return true
}
