package util

import (
	"github.com/axgle/mahonia"
	"regexp"
	"unicode"
	"math"
)

/***
 * 将utf-8 转码 gbk
 */
func ConvertToGbk(src string) string {
	enc := mahonia.NewEncoder("gbk")
	str := enc.ConvertString(src)
	return str
}


// GetStrLength 返回输入的字符串的字数，汉字和中文标点算 2个字数，英文和其他字符 1 个字数，不足 1 个算 1个
func GetStrLength(str string) float64 {
	var total float64

	reg := regexp.MustCompile("/·|，|。|《|》|‘|’|”|“|；|：|【|】|？|（|）|、/")

	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || reg.Match([]byte(string(r))) {
			total = total + 2
		} else {
			total = total + 1
		}
	}

	return math.Ceil(total)
}