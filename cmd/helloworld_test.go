package main

import "testing"

/*import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

/*func TestHeader_Values(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		h    Header
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			h:    map[string][]string{"1": {"1", "2"}, "2": {"3", "4"}},
			args: args{key: "1"},
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Values(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAntiParse(t *testing.T) {
	type args struct {
		request Request
	}
	// 使用 buf 作为 io.ReadCloser 对象
	resp, err := http.Get("http://www.example.com")
	if err != nil {
		// 处理错误
	}
	defer resp.Body.Close()

	// 使用 resp.Body 作为 io.ReadCloser 对象

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{request: Request{
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
				Proto:            "HTTP/1.1",
				ProtoMajor:       0,
				ProtoMinor:       0,
				Header:           map[string][]string{"Host": {"huangyijian.com"}, "Connetion": {"Keep-alive"}},
				Body:             resp.Body,
				GetBody:          nil,
				ContentLength:    0,
				TransferEncoding: nil,
				Close:            false,
				Host:             "huangyijian.com",
				Form:             nil,
				PostForm:         nil,
				MultipartForm:    nil,
				Trailer:          nil,
				RemoteAddr:       "",
				RequestURI:       "",
				TLS:              nil,
				Cancel:           nil,
				Response:         nil,
				ctx:              nil,
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

func TestFrame_Fin1(t *testing.T) {
	type fields struct {
		Head       int
		PayloadLen int
		Payload    []int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				Head:       0b10101010,
				PayloadLen: 7,
				Payload:    []int{1, 2, 3},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := &Frame{
				Head:       tt.fields.Head,
				MaskAndLen: tt.fields.PayloadLen,
				Payload:    tt.fields.Payload,
			}
			if got := frame.Fin(); got != tt.want {
				t.Errorf("Fin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBit(t *testing.T) {
	type args struct {
		num int
		i   uint
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				num: 0b10101010,
				i:   8,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBit(tt.args.num, tt.args.i); got != tt.want {
				t.Errorf("getBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrame_Opcode(t *testing.T) {
	type fields struct {
		Head       int
		PayloadLen int
		Payload    []int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				Head:       0b10101010,
				PayloadLen: 0,
				Payload:    []int{1},
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := Frame{
				Head:       tt.fields.Head,
				MaskAndLen: tt.fields.PayloadLen,
				Payload:    tt.fields.Payload,
			}
			if got := frame.Opcode(); got != tt.want {
				t.Errorf("Opcode() = %v, want %v", got, tt.want)
			}
		})
	}
}
