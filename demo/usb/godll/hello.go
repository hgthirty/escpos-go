package main

//export add
func add( a , b int) int {
	return a+b
}


func main() {
	// Need a main function to make CGO compile package as C shared library
}