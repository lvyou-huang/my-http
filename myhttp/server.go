package myhttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 报文
type Response struct {
	Status        string
	StatusCode    int
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Header        Header
	Body          io.ReadCloser
	ContentLength int64
	Close         bool
	Uncompressed  bool
	Trailer       Header
	Request       *Request
	TLS           *tls.ConnectionState
}

func NewResponse() *Response {
	return &Response{
		Status:        "",
		StatusCode:    0,
		Proto:         "",
		ProtoMajor:    0,
		ProtoMinor:    0,
		Header:        make(Header),
		Body:          nil,
		ContentLength: 0,
		Close:         false,
		Uncompressed:  false,
		Trailer:       nil,
		Request:       nil,
		TLS:           nil,
	}
}

type Request struct {
	Method        string
	URL           *url.URL
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Header        Header
	Body          io.ReadCloser
	GetBody       func() (io.ReadCloser, error)
	ContentLength int64
	Close         bool
	Trailer       Header
	RemoteAddr    string
	RequestURI    string
	TLS           *tls.ConnectionState
	Cancel        <-chan struct{}
	Response      *Response
	ctx           context.Context
}
type Client struct {
	Cookie   []byte
	Conn     net.Conn
	DeadLine time.Time
	Reqs     []*Request
}
type Handlers []Handler
type Handler func(request *Request, response *Response)
type Router func(request *Request, response *Response, conn net.Conn)
type MethodHandler struct {
	Address string
	Handler Handler
}
type MethodRouter struct {
	Method  string
	Address string
	Router  Router
}
type Engine struct {
	MethodHandlers []MethodHandler
	MethodRouters  []MethodRouter
	Close          bool
}

var Nobody noBody

type noBody struct{}

func (noBody) Read([]byte) (int, error)         { return 0, io.EOF }
func (noBody) Close() error                     { return nil }
func (noBody) WriteTo(io.Writer) (int64, error) { return 0, nil }

func NewEngine() *Engine {

	engine := new(Engine)
	engine.Close = false
	engine.MethodHandlers = make([]MethodHandler, 1)
	return engine
}

func (handler Handler) AddRouter(engine *Engine, path string) {
	handlers := append(engine.MethodHandlers, MethodHandler{
		Address: path,
		Handler: handler,
	})
	engine.MethodHandlers = handlers
}
func (engine *Engine) Run(network string, address string) {
	listener, err := net.Listen(network, address)
	defer listener.Close()
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, engine.MethodHandlers, engine.MethodRouters)
	}
}

func (engine *Engine) POST(address string, routers ...Router) {
	for _, r := range routers {
		r.Add(engine, "POST", address)
	}
}
func (router Router) Add(engine *Engine, method string, path string) {
	routers := append(engine.MethodRouters, MethodRouter{
		Method:  method,
		Address: path,
		Router:  router,
	})
	engine.MethodRouters = routers
}
func handleConnection(conn net.Conn, handlers []MethodHandler, routers []MethodRouter) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		panic("错误")
	}
	req := NewRequest()
	resp := NewResponse()
	req.Response = resp
	resp.Request = req
	resp.Proto = req.Proto
	resp.ProtoMajor = req.ProtoMajor
	resp.ProtoMinor = req.ProtoMinor
	req.ReadRequest(buf)
	path := req.URL.Path
	relativepaths := processpath(path)
	for _, path := range relativepaths {
		for _, h := range handlers {
			if h.Address == path {
				h.Handler(req, resp)
			}
		}
	}
	for _, r := range routers {
		if path == r.Address && req.Method == r.Method {
			r.Router(req, resp, conn)
		}
	}
}
func processpath(path string) []string {
	split := strings.Split(path, "/")
	for i, _ := range split {
		lpath := ""
		for j := 0; j <= i; j++ {
			lpath += split[j] + "/"
		}
		split[i] = lpath[:len(lpath)-1]
	}
	return split
}
func Send(conn net.Conn, resp *Response) error {
	unparse := resp.Unparse()
	_, err := fmt.Fprint(conn, unparse)
	if err != nil {
		return err
	}
	return nil
}
func (resp *Response) Unparse() string {
	var body []byte
	if resp.Body != nil {
		length, err := resp.Body.Read(body)
		if err != nil {
			panic("错误")
		}
		resp.ContentLength = int64(length)
	}
	var respline, respheader, blank, respbody, respstring = "", "", "", "", ""
	respline = resp.Proto + " " + strconv.Itoa(resp.StatusCode) + resp.Status
	for k, v := range resp.Header {
		value := ""
		for i, val := range v {
			value += val
			if i != len(v)-1 {
				value += ","
			}
		}
		respheader += k + ":" + value + "\r\n"
	}
	if resp.ContentLength != 0 {
		respbody = string(body)
		blank = "\r\n"
	}
	respstring = respline + "\r\n" + respheader + blank + respbody
	return respstring
}

func main() {
	/*engine := NewEngine("tcp", ":8080")
	engine.Get("/", func(request *Request, response *Response) {

	})*/
}

/*
	func main() {
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}

			go handleConnection(conn)
		}
	}
*/

func NewRequest() *Request {
	return &Request{
		Method:        "",
		URL:           nil,
		Proto:         "",
		ProtoMajor:    0,
		ProtoMinor:    0,
		Header:        make(Header),
		Body:          nil,
		GetBody:       nil,
		ContentLength: 0,
		Close:         false,
		Trailer:       nil,
		RemoteAddr:    "",
		RequestURI:    "",
		TLS:           nil,
		Cancel:        nil,
		Response:      nil,
		ctx:           context.Background(),
	}
}

func (req *Request) ReadRequest(buf []byte) {
	split := bytes.Split(buf, []byte{'\r', '\n'})
	req.ParseLine(split[0])
	var (
		i    int
		v    []byte
		body = []byte{}
	)
	for i, v = range split[1:] {
		if len(v) == 0 {
			break
		}
		req.ParseHeader(v)
	}
	for i += 1; i < len(split); i++ {
		body = bytes.Join([][]byte{body, split[i]}, []byte{})
	}
	req.ParseBody(body)
}

func (req *Request) ParseLine(line []byte) {
	split := bytes.Split(line, []byte{' '})
	req.Method = string(split[0])
	parseurl, err := url.Parse(string(split[1]))
	if err != nil {
		panic("请求行解析错误")
	}
	req.URL = parseurl
	req.Proto = string(split[2])
	proto := bytes.Split(split[2], []byte{'/'})
	majorMinor := bytes.Split(proto[1], []byte{'.'})
	major, err := strconv.Atoi(string(majorMinor[0]))
	if err != nil {
		req.ProtoMajor = 1
	}
	minor, err := strconv.Atoi(string(majorMinor[1]))
	if err != nil {
		req.ProtoMinor = 0
	}
	req.ProtoMajor = major
	req.ProtoMinor = minor

}
func (req *Request) ParseHeader(line []byte) {
	header := bytes.Split(line, []byte{':'})
	header[1] = bytes.TrimSpace(header[1])
	req.Header.Add(string(header[0]), string(header[1]))
}
func (req *Request) ParseBody(line []byte) {
	req.ContentLength = int64(len(line))
	reader := bytes.NewBuffer(line)
	readCloser := io.NopCloser(reader)
	req.Body = readCloser
}
