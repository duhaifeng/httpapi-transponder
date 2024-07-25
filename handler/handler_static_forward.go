/**
 * 转发静态资源的的handler，用于为web提供网络服务
 * @author duhaifeng
 * @date   2021/07/05
 */
package handler

import (
	"cv-api-gw/common"
	"cv-api-gw/common/busierr"
	"io/ioutil"
	"path"
	"strings"

	api "github.com/duhaifeng/simpleapi"
)

type StaticForwardHandler struct {
	GatewayHandlerBase
}

func (this *StaticForwardHandler) HandleRequest(r *api.Request) (interface{}, error) {
	this.DoNotWrapResponse()
	url := r.GetUrl().String()
	url = strings.Split(url, "?")[0]

	filePath := r.GetUrl().String()[strings.Index(url, "/static/"):]
	filePath = strings.Replace(filePath, "/static/", "", -1)
	filePath = common.Configs.System.WebRoot + filePath
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, busierr.WrapGoError(err)
	}
	fileExt := path.Ext(filePath)
	if fileExt == ".html" || fileExt == ".htm" {
		this.GetResponse().SetHeader("Content-Type", "text/html")
	} else if fileExt == ".css" {
		this.GetResponse().SetHeader("Content-Type", "text/css")
	} else if fileExt == ".js" {
		this.GetResponse().SetHeader("Content-Type", "application/x-javascript")
	} else if fileExt == ".svg" {
		this.GetResponse().SetHeader("Content-Type", "text/xml")
	} else if fileExt == ".jpg" || fileExt == ".jpeg" {
		this.GetResponse().SetHeader("Content-Type", "image/jpeg")
	} else if fileExt == ".png" {
		this.GetResponse().SetHeader("Content-Type", "image/png")
	} else if fileExt == ".gif" {
		this.GetResponse().SetHeader("Content-Type", "image/gif")
	}
	_, err = this.GetResponse().Write(fileBytes)
	return nil, err
}
