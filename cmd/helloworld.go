package main

import (
	"fmt"
	"k8s/myhttp"
	"log"
	"net"
)

// 头部字段

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	request := myhttp.NewReqMini()
	err = request.SetMethod("GET")
	if err != nil {
		log.Println(err)
		return
	}
	err = request.SetProto("HTTP/1.0")
	if err != nil {
		log.Println(err)
		return
	}
	err = request.SetURL("http://localhost:8080/huangyijian/hello")
	fmt.Println(request.URL)
	if err != nil {
		log.Println(err)
		return
	}
	var SecWebSocketKey = "hj0eNqbhE/A0GkBXDRrYYw=="
	request.Header.Add("Upgrade", "websocket")
	request.Header.Add("Connection", "Upgrade")
	request.Header.Add("Sec-WebSocket-Key", SecWebSocketKey)
	req, err := request.AntiParse()
	if err != nil {
		log.Println(err)
		return
	}
	//fmt.Fprintf(conn, "POST / HTTP/1.0\r\nHost: example.com\r\nContent-Length: %d\r\n\r\n%s", len(message), message)

	fmt.Fprintf(conn, req)

	buf := make([]byte, 1024)
	//for {
	_, err = conn.Read(buf)

	response := myhttp.NewResponse()
	response.ReadResponse(buf)
	fmt.Printf("%+v", response)
	accept := response.Header.Get("Sec-Websocket-Accept")
	if accept == myhttp.GenerateWebSocketAccept(SecWebSocketKey) {
		go func() {
			for {
				n, err := conn.Read(buf)
				if n == 0 || err != nil {
					break
				}
				fmt.Println(buf[:n])
				frame := myhttp.ParseFrame(buf[:n])
				fmt.Println(frame)
				fmt.Println(string(frame.Payload))
			}
		}()
		for {
			myhttp.Do(conn)
		}
	}
}
