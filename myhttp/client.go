package myhttp

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 请求 抄了http包的，
type ReqMini struct {
	Method        string
	URL           *url.URL
	Host          string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Header        Header
	Body          io.ReadWriter //原是read close
	ContentLength int64
}

// 主体
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
type Header map[string][]string

func NewReqMini() *ReqMini {
	return &ReqMini{
		Method:        "",
		Host:          "",
		Proto:         "",
		ProtoMajor:    0,
		ProtoMinor:    0,
		Header:        nil,
		Body:          nil,
		ContentLength: 0,
	}
}
func (req *ReqMini) SetBodyString(body string) {
	req.Body = bytes.NewBufferString("")
	length, _ := req.Body.Write([]byte(body))
	req.ContentLength = int64(length)
}
func (req *ReqMini) SetBodyByte(body []byte) {
	req.Body = bytes.NewBufferString("")
	length, _ := req.Body.Write(body)
	req.ContentLength = int64(length)
}
func (req *ReqMini) SetProto(proto string) error {
	_, protomajor, protominor, err := parseProto(proto)
	if err != nil {
		return err
	}
	req.Proto = proto
	req.ProtoMajor = protomajor
	req.ProtoMinor = protominor
	return nil
}
func parseProto(proto string) (protoName string, protomajor int, protominor int, err error) {

	pro := strings.Split(proto, "/")
	if len(pro) != 2 {
		return "", 0, 0, errors.New("格式错误1")
	}
	majorMinor := strings.Split(pro[1], ".")
	if len(majorMinor) != 2 {
		return "", 0, 0, errors.New("格式错误2")
	}
	major, err := strconv.Atoi(majorMinor[0])
	if err != nil {
		return "", 0, 0, errors.New("格式错误3")
	}
	minor, err := strconv.Atoi(majorMinor[1])
	if err != nil {
		return "", 0, 0, errors.New("格式错误4")
	}
	return pro[0], major, minor, nil
}
func (req *ReqMini) SetMethod(method string) error {
	if ismethodcontain(method) {
		req.Method = method
		return nil
	}
	return errors.New("未知方法")
}
func ismethodcontain(method string) bool {
	var methods = []string{"GET", "POST", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	for _, v := range methods {
		if method == v {
			return true
		}
	}
	return false
}
func (req *ReqMini) SetHost(host string) {
	req.URL.Host = host
}
func (req *ReqMini) SetURL(urlstr string) error {
	urlParse, err := url.Parse(urlstr)
	if err != nil {
		log.Println(err)
		return err
	}
	req.URL = urlParse
	req.AddHeader("Host", []string{req.URL.Host})
	return nil
}
func (req *ReqMini) AddHeader(key string, val []string) {
	if req.Header == nil {
		req.Header = map[string][]string{}
	}
	for i := 0; i < len(val); i++ {
		req.Header.Add(key, val[i])
	}
}
func (request ReqMini) AntiParse() (string, error) {
	if request.URL == nil || request.Method == "" || request.Proto == "" {
		return "", errors.New("NO HOST")
	}
	if request.URL.Path == "" {
		request.URL.Path = "/"
	}
	reqLine := request.Method + " " + request.URL.Path + " " + request.Proto
	var reqHead string
	if len(request.Header) == 0 {
		reqHead = ""
	} else {
		for key, val := range request.Header {
			var (
				i int
				v = ""
			)
			for i, _ = range val {
				if i+1 == len(val) {
					v += val[i]
				} else {
					v += val[i] + ","
				}

			}
			reqHead += key + ":" + v + "\r\n"
		}
	}

	blank := "\r\n"
	var body string = ""
	buf := make([]byte, 1024)
	if request.Body != nil {
		for {
			_, err := request.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				// 处理错误
			}
			// 使用读取到的数据
			body += string(buf)
		}
		blank = "\r\n"
	}
	var req string
	req = reqLine + "\r\n" + reqHead + blank + body
	return req, nil
}
func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

// Set sets the header entries associated with key to the
// single element value. It replaces any existing values
// associated with key. The key is case insensitive; it is
// canonicalized by textproto.CanonicalMIMEHeaderKey.
// To use non-canonical keys, assign to the map directly.
func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

// Get gets the first value associated with the given key. If
// there are no values associated with the key, Get returns "".
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is
// used to canonicalize the provided key. Get assumes that all
// keys are stored in canonical form. To use non-canonical keys,
// access the map directly.
func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

// Values returns all values associated with the given key.
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is
// used to canonicalize the provided key. To use non-canonical
// keys, access the map directly.
// The returned slice is not a copy.
func (h Header) Values(key string) []string {
	return textproto.MIMEHeader(h).Values(key)
}

func (resp *Response) ReadResponse(buf []byte) {
	split := bytes.Split(buf, []byte{'\r', '\n'})
	resp.ParseLine(split[0])
	var (
		i    int
		v    []byte
		body = []byte{}
	)
	for i, v = range split[1:] {
		if len(v) == 0 {
			break
		}
		resp.ParseHeader(v)
	}
	for i += 1; i < len(split); i++ {
		body = bytes.Join([][]byte{body, split[i]}, []byte{})
	}
	resp.ParseBody(body)
}
func (resp *Response) ParseHeader(line []byte) {
	header := bytes.SplitN(line, []byte{':'}, 2)
	if len(header) == 1 {
		return
	}
	header[1] = bytes.TrimSpace(header[1])
	resp.Header.Add(string(header[0]), string(header[1]))
}
func (resp *Response) ParseLine(i []byte) {
	split := bytes.Split(i, []byte{' '})
	resp.Proto = string(split[0])
	proto := bytes.Split(split[0], []byte{'/'})
	majorMinor := bytes.Split(proto[1], []byte{'.'})
	major, err := strconv.Atoi(string(majorMinor[0]))
	if err != nil {
		resp.ProtoMajor = 1
	}
	minor, err := strconv.Atoi(string(majorMinor[1]))
	if err != nil {
		resp.ProtoMinor = 0
	}
	resp.ProtoMajor = major
	resp.ProtoMinor = minor
	statecode, _ := strconv.Atoi(string(split[1]))
	resp.StatusCode = statecode
	resp.Status = string(split[2])
}
func (resp *Response) ParseBody(line []byte) {
	resp.ContentLength = int64(len(line))
	reader := bytes.NewBuffer(line)
	readCloser := io.NopCloser(reader)
	resp.Body = readCloser
}
