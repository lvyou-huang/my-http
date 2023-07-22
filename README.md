# my-http
# 快速开始
- go get
- 在一个窗口运行cmd/helloworld.go
- 在另一个窗口运行cmd.go
- 然后就可以快乐的使用websocket聊天了（本地）（没有精美的界面十分丑陋，但是功能是有的）
# 架构
client是用户端，server是服务端，websocket专门处理websocket
由于开始写的时候架构没分清，在client可能有sever的内容，sever可能有client的内容，考虑到代码的稳定性，还是不改了
在基本使用方式仿照gin，但是没有看gin的代码，绝对gin写起来好看就用了他的样子。
- 服务端
  engine := myhttp.NewEngine()

	/*engine.POST("/huangyijian/hello", func(request *myhttp.Request, response *myhttp.Response, conn net.Conn) {
		response.Status = "ok"
		response.StatusCode = 200
		fmt.Printf("%+v\n", request)
		fmt.Printf("%+v\n", response)
		response.Header.Add("huangyijian", "nibaba")
		response.Header.Add("nihao", "hhh")
		fmt.Fprint(conn, response.Unparse())
	})*/
	engine.GET("/huangyijian/hello", myhttp.Websocket)
	engine.Run("tcp", ":8080")
- 客户端
  通过set设置request字段的值，在通过antiparse()方法返回反解析后的字符串，直接往conn里写
- 添加中间件
engine.handler.AddHandler(handler,address)
- 上下文
  在request结构里有context.Context,所有的方法都可以使用
- websocket
  定义了一个路由，可以通过这个路由实现websocket  
# 借鉴
想newbing问了问题
