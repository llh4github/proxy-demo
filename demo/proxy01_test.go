package demo

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

type proxyHandler struct {
}

func (proxyHandler) ServeHTTP(write http.ResponseWriter, req *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			write.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	}()
	url := req.URL.Path
	if url == "/a" {
		nreq, _ := http.NewRequest(req.Method, "http://localhost:9091", req.Body)
		nresp, _ := http.DefaultClient.Do(nreq)
		defer nresp.Body.Close()
		content, _ := ioutil.ReadAll(nresp.Body)
		write.Write(content)
		return
	}
	write.Write([]byte("default"))
}

// 测试案例服务器启动情况
func TestDemoServeStatus(t *testing.T) {
	StartDemoServe()
	time.Sleep(time.Second)
	// 第一个服务器
	resp, err := http.Get("http://localhost:9091")
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "web1", string(body))
	// 第二个服务器
	resp, err = http.Get("http://localhost:9092")
	body, err = ioutil.ReadAll(resp.Body)
	assert.Equal(t, "web2", string(body))
	if err != nil {
		t.Log(err)
	}
	t.Log(string(body))
}

// 测试反向代理最简单的功能
func TestProxyDemo01(t *testing.T) {
	go http.ListenAndServe(":8080", proxyHandler{})
	StartDemoServe()
	time.Sleep(time.Second)
	// 是否代理到9091上
	resp, _ := http.Get("http://localhost:8080/a")
	body, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(body))
	assert.Equal(t, "web1", string(body))

	// 没有配置代理的URL
	resp, _ = http.Get("http://localhost:8080/ccc")
	body, _ = ioutil.ReadAll(resp.Body)
	t.Log(string(body))
	assert.Equal(t, "default", string(body))

}
