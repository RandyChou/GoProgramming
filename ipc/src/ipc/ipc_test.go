package ipc

import (
	"testing"
)

type EchoServer struct {
}

func (server *EchoServer) Handle(method, params string) *Response {
	return &Response{method, params} 
}

func (server *EchoServer) Name() string {
	return "EchoServer"
}

func TestIpc(t *testing.T) {
	server := NewIpcServer(&EchoServer{})

	client1 := NewIpcClient(server)
	client2 := NewIpcClient(server)

	resp1, err1 := client1.Call("From Client1", "client1")
	resp2, err2 := client2.Call("From Client2", "client2")

	if err1 != nil || err2 != nil {
		t.Error("IpcClient.Call failed.")		
	}

	if resp1.Body != "client1" || resp1.Code != "From Client1" ||
		resp2.Body != "client2" || resp2.Code != "From Client2" {
			t.Error("IpcClient.Call failed. resp1:", resp1, "resp2:", resp2)
	}

	client1.Close()
	client2.Close()
}