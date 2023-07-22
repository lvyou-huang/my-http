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
# 借鉴
想newbing问了问题
