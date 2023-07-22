package myhttp

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net"
)

func Websocket(req *Request, resp *Response, conn net.Conn) {
	if req.Method == "GET" && req.Header.Get("Upgrade") == "websocket" {
		OpenWebsocket(req, resp, conn)
	}
}
func OpenWebsocket(req *Request, resp *Response, conn net.Conn) {
	Key := req.Header.Get("Sec-WebSocket-Key")
	if Key == "" {
		panic("不能为空")
	}
	accept := GenerateWebSocketAccept(Key)
	resp.StatusCode = 101
	resp.Status = "Switching Protocols"
	resp.Header.Add("Upgrade", "websocket")
	resp.Header.Add("Sec-Websocket-Accept", accept)
	resp.Header.Add("Connection", " Upgrade")
	respunparse := resp.Unparse()
	fmt.Fprint(conn, respunparse)
}
func GenerateWebSocketAccept(key string) string {
	h := sha1.New()
	io.WriteString(h, key+"258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
