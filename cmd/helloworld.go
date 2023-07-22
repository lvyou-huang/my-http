package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"k8s/myhttp"
	"log"
	"net"
	"os"
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
				frame := ParseFrame(buf[:n])
				fmt.Println(frame)
				fmt.Println(string(frame.Payload))
			}
		}()
		for {
			Do(conn)
		}
	}
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
