package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/sunliang711/aliddns2/recordOperation"
	"github.com/sunliang711/aliddns2/types"
)

// UpdateRecordReq 存储请求更新的主机记录信息
type UpdateRecordReq struct {
	NewIP        string `json:"new_ip"`
	RR           string `json:"rr"`     // 二级域名
	Domain       string `json:"domain"` // 一级域名
	Type         string `json:"type"`
	TTL          string `json:"ttl"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
}

// UpdateRecord 根据请求来更新主机记录，如果不存在，则会新增主机记录
func UpdateRecord(c *gin.Context) {
	var req UpdateRecordReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("bad request,not json?")
		c.JSON(1, "bad request")
		return
	}

	if req.NewIP == "" || req.RR == "" || req.Domain == "" || req.Type == "" || req.AccessKey == "" || req.AccessSecret == "" {
		logrus.Warnf("some request field is empty")
		c.JSON(1, gin.H{"message": types.ErrReqFieldEmpty})
		return
	}
	regionID := viper.GetString("server.regionID")
	recordOperator, err := recordOperation.NewOperator(regionID, req.AccessKey, req.AccessSecret, "", "")
	err = recordOperator.DoUpdate(req.NewIP, req.RR, req.Domain, req.Type, req.TTL)

	if err != nil {
		c.JSON(1, gin.H{"message": err})
	} else {
		c.JSON(0, gin.H{"message": "OK"})
	}
}
