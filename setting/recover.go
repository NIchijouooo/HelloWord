package setting

import (
	"archive/zip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gateway/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type BackupFileTemplate struct {
	BackupId   int    `json:"backupId"`
	BackupTime string `json:"backupTime"`
}

func UnZipFiles(zipFile string, destDir string) error {

	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		ZAPS.Errorf("OpenReader err,%v", err)
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		func() error {
			fPath := filepath.Join(destDir, f.Name)
			//log.Println("fPath ", fPath)
			if f.FileInfo().IsDir() {
				_ = os.MkdirAll(fPath, os.ModePerm)
			} else {
				if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
					log.Println("mkdir err", err)
					return err
				}
				inFile, err := f.Open()
				if err != nil {
					log.Println("open err,", err)
					return err
				}
				defer inFile.Close()

				outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					log.Println("openFile err,", err)
					return err
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, inFile)
				if err != nil {
					log.Println("copy err,", err)
					return err
				}
			}
			return nil
		}()
	}

	return nil
}

func RecoverFiles(name string) error {

	fileName := "./tmp/" + name
	ZAPS.Debugf("fileName %v", fileName)

	err, _, _ := RunShellCmd("tar", "-xvf", fileName, "-C", "../")
	if err != nil {
		return err
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

//func RecoverFiles(name string) bool {
//
//	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
//
//	fileName := exeCurDir + "/selfpara/" + name
//	fileAbsoluteDir := exeCurDir + "/"
//	ZAPS.Debugf("fileName %v", fileName)
//	if err := UnZipFiles(fileName, fileAbsoluteDir); err != nil {
//		ZAPS.Errorf("err %v", err)
//		return false
//	}
//	err := os.Remove(fileName)
//	if err != nil {
//		log.Printf("removeFile err,%s\n", fileName)
//	}
//
//	return true
//}

func GetRecoverFileListFromRemote(sn string) (error, []BackupFileTemplate) {

	snCodeStr := sn

	unixStr := strconv.FormatInt(time.Now().UnixMilli(), 10)

	str := snCodeStr + unixStr + "rtCloud"

	sign := md5.Sum([]byte(str))
	signStr := fmt.Sprintf("%x", sign)

	urls := url.Values{}
	urls.Add("snCode", snCodeStr)
	urls.Add("time", unixStr)
	urls.Add("sign", signStr)
	urlsParam := "http://" + "cloud.reatgreen.com" + "/gateway/device-maintenance/deviceBackup/getBackUpFileList?" + urls.Encode()
	//ZAPS.Debugf("urlsParam %s", urlsParam)

	response, err := http.Get(urlsParam)
	if err != nil {
		ZAPS.Errorf("http Get获取失败 %v", err)
		return err, nil
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			ZAPS.Errorf("http Get关闭失败 %v", err)
		}
	}(response.Body)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ZAPS.Errorf("http Get获取Response失败 %v", err)
		return err, nil
	}
	ZAPS.Debugf("responseBody,%v", string(responseBody))

	fileInfo := &struct {
		Code int                  `json:"code"`
		Msg  string               `json:"Msg"`
		Data []BackupFileTemplate `json:"data"`
	}{}

	err = json.Unmarshal(responseBody, fileInfo)
	if err != nil {
		ZAPS.Errorf("unmarshall responseBody json err,%v", err)
		return err, nil
	}

	return nil, fileInfo.Data
}

func GetRecoverFileFromRemote(sn string, backupId int) error {
	snCodeStr := sn

	//BeiJingZone := time.FixedZone("CST", 8*3600)
	unixStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
	ZAPS.Debugf("unixStr %v", unixStr)
	str := snCodeStr + unixStr + "rtCloud"

	sign := md5.Sum([]byte(str))
	signStr := fmt.Sprintf("%x", sign)
	ZAPS.Debugf("signStr %v", signStr)

	urls := url.Values{}
	urls.Add("snCode", snCodeStr)
	urls.Add("time", unixStr)
	urls.Add("sign", signStr)
	urls.Add("backupId", strconv.Itoa(backupId))
	urlsParam := "http://" + "cloud.reatgreen.com" + "/gateway/device-maintenance/deviceBackup/downloadFile?" + urls.Encode()
	ZAPS.Debugf("urlsParam %s", urlsParam)

	response, err := http.Get(urlsParam)
	if err != nil {
		ZAPS.Errorf("http Get获取失败 %v", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			ZAPS.Errorf("http Get关闭失败 %v", err)
		}
	}(response.Body)

	utils.DirIsExist("./tmp")

	fileName := "./tmp/" + "backup.tar"
	f, _ := os.Create(fileName)
	_, err = io.Copy(f, response.Body)
	if err != nil {
		ZAPS.Errorf("copy file err,%v", err)
		return err
	}

	_, _, _ = RunShellCmd("tar", "-xvf", fileName, "../")

	return nil
}
