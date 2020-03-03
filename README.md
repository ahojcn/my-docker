branch 18:

在 branch17 的基础上，实现进入容器 Namespace。


tags:

18-1:
test cgo

> cgo
> 可以直接在 Go 源代码里写 c 代码
> 也可以直接把 c 代码放到 go 代码里面并打上注释，紧接着加一个 `import "c"`就可以了
>
> __attribute__ constructor/destructor 若函数被设定为 constructor 属性，则该函数会在 main 函数执行之前被自动执行。
> 拥有此类属性的函数经常饮食的用在程序的初始化数据方面。

---

EOF
