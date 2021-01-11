package demo

import "net/http"

type web1 struct {
}

type web2 struct {
}

// ServeHTTP1 1
func (web1) ServeHTTP(write http.ResponseWriter, req *http.Request) {

	write.Write([]byte("web1"))
}

// ServeHTTP2 2
func (web2) ServeHTTP(write http.ResponseWriter, req *http.Request) {

	write.Write([]byte("web2"))
}

// StartDemoServe 开启测试服务器
// 端口号：9091，9092
func StartDemoServe() {
	go func() {
		http.ListenAndServe(":9091", web1{})
	}()

	go func() {
		http.ListenAndServe(":9092", web2{})
	}()
}
