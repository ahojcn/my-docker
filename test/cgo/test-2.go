package main

/*
#include <stdio.h>
__attribute__((constructor)) void before_main() {
	printf("before main\n");
}
*/
import "C"

import "log"

func main() {
	log.Println("hello world!")
}
