package tool

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ZipDirectory(srcPath string) ([]byte, error) {
	// 压缩文件
	zipbuf := &bytes.Buffer{}
	zipWriter := zip.NewWriter(zipbuf)
	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// root目录直接跳过
		if path == srcPath {
			return nil
		}
		// 替换为相对路径
		path = strings.TrimPrefix(path, fmt.Sprintf("%v/", srcPath))
		// 添加文件
		if info.IsDir() {
			_, err := zipWriter.Create(fmt.Sprintf("%v/", path))
			return err
		} else {
			zipFile, err := zipWriter.Create(path)
			if err != nil {
				return err
			}

			binary, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", srcPath, path))
			if err != nil {
				return err
			}

			if _, err := zipFile.Write(binary); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return zipbuf.Bytes(), nil
}
