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

func TestDemo2(t *testing.T) {
	StartDemoServe()
	time.Sleep(time.Second)
	resp, err := http.Get("http://localhost:9091")
	if err != nil {
		t.Log(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(body))
}
func TestDemo01(t *testing.T) {
	go http.ListenAndServe(":8080", proxyHandler{})
	StartDemoServe()

	time.Sleep(time.Second)
	resp, _ := http.Get("http://localhost:8080/a")
	body, _ := ioutil.ReadAll(resp.Body)

	t.Log(string(body))
	assert.Equal(t, "web1", string(body))

	resp, _ = http.Get("http://localhost:8080/ccc")
	body, _ = ioutil.ReadAll(resp.Body)

	t.Log(string(body))
	assert.Equal(t, "default", string(body))

}
