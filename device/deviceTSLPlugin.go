package device

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"gateway/setting"
	"gateway/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type TSLPluginTemplate struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

var tslPlugins = make(map[string]TSLPluginTemplate)

func ImportTSLPlugin(name string) error {

	fileFullName := "./plugin/" + name
	fileName := strings.TrimSuffix(name, ".zip")
	fileDir := "./plugin/" + fileName

	err := unZip(fileFullName, fileDir)
	if err != nil {
		return errors.New("解压" + fileFullName + "失败")
	}
	err = os.Remove(fileFullName)
	if err != nil {
		setting.ZAPS.Errorf("删除%s失败 %v", fileFullName, err)
	}

	//打开json
	fileJsonName := "./plugin/" + fileName + "/" + fileName + ".json"
	jsonBytes, err := utils.FileRead(fileJsonName)
	if err != nil {
		setting.ZAPS.Errorf("打开%s文件失败 %v", fileFullName, err)
		return errors.New("打开" + fileJsonName + "失败")
	}

	//解析json
	pluginJson := TSLPluginTemplate{}
	err = json.Unmarshal(jsonBytes, &pluginJson)
	if err != nil {
		return errors.New("格式化" + fileJsonName + "失败")
	}

	tslPlugins[pluginJson.Name] = pluginJson

	return nil
}

func unZip(zipFile string, destDir string) error {

	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		setting.ZAPS.Errorf("打开压缩包头失败 %v", err)
		return err
	}
	defer zipReader.Close()

	for k, f := range zipReader.File {
		setting.ZAPS.Infof("压缩包内文件[%v][%v]", k, f.Name)
		if f.FileInfo().IsDir() {
			setting.ZAPS.Debugf("压缩包内文件[%v][%v]是路径", k, f.Name)
			continue
		} else {
			inFile, err := f.Open()
			if err != nil {
				setting.ZAPS.Errorf("压缩包内文件[%v][%v]打开失败 %v", k, f.Name, err)
				continue
			}

			outFileName := filepath.Base(f.Name)

			outFile, err := os.OpenFile(destDir+"/"+outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				setting.ZAPS.Errorf("压缩包内文件[%v][%v]创建输出文件失败%v", k, f.Name, err)
				inFile.Close()
				continue
			}
			setting.ZAPS.Infof("压缩包内文件[%v][%v]创建输出文件[%v]成功", k, f.Name, outFileName)

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				setting.ZAPS.Errorf("压缩包内文件[%v][%v]拷贝文件失败%v", k, f.Name, err)
				inFile.Close()
				continue
			}

			inFile.Close()
			outFile.Close()
		}
	}

	return nil
}

func ExportTSLPlugin(pluginName string) (bool, string) {

	//遍历文件
	pluginPath := "./plugin/" + pluginName
	fileNameMap := make([]string, 0)
	fileNameMap, _ = TraverseTSLPlugin(pluginPath, fileNameMap)

	utils.DirIsExist("./tmp")
	_ = utils.CompressFilesToZip(fileNameMap, "./tmp/"+pluginName+".zip")

	return true, "./tmp/" + pluginName + ".zip"
}

//遍历plugin
func TraverseTSLPlugin(path string, fileName []string) ([]string, error) {

	rd, err := ioutil.ReadDir(path)
	if err != nil {
		setting.ZAPS.Errorf("打开目录实错误 %v", err)
		return fileName, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := path + "/" + fi.Name()
			fileName, _ = DeviceTSLTraversePlugin(fullDir, fileName)
		} else {
			fullName := path + "/" + fi.Name()
			if strings.Contains(fi.Name(), ".json") {
				fileName = append(fileName, fullName)
			} else if strings.Contains(fi.Name(), ".lua") {
				fileName = append(fileName, fullName)
			}
		}
	}

	return fileName, nil
}

func GetTSLPluginParam(name string) (TSLPluginTemplate, error) {

	plugin := TSLPluginTemplate{}

	//打开json
	fileJsonName := "./plugin/" + name + "/" + name + ".json"
	jsonBytes, err := utils.FileRead(fileJsonName)
	if err != nil {
		setting.ZAPS.Errorf("打开%s文件失败 %v", name, err)
		return plugin, errors.New("打开" + fileJsonName + "失败")
	}

	//解析json
	err = json.Unmarshal(jsonBytes, &plugin)
	if err != nil {
		return plugin, errors.New("格式化" + fileJsonName + "失败")
	}

	return plugin, nil
}
