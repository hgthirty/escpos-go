package main

/*
#cgo LDFLAGS: -L. -llibhello
#include <stdio.h>
#include <stdlib.h>
#include "hello.h"
void printStr(char *p){
 fprintf(stdout, "imgPath = %s\n", p);
};
GoStringArr testStringArr(int len);
void printArr(GoStringArr strArr){
	_GoString *p = strArr.data;
	int i = 0;
	for (i = 0; i < strArr.len; i++) {
		//printf("Str : %s \n", (p + i)->p);
	}
}
const char* getOffsetData(GoStringArr strArr , int offset){
	_GoString *head = strArr.data;
	return  (head + offset)->p;
}
*/
import "C"
import "fmt"

func main() {
	//shell := syscall.MustLoadDLL("libhello.dll")
	//testStringArr := shell.MustFindProc("testStringArr")
	//var str1 *C.GoStr
	//strArr, _, err := testStringArr.Call(10)
	//fmt.Println( err)
	strArr := C.testStringArr(10)
	C.printArr(strArr)
	fmt.Println(C.GoString(C.getOffsetData(strArr ,9)))

	//print()
	//b := (C.testCharP())
	//C.printStr(b)
	//println( C.GoString(b) )
	//testCharP := shell.MustFindProc("testCharP")
	//str,_,err =testCharP.Call()
	//
	//fmt.Println((*byte)(str),err)
}
