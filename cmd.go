package main

import (
	"awesomeProject1/myhttp"
	"bufio"
	"bytes"
	"io"
)

func main() {
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
}

func SetBody(response *myhttp.Response, s string) {
	response.ContentLength = int64(len(s))
	bufferString := bytes.NewBufferString(s)
	reader := bufio.NewReader(bufferString)
	readCloser := io.NopCloser(reader)
	response.Body = readCloser
}
