package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DirIsExist(dir string) {
	//exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	_, err := os.Stat(dir)
	//文件夹或者文件不存在
	if err != nil {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Printf("创建%s文件夹失败 %v", dir, err)
		}
		_ = os.Chmod(dir, 0777)
	}
}

func FileIsExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func FileCreate(name string) error {
	fp, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(fp)

	_ = fp.Sync()
	return nil
}

func FileWrite(dir string, data []byte) error {
	fp, err := os.OpenFile(dir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(fp)

	_, err = fp.Write(data)
	if err != nil {
		return err
	}
	_ = fp.Sync()
	return nil
}

func FileRead(dir string) ([]byte, error) {
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func FileRemove(dir string) error {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	err := os.Remove(exeCurDir + dir)
	if err != nil {
		return err
	}
	return nil
}

func FileCopy(src, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
	}
	return nil
}

// 压缩多个路径到一个zip文件里面
// Param 1: 需要添加到zip文件里面的路径
// Param 2: 打包后zip文件名称
func CompressDirsToZip(srcDirMap []string, fileName string) error {

	//创建zip文件
	newZipFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(newZipFile *os.File) {
		err = newZipFile.Close()
		if err != nil {

		}
	}(newZipFile)

	//打开zip文件
	zipWriter := zip.NewWriter(newZipFile)
	defer func(zipWriter *zip.Writer) {
		err = zipWriter.Close()
		if err != nil {

		}
	}(zipWriter)

	for _, v := range srcDirMap {
		//遍历文件夹
		err = filepath.Walk(v, func(path string, info os.FileInfo, _ error) error {

			//如果是源路径，提前进行下一个遍历(没明白这个判断)
			//if path == v {
			//	return nil
			//}

			// 获取文件头信息
			header, _ := zip.FileInfoHeader(info)
			header.Name = strings.TrimPrefix(path, v+`/`)

			// 判断文件是不是文件夹
			if info.IsDir() {
				header.Name += `/`
			} else {
				// 设置：zip的文件压缩算法
				header.Method = zip.Deflate
			}

			// 创建压缩包头部信息
			writer, _ := zipWriter.CreateHeader(header)
			if !info.IsDir() {
				file, _ := os.Open(path)
				defer func(file *os.File) {
					err = file.Close()
					if err != nil {
					}
				}(file)
				_, err = io.Copy(writer, file)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	return err
}

// 压缩单个路径到一个zip文件里面
// Param 1: 需要添加到zip文件里面的路径
// Param 2: 打包后zip文件名称
func CompressDirToZip(srcDirMap string, fileName string) error {

	//创建zip文件
	newZipFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(newZipFile *os.File) {
		err = newZipFile.Close()
		if err != nil {

		}
	}(newZipFile)

	//打开zip文件
	zipWriter := zip.NewWriter(newZipFile)
	defer func(zipWriter *zip.Writer) {
		err = zipWriter.Close()
		if err != nil {

		}
	}(zipWriter)

	//遍历文件夹
	err = filepath.Walk(srcDirMap, func(path string, info os.FileInfo, _ error) error {
		//如果是源路径，提前进行下一个遍历(没明白这个判断)
		//if path == v {
		//	return nil
		//}

		// 获取文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcDirMap+`\`)

		// 判断文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建压缩包头部信息
		writer, _ := zipWriter.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
				}
			}(file)
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// 压缩单个路径到一个zip文件里面(压缩后不带路径名)
// Param 1: 需要添加到zip文件里面的路径
// Param 2: 打包后zip文件名称
func NewCompressDirToZip(srcDir string, dstFileName string) error {
	var srcFiles []string
	DirIsExist(srcDir)
	// 读取目录下全部文件
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {
		if !info.IsDir() {
			srcFiles = append(srcFiles, path)
		}
		return nil
	})
	if err != nil || len(srcFiles) == 0 {
		log.Fatalln("NewCompressDirToZip error or no files err = ", err)
		return err
	}
	return CompressFilesToZip(srcFiles, dstFileName)
}

// 压缩单个路径到一个zip文件里面
// Param 1: 需要添加到zip文件里面的路径
// Param 2: 打包后zip文件名称
func CompressFilesToZip(srcFiles []string, dstFileName string) error {
	//创建zip文件
	newZipFile, err := os.Create(dstFileName)
	if err != nil {
		return err
	}
	defer func(newZipFile *os.File) {
		err = newZipFile.Close()
		if err != nil {

		}
	}(newZipFile)

	//打开zip文件
	zipWriter := zip.NewWriter(newZipFile)
	defer func(zipWriter *zip.Writer) {
		err = zipWriter.Close()
		if err != nil {

		}
	}(zipWriter)

	for _, file := range srcFiles {
		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// Get the file information
		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = filepath.Base(file)

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fileToZip)

	}

	return err
}

var crc16tab = []uint16{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
	0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
	0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
	0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
	0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
	0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
	0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
	0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
	0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
	0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
	0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
	0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
	0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
	0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
	0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
	0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
	0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
	0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
	0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
	0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
	0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
	0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
	0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
	0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
	0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
	0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
	0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
	0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
	0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
	0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
	0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
	0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0,
}

func crc16_ccitt(crc uint16, buf []byte) uint16 {
	for _, b := range buf {
		crc = (crc << 8) ^ crc16tab[((crc>>8)^uint16(b))&0xFF]
	}
	return crc
}

const (
	CRC_READ_BUFF_SIZE = 1024
)

func CalculateFileCRC16(filePath string) string {
	var file_crc uint16
	crc_read_buff := make([]byte, CRC_READ_BUFF_SIZE)

	fd, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	// get the size of file
	fi, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	length := fi.Size()

	for length > 0 {
		n, err := fd.Read(crc_read_buff)
		if err != nil && err != io.EOF {
			panic(err)
		}
		file_crc = crc16_ccitt(file_crc, crc_read_buff[:n])
		length -= int64(n)
	}

	return fmt.Sprintf("%d", file_crc)
}

func GetAllFileFormDir(dir string) ([]string, error) {

	var fileList []string

	dirInfo, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	fileInfolist, err := dirInfo.Readdir(-1)
	dirInfo.Close()
	if err != nil {
		return nil, err
	}

	for _, fileinfo := range fileInfolist {
		fileList = append(fileList, fmt.Sprintf("%s/%s", dir, fileinfo.Name()))
	}

	return fileList, nil
}

func Unzip(zipFilePath, destFolder string) bool {
	exist := FileIsExist(zipFilePath)
	if !exist {
		log.Fatalln("unzip file error file is not exist path = [%s]", zipFilePath)
		return false
	}
	DirIsExist(destFolder)
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		log.Fatalln("open zip error", err)
		return false
	}
	defer r.Close()
	for _, f := range r.File {
		filePath := filepath.Join(destFolder, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, f.Mode())
			continue
		}
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Fatalln("mk dir all error", err)
			return false
		}
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatalln("open file error", err)
			return false
		}
		defer outFile.Close()
		rc, err := f.Open()
		if err != nil {
			log.Fatalln("open error", err)
			return false
		}
		defer rc.Close()
		if _, err = io.Copy(outFile, rc); err != nil {
			log.Fatalln("copy file error", err)
			return false
		}
	}
	return true
}
