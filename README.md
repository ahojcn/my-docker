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

## branch 7:

在 branch6 的基础上，梳理代码，使其更清晰，也为了方便拓展加入其他的 subsystem 限制。

tags:

7-1: 实现。

注意:

循环引用问题！
注意 memory.swappiness 的限制 memory.swappiness = 0 时候是没有交换分区的
参考：
https://segmentfault.com/a/1190000008125116
https://segmentfault.com/a/1190000008125359

---

## branch 8:

在 branch7 的基础上，加入管道功能。

通常进程间通信使用管道进行通信，所以本文将对之前进程间传输的 command 用管道的方式来执行。

branch7 实现的那个简单版本的 run 命令有一个缺陷，就是传递参数。 在父进程和子进程之间传参，是通过调用命令后面跟上参数，也就是 /proc/self/exec init args 这种方式进行的 然后在 init 进程内去解析这个参数，执行响应的命令。 缺点是：如果用户输入的参数过长，或者其中带有一些特殊字符，那么这种方案就会失败了。 其实 runC 实现的方案是通过匿名管道来实现父子进程之间通信的。 branch8 就是实现这个功能。

tags:

8-1: 具体实现。

---

## branch 11:

在 branch8 的基础上，一步步实现使用 busybox 创建容器。

tags:

11-1: 实现改变 init 程序执行路径。 给 cmd 加入一些参数， cmd.Dir = "/root"，在执行用户程序的时候可以设置该程序在哪个目录下执行。

11-2: 实现用 busybox 作为容器的跟目录。


---


EOF.
