branch 22:

在 branch21 的基础上，实现通过容器制作镜像。


tags:

22-1:
branch 21存在问题：在以 -d 方式运行容器创建进程时无法将 mount 进来的文件持久化。

（当 -d 运行起来后，再 exec 进入容器时候看不到挂载的 containerVolume）

原因：-d 运行后在父进程退出时会调用 ClearWorkDir，所以再次进入到容器后看不到 mount 的文件。

解决：
1. 在 Run 方法中判断如果是以后台运行的形式启动容器则不调用 ClearWorkDir 方法，在 stop 命令的时候需要调用 ClearWorkDir 方法。
2. 如果是以 tty 形式启动容器则在子进程运行结束后应该调用 ClearWorkDir 方法。

22-2:
通过容器制作镜像。
镜像其实就是由一些文件组成，因此直接将运行中的容器打包就可以组成一个新的镜像。

---

EOF
