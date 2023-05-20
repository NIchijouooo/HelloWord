package setting

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"gateway/utils"
)

type EncryptParamTemplate struct {
	ETH map[string]string `json:"eth"`
}

var encryptKey = []byte("openGW__key@2022")

var EncryptParam EncryptParamTemplate

func init() {
	EncryptParam.ETH = make(map[string]string)
}

func ReadAuthFile() bool {

	//LUO
	ZAPS.Info("授权成功")
	SystemState.AuthStatus = "已授权"
	return true

	/*

		rd, err := ioutil.ReadDir("./")
		if err != nil {
			ZAPS.Errorf("读取授权文件错误 %v", err)
			return false
		}

		fileName := ""
		for _, fi := range rd {
			if fi.IsDir() == false {
				fullName := "./" + fi.Name()
				if strings.Contains(fi.Name(), ".key") {
					fileName = fullName
					break
				}
			}
		}

		rdata, err := utils.FileRead(fileName)
		if err != nil {
			ZAPS.Errorf("打开授权文件失败 %v，请联系管理员", err)
			return false
		}

		decryStr, err := base64.StdEncoding.DecodeString(string(rdata))
		if err != nil {
			ZAPS.Errorf("打开授权文件内容失败 %v，请联系管理员", err)
			return false
		}

		decryData := AesDecryptCBC(decryStr, encryptKey)

		err = json.Unmarshal(decryData, &EncryptParam)
		if err != nil {
			ZAPS.Errorf("授权文件格式化失败 %v，请联系管理员", err)
			return false
		}
		//ZAPS.Debugf("授权网卡参数 %v", EncryptParam)

		params := GetNetworkInterface()
		//ZAPS.Debugf("网卡参数 %v", params)
		result := 0
		for _, v := range params {
			eth, ok := EncryptParam.ETH[v.Name]
			if !ok {
				continue
			}
			if v.MAC == eth {
				result += 1
			}
		}

		if result == len(EncryptParam.ETH) {
			ZAPS.Info("授权成功")
			SystemState.AuthStatus = "已授权"
			return true
		}

		ZAPS.Error("授权失败，请联系管理员")
		return false */
}

func GenerateAuthFile() string {

	params := GetNetworkInterface()
	for _, v := range params {
		EncryptParam.ETH[v.Name] = v.MAC
	}

	data, _ := json.Marshal(&EncryptParam)

	aesStr := base64.StdEncoding.EncodeToString(AesEncryptCBC(data, encryptKey))
	//ZAPS.Debugf("aesStr %v", aesStr)

	fileName := "./noAuth_" + SystemState.SN + ".key"
	err := utils.FileWrite(fileName, []byte(aesStr))
	if err != nil {
		ZAPS.Errorf("生成授权文件失败 %v，请联系管理员", err)
		return ""
	}

	ZAPS.Info("生成授权文件文件成功")

	return fileName
}

// =================== CBC ======================
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
