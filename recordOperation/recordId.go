package recordOperation

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/aliddns2/types"
)

type SubDomainRecordResponse struct {
	PageNumber    int    `json:"PageNumber"`
	TotalCount    int    `json:"TotalCount"`
	PageSize      int    `json:"PageSize"`
	RequestId     string `json:"RequestId"`
	DomainRecords struct {
		Record []struct {
			RR         string `json:"RR"`
			Status     string `json:"Status"`
			Value      string `json:"Value"`
			Weight     int    `json:"Weight"`
			RecordId   string `json:"RecordId"`
			Type       string `json:"Type"`
			DomainName string `json:"DomainName"`
			Locked     bool   `json:"Locked"`
			Line       string `json:"Line"`
			TTL        int    `json:"TTL"`
		} `json:"Record"`
	} `json:"DomainRecords"`
}

//查询子域名记录Id
//return : id,current ip,err
func (o *Operator) GetRecordId(subDomain string, pageSize string, pageNumber string) (string, string, error) {
	logrus.Printf("GetRecordId(): subDomain: %v", subDomain)
	defer func() {
		logrus.Printf("Leave GetRecordId()")
	}()

	request := alidns.CreateDescribeSubDomainRecordsRequest()

	request.PageSize = requests.Integer(pageSize)
	request.PageNumber = requests.Integer(pageNumber)
	request.SubDomain = subDomain

	response, err := o.client.DescribeSubDomainRecords(request)
	if err != nil {
		logrus.Printf(">>DescribeSubDomainRecords error:%s", err)
		return "", "", err
	}
	if response.GetHttpStatus() != http.StatusOK {
		logrus.Printf(">>%v", types.ErrHttpStatusNotOK)
		return "", "", types.ErrHttpStatusNotOK
	}
	logrus.Println(">> Response Content: ", response.GetHttpContentString())
	var res SubDomainRecordResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &res)
	if err != nil {
		logrus.Printf(">>json.Unmarshal error:%v", err)
		return "", "", err
	}
	if res.TotalCount == 0 {
		logrus.Printf(">>%v", types.ErrNoSubDomain)
		return "", "", types.ErrNoSubDomain
	}
	//RR.DomainName === subDomain
	if strings.Compare(res.DomainRecords.Record[0].RR+"."+res.DomainRecords.Record[0].DomainName, subDomain) != 0 {
		logrus.Printf(">>%v", types.ErrSubDomainNotMatch)
		return "", "", types.ErrSubDomainNotMatch
	}
	return res.DomainRecords.Record[0].RecordId, res.DomainRecords.Record[0].Value, nil
}
