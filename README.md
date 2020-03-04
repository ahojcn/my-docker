<<<<<<< HEAD
branch 25:

在 branch24 的基础上，实现容器的网络并测试。

---

EOF
=======
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

5-1: 实现 run 方法。

5-2: 实现 init 方法（在容器起来前先执行 mount -t proc proc /proc）。

5-3:

实现用户进程为1号进程
（因为5-2的程序运行后，1号进程并不是我们要的 sh，而是 mount 操作使用的一个进程）。


## branch 6:

主要实现 memory 限制，在 branch5 基础上。

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

6-4: 实现容器资源隔离。

6-5:

实现资源删除。
资源删除其实是在进程结束的时候把限制解除，其实就是把对应的文件夹给删除。
Remove 函数。



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



## branch 8:

在 branch7 的基础上，加入管道功能。

通常进程间通信使用管道进行通信，所以本文将对之前进程间传输的 command 用管道的方式来执行。

branch7 实现的那个简单版本的 run 命令有一个缺陷，就是传递参数。 在父进程和子进程之间传参，是通过调用命令后面跟上参数，也就是 /proc/self/exec init args 这种方式进行的 然后在 init 进程内去解析这个参数，执行响应的命令。 缺点是：如果用户输入的参数过长，或者其中带有一些特殊字符，那么这种方案就会失败了。 其实 runC 实现的方案是通过匿名管道来实现父子进程之间通信的。 branch8 就是实现这个功能。

tags:

8-1: 具体实现。



## branch 11:

在 branch8 的基础上，一步步实现使用 busybox 创建容器。

tags:

11-1: 实现改变 init 程序执行路径。 给 cmd 加入一些参数， cmd.Dir = "/root"，在执行用户程序的时候可以设置该程序在哪个目录下执行。

11-2: 实现用 busybox 作为容器的跟目录。




## branch 12:

在 branch11 的基础上，一步步实现使用AUFS包装busybox。

branch 11 的存在问题： 利用 busybox 创建的容器, 创建文件夹并且创建文件。 退出容器后, 查看宿主机的内容, 返现内容在宿主机中也存在。 这样会有一个问题, 其实 busybox 就是容器的镜像层, 如果多个容器共享该镜像层, 那就会造成容器之间互相看到对方文件, 并且文件覆盖等等问题。 branch12 就是利用 AUFS 解决此问题

https://www.jianshu.com/p/ecbdcc98db76

tags:

12:

根据 busybox 镜像生成容器，其实就是解压 busybox.tar 生成 rootPath/busybox。
创建挂载点、可写层，并将 writeLayer 和 busybox 挂载到 mnt 下。
当退出时候执行 umount mnt/ 并删除 writeLayer。



## branch 13:

在 branch12 的基础上，一步步实现 volume 操作。

branch12 中的容器内增删文件都不会保存。 如果用户需要保存则需要 -v 参数把宿主机的目录挂载到容器内。

tags:

13-1: 实现单个 volume 挂载到容器中。注意 umount 的时候先 umount volume path，再 umount mnt/。

13-2: 实现多个 -v。



## branch 14:

在 branch13 的基础上，实现保存镜像。

直接把容器运行时的整个目录保存起来即可。

tags:

14: 具体实现。



## branch 15:

在 branch14 的基础上，实现容器的后台运行。

可以允许程序在后台运行。 因为容器再后台运行，实际上就是不采用交互式，也就是父进程会退出，然后用户进程会继续执行。 父进程退出后，子进程会被 1 号进程接管。

容器再退出时需要深处对应的中间文件夹 writeLayer，但是 branch15 的代码会提示删除失败 error 因为 /bin/top 进程与该文件夹有关。什么原因没搞清楚……

tags:

15: 实现容器的后台运行。



## branch 16:

在 branch15 的基础上，实现查看运行中的容器。

tags:

16-1: 测试代码，实现查看运行中的容器。
16-2: 根据测试代码实现。



## branch 17:

在 branch16 的基础上，实现产看容器的日志。

将后台程序日志持久化到文件中。
如果容器是后台运行，就把标准输出重定向到某一个文件中。
查看容器日志时读取文件内容并显示。


tags:

17: 具体实现。



## branch 18:

在 branch17 的基础上，实现进入容器 Namespace。

tags:

18-1: test cgo。

cgo 可以直接在 Go 源代码里写 c 代码 也可以直接把 c 代码放到 go 代码里面并打上注释，紧接着加一个 import "c"就可以了

attribute constructor/destructor 若函数被设定为 constructor 属性，则该函数会在 main 函数执行之前被自动执行。 
拥有此类属性的函数经常饮食的用在程序的初始化数据方面。

18-2: 命令实现。

setns 是一个系统调用，可以根据提供的 PID 再次进入到指定的 Namespace 中。 
它需要打开 /proc/[pid]/ns 文件夹下对应的文件，然后使当前进程进入到指定的 Namespace 中。 
系统调用描述非常简单，但是有一点对于 Go 来说很麻烦。 
对于 Mount Namespace 来说，一个具有多线程的进程无法使用 setns 调用进入到对应的命名空间的。 
但是 Go 没启动一个程序就会进入多线程状态，因此无法简单的在 Go 中直接调用系统调用，使当前的进程进入对应的 Mount Namespace 这里需要借助 C 来实现这个功能。

根据容器名去 /var/run/mydocker/容器名/config.json 中找 pid
从命令行中获取 exec 需要执行的命令 command
执行 c 代码进入 pid 中的 Namespace 并调用 command 执行，最终 command fork 出来的进程会与 pid 拥有一样的 Namespace（mnt ipc uts net pid）




## branch 19:

在 branch18 的基础上，实现停止容器。

1. 找到容器的 pid
2. 删除该 pid 对应的进程
3. 将该容器的 metadata 中的状态设置为 stopped


tags: 增加了记录 log 的行数信息。



## branch 20:

在 branch19 的基础上，实现删除容器。

实现删除一个容器就是把这个容器对应的文件夹删除， 这个容器必须是一个已经被停止的容器； 
因为正常逻辑是先停止运行的容器，然后再能删除。



## branch 21:

在 branch20 的基础上，实现容器层隔离。

由于 volume 的实现会让多个容器共用容器层，这样会导致多个容器层之间的数据不隔离，互相可以访问、修改数据。

tags:

21-1: 实现容器层隔离。

21-2: 实现镜像参数化。

之前起的容器都是在代码里写的 busybox，所以这里是将该部分参数化，就是在命令行中指定使用什么镜像
但是 /rootPath 必须要有该镜像的压缩包（比如 busybox.tar）



## branch 22:

在 branch21 的基础上，实现通过容器制作镜像。

tags:

22-1: branch 21存在问题：在以 -d 方式运行容器创建进程时无法将 mount 进来的文件持久化。

（当 -d 运行起来后，再 exec 进入容器时候看不到挂载的 containerVolume）

原因：-d 运行后在父进程退出时会调用 ClearWorkDir，所以再次进入到容器后看不到 mount 的文件。

解决：

在 Run 方法中判断如果是以后台运行的形式启动容器则不调用 ClearWorkDir 方法，在 stop 命令的时候需要调用 ClearWorkDir 方法。
如果是以 tty 形式启动容器则在子进程运行结束后应该调用 ClearWorkDir 方法。


22-2: 通过容器制作镜像。 镜像其实就是由一些文件组成，因此直接将运行中的容器打包就可以组成一个新的镜像。


---

EOF.
>>>>>>> c66c859454d3c0542e2d13fc80fa1e711bd10008
