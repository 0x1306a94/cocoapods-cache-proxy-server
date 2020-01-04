package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func ZipDir(dir, zipFile string) bool {

	fz, err := os.Create(zipFile)
	if err != nil {
		fmt.Println("Create zip file failed: ", err)
		return false
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 忽略 git
		if strings.Contains(path, ".git") {
			return nil
		}
		if !info.IsDir() {
			dst, err := w.Create(path[len(dir)+1:])
			if err != nil {
				fmt.Println("Create failed:", err)
				return nil
			}
			src, err := os.Open(path)
			if err != nil {
				fmt.Println("Open failed: ", err)
				return nil
			}
			defer src.Close()
			_, err = io.Copy(dst, src)
			if err != nil {
				fmt.Println("Copy failed: %s", err)
				return nil
			}
		}
		return nil
	})
	return true
}
