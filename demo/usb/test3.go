package main

/*
#cgo LDFLAGS: -L. -lAppendPrint
#include <stdio.h>
int HelloWorld(int a, int b);

int add(int a ,int b){
    return a +b ;
}

typedef int (*intFunc) ();

int bridge_int_func(intFunc f)
{
     return f();
}

int fortytwo()
{
     return 42;
}
*/
import "C"
import "fmt"

func main() {
	C.puts(C.CString("Hello, world\n"))
	fmt.Println(C.add(1,2))
	test()

	C.HelloWorld(1,2)
}

func test() {
	f := C.intFunc(C.fortytwo)
	fmt.Println(int(C.bridge_int_func(f)))
}