#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define MB (1024 * 1024)

/*
  申请内存空间
  测试 memory.oom_control：
  oom_kill_disable 0
  under_oom 0
  */

int main() {

    char *p;
    int i = 0;
    while (1) {
        p = (char *)malloc(MB);
        memset(p, 0, MB);
        printf("%dM memory allocated\n", ++i);
        sleep(1);
    }

    return 0;
}