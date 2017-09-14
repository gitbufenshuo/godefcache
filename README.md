# -- godef doesnt support cache, so big package searching is slow.
# -- this is only a cache tool based on godef. so you should install godef first.
# -- godefcache will use mongodb, so mongodb is needed.

----

    - [install] go get github.com/gitbufenshuo/godefcache

    - [start mongo] google this or you've already done this

    - [replace step_1] $GOPATH/bin: mv godef godef_raw // any good name you like

    - [replace step_2] $: godefcache -s godef_raw    // the previous good name you like

    - [replace step_3] $GOPATH/bin: mv godefcache godef // now (godef) is (godefcache) and (godef_raw) is the (real godef)....

    - [tested in vscode] runs well

    - [CAUTION] remember (the previous replace steps) when hack this code. And remember the mongodb。
----
----
----
----

# -- godef 不支持缓存，所以大型 package 搜索起来很慢。
# -- godefcache 基于 godef，所以得先装 godef 。
# -- godefcache 用 mongodb 缓存数据，所以得先装 mongodb 。

----

    - [安装] go get github.com/gitbufenshuo/godefcache

    - [启动 mongo] google 这个，或者你已经启动了。

    - [替换步骤_1] 进入$GOPATH/bin目录: mv godef godef_raw // 换成任何你喜欢的名字

    - [替换步骤_2] $: godefcache -s godef_raw    // 上面的你喜欢的名字

    - [替换步骤_3] 进入$GOPATH/bin目录: mv godefcache godef // 现在 godef 是 godefcache，而 godef_raw 才是真正的 godef 。

    - [在 vscode 中测试过] 跑起来很好

    - [注意] 修改 godefcache 代码之后，记得上面的替换步骤，还有 mongodb 记得开。
