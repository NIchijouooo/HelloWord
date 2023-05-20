package setting

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
)

type HTTPDownFileTemplate struct {
	io.Reader
	Total      int64
	Current    int64
	Percent    int
	PrePercent int
	PercentOut chan int
	Finish     chan bool
	Error      chan error
}

func NewHTTPDownFile() *HTTPDownFileTemplate {
	return &HTTPDownFileTemplate{
		PercentOut: make(chan int),
		Finish:     make(chan bool),
		Error:      make(chan error),
	}
}

func HttpDownFileFromURL(url string, fileName string, reader *HTTPDownFileTemplate) {
	ZAPS.Debugf("httpDownFile开始获取URL %v", url)
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	response, err := http.Get(url)
	if err != nil {
		ZAPS.Errorf("httpGet请求失败 %v", err)
		reader.Error <- err
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			ZAPS.Errorf("httpGet关闭body失败 %v", err)
		}
	}(response.Body)

	f, _ := os.Create(fileName)

	ZAPS.Debugf("resonseBody %v", response.Body)
	reader.Reader = response.Body
	reader.Total = response.ContentLength

	_, err = io.Copy(f, reader)
	if err != nil {
		ZAPS.Errorf("httpGet从httpResponse拷贝文件错误 %v", err)
		reader.Error <- err
	}
	reader.Finish <- true
}

//func (h *HTTPDownFileTemplate) Read(p []byte) (int, error) {
//	n, err := h.Reader.Read(p)
//	if err != nil {
//		ZAPS.Errorf("httpGet从httpResponse读取文件错误 %v", err)
//	}
//
//	h.Current += int64(n)
//	h.Percent = int(h.Current * 100 / h.Total)
//	if (h.Percent%10 == 0) && (h.Percent != h.PrePercent) {
//		ZAPS.Debugf("percent %v", h.Percent)
//		h.PercentOut <- h.Percent
//	}
//	h.PrePercent = h.Percent
//	return n, err
//}
