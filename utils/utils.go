package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strings"
)

//FormatGrams Format nanograms to grams float
func FormatGrams(n int64) string {
	belowZero := n < 0
	if belowZero {
		n = -n
	}

	in := fmt.Sprintf("%010d", n)
	numOfDigits := len(in)

	numOfCommas := (numOfDigits - 1) / 9

	out := make([]byte, len(in)+numOfCommas)

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			if belowZero {
				return "-" + string(out)
			}
			return string(out)
		}
		if k++; k == 9 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

//Hex2int hex to int conversion
func Hex2int(hexStr string) *big.Int {
	i := new(big.Int)
	i.SetString(hexStr, 16)
	return i
}

//FileExists check file existence
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

//CmdExec Command execution
func CmdExec(path string, verbose bool, args ...string) (string, error) {
	cmd := exec.Command(path, args...)

	cmd.StdinPipe()
	output, err := cmd.CombinedOutput()
	if verbose {
		log.Println(cmd.Args)
		log.Println(string(output))
	}
	if err != nil {
		log.Println(path, "failed:", err)
		return string(output), err
	}
	return string(output), err
}

//AppendToFile append to file
func AppendToFile(path, key, value string) {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(key + "=" + value + "\n"); err != nil {
		log.Println(err)
	}
}

//PubKeyToHex pulibc key to hex conversion
func PubKeyToHex(pubKey string) string {
	p, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		log.Fatalln(err)
	}
	h := hex.EncodeToString(p)
	h = strings.TrimPrefix(h, "c6b41348")
	return h
}
