package handler

import (
	"cocoapods-cache-proxy-server/model"
	"cocoapods-cache-proxy-server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DownloadGitHandler(ctx *gin.Context, params model.ReposParams, cacheDir string) {

	dir, err := ioutil.TempDir("", "cocoapods-proxy-temp"+params.Name+fmt.Sprint(time.Now().Unix()))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		fmt.Println(err)
		return
	}
	savePath := filepath.Join(dir, params.Name)
	defer func() {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			os.RemoveAll(savePath)
		}()
	}()

	rep, err := git.PlainClone(savePath, false, &git.CloneOptions{
		URL:      params.Git,
		Progress: os.Stdout,
		Depth:    1,
		//NoCheckout:   true,
		SingleBranch: true,
		RemoteName:   params.Tag,
		Tags:         git.NoTags,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	work, err := rep.Worktree()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	if params.Submodules {
		sm, err := work.Submodules()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 0,
				"msg":  err.Error(),
			})
			fmt.Println(err)
			return
		}
		for _, v := range sm {
			options := &git.SubmoduleUpdateOptions{
				Init: true,
			}
			if err := v.Update(options); err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code": 0,
					"msg":  err.Error(),
				})
				fmt.Println(err)
				return
			}
		}
	}

	// zip
	tgz_file_path := filepath.Join(dir, params.Name+"-"+params.Tag+".tgz")
	if err := util.TarGzDir(savePath, tgz_file_path); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		fmt.Println(err)
		return
	}
	cache_tgz_file_path := filepath.Join(cacheDir, "repos", params.Name, params.Name+"-"+params.Tag+".tgz")
	cache_base_path := filepath.Join(cacheDir, "repos", params.Name)
	fmt.Println("cache_base_path", cache_base_path)
	if err := os.MkdirAll(cache_base_path, os.ModePerm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		fmt.Println(err)
		return
	}
	if err := os.Rename(tgz_file_path, cache_tgz_file_path); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		fmt.Println(err)
		return
	}
	RedirectToCacheFile(ctx, params, cacheDir)
}
