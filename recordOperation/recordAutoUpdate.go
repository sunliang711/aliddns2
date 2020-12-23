package recordOperation

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/sunliang711/aliddns2/types"
)

// LoopUpdate update record every 'updateInterval' seconds
func (o *Operator) LoopUpdate() {
	var (
		RR             = viper.GetString("local.RR")
		DomainName     = viper.GetString("local.DomainName")
		Type           = viper.GetString("local.Type")
		TTL            = viper.GetString("local.TTL")
		updateInterval = viper.GetUint("local.updateInterval")
	)
	if updateInterval > 0 {
		tick := time.NewTicker(time.Duration(updateInterval) * time.Second)
		for {
			select {
			case <-tick.C:
				newIP, err := o.GetNewIP()
				if err != nil {
					log.Printf(">>GetNewIp error:%s", err)
					continue
				}
				newIP = strings.TrimSpace(newIP)
				o.DoUpdate(newIP, RR, DomainName, Type, TTL)
			}
		}
	}
}

// DoUpdate update record
func (o *Operator) DoUpdate(newIP, RR, DomainName, Type, TTL string) error {
	//1. getRecordId
	subDomain := fmt.Sprintf("%v.%v", RR, DomainName)
	pageNumber := viper.GetString("server.pageNumber")
	pageSize := viper.GetString("server.pageSize")
	recordID, currentDNSIP, err := o.GetRecordId(subDomain, pageSize, pageNumber)
	if err != nil {
		if err == types.ErrNoSubDomain {
			recordID, err = o.AddRecord(DomainName, Type, RR, newIP, TTL)
			return nil
		}
		logrus.Printf(">>Exist such subDomain,but cann't get recordID")
		return err
	}
	currentDNSIP = strings.TrimSpace(currentDNSIP)
	logrus.Printf("Current ip: %s", currentDNSIP)
	if currentDNSIP != newIP {
		//2. update
		res, err := o.UpdateRecord(recordID, RR, Type, newIP, TTL)
		if err != nil {
			return err
		}
		logrus.Printf(">>update OK:%v", res)
		return nil
	} else {
		return types.ErrNotNeedUpdate
	}
}
