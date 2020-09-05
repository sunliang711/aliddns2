package recordOperation

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/sunliang711/aliddns2/types"
)

type DeleteRecordResponse struct {
	RequestId string `json:"RequestId"`
	RecordId  string `json:"RecordId"`
}

func (o *Operator) DeleteRecord(recordId string) (string, error) {
	log.Printf("DeleteRecord(): recordId: %v", recordId)
	defer func() {
		log.Printf("Leave DeleteRecord()")
	}()

	request := alidns.CreateDeleteDomainRecordRequest()

	request.RecordId = recordId

	response, err := o.client.DeleteDomainRecord(request)
	if err != nil {
		log.Printf(">>DeleteDomainRecord() error:%v", err)
		return "", err
	}
	if response.GetHttpStatus() != http.StatusOK {
		log.Printf(">>%v", types.ErrHttpStatusNotOK)
		return "", types.ErrHttpStatusNotOK
	}

	var res DeleteRecordResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &res)
	if err != nil {
		log.Printf(">>json.Unmarshal error: %v", err)
		return "", err
	}
	if res.RecordId != recordId {
		log.Printf(">>%v", types.ErrResponseIdNotMatchRequestId)
		log.Printf(">>response id:%v, request id:%v", res.RecordId, recordId)
		return "", types.ErrResponseIdNotMatchRequestId
	}
	return response.GetHttpContentString(), nil
}
