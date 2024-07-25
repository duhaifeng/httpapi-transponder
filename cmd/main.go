/**
 * 系统启动入口，负责静态Handler路由、Web静态文件路由、以及动态配置路由配置
 * @author duhaifeng
 * @date   2021/04/14
 */
package main

import (
	"cv-api-gw/common"
	"cv-api-gw/handler"
	"cv-api-gw/interceptor"
	"cv-api-gw/service"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	api "github.com/duhaifeng/simpleapi"
)

const BANNER = `
                _____    _____         _____  __          __
        /\     |  __ \  |_   _|       / ____| \ \        / /
       /  \    | |__| |   | |        | |  __   \ \  /\  / /
      / /\ \   |  ___/    | |        | | |_ |   \ \/  \/ /
     / ____ \  | |       _| |_       | |__| |    \  /\  /
    /_/    \_\ |_|      |_____|       \_____|     \/  \/
                                                                            
`

var (
	BuildTime     string
	BuildBranch   string
	BuildCommitID string
)

func printBuildInfo() {
	fmt.Printf("build branch:\t%s\n", BuildBranch)
	fmt.Printf("build time:\t%s\n", BuildTime)
	fmt.Printf("build commit:\t%s\n", BuildCommitID)
	fmt.Printf("build version:\t%s\n", common.Configs.System.Version)
}

type NotFoundHandler struct{}

func (this *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "Application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"code\":404, \"message\":\"request route not exist\"}\n"))
}

//为Web静态文件自动生成路由
func findWebFileRoutes() []string {
	var webFileRoutes []string
	if common.Configs.System.WebRoot == "" {
		return nil
	}
	if common.Configs.System.WebRoot == "." {
		common.Configs.System.WebRoot = common.GetCurrentDirectory() + "/static/"
	}
	webFiles, err := common.GetAllFiles(common.Configs.System.WebRoot)
	if err != nil {
		common.Log.Error("load web file failed: %s", err.Error())
		return nil
	}
	for _, webFile := range webFiles {
		webFileRoute := strings.Replace(webFile, common.Configs.System.WebRoot, "/static/", -1)
		webFileRoute = strings.Replace(webFileRoute, "\\", "/", -1)
		webFileRoute = strings.Replace(webFileRoute, "//", "/", -1)
		webFileRoutes = append(webFileRoutes, webFileRoute)
	}
	return webFileRoutes
}

func StartApiServer() {
	//注册配置的自动转发接口路由
	server := new(api.ApiServer)
	server.Init()
	server.AllowCrossDomainRequest(true)
	server.SetNotFoundHandler(&NotFoundHandler{})

	for _, methodCategoryConfigs := range common.Configs.Endpoints {
		for _, endpointConf := range methodCategoryConfigs {
			server.RegisterHandler(endpointConf.Method, endpointConf.Endpoint, handler.AutoForwardHandler{})
		}
	}

	webFileRoutes := findWebFileRoutes()
	for _, webFileRoute := range webFileRoutes {
		server.RegisterHandler(http.MethodGet, webFileRoute, handler.StaticForwardHandler{})
	}

	server.RegisterHandler(http.MethodPost, "/v1/user", handler.AddUserHandler{})               //注册一个新用户，其中admin用户为系统内置，不可新增或删除
	server.RegisterHandler(http.MethodGet, "/v1/user/{userID}", handler.UserDetailHandler{})    //获取指定用户的明细信息
	server.RegisterHandler(http.MethodGet, "/v1/users", handler.UserListHandler{})              //获取系统内所有用户信息，只允许admin用户调用本接口
	server.RegisterHandler(http.MethodDelete, "/v1/user/{userID}", handler.UserDeleteHandler{}) //删除指定ID对应的用户
	server.RegisterHandler(http.MethodPut, "/v1/user/{userID}", handler.UpdateUserHandler{})    //更新指定ID对应的用户信息

	server.RegisterHandler(http.MethodPost, "/v1/token", handler.GetBoxTokenHandler{})
	server.RegisterHandler(http.MethodPost, "/v1/token/{userID}", handler.CreateTokenHandler{})
	server.RegisterHandler(http.MethodGet, "/v1/token/{userID}", handler.TokenDetailHandler{})
	server.RegisterHandler(http.MethodDelete, "/v1/token/{userID}", handler.TokenDeleteHandler{})
	server.RegisterHandler(http.MethodPost, "/v1/aksk/{userID}", handler.CreateAkskHandler{})
	server.RegisterHandler(http.MethodGet, "/v1/aksk/{userID}", handler.AkskDetailHandler{})
	server.RegisterHandler(http.MethodDelete, "/v1/aksk/{userID}", handler.AkskDeleteHandler{})

	server.RegisterHandler(http.MethodGet, "/v1/health", handler.HealthCheckHandler{})

	//注册拦截器（目前拦截器只对StructHandler有效）
	server.RegisterInterceptor(new(interceptor.FormatResponseInterceptor))
	server.RegisterInterceptor(new(interceptor.RequestIdInterceptor))
	server.RegisterInterceptor(new(interceptor.ExceptionInterceptor))
	server.RegisterInterceptor(new(interceptor.ValidateAccessAuthInterceptor))
	server.RegisterInterceptor(new(interceptor.EnvValuesInterceptor))

	//设置限流沙漏为每秒1000QPS
	server.GetTokenFunnel().SetDefaultTokenQuota(1000)
	//打开数据库链接池
	err := server.OpenMySQLOrmConn(common.Configs.Mysql.Host, common.Configs.Mysql.Port, common.Configs.Mysql.User, common.Configs.Mysql.Password, common.Configs.Mysql.Schema)
	if err != nil {
		common.Log.Error("!!! open mysql failed: %s !!!", err.Error())
		//return
	}
	server.SetLogger(common.Log.GetOriginalLogger())
	server.StartListen(common.Configs.System.ListenAddr, common.Configs.System.ListenPort)
}

func main() {
	fmt.Print(BANNER)
	err := common.ParseConfigFile()
	if err != nil {
		fmt.Printf("parse config error: %s\n", err.Error())
		os.Exit(1)
	}
	common.InitLog()
	printBuildInfo()
	fmt.Println()
	common.Log.Debug("config file parse finish.")
	service.AccessAuthManager.Init()
	StartApiServer()
	//停留一秒钟确保日志全部输出
	time.Sleep(time.Second)
}
