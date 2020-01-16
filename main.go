package main

import (
	"cocoapods-cache-proxy-server/config"
	"cocoapods-cache-proxy-server/handler"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
)

var (
	authConfig  *config.AuthorizationConfig = &config.AuthorizationConfig{}
	lauchConfig *LauchConfig                = &LauchConfig{}
	_version_ = ""
	_commit_ = ""
)

type LauchConfig struct {
	User     string
	Password string
	Port     int64
	Verbose bool
	CacheDir string
}

func parseLauchConfig(cnf *LauchConfig) {
	flag.StringVar(&cnf.User, "user", "", "http basic auth user name, 3-10 characters, or environment variable COCOAPODS_CACHE_PROXY_USER")
	flag.StringVar(&cnf.Password, "password", "", "http basic auth password, 6-10 characters, or environment variable COCOAPODS_CACHE_PROXY_PASSWORD")
	flag.Int64Var(&cnf.Port, "port", 9898, "监听端口 1024 ~ 65535 之间")
	flag.StringVar(&cnf.CacheDir, "dir", "./repo/cache", "缓存文件存放目录, or environment variable COCOAPODS_CACHE_PROXY_CACHE_DIR")
	flag.BoolVar(&cnf.Verbose, "verbose", false, "是否开启请求日志")
	version := false
	flag.BoolVar(&version, "version", false, "显示版本信息")
	flag.Parse()
	if version {
		fmt.Println("version: ", _version_)
		fmt.Println("commit: ", _commit_)
		os.Exit(0)
	}
	if len(cnf.User) == 0 {
		cnf.User = os.Getenv("COCOAPODS_CACHE_PROXY_USER")
	}

	if len(cnf.Password) == 0 {
		cnf.Password = os.Getenv("COCOAPODS_CACHE_PROXY_PASSWORD")
	}

	if len(cnf.User) < 3 || len(cnf.User) > 10 {
		flag.Usage()
		os.Exit(-1)
	}

	if len(cnf.Password) < 6 || len(cnf.Password) > 10 {
		flag.Usage()
		os.Exit(-1)
	}

	if cnf.Port < 1024 || cnf.Port > 65535 {
		flag.Usage()
		os.Exit(-1)
	}

	if len(cnf.CacheDir) == 0 {
		cnf.CacheDir = os.Getenv("COCOAPODS_CACHE_PROXY_CACHE_DIR")
	}

	if len(cnf.CacheDir) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	path, err := filepath.Abs(cnf.CacheDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("cache dir:", path)
	cnf.CacheDir = path
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	_, err = os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = os.Chmod(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	parseLauchConfig(lauchConfig)
	tips := `
如何使用: 
	1. 请先安装 cocoapods 插件 https://github.com/0x1306a94/cocoapods-cache-proxy
	2. 安装完插件后 执行 pod cache proxy add NAME http://domain/cocoapods/proxy/repos USER PASSWORD (USER PASSWORD 为 http basic auth user and password)
	3. 在Podfile 中 添加 plugin "cocoapods-cache-proxy", :proxy => "NAME" NAME 为第二步中的 NAME
	4. 然后直接 pod install 就行
	`
	fmt.Println(tips)
	if lauchConfig.Verbose {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	authConfig.SetupUser(lauchConfig.User, lauchConfig.Password)

	router := setupRouter()

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", lauchConfig.Port),
		Handler: router,
	}
	s.ListenAndServe()
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/cache", lauchConfig.CacheDir).Use(handler.BasicAuthMiddleware(authConfig))
	reposRouter := router.Group("/cocoapods/proxy/repos")
	reposRouter.Use(handler.BasicAuthMiddleware(authConfig))
	//reposRouter.Use(func(ctx *gin.Context) {
	//	ctx.Set("authConfig", authConfig)
	//	ctx.Next()
	//})
	reposRouter.GET("/:name", handler.ReposIndexHandler(authConfig, lauchConfig.CacheDir))
	return router
}
