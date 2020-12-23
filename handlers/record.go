package handlers

import (
	"fmt"
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
	remoteAddr := c.Request.RemoteAddr
	logrus.Infof("Remote address: %v",remoteAddr)
	var req UpdateRecordReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("bad request,not json?")
		c.JSON(200, gin.H{"message": "bad request,invalid json"})
		return
	}

	if req.NewIP == "" || req.RR == "" || req.Domain == "" || req.Type == "" || req.AccessKey == "" || req.AccessSecret == "" {
		logrus.Warnf("some request field is empty")
		c.JSON(200, gin.H{"message": types.ErrReqFieldEmpty, "format": "need new_ip,rr,domain,type,ttl,access_key,access_secret"})
		return
	}
	regionID := viper.GetString("server.regionID")
	recordOperator, err := recordOperation.NewOperator(regionID, req.AccessKey, req.AccessSecret, "", "")
	if err != nil {
		logrus.Infof("NewOperator error: %v",err)
		c.JSON(200,gin.H{"message":fmt.Sprintf("Internal error: %v",err)})
		return
	}
	err = recordOperator.DoUpdate(req.NewIP, req.RR, req.Domain, req.Type, req.TTL)

	if err != nil {
		msg := fmt.Sprintf("update error: %v",err)
		logrus.Info(msg)
		c.JSON(200, gin.H{"message": msg})
	} else {
		logrus.Info("updated")
		c.JSON(200, gin.H{"message": "updated"})
	}
}
