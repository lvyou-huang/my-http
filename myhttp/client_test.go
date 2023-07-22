package myhttp

import (
	"crypto/tls"
	"fmt"
	"io"
	"testing"
)

/*func TestAntiParse(t *testing.T) {
	type args struct {
		request ReqMini
	}
	var buf bytes.Buffer
	buf.WriteString("123:huangyijian")
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{request: ReqMini{
				Method: "GET",
				URL: &url.URL{
					Scheme:      "",
					Opaque:      "",
					User:        &url.Userinfo{},
					Host:        "",
					Path:        "/huangyijian/nmsl",
					RawPath:     "",
					OmitHost:    false,
					ForceQuery:  false,
					RawQuery:    "",
					Fragment:    "",
					RawFragment: "",
				},
				Proto:         "HTTP/1.1",
				ProtoMajor:    0,
				ProtoMinor:    0,
				Header:        map[string][]string{"Host": {"huangyijian.com"}, "Connetion": {"Keep-alive"}},
				Body:          &buf,
				ContentLength: 0,
			}},
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AntiParse(tt.args.request); got != tt.want {
				t.Errorf("AntiParse() = %v, want %v", got, tt.want)
			}
		})
	}
}*/

func TestResponse_ReadRequest(t *testing.T) {
	type fields struct {
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
	type args struct {
		buf []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{buf: []byte("HTTP/1.0 200 ok\r\nHuangyijian:nibaba\r\nNihao:hhh\r\n\r\nnihao hello\r\n")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewResponse()
			resp.ReadResponse(tt.args.buf)
			fmt.Printf("%+v", resp)
		})
	}
}
