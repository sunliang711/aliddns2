package recordOperation

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/spf13/viper"
	"github.com/sunliang711/aliddns2/types"
)

func (o *Operator) GetNewIP() (string, error) {
	log.Printf("GetNewIP()")
	defer func() {
		log.Printf("Leave GetNewIP()")
	}()
	newIPCommand := viper.GetString("local.newIPCommand")
	newIPSource := viper.GetString("local.newIPSource")
	masterIndex := viper.GetUint("local.masterIndex")
	slaveIndex := viper.GetUint("local.slaveIndex")
	if len(newIPCommand) > 0 {
		result, err := exec.Command("sh", "-c", newIPCommand).Output()
		if err != nil {
			return "", types.ErrCannotGetIpFromIpCommnad
		}
		return string(result), nil
	}
	return getIPByRegex(newIPSource, o.filterIPRegex, masterIndex, slaveIndex)
}

func getIPByRegex(url string, re *regexp.Regexp, masterIndex, slaveIndex uint) (string, error) {
	log.Printf("getIpByRegex(): url:%s,masterIndex: %d,slaveIndex: %d", url, masterIndex, slaveIndex)

	res, err := http.Get(url)
	if err != nil {
		log.Printf(">>http.Get error: %s", err)
		return "", err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		log.Printf(">>io.Copy error: %s", err)
		return "", err
	}
	result := re.FindAllStringSubmatch(buf.String(), -1)
	if uint(len(result)) <= masterIndex {
		return "", types.ErrCannotGetIpFromIpSource
	}
	if uint(len(result[masterIndex])) <= slaveIndex {
		return "", types.ErrCannotGetIpFromIpSource
	}
	return result[masterIndex][slaveIndex], nil
}
