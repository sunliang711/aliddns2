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
	logrus.Infof(">> Remote address: %v", remoteAddr)
	var req UpdateRecordReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("bad request,not json?")
		c.JSON(200, gin.H{"message": "bad request,invalid json"})
		return
	}
	newIP := req.NewIP
	if newIP == "" {
		logrus.Infof("use client ip instead of request body")
		newIP = c.ClientIP()
	}

	accessKey := req.AccessKey
	if accessKey == "" {
		logrus.Infof("use default access key")
		accessKey = viper.GetString("server.defaultAccessKey")
	}
	accessSecret := req.AccessSecret
	if accessSecret == "" {
		logrus.Infof("use default access secret")
		accessSecret = viper.GetString("server.defaultAccessSecret")
	}

	rr := req.RR
	domain := req.Domain
	domainType := req.Type
	if domainType == "" {
		domainType = "A"
	}
	ttl := req.TTL
	if ttl == "" {
		ttl = "600"
	}

	if newIP == "" || rr == "" || domain == "" || domainType == "" || accessKey == "" || accessSecret == "" {
		logrus.Warnf("some request field is empty")
		c.JSON(200, gin.H{"message": types.ErrReqFieldEmpty, "format": "need new_ip(optional),rr,domain,type,ttl,access_key(use default when empty),access_secret(use default when empty)"})
		return
	}
	regionID := viper.GetString("server.regionID")
	recordOperator, err := recordOperation.NewOperator(regionID, accessKey, accessSecret, "", "")
	if err != nil {
		logrus.Infof("NewOperator error: %v", err)
		c.JSON(200, gin.H{"message": fmt.Sprintf("Internal error: %v", err)})
		return
	}
	err = recordOperator.DoUpdate(newIP, rr, domain, domainType, ttl)

	if err == types.ErrNotNeedUpdate {
		msg := fmt.Sprintf("%v", err)
		logrus.Info(msg)
		c.JSON(200, gin.H{"code": 0, "message": msg})
	} else if err != nil {
		msg := fmt.Sprintf("Update error: %v", err)
		logrus.Info(msg)
		c.JSON(200, gin.H{"code": 1, "message": msg})
	} else {
		msg := fmt.Sprintf("Update %v.%v to %v", rr, domain, newIP)
		logrus.Info(msg)
		c.JSON(200, gin.H{"code": 0, "message": msg})
	}
}
