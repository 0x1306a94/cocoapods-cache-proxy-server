package util

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func ZipCompressDir(dir, zipFile string) error {

	fz, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 忽略 git
		if strings.Contains(path, ".git") {
			return nil
		}
		if !info.IsDir() {
			dst, err := w.Create(path[len(dir)+1:])
			if err != nil {
				return err
			}
			src, err := os.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()
			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func ZipDeCompress(zipFile, dst string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		if strings.HasSuffix(file.Name, "/") {
			continue
		}
		//split := strings.Split(file.Name, "/")
		//s := split[len(split) - 3:]
		filename := filepath.Join(dst, file.Name)
		dir := filepath.Dir(filename)
		fmt.Println(file.Name)
		//if strings.HasSuffix(file.Name, "/") {
		//	dir = filepath.Join(dst,file.Name)
		//}
		err = os.MkdirAll(dir, 0755)

		if err != nil && os.IsNotExist(err) {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		rc.Close()
		w.Close()
	}
	return nil
}

func TarGzDir(dir, tarFile string) error {
	d, err := os.Create(tarFile)
	if err != nil {
		return nil
	}
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 忽略 git
		if strings.Contains(path, ".git") {
			return nil
		}
		if !info.IsDir() {
			var link string
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(path); err != nil {
					return err
				}
			}
			header, err := tar.FileInfoHeader(info, link)
			if err != nil {
				return err
			}
			header.Name = path[len(dir)+1:]
			err = tw.WriteHeader(header)
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() { //nothing more to do for non-regular
				return nil
			}
			src, err := os.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()
			_, err = io.Copy(tw, src)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
