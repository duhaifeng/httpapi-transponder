/**
 * 文件操作工具
 * @author duhaifeng
 * @date   2021/04/15
 */
package common

import (
	"cv-api-gw/common/busierr"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 判断所给路径文件/文件夹是否存在
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

/**
 * 获取程序运行路径
 */
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "/", "/", -1)
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, busierr.WrapGoError(err)
	}
	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, file := range dir {
		if file.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+file.Name())
			GetAllFiles(dirPth + PthSep + file.Name())
		} else {
			files = append(files, dirPth+PthSep+file.Name())
		}
	}

	// 读取子目录下文件
	for _, subDir := range dirs {
		subDirFiles, _ := GetAllFiles(subDir)
		for _, subDirFile := range subDirFiles {
			files = append(files, subDirFile)
		}
	}
	return files, nil
}
