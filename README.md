# my-docker

@a:<ahojcn@gmail.com>

---


## branch 5:

主要实现如下命令：
```shell
宿主机 $ ./mydocker run -it /bin/sh
容器内 # ps -ef
......
```

tags:

5-1:实现 run 方法。

5-2:实现 init 方法（在容器起来前先执行 mount -t proc proc /proc）。

5-3:实现用户进程为1号进程
（因为5-2的程序运行后，1号进程并不是我们要的 sh，而是 mount 操作使用的一个进程）。


## branch 6:

主要实现 memory 限制，在 branch5 基础上。

---

tags:

6-1:
修改 run 命令。
加入参数 -m 表示接受 memory 限制。

6-2:
实现一些 cgroups utils 函数。
找到当前进程的 cgroup 的路径。

6-3:
实现资源限制。
memory 的 Set 和 Apply 函数将内存限制写入文件。

6-4:
实现容器资源隔离。

6-5:
实现资源删除。
资源删除其实是在进程结束的时候把限制解除，其实就是把对应的文件夹给删除。
Remove 函数。

---


EOF.
