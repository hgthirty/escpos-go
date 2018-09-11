package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
//创建 _GoString 字符串 (给C程序调用)
_GoString_  NewGoString(const char * str) {
	size_t len = strlen(str)+1;
	_GoString_ temp;
	char *buffer;
	buffer = (char*)malloc(len);
	memset(buffer, 0, len);
	memcpy(buffer, str, len-1);
	temp.p = buffer;
	temp.n = len;
	return temp;
}

//测试创建字符串 , (经过测试  _GoString_ 对于go 的string 可以直接使用)
_GoString_  testGoString() {
	_GoString_ temp = NewGoString("hello world 100010101001010lllsss");
	return temp;
}

//打印字符串
void printStr(char *p){
 fprintf(stdout, "imgPath = %s\n", p);
};

//测试字符串指针
char* testCharPoint(){
	char *p = "hello testCharPoint!\n";
	return p;
}

*/
import "C"
import "fmt"

func main() {
	str := C.testGoString()
	fmt.Println(str)
	fmt.Println("字符串内容", str, "len:", len(string(str)))

	str1 := C.testCharPoint()
	fmt.Print(C.GoString(str1))
}
