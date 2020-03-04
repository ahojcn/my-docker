package main

/*
int add(int a, int b) {
	return a + b;
}
*/
import "C"  // 注意是紧跟着上面的注释代码，不能有空行
import "fmt"

func main() {
	a := C.int(1)
	b := C.int(2)
	value := C.add(a, b)
	fmt.Printf("value: %d\n", value)
}