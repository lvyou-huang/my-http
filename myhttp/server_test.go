package myhttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"reflect"
	"testing"
)

/*func TestNewRequestWithContext(t *testing.T) {
	type args struct {
		ctx    context.Context
		method string
		url    string
		body   io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Request
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				ctx:    context.Background(),
				method: "GET",
				url:    "",
				body:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRequestWithContext(tt.args.ctx, tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequestWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequestWithContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}*/

func TestRequest_ReadRequest(t *testing.T) {
	type fields struct {
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
			name:   "1",
			fields: fields{},
			args:   args{buf: []byte("POST / HTTP/1.0\r\nHost: example.com\r\nContent-Length: 13\r\n\r\nhello,world!")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest()
			req.ReadRequest(tt.args.buf)
			fmt.Printf("%+v", req)
		})
	}
}

func TestNewResponse(t *testing.T) {
	tests := []struct {
		name string
		want *Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processpath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{path: "huangyijian/nmsl/hhh"},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processpath(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processpath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEngine(t *testing.T) {
	type args struct {
		network string
		address string
	}
	tests := []struct {
		name string
		args args
		want *Engine
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				network: "tcp",
				address: ":8080",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEngine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_POST(t *testing.T) {

	type fields struct {
		Listener       net.Listener
		MethodHandlers []MethodHandler
		MethodRouters  []MethodRouter
		Close          bool
	}
	type args struct {
		address string
		routers []Router
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				address: "/hello/huangyijian",
				routers: []Router{func(request *Request, response *Response, conn net.Conn) {
					fmt.Println(request.Proto)
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine()
			engine.POST(tt.args.address, tt.args.routers...)
		})
	}
}
