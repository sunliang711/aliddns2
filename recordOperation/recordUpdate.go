package recordOperation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/sunliang711/aliddns2/types"
)

type UpdateRecordResponse struct {
	RequestId string `json:"RequestId"`
	RecordId  string `json:"RecordId"`
}

func (o *Operator) UpdateRecord(recordId, RR, Type, Value, TTL string) (string, error) {
	log.Printf("UpdateRecord(): recordId:%v, RR:%v, Type:%v, Value:%v, TTL:%v", recordId, RR, Type, Value, TTL)
	defer func() {
		log.Printf("Leave UpdateRecord()")
	}()

	request := alidns.CreateUpdateDomainRecordRequest()

	request.RecordId = recordId
	request.RR = RR
	request.Type = Type
	request.Value = Value
	request.TTL = requests.Integer(TTL)

	response, err := o.client.UpdateDomainRecord(request)
	if err != nil {
		log.Printf(">>UpdateDomainRecord error:%v", err)
		fmt.Print(err.Error())
	}

	if response.GetHttpStatus() != http.StatusOK {
		log.Printf(">>%v", types.ErrHttpStatusNotOK)
		return "", types.ErrHttpStatusNotOK
	}
	var res UpdateRecordResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &res)
	if err != nil {
		log.Printf(">>json.Unmarshal error:%v", err)
		return "", err
	}

	if res.RecordId != recordId {
		log.Printf(">>%v", types.ErrResponseIdNotMatchRequestId)
		return "", types.ErrResponseIdNotMatchRequestId
	}

	return response.GetHttpContentString(), nil
}
