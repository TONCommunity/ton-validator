package liteclient

import (
	"github.com/mercuryoio/ton-validator/utils"
	"log"
	"strconv"
	"strings"
)

//Config config
type Config struct {
	LiteClient       *string
	LiteclientConfig *string
	Verbose          *bool
}

//NewClient new client
func NewClient(config *Config) *Config {
	return config
}

//GetCurrentElectorAddress get current elector address
func (c *Config) GetCurrentElectorAddress() (string, error) {
	args := []string{"-C", *c.LiteclientConfig, "-v 0", "-r", "-c", "getconfig 1"}

	output, err := utils.CmdExec(*c.LiteClient, *c.Verbose, args...)
	if err != nil {
		log.Println("getCurrentElectorAddress err:", err)
		return string(output), err
	}
	i := strings.Index(output, "x{")
	output = output[i:]
	output = strings.TrimPrefix(output, "x{")
	output = strings.Replace(output, "}", "", 1)
	output = strings.TrimSpace(output)
	output = "-1:" + output

	return string(output), err
}

//ElectionPeriods election periods
type ElectionPeriods struct {
	ValidatorsElectedFor int64
	ElectionsStartBefore int64
	ElectionsEndBefore   int64
	StakeHeldFor         int64
}

//StakeConfig network stake config
type StakeConfig struct {
	MinStake       int64
	MaxStake       int64
	MinTotalStake  int64
	MaxStakeFactor int64
}

//GetElectionConfig get election config
func (c *Config) GetElectionConfig() (ElectionPeriods, error) {
	args := []string{"-C", *c.LiteclientConfig, "-v 0", "-r", "-c", "getconfig 15"}

	output, err := utils.CmdExec(*c.LiteClient, *c.Verbose, args...)
	if err != nil {
		log.Println("getCurrentElectorAddress err:", err)
		return ElectionPeriods{}, err
	}
	i := strings.Index(output, "validators_elected_for")
	output = output[i:]

	output = strings.Replace(output, ")", "", 1)
	lines := strings.Split(output, "\n")
	output = lines[0]
	lines = strings.Split(output, " ")
	var values []string
	for _, s := range lines {
		strs := strings.Split(s, ":")
		values = append(values, strs[1])
	}

	var periods ElectionPeriods
	periods.ValidatorsElectedFor, _ = strconv.ParseInt(values[0], 10, 64)
	periods.ElectionsStartBefore, _ = strconv.ParseInt(values[1], 10, 64)
	periods.ElectionsEndBefore, _ = strconv.ParseInt(values[2], 10, 64)
	periods.StakeHeldFor, _ = strconv.ParseInt(values[3], 10, 64)
	output = strings.TrimSpace(output)

	return periods, err

}

//GetStakeConfig get stake config
func (c *Config) GetStakeConfig() (StakeConfig, error) {
	args := []string{"-C", *c.LiteclientConfig, "-v 0", "-r", "-c", "getconfig 17"}

	output, err := utils.CmdExec(*c.LiteClient, *c.Verbose, args...)
	if err != nil {
		log.Println("getStakeConfig err:", err)
		return StakeConfig{}, err
	}
	i := strings.Index(output, "min_stake")
	output = output[i:]
	lines := strings.Split(output, "\n")
	var values []string
	for _, s := range lines {
		i = strings.Index(s, "value")
		if i > 0 {
			output = s[i:]
			output = strings.TrimPrefix(output, "value:")
			output = strings.Replace(output, "))", "", 1)
			output = strings.TrimRight(output, " ")
			output = strings.TrimSpace(output)
			i = strings.Index(output, "max_stake_factor")
			if i > 0 {
				minTotalStake := strings.Split(output, " ")
				if minTotalStake[0] != "" {
					values = append(values, minTotalStake[0])
				}
				output = output[i:]
				output = strings.TrimPrefix(output, "max_stake_factor:")
				output = strings.Replace(output, ")", "", 1)
				output = strings.TrimSpace(output)
				if output != "" {
					values = append(values, output)
				}
			} else if output != "" {
				values = append(values, output)
			}
		}
	}

	var stakeConfig StakeConfig
	stakeConfig.MinStake, _ = strconv.ParseInt(values[0], 10, 64)
	stakeConfig.MaxStake, _ = strconv.ParseInt(values[1], 10, 64)
	stakeConfig.MinTotalStake, _ = strconv.ParseInt(values[2], 10, 64)
	stakeConfig.MaxStakeFactor, _ = strconv.ParseInt(values[3], 10, 64)

	return stakeConfig, nil
}
