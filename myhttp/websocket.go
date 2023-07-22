package myhttp

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
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

func ParseFrame(frame []byte) Frame {
	f := Frame{
		Head:       0,
		MaskAndLen: 0,
		Mask:       nil,
		Payload:    nil,
	}
	f.Head = int(frame[0])
	f.MaskAndLen = int(frame[1])
	f.Mask = frame[2:6]
	f.Payload = frame[6:]
	return f
}

func Do(conn net.Conn) {
	var input string
	fmt.Scan(&input)
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	frames := GenerateFrame(input)
	for _, f := range frames {
		/*h := byte(f.Head)
		mal := byte(f.MaskAndLen)
		m := f.Mask
		pl := f.Payload*/
		bytes := append([]byte{byte(f.Head), byte(f.MaskAndLen)}, f.Mask...)
		fm := append(bytes, f.Payload...)
		conn.Write(fm)
		//fmt.Fprint(conn, f.Head, f.MaskAndLen, f.Mask, f.Payload)
	}
}

type Frame struct {
	Head       int
	MaskAndLen int
	Mask       []byte
	Payload    []byte
}

// Mask masks the data with the given masking key.
func Mask(data []byte, mask []byte) {
	for i := 0; i < len(data); i++ {
		data[i] ^= mask[i%4]
	}
}

// Unmask unmasks the data with the given masking key.
func Unmask(data []byte, mask []byte) {
	Mask(data, mask)
}

func generateMask() []byte {
	mask := make([]byte, 4)
	rand.Read(mask)
	return mask
}

func GenerateFrame(s string) []Frame {
	var f Frame
	if len(s) <= 255 {
		f.Head = 0b10000001
		f.MaskAndLen = 128 + len(s)
		f.Mask = generateMask()
		f.Payload = []byte(s)
		return []Frame{f}
	}
	return nil
}
func (frame Frame) Fin() int {
	return getBit(frame.Head, 1)
}
func (frame Frame) RSV1() int {
	return getBit(frame.Head, 2)
}
func (frame Frame) RSV2() int {
	return getBit(frame.Head, 3)
}
func (frame Frame) RSV3() int {
	return getBit(frame.Head, 4)
}
func (frame Frame) Opcode() int {
	return frame.Head & 0x0F
}
func (frame Frame) IsMask() int {
	return getBit(frame.MaskAndLen, 1)
}
func (frame Frame) Payloadlen() int {
	return frame.Head & 0x7F
}

func getBit(num int, i uint) int {
	return (num >> i) & 1
}
