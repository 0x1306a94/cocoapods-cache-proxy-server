package handler

import (
	"cocoapods-cache-proxy-server/model"
	"cocoapods-cache-proxy-server/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func RedirectToCacheFile(ctx *gin.Context, params model.ReposParams, cacheDir string) {
	cache_zip_file_path := filepath.Join(cacheDir, params.Repo, params.Name, params.Name+"-"+params.Tag+".zip")
	if !util.FileExists(cache_zip_file_path) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "The file does not exist",
		})
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, filepath.Join("/cache", params.Repo, params.Name, params.Name+"-"+params.Tag+".zip"))
}
