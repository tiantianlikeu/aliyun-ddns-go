package main

import (
	myUtil "aliyun-ddns-go/util"
	"encoding/json"
	"fmt"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"time"
)

var Client *alidns20150109.Client

var lastIp string

func CreateClient() (_result *alidns20150109.Client) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(myUtil.Config.AccessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(myUtil.Config.AccessKeySecret),
	}
	// 访问的域名
	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	_result = &alidns20150109.Client{}
	if _result, _err := alidns20150109.NewClient(config); _err == nil {
		return _result
	}
	return nil
}

func UpdateDomainRecord(recordId, newIp string) {

	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{}
	updateDomainRecordRequest.SetType("A")
	updateDomainRecordRequest.SetValue(newIp)
	updateDomainRecordRequest.SetRR(myUtil.Config.Rr)
	updateDomainRecordRequest.SetTTL(600)
	updateDomainRecordRequest.SetRecordId(recordId)
	runtime := &util.RuntimeOptions{}

	result, err := Client.UpdateDomainRecordWithOptions(updateDomainRecordRequest, runtime)
	fmt.Println(err)
	fmt.Println(result)
	if err == nil {
		lastIp = newIp
		fmt.Println(time.Now(), "DDNS ip 修改成功~    新解析的IP为：", newIp)
	}

}

func main() {
	if Client == nil {
		Client = CreateClient()
	}
	c := cron.New()
	EntryID, err := c.AddFunc("*/1 * * * *", ddns)
	fmt.Println(time.Now(), EntryID, err)
	c.Start()
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}

func ddns() {
	newIp := GetNewIp()
	if newIp != "" && newIp != lastIp {
		recordId, oldIp := GetDnsRecordId()
		if newIp != *oldIp {
			UpdateDomainRecord(*recordId, newIp)
		} else {
			lastIp = *oldIp
			fmt.Println(time.Now(), "ip未变化，oldIp is: "+*oldIp+", newIp is: "+newIp)
		}
	} else {
		fmt.Println(time.Now(), "ip未变化，lastIp is: "+lastIp+", newIp is: "+newIp)
	}
}

func GetDnsRecordId() (recordId *string, oldIp *string) {
	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{}
	describeDomainRecordsRequest.SetDomainName(myUtil.Config.DomainName)
	describeDomainRecordsRequest.SetRRKeyWord(myUtil.Config.RrKeyword)
	runtime := &util.RuntimeOptions{}

	result, recordErr := Client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, runtime)
	if recordErr != nil {
		return nil, nil
	}
	records := result.Body.DomainRecords.Record
	if len(records) > 0 {
		recordId = records[0].RecordId
		oldIp = records[0].Value
	}
	return
}

func GetNewIp() string {
	var resp *http.Response
	var err error
	// 接口调用
	resp, err = http.Get("https://jsonip.com/")
	if err != nil {
		log.Println("err:", err)
		return ""
	}
	rs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("readAll,err:", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("http请求，流关闭失败")
		}
	}(resp.Body)
	var jsonip JsonIp
	if err = json.Unmarshal(rs, &jsonip); err == nil {
		return jsonip.Ip
	}
	return ""
}

type JsonIp struct {
	Ip      string
	Country string
}
