package webserver

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRunHTTPServerPortConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 先在一个随机端口上建立监听
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("failed to split host port: %v", err)
	}
	port, _ := strconv.Atoi(portStr)

	eng := gin.New()
	s := &Server{
		Address: "127.0.0.1",
		Port:    port,
		Proto:   "http",
	}

	// 此时调用 RunHTTPServer 应该立即返回端口冲突错误，而不是无声无息或卡住
	err = RunHTTPServer(eng, s)
	if err == nil {
		t.Fatalf("expected error due to port conflict, got nil")
	}
	t.Logf("successfully caught port conflict error: %v", err)
}

func TestServiceRunAndClose(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 获取两个未被占用的端口
	ln1, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	port1 := ln1.Addr().(*net.TCPAddr).Port
	ln1.Close()

	s := New()
	gConf.Server = []Server{
		{
			Address: "127.0.0.1",
			Port:    port1,
			Proto:   "http",
		},
	}

	err = s.Run()
	if err != nil {
		t.Fatalf("service run failed: %v", err)
	}

	// 等待启动完成
	time.Sleep(50 * time.Millisecond)

	// 测试 Close 是否能正常优雅关闭
	s.Close()

	// 再次监听该端口，验证端口已被释放
	lnTest, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port1)))
	if err != nil {
		t.Fatalf("expected port %d to be released, but listen failed: %v", port1, err)
	}
	lnTest.Close()
}
