branch 12:

在 branch11 的基础上，一步步实现使用AUFS包装busybox。

> branch 11 的存在问题：
> 利用 busybox 创建的容器, 创建文件夹并且创建文件。
> 退出容器后, 查看宿主机的内容, 返现内容在宿主机中也存在。
> 这样会有一个问题, 其实 busybox 就是容器的镜像层, 如果多个容器共享该镜像层, 那就会造成容器之间互相看到对方文件, 并且文件覆盖等等问题。
> branch12 就是利用 AUFS 解决此问题
>
> https://www.jianshu.com/p/ecbdcc98db76

---

tags:

12:

1. 根据 busybox 镜像生成容器，其实就是解压 busybox.tar 生成 rootPath/busybox。
2. 创建挂载点、可写层，并将 writeLayer 和 busybox 挂载到 mnt 下。
3. 当退出时候执行 umount mnt/ 并删除 writeLayer。

---

EOF
