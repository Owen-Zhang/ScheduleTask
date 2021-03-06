package system

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/mholt/archiver"
	uuid "github.com/satori/go.uuid"

	"os/exec"
)

//FileExist 验证文件是否存在
func FileExist(filename string) bool {
	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

//DirExist 验证目录是否存在
func DirExist(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

//FileCopy 复制文件
func FileCopy(source string, dest string) (int64, error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return 0, err
	}

	defer sourcefile.Close()
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return 0, err
	}

	destfile, err := os.Create(dest)
	if err != nil {
		return 0, err
	}

	w, err := io.Copy(destfile, sourcefile)
	if err != nil {
		destfile.Close()
		return 0, err
	}

	if err := os.Chmod(dest, sourceinfo.Mode()); err != nil {
		destfile.Close()
		return 0, err
	}
	destfile.Close()
	return w, nil
}

//DirectoryCopy 复制目录
func DirectoryCopy(source string, dest string) error {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, sourceinfo.Mode()); err != nil {
		return err
	}

	directory, err := os.Open(source)
	if err != nil {
		return err
	}

	defer directory.Close()
	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			if err := DirectoryCopy(sourcefilepointer, destinationfilepointer); err != nil {
				return err
			}
		} else {
			if _, err := FileCopy(sourcefilepointer, destinationfilepointer); err != nil {
				return err
			}
		}
	}
	return nil
}

//FileNameWithExt 获取文件名带后缀
func FileNameWithExt(filepath string) string {
	if filepath == "" {
		return ""
	}
	return path.Base(filepath)
}

//Ext 获取文件后缀
func Ext(filepath string) string {
	fileName := FileNameWithExt(filepath)
	if fileName == "" {
		return ""
	}
	return path.Ext(fileName)
}

//FileName 获取文件名,不包含后缀
func FileName(path string) string {
	fileName := FileNameWithExt(path)
	if fileName == "" {
		return ""
	}
	return strings.TrimSuffix(fileName, Ext(path))
}

//CreateUuidFile 生成UUID的文件名
func CreateUuidFile(filepath string) string {
	ext := Ext(filepath)
	if ext == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", uuid.NewV4().String(), ext)
}

//GetUuid 生成UUID字符串
func GetUuid() string {
	return uuid.NewV4().String()
}

//CheckFileExt 判断文件类型是否为想要类型
func CheckFileExt(exts []string, filepath string) bool {
	if len(exts) == 0 {
		return true
	}
	ext := strings.TrimLeft(Ext(filepath), ".")

	for _, value := range exts {
		if ext == value {
			return true
		}
	}
	return false
}

//IsExist 判断文件是否存在
func IsExist(filepath string) bool {
	if filepath == "" {
		return false
	}
	_, err := os.Stat(filepath)
	return err == nil || os.IsExist(err)
}

//UrlFileName 找出url中的文件名，如http://www.baidu.com/aaa/12.zip?name=sdfasd 要取出12.zip文件名
func UrlFileName(url string) string {
	array := strings.Split(url, "/")
	length := len(array)
	if length <= 0 {
		return ""
	}

	filename := array[length-1]
	arrry2 := strings.Split(filename, "?")

	return arrry2[0]
}

//UnzipFile 解压
func UnzipFile(filePath, rundatafolder string) error {
	if filePath == "" {
		return errors.New("filePath is empty")
	}
	if !IsExist(rundatafolder) {
		if err := os.MkdirAll(rundatafolder, 0777); err != nil {
			return err
		}
	}
	if err := archiver.NewZip().Unarchive(filePath, rundatafolder); err != nil {
		return err
	}

	return nil
}

//GetCurrentPath 获取目录的根目录如：D:/ddd/
func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Println(err.Error())
	}
	s = strings.Replace(s, "\\", "/", -1)
	s = strings.Replace(s, "\\\\", "/", -1)
	i := strings.LastIndex(s, "/")
	path := string(s[0:i])
	return path
}
