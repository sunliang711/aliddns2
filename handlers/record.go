package handlers

import (
	"fmt"
	"strings"

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

func (req *UpdateRecordReq) check(ctx *gin.Context) error {
	var err error
	if req.NewIP == "" {
		logrus.Infof("use client ip instead of request body")
		req.NewIP = ctx.ClientIP()
	}

	if req.RR == "" {
		return fmt.Errorf("require rr")
	}

	if req.Domain == "" {
		return fmt.Errorf("require domain")
	}

	if req.AccessKey == "" || req.AccessSecret == "" {
		req.AccessKey, req.AccessSecret, err = getAccessKey(req.Domain)
		if err != nil {
			return err
		}
	}

	if req.TTL == "" {
		req.TTL = "600"
	}

	if req.Type == "" {
		req.Type = "A"
	}

	return nil
}

func getAccessKey(domain string) (string, string, error) {
	accessKeys := strings.Split(viper.GetString("server.accessKeys"), ";")
	// accessKey format: <domain>:<access_key>:<access_secret>
	for _, accessKey := range accessKeys {
		parts := strings.Split(accessKey, ":")
		if len(parts) != 3 {
			logrus.Warnf("invalid format or accessKeys in config")
			break
		}
		dm := parts[0]
		key := parts[1]
		secret := parts[2]
		if dm == domain {
			logrus.Infof("got access key for domain: %s", domain)
			return key, secret, nil
		}
	}

	defaultAccessKey := viper.GetString("server.defaultAccessKey")
	defaultAccessSecret := viper.GetString("server.defaultAccessSecret")
	if defaultAccessKey == "" || defaultAccessSecret == "" {
		return "", "", fmt.Errorf("no default access key or secret")
	}

	logrus.Infof("use default access key and secret")
	// return defaultAccessKey defaultAccessSecret
	return defaultAccessKey, defaultAccessSecret, nil
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

	if err = req.check(c); err != nil {
		logrus.Errorf("%s", err.Error())
		c.JSON(200, gin.H{"message": err.Error()})
		return
	}
	regionID := viper.GetString("server.regionID")
	recordOperator, err := recordOperation.NewOperator(regionID, req.AccessKey, req.AccessSecret, "", "")
	if err != nil {
		logrus.Infof("NewOperator error: %v", err)
		c.JSON(200, gin.H{"message": fmt.Sprintf("Internal error: %v", err)})
		return
	}
	err = recordOperator.DoUpdate(req.NewIP, req.RR, req.Domain, req.Type, req.TTL)

	if err == types.ErrNotNeedUpdate {
		msg := fmt.Sprintf("%v", err)
		logrus.Info(msg)
		c.JSON(200, gin.H{"code": 0, "message": msg})
	} else if err != nil {
		msg := fmt.Sprintf("Update error: %v", err)
		logrus.Info(msg)
		c.JSON(200, gin.H{"code": 1, "message": msg})
	} else {
		msg := fmt.Sprintf("Update %v.%v to %v", req.RR, req.Domain, req.NewIP)
		logrus.Info(msg)
		c.JSON(200, gin.H{"code": 0, "message": msg})
	}
}
