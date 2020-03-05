#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define MB (1024 * 1024)

int main() {
    char *p;
    int i = 0;
    while (1) {
        p = (char *)malloc(MB);
        memset(p, 0, MB);
        printf("%dMB memory allocated\n", ++i);
        sleep(1);
    }

    return 0;
}
