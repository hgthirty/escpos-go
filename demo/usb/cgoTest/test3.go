package main

/*
#cgo LDFLAGS: -L"D:/install/pkg-config/bin/libhello" -llibhello
#include <stdio.h>
#include <stdlib.h>
#include "hello.h"
void printStr(char *p){
 fprintf(stdout, "imgPath = %s\n", p);
};
GoStringArr testStringArr(int len);
GoStringArr printer();

void printArr(GoStringArr strArr){
	_GoString *p = strArr.data;
	int i = 0;
	for (i = 0; i < strArr.len; i++) {
		printf("Str : %s \n", (p + i)->p);
	}
}

_GoString_ getOffset(GoStringArr strArr , int offset){
	_GoString_ *head = (_GoString_*)strArr.data;
	return  *(head + offset);
}

const char* getOffsetData(GoStringArr strArr , int offset){
	_GoString *head = strArr.data;
	return  (head + offset)->p;
}


*/
import "C"
import (
	"fmt"
	"os"
)

func main() {

	strArr := C.testStringArr(10)
	//C.printArr(strArr)
	fmt.Println(C.GoString(C.getOffsetData(strArr ,9)))
	paths :=C.printer()
	//fmt.Println(paths)
	fmt.Println("=================");
	//fmt.Print(C.getOffset(paths ,1))

	for i:=0 ;i<4 ;i++{
		//fmt.Println(C.getOffset(paths ,C.int(i) ))
		path :=C.getOffset(paths ,C.int(i))
		fileHandel,err := os.OpenFile(path, os.O_RDWR, 0660)
		if err != nil{
			fmt.Println(err)
		}
		fileHandel.Write([]byte{0x1B ,0x76});
		status  := []byte{0}
		fileHandel.Read(status);
		fileHandel.Close()
		fmt.Printf("status: %x\n" ,status)

	}

}
