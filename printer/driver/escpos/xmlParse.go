package escpos

import (
	"io"
	"io/ioutil"
	"fmt"
	"encoding/xml"
	"os"
	"net/http"
)

/***
 * 解析基础的xml io流
 */
func ParseXml(io io.Reader) (root Root, err error) {

	content, err := ioutil.ReadAll(io)
	if err != nil {
		fmt.Println(err)
		return root, err
	}

	err = xml.Unmarshal(content, &root)

	return root, err
}

/***
 * 解析基础的xml io流
 */
func ParseString(content string) (root Root, err error) {

	err = xml.Unmarshal([]byte(content), &root)

	return root, err
}

/***
 * 解析基础的xml 文件
 */
func ParseLocalXml(filename string) (root Root, err error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return root, err
	}
	root, err = ParseXml(file)
	return root, err
}

/***
 * 解析基础的xml url (get请求)
 */
func ParseRemoteXml(url string) (root Root, err error) {
	response, err := http.Get(url)
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err)
		return root, err
	}

	root, err = ParseXml(response.Body)
	return root, err
}
