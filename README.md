branch 18:

在 branch17 的基础上，实现进入容器 Namespace。


tags:

18-1:
test cgo。

> cgo
> 可以直接在 Go 源代码里写 c 代码
> 也可以直接把 c 代码放到 go 代码里面并打上注释，紧接着加一个 `import "c"`就可以了
>
> __attribute__ constructor/destructor 若函数被设定为 constructor 属性，则该函数会在 main 函数执行之前被自动执行。
> 拥有此类属性的函数经常饮食的用在程序的初始化数据方面。


18-2:
命令实现。

> `setns` 是一个系统调用，可以根据提供的 PID 再次进入到指定的 Namespace 中。
> 它需要打开 `/proc/[pid]/ns` 文件夹下对应的文件，然后使当前进程进入到指定的 Namespace 中。
> 系统调用描述非常简单，但是有一点对于 Go 来说很麻烦。
> 对于 Mount Namespace 来说，一个具有多线程的进程无法使用 setns 调用进入到对应的命名空间的。
> 但是 Go 没启动一个程序就会进入多线程状态，因此无法简单的在 Go 中直接调用系统调用，使当前的进程进入对应的 Mount Namespace
> 这里需要借助 C 来实现这个功能。
>
> 1. 根据容器名去 /var/run/mydocker/容器名/config.json 中找 pid
> 2. 从命令行中获取 exec 需要执行的命令 command
> 3. 执行 c 代码进入 pid 中的 Namespace 并调用 command 执行，最终 command fork 出来的进程会与 pid 拥有一样的 Namespace（mnt ipc uts net pid）
>

---

EOF
