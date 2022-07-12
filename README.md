gprc在注册时需要开启反射服务，不然将无法使用，golang开启反射如下:

```golang
func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // 创建 RPC 服务容器
    grpcServer := grpc.NewServer()

    // 为 User 服务注册业务实现 将 User 服务绑定到 RPC 服务容器上
    user.RegisterUserServer(grpcServer, &UserService{})
    // 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系

    reflection.Register(grpcServer)

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

  go-zero框架只需要在etc文件里加入 `Mode: dev` 就能开启反射



**添加grpc的地址，会得到服务名字和其方法**

![](https://github.com/hxoreyer/grpc-tool/png/1.png)

**点击方法将会在request区域显示其参数**

![](https://github.com/hxoreyer/grpc-tool/png/2.png)

**中间按钮点击即可发送请求**

![](https://github.com/hxoreyer/grpc-tool/png/3.png)

> 工具还不完善，大家发现bug啥的可以留言哈
>
> 下载地址:  https://github.com/hXoreyer/grpc-tool/releases/tag/v1.1
>
> 视频演示:  https://www.bilibili.com/video/BV1JZ4y1Y7yr/

