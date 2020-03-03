branch 19:

在 branch18 的基础上，实现停止容器。

> 1. 找到容器的 pid
> 2. 删除该 pid 对应的进程
> 3. 将该容器的 metadata 中的状态设置为 stopped


tags:
增加了记录 log 的行数信息。

---

EOF
