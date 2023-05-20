package setting

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/utils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func BackupFiles(names []string) (error, string) {
	utils.DirIsExist("./backup")

	//判断plugin是否要备份
	for _, name := range names {
		if name == "plugin" {
			_, _, _ = RunShellCmd("cp", "-r", "./plugin/", "./backup/plugin")
		}
	}

	//遍历selfpara文件夹
	fileMap, err := ioutil.ReadDir("./selfpara")
	if err != nil {
		ZAPS.Errorf("打开selfpara目录失败 %v", err)
		return err, ""
	}
	utils.DirIsExist("./backup/selfpara")

	//判断通信接口、采集接口、上报接口是否要备份
	for _, name := range names {
		for _, v := range fileMap {
			if strings.Contains(v.Name(), ".json") && strings.Contains(v.Name(), name) {
				dstFile := "./backup/selfpara/" + v.Name()
				srcFile := "./selfpara/" + v.Name()

				dFile, _ := os.Create(dstFile)
				sFile, _ := os.Open(srcFile)
				_, err = io.Copy(dFile, sFile)
				if err != nil {
					ZAPS.Debugf("拷贝文件错误 %v", err)
				}
			}
		}
	}

	zipFileName := "./tmp/backup" + SystemState.SN + ".tar"
	_, _, _ = RunShellCmd("tar", "cvf", zipFileName, "./backup")

	_, _, _ = RunShellCmd("rm", "-r", "./backup/")

	return nil, zipFileName
}

//func BackupFiles() (error, string) {
//	utils.DirIsExist("./tmp")
//	dirMap := []string{"./plugin", "./selfpara"}
//
//	zipFileName := "./tmp/backup_" + SystemState.SN + ".zip"
//
//	_ = utils.CompressDirsToZip(dirMap, zipFileName)
//
//	return nil, zipFileName
//}

func BackupFilesToRemote() (error, string) {
	utils.DirIsExist("./tmp")
	//dirMap := []string{"./plugin", "./selfpara"}

	zipFileName := "./tmp/backup_" + SystemState.SN + ".tar"

	//_ = utils.CompressDirsToZip(dirMap, zipFileName)

	RunShellCmd("tar", "cvf", zipFileName, "./plugin", "./selfpara")

	unixStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
	ZAPS.Debugf("unixStr %v", unixStr)
	str := SystemState.SN + unixStr + "rtCloud"

	sign := md5.Sum([]byte(str))
	signStr := fmt.Sprintf("%x", sign)
	ZAPS.Debugf("signStr %v", signStr)

	urls := url.Values{}
	urls.Add("snCode", SystemState.SN)
	urls.Add("time", unixStr)
	urls.Add("sign", signStr)
	//urlsParam := "http://" + "cloud.reatgreen.com" + "/gateway/deviceBackup/upload?" + urls.Encode()
	urlsParam := "http://" + "cloud.reatgreen.com" + "/gateway/device-maintenance/deviceBackup/upload?"
	ZAPS.Debugf("urlsParam %s", urlsParam)

	//创建一个缓冲区对象,后面的要上传的body都存在这个缓冲区里
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	err := bodyWriter.WriteField("snCode", SystemState.SN)
	if err != nil {
		ZAPS.Errorf("http WriteField err,%v", err)
	}
	err = bodyWriter.WriteField("time", unixStr)
	if err != nil {
		ZAPS.Errorf("http WriteField err,%v", err)
	}
	err = bodyWriter.WriteField("sign", signStr)
	if err != nil {
		ZAPS.Errorf("http WriteField err,%v", err)
	}
	//创建需要上传的文件
	fileWriter, err := bodyWriter.CreateFormFile("file", filepath.Base(zipFileName))
	if err != nil {
		ZAPS.Errorf("CreateFormFile err,%v", err)
		return err, zipFileName
	}
	//打开文件
	fd, err := os.Open(zipFileName)
	if err != nil {
		ZAPS.Errorf("open file err,%v", err)
		return err, zipFileName
	}
	defer fd.Close()
	_ = os.Remove(zipFileName)
	//把文件写入到缓冲区里
	_, err = io.Copy(fileWriter, fd)
	if err != nil {
		ZAPS.Errorf("copy file err,%v", err)
		return err, ""
	}
	_ = bodyWriter.Close()
	//ZAPS.Debugf("http Post body %+v", bodyBuf)
	response, err := http.Post(urlsParam, bodyWriter.FormDataContentType(), bodyBuf)
	if err != nil {
		ZAPS.Errorf("http Post err,%v", err)
		return err, ""
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(response.Body)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ZAPS.Errorf("read httpResponse err,%v", err)
		return err, ""
	}
	ZAPS.Debugf("responseBody,%v", string(responseBody))

	res := struct {
		Code  int         `json:"code"`
		Msg   string      `json:"msg"`
		Total int         `json:"total"`
		Data  interface{} `json:"data"`
	}{}

	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		ZAPS.Errorf("responseBody Json格式化错误 %v", err)
		return err, ""
	}
	if res.Code == 0 {
		return nil, zipFileName
	} else {
		return errors.New("返回错误"), ""
	}

}
