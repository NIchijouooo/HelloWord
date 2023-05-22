package setting

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gateway/utils"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

//func DownloadFileFeisjy(fileSavePath string, fileName string, url string) error {
//	var req *http.Request = nil
//	var err error
//
//	// 实例化一个http client对象
//	client := &http.Client{}
//	// 实例化一个http Get Request对象
//	req, err = http.NewRequest("GET", url, nil)
//
//	if err != nil {
//		return errors.New(fmt.Sprintf("http.NewRequest Fail! %v", err))
//	}
//
//	req.Header.Add("Cookie", "SHAREJSESSIONID=f0bae696-f1be-4df8-a5e1-10ba646b7ff0")
//
//	// 将http Get 请求对象发送服务器上 若发送成功则会得到响应
//	res, err := client.Do(req)
//	if err != nil {
//		return errors.New(fmt.Sprintf("client.Do Fail! %v", err))
//	}
//	defer res.Body.Close()
//
//	// 读取响应体
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return errors.New(fmt.Sprintf("ioutil.ReadAll Fail! %v", err))
//	}
//
//	// 判断是否路径是否存在，不存在则递归创建
//	utils.DirIsExist(fileSavePath)
//
//	filePath := fmt.Sprintf("%s%s", fileSavePath, fileName)
//
//	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
//	defer file.Close()
//	if err != nil {
//		return errors.New(fmt.Sprintf("文件打开失败[%s]! %v", filePath, err))
//	} else {
//		// 创建一个带缓冲区的写入器
//		write := bufio.NewWriter(file)
//		// 将响应体写入文件中
//		_, err = write.WriteString(string(body))
//		fmt.Printf(">>>>>>")
//		if err != nil {
//			return err
//		}
//		write.Flush()
//	}
//
//	return nil
//}

func DownloadFileFeisjy(fileSavePath string, fileName string, url string, maxRetries int) error {

	// 实例化一个http client对象
	client := &http.Client{}
	// 实例化一个http Get Request对象
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("http.NewRequest Fail! %v", err))
	}

	req.Header.Add("Cookie", "SHAREJSESSIONID=f0bae696-f1be-4df8-a5e1-10ba646b7ff0")

	// 设置 Range 头字段，指定要下载的文件块的起始和结束位置
	blockSize := 1024 * 1024 // 每次下载1MB的数据
	var start, end int64
	var fileSize int64
	for retries := 0; retries < maxRetries; retries++ {
		fmt.Printf("[%d]", retries+1)

		// 获取文件大小
		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("获取文件大小失败，正在进行第 %d 次重试，错误信息：%v\n", retries+1, err)
			continue
		}
		defer res.Body.Close()

		fileSize = res.ContentLength
		if fileSize == -1 {
			return errors.New(fmt.Sprintf("无法获取文件大小!"))
		}

		// 判断是否路径是否存在，不存在则递归创建
		utils.DirIsExist(fileSavePath)

		filePath := fmt.Sprintf("%s%s", fileSavePath, fileName)

		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return errors.New(fmt.Sprintf("文件打开失败[%s]! %v", filePath, err))
		}
		defer file.Close()

		// 创建一个带缓冲区的写入器
		write := bufio.NewWriter(file)

		// 分块下载文件
		for start < fileSize {

			fmt.Printf("□□□□□□□□")

			end = start + int64(blockSize) - 1
			if end >= fileSize {
				end = fileSize - 1
			}

			// 设置 Range 头字段，指定要下载的文件块的起始和结束位置
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

			// 将http Get 请求对象发送服务器上 若发送成功则会得到响应
			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("下载失败，正在进行第 %d 次重试，错误信息：%v\n", retries+1, err)
				break
			}
			defer res.Body.Close()

			// 将响应体直接写入文件中
			_, err = io.Copy(write, res.Body)
			if err != nil {
				fmt.Printf("写入文件失败，正在进行第 %d 次重试，错误信息：%v\n", retries+1, err)
				break
			}

			// 更新起始位置
			start = end + 1
		}
		fmt.Printf("\n")
		// 确保文件被正确关闭
		err = write.Flush()
		if err != nil {
			return errors.New(fmt.Sprintf("文件关闭失败[%s]! %v", filePath, err))
		}

		return nil
	}

	return errors.New(fmt.Sprintf("下载失败，已达到最大重试次数：%d", maxRetries))
}

/*
*
上传网关配置到云平台,先将配置文件压缩成一个压缩包再上传
压缩的文件列表暂时写死测试
url=上传接口
deviceId = 云平台下发的设备id,再返回给云平台
*/
func UpLoadGatewayConfigFeisjy(url string, deviceId string) (string, error) {
	fileName := "current.zip"
	dir := "./tmp/"
	filePath := dir + fileName
	err := utils.CompressFilesToZip([]string{"./selfpara/collInterface.json",
		"./selfpara/reportModel.json",
		"./selfpara/reportServiceParamListFeisjyIot.json",
		"./selfpara/commInterfaceProtocol.json"}, filePath)
	if err != nil {
		return "", err
	}

	//创建一个新的 multipart form
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	_ = writer.WriteField("deviceId", deviceId)
	var file fs.File
	var part io.Writer
	file, err = os.Open(filePath)
	defer file.Close()
	if err != nil {
		return "", errors.New(fmt.Sprintf("UpLoadFileFeisjy -> os.Open() Faile! : %s", err.Error()))
	}
	defer file.Close()

	// 创建一个新的表单字段，并将文件内容写入其中
	part, err = writer.CreateFormFile("fileName", filepath.Base(filePath))
	if err != nil {
		return "", errors.New(fmt.Sprintf("UpLoadFileFeisjy -> writer.CreateFormFile() Faile! : %s", err.Error()))
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", errors.New(fmt.Sprintf("UpLoadFileFeisjy -> io.Copy() Faile! : %s", err.Error()))
	}

	// 关闭 multipart form 并设置 Content-Type
	err = writer.Close()
	contentType := writer.FormDataContentType()
	if err != nil {
		return "", errors.New(fmt.Sprintf("UpLoadFileFeisjy -> writer.Close() Faile! : %s", err.Error()))
	}

	// 发送 HTTP POST 请求并上传文件
	response, err := http.Post(url, contentType, requestBody)
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return "", errors.New(fmt.Sprintf("UpLoadFileFeisjy -> http.Post() Faile! : %s", err.Error()))
	}
	defer response.Body.Close()
	return fileName, nil
}

func UpLoadFileFeisjy(filePath string, url string, uniqueCode string, deviceAddr string) error {

	var file fs.File
	var err error
	var part io.Writer

	//创建一个新的 multipart form
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	// 打开上传文件
	if len(uniqueCode) > 0 {
		_ = writer.WriteField("uniqueCode", uniqueCode)
	}
	if len(deviceAddr) > 0 {
		_ = writer.WriteField("deviceAddr", deviceAddr)
	}
	file, err = os.Open(filePath)
	defer file.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("UpLoadFileFeisjy -> os.Open() Faile! : %s", err.Error()))
	}
	defer file.Close()

	// 创建一个新的表单字段，并将文件内容写入其中
	part, err = writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return errors.New(fmt.Sprintf("UpLoadFileFeisjy -> writer.CreateFormFile() Faile! : %s", err.Error()))
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return errors.New(fmt.Sprintf("UpLoadFileFeisjy -> io.Copy() Faile! : %s", err.Error()))
	}

	// 关闭 multipart form 并设置 Content-Type
	err = writer.Close()
	contentType := writer.FormDataContentType()
	if err != nil {
		return errors.New(fmt.Sprintf("UpLoadFileFeisjy -> writer.Close() Faile! : %s", err.Error()))
	}

	// 发送 HTTP POST 请求并上传文件
	response, err := http.Post(url, contentType, requestBody)
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return errors.New(fmt.Sprintf("UpLoadFileFeisjy -> http.Post() Faile! : %s", err.Error()))
	}
	defer response.Body.Close()

	// 打印响应结果
	//fmt.Println("Upload successful!")
	//fmt.Println("Response:", response.Status)

	return nil
}
