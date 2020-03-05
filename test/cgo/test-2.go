package main

///*
//#include <stdio.h>
//__attribute__((constructor)) void before_main() {
//	printf("before main\n");
//}
//*/
//import "C"
//
//import "log"
//
//func main() {
//	log.Println("hello world!")
//}

/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
void test() {
	char *mydocker_pid;
	mydocker_pid = getenv("mydocker_pid");
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };
	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", mydocker_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);
		if (setns(fd, 0) == -1) {
			fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	int res = system("/bin/sh");
	exit(0);
	return;
}
*/
import "C"

func main() {
	C.test()
}
