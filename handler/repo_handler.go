package handler

import (
	"cocoapods-cache-proxy-server/config"
	"cocoapods-cache-proxy-server/model"
	"cocoapods-cache-proxy-server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func ReposIndexHandler(authConfig *config.AuthorizationConfig, cacheDir string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params model.ReposParams
		if err := ctx.ShouldBindUri(&params); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
			fmt.Println(err)
			return
		}
		if err := ctx.ShouldBindQuery(&params); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
			fmt.Println(err)
			return
		}
		fmt.Println(params)
		cache_zip_file_path := filepath.Join(cacheDir, params.Repo, params.Name, params.Name+"-"+params.Tag+".zip")
		if util.FileExists(cache_zip_file_path) {
			RedirectToCacheFile(ctx, params, cacheDir)
			return
		}
		DownloadGitHandler(ctx, params, cacheDir)
	}
}
