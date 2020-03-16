package utils

import "fmt"

// RedisGetStringResult 将Redis查询获得的数据，转换为字符串
func RedisGetStringResult(strRes interface{}) string {
	resArray, ok := strRes.([]uint8)
	if !ok {
		fmt.Println("redis查询过程错误！")
		return ""
	}
	strByte := []byte{}
	for _, b := range resArray {
		strByte = append(strByte, byte(b))
	}
	res := string(strByte)
	return res
}

// JointKey 拼接redis的key
func JointKey(args ...interface{}) string {
	var result string
	for k, s := range args {
		if k == 0 {
			result = fmt.Sprintf("%v", s)
		} else {
			result = fmt.Sprintf("%s::%v", result, s)
		}
	}
	return result
}
