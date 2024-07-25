/**
 * HTTPClient工具
 * @author duhaifeng
 * @date   2021/04/15
 */
package common

import (
	"cv-api-gw/common/busierr"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/**
 * Http请求通用方法
 * @param url: 请求url，需要有http://, method: Get,Post等, header: 请求的header, form: 请求的form, retry：失败时重试的次数，为0时不重试
 * @return 统一的apiResponse返回信息 *ApiResponse, 错误信息 error
 */
func RequestHttp(method, url string, header map[string]string, body string) (int, []byte, error) {
	//如果url中包含json字样，说明此时是测试请求，则读取本地配置信息
	if Configs != nil && Configs.System.FakeBackend {
		return getLocalFakeResponse(url, method)
	}
	Log.Debug("request api: <%s> %s", method, url)
	var req *http.Request
	var err error
	client := &http.Client{}
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		Log.Error("create api request failed: %s", err.Error())
		return http.StatusInternalServerError, nil, busierr.WrapGoError(err)
	}
	req.Body = ioutil.NopCloser(strings.NewReader(body))
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Type", "application/json")
	//设置header，增加认证信息
	for k, v := range header {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, busierr.WrapGoError(err)
	}
	defer resp.Body.Close()
	//if resp.StatusCode != http.StatusOK {
	//	return nil, fmt.Errorf("response code from %s is not succeed: %d", url, resp.StatusCode)
	//}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Error("read api body data failed: %s", err.Error())
		return resp.StatusCode, nil, busierr.WrapGoError(err)
	}
	return resp.StatusCode, data, nil
}

func getLocalFakeResponse(url, method string) (int, []byte, error) {
	url = strings.Replace(url, "http://", "", -1)
	url = strings.Replace(url, "https://", "", -1)
	url = "/" + strings.Join(strings.Split(url, "/")[1:], "/")
	jsonFakeFiles, err := GetAllFiles(Configs.System.FakeBackendFilePath)
	if err != nil {
		return http.StatusInternalServerError, nil, busierr.WrapGoError(err)
	}
	for _, jsonFakeFile := range jsonFakeFiles {
		jsonFakeBytes, err := readLocalJsonFile(jsonFakeFile)
		if err != nil {
			continue
		}
		jsonFakeMap := make(map[string]interface{})
		err = json.Unmarshal(jsonFakeBytes, &jsonFakeMap)
		if err != nil {
			continue
		}
		fakeUrl, ok := jsonFakeMap["_url"]
		if !ok {
			continue
		}
		fakeMethod, ok := jsonFakeMap["_method"]
		if !ok {
			fakeMethod = "GET"
		}
		if strings.ToLower(url) == strings.ToLower(fakeUrl.(string)) && method == fakeMethod {
			delete(jsonFakeMap, "_url")
			delete(jsonFakeMap, "_method")
			jsonFakeBytes, _ = json.MarshalIndent(jsonFakeMap, "", "  ")
			return http.StatusOK, jsonFakeBytes, nil
		}
	}
	return http.StatusNotFound, nil, nil
}

/**
 * 读取本地JSON文件中定义的假数据，本方法用于支持本地离线测试
 */
func readLocalJsonFile(fileName string) ([]byte, error) {
	//fileName = strings.Replace(fileName, "/", ".", -1)
	//fileName = strings.TrimPrefix(fileName, ".")
	//fileName = common.Configs.System.FakeBackendFilePath + "/" + fileName + ".json"
	//common.Log.Debug("request local api: file name: %s", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		Log.Error("read local json data file<%s> failed: %s", fileName, err.Error())
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadFile(fileName)
}
