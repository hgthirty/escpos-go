package util

import (
	"github.com/axgle/mahonia"
	"regexp"
	"unicode"
	"math"
	"os/exec"
	"fmt"
	"strings"
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

/***
  获取主机所在网络列表
 */
func GetHostNet() []string {
	cmd := exec.Command("ip", "addr")
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Fprintf(os.Stdout, "Result: %s", buf)
	r, _ := regexp.Compile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.\d{1,3}/\d{1,2}`)
	//fmt.Println(r.FindAllString(string(buf),-1))
	temp := r.FindAllString(string(buf), -1)
	ips := make([]string, 0)
	for _, ip := range temp {
		ipInfo := strings.Split(ip, "/")
		//忽略本地ip
		if  len(ipInfo) != 2 || ipInfo[0] == "127.0.0.1" {
			continue
		}
		ips = append(ips ,ipInfo[0])
	}
	return ips
}
