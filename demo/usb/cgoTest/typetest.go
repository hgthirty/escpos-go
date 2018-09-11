package main

/*
#include <stdio.h>
#include <stdlib.h>

char ch = 'M';
unsigned char uch = 253;
short st = 233;
int i = 257;
long lt = 11112222;
float f = 3.14;
double db = 3.15;
void * p;
char *str = "const string";
char str1[64] = "char array";

void printI(void *i)
{
    printf("print i = %d\n", (*(int *)i));
}

struct ImgInfo {
    char *imgPath;
    int format;
    unsigned int width;
    unsigned int height;
};

void printStruct(struct ImgInfo *imgInfo)
{
    if(!imgInfo) {
        fprintf(stderr, "imgInfo is null\n");
        return ;
    }

    fprintf(stdout, "imgPath = %s\n", imgInfo->imgPath);
    fprintf(stdout, "format = %d\n", imgInfo->format);
    fprintf(stdout, "width = %d\n", imgInfo->width);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
	"reflect"
)

func main() {

	testtypes()
}

func testtypes() {
	fmt.Println("----------------Go to C---------------")
	fmt.Println(C.char('Y'))
	fmt.Printf("%c\n", C.char('Y'))
	fmt.Println(C.uchar('C'))
	fmt.Println(C.short(254))
	fmt.Println(C.long(11112222))
	var goi int = 2
	// unsafe.Pointer --> void *
	cpi := unsafe.Pointer(&goi)
	C.printI(cpi)
	fmt.Println("----------------C to Go---------------")
	fmt.Println(C.ch)
	fmt.Println(C.uch)
	fmt.Println(C.st)
	fmt.Println(C.i)
	fmt.Println(C.lt)
	f := float32(C.f)
	fmt.Println(reflect.TypeOf(f))
	fmt.Println(C.f)
	db := float64(C.db)
	fmt.Println(reflect.TypeOf(db))
	fmt.Println(C.db)
	// 区别常量字符串和char数组，转换成Go类型不一样
	str := C.GoString(C.str)

	fmt.Println(str)

	fmt.Println(reflect.TypeOf(C.str1))
	var charray []byte
	for i := range C.str1 {
		if C.str1[i] != 0 {
			charray = append(charray, byte(C.str1[i]))
		}
	}

	fmt.Println(charray)
	fmt.Println(string(charray))

	for i := 0; i < 10; i++ {
		imgInfo := C.struct_ImgInfo{imgPath: C.CString("../images/xx.jpg"), format: 0, width: 500, height: 400}
		defer C.free(unsafe.Pointer(imgInfo.imgPath))
		C.printStruct(&imgInfo)
	}

	fmt.Println("----------------C Print----------------")
}
