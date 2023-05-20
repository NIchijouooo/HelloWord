package utils

import "errors"

func TypeConvertToInt32(value interface{}) (int32, error) {

	switch value.(type) {
	case uint8:
		return int32(value.(uint8)), nil
	case uint16:
		return int32(value.(uint16)), nil
	}

	return 0, errors.New("变量类型不支持")
}

func TypeConvertToUint32(value interface{}) (uint32, error) {

	switch value.(type) {
	case uint8:
		return uint32(value.(uint8)), nil
	}

	return 0, errors.New("变量类型不支持")
}
