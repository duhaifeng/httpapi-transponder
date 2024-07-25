/**
 * 解析和缓存全局与自动转发配置的工具
 * @author duhaifeng
 * @date   2021/04/15
 */
package common

import (
	"cv-api-gw/common/busierr"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var Configs *AllConfig

type AllConfig struct {
	System      *SystemConfig
	Log         *LogConfig
	Mysql       *MySQLConfig
	GrpcForward *GrpcForward
	Endpoints   map[string]map[string]*ForwardEndpoint
}

type SystemConfig struct {
	Version             string
	ListenAddr          string
	ListenPort          string
	FakeBackend         bool
	FakeBackendFilePath string
	WebRoot             string
	DefaultBackend      map[string]string
	LogLevel            string
}

type LogConfig struct {
	LogLevel      string
	LogOutput     string
	LogFilePath   string
	LogFileSize   string
	LogFileNumber string
}

type GrpcForward struct {
	ListenPort  string
	PushAddress []string
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Schema   string
}

type ForwardEndpoint struct {
	Endpoint               string
	Method                 string
	ValidateRule           string
	PostProcessor          string
	OnBackendErrorContinue bool
	DoNotWrapResponse      bool
	Backend                []*EndpointBackend
}

type EndpointBackend struct {
	Alias                        string
	UrlPattern                   string
	ResponseProcessor            string
	AttachUrlParamPart           bool
	DoNotPrintDebug              bool
	Method                       string
	Host                         []string
	DropResponse                 bool
	BodyFieldOperationSameAs     string
	ResponseFieldOperationSameAs string
	BodyFieldOperation           []*FieldOperation
	ResponseFieldOperation       []*FieldOperation
	OnErrorCall                  []*EndpointBackend //当发生错误时的回滚处理，一般也是调用后端的接口
}

type FieldOperation struct {
	FieldName           string
	DestName            string
	Operation           string
	ArrayIndexAlignMode string
	FormatConverter     string
	Value               interface{}
}

/**
 * 寻找配置文件的路径
 */
func findConfigFile(fileName string) string {
	// app path
	curDirPath := GetCurrentDirectory() + "/" + fileName
	if IsFileExists(curDirPath) {
		return curDirPath
	}
	// dev path
	curDirPath = GetCurrentDirectory() + "/config/" + fileName
	if IsFileExists(curDirPath) {
		return curDirPath
	}
	curDirPath = GetCurrentDirectory() + "/../config/" + fileName
	if IsFileExists(curDirPath) {
		return curDirPath
	}
	if IsFileExists("C:/" + fileName) {
		return "C:/" + fileName
	}
	if IsFileExists("/etc/" + fileName) {
		return "/etc/" + fileName
	}
	return ""
}

/**
 * 将配置文件中的json格式，解析到结构体中
 */
func ParseConfigFile() error {
	confFilePath := findConfigFile("cv_api_gw_conf.json")
	if confFilePath == "" {
		return errors.New("can not find system config file")
	}
	fileContent, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return busierr.WrapGoError(err)
	}
	//fmt.Println(string(fileContent))
	err = json.Unmarshal(fileContent, &Configs)
	if err != nil {
		return busierr.WrapGoError(err)
	}

	ipFlag := flag.String("config_ip", "", "替换配置文件里的 {IP}")
	flag.Parse()
	// 替换模板里的 IP：优先使用 flag 里的 ip，供开发使用
	template := "{IP}"
	ip := "NOT_EXIST"
	if *ipFlag != "" {
		ip = *ipFlag
	} else {
		ip, err = ExternalIP()
		if err != nil {
			fmt.Sprint("get ip with err: ", err)
		}
	}
	fmt.Sprint("replace {IP} in config file with: ", ip)
	ReplaceFieldValueRecursively(Configs, template, ip)

	endpointConfigDirPath := strings.Replace(confFilePath, "cv_api_gw_conf.json", "", -1)
	endpointConfigDirPath = strings.Replace(endpointConfigDirPath, "config", "config_of_endpoint", -1)
	confFiles, err := GetAllFiles(endpointConfigDirPath)
	if err != nil {
		return busierr.WrapGoError(err)
	}

	backendAlias := make(map[string]*EndpointBackend)
	allEndpointConf := make(map[string]map[string]*ForwardEndpoint) //map[METHOD]map[url]Config
	for _, confFile := range confFiles {
		if strings.Contains(confFile, "endpoint") && strings.HasSuffix(confFile, ".json") {
			fileContent, err = ioutil.ReadFile(confFile)
			if err != nil {
				return busierr.WrapGoError(err)
			}
			var tmpConf []*ForwardEndpoint
			err = json.Unmarshal(fileContent, &tmpConf)
			if err != nil {
				return busierr.WrapGoError(err)
			}

			for _, endpointConf := range tmpConf {
				methodCategoryConf := allEndpointConf[endpointConf.Method]
				if methodCategoryConf == nil {
					methodCategoryConf = make(map[string]*ForwardEndpoint)
				}
				methodCategoryConf[endpointConf.Endpoint] = endpointConf
				allEndpointConf[endpointConf.Method] = methodCategoryConf
				//同步按照alias收集每个Backend配置，用于下面引用的修改
				for _, backendConf := range endpointConf.Backend {
					var backendHosts []string
					for _, backendHost := range backendConf.Host {
						if strings.Contains(backendHost, "{defaultBackend.") {
							backendHostKey := strings.Replace(backendHost, "{defaultBackend.", "", -1)
							backendHostKey = strings.Replace(backendHostKey, "}", "", -1)
							defaultHost, ok := Configs.System.DefaultBackend[backendHostKey]
							if ok {
								backendHosts = append(backendHosts, defaultHost)
							} else {
								backendHosts = append(backendHosts, backendHost)
							}
						} else {
							backendHosts = append(backendHosts, backendHost)
						}
					}
					backendConf.Host = backendHosts
					if backendConf.Alias != "" {
						backendAlias[backendConf.Alias] = backendConf
					}
				}
			}
		}
	}
	//将backend引用的别名（SameAs），替换为真正的配置
	for _, methodCategoryConf := range allEndpointConf {
		for _, forwardConf := range methodCategoryConf {
			for _, backendConf := range forwardConf.Backend {
				if backendConf.BodyFieldOperationSameAs != "" {
					referedForwardConf, ok := backendAlias[backendConf.BodyFieldOperationSameAs]
					if ok {
						//注意这里要用append方式，而不能直接赋值（因为有可能是混合配置）
						for _, referedBodyOperation := range referedForwardConf.BodyFieldOperation {
							backendConf.BodyFieldOperation = append(backendConf.BodyFieldOperation, referedBodyOperation)
						}
					}
				}
				if backendConf.ResponseFieldOperationSameAs != "" {
					referedForwardConf, ok := backendAlias[backendConf.ResponseFieldOperationSameAs]
					if ok {
						for _, referedRespOperation := range referedForwardConf.ResponseFieldOperation {
							backendConf.ResponseFieldOperation = append(backendConf.ResponseFieldOperation, referedRespOperation)
						}
					}
				}
			}
		}
	}

	Configs.Endpoints = allEndpointConf
	return nil
}

func GetApiForwardConf(url, method string) *ForwardEndpoint {
	//TODO：测试完毕后删除
	//ParseConfigFile()
	methodCategoryConf, ok := Configs.Endpoints[method]
	if !ok {
		return nil
	}
	for _, endpointConf := range methodCategoryConf {
		if isUrlMatch(endpointConf.Endpoint, url) {
			return endpointConf
		}
	}
	return nil
}

func isUrlMatch(route, url string) bool {
	routeItems := strings.Split(route, "/")
	urlItems := strings.Split(url, "/")
	if len(routeItems) != len(urlItems) {
		return false
	}
	for i := 0; i < len(routeItems); i++ {
		routeItem := routeItems[i]
		urlItem := urlItems[i]
		//如果URL Route的定义中有{xxx}的部分，则说明这部分是一个通配符，算作匹配
		if strings.HasPrefix(routeItem, "{") && strings.HasSuffix(routeItem, "}") {
			continue
		}
		//如果当前不是通配符，则按照
		if routeItem != urlItem {
			return false
		}
	}
	return true
}
