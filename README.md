branch 13:

在 branch12 的基础上，一步步实现 volume 操作。

> branch12 中的容器内增删文件都不会保存。
> 如果用户需要保存则需要 -v 参数把宿主机的目录挂载到容器内。

---

tags:

13-1:
实现单个 volume 挂载到容器中。注意 umount 的时候先 umount volume path，再 umount mnt/。

13-2:
实现多个 -v。


---

EOF
