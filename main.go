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
	user     = ""
	password = ""
)

var (
	authConfig  *config.AuthorizationConfig = &config.AuthorizationConfig{}
	lauchConfig *LauchConfig                = &LauchConfig{}
)

type LauchConfig struct {
	AdminUser     string
	AdminPassword string
	Port          int64
	BaseURL       string
	CacheDir      string
}

func parseLauchConfig(cnf *LauchConfig) {
	flag.StringVar(&cnf.AdminUser, "user", "admin", "admin user name, 3-10 characters")
	flag.StringVar(&cnf.AdminPassword, "password", "123456dd", "admin user password, 6-10 characters")
	flag.Int64Var(&cnf.Port, "port", 9898, "监听端口 1024 ~ 65535 之间")
	flag.StringVar(&cnf.BaseURL, "baseurl", "", "base url 默认为 http://127.0.0.1:port")
	flag.StringVar(&cnf.CacheDir, "dir", "./repo/cache", "缓存文件存放目录")
	flag.Parse()

	if len(cnf.AdminUser) < 3 || len(cnf.AdminUser) > 10 {
		flag.Usage()
		os.Exit(-1)
	}

	if len(cnf.AdminPassword) < 6 || len(cnf.AdminPassword) > 10 {
		flag.Usage()
		os.Exit(-1)
	}

	if cnf.Port < 1024 || cnf.Port > 65535 {
		flag.Usage()
		os.Exit(-1)
	}

	if len(cnf.BaseURL) == 0 {
		cnf.BaseURL = fmt.Sprintf("http://127.0.0.1:%v", cnf.Port)
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
	authConfig.SetupAdminUser(lauchConfig.AdminUser, lauchConfig.AdminPassword)

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
	reposRouter.GET("/:repo/:name", handler.ReposIndexHandler(authConfig, lauchConfig.CacheDir))
	return router
}
