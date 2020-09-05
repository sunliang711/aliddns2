package recordOperation

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type Operator struct {
	// Config        *config.Config
	client        *alidns.Client
	filterIPRegex *regexp.Regexp
}

func NewOperator(regionID, accessKey, accessSecret, newIPCommand, filterIPRegex string) (*Operator, error) {
	o := &Operator{
		// Config: cfg,
	}
	var err error
	// regionID := viper.GetString("server.regionID")
	// accessKey := viper.GetString("local.accessKey")
	// accessSecret := viper.GetString("local.accessSecret")

	// newIPCommand := viper.GetString("local.newIPCommand")
	// filterIPRegex := viper.GetString("local.filterIPRegex")
	o.client, err = alidns.NewClientWithAccessKey(regionID, accessKey, accessSecret)
	if err != nil {
		return nil, err
	}
	if len(newIPCommand) == 0 {
		//compile regex
		re, err := regexp.Compile(filterIPRegex)
		if err != nil {
			msg := fmt.Sprintf("NewipCommand is null and compile filterIpRegex failed:%v", err)
			log.Println(msg)
			return nil, fmt.Errorf(msg)
		}
		o.filterIPRegex = re
	}
	return o, nil
}
