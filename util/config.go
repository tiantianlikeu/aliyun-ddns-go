package util

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var Config *config

type config struct {
	DomainName      string
	Rr              string
	RrKeyword       string
	AccessKeyId     string
	AccessKeySecret string
}

func NewConfig() *config {
	cf := new(config)
	cf.loadConfigYaml()
	return cf
}

func init() {
	Config = NewConfig()
}

func (cf *config) loadConfigYaml() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Println("conf路径未找到config.yaml，开始在根目录查询文件~")
		viper.AddConfigPath(".")
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("配置文件读取失败: %v \n", err))
		}
	}
	if err := viper.Unmarshal(cf); nil != err {
		log.Fatalf("赋值配置对象失败，异常信息：%v", err)
	}
}
