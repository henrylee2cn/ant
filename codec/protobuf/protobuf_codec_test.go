package protobuf

import (
	"testing"
	"time"

	"github.com/henrylee2cn/rpc2"
	"github.com/henrylee2cn/rpc2/codec"
)

type ProtoArith int

func (t *ProtoArith) Mul(args *ProtoArgs, reply *ProtoReply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *ProtoArith) Error(args *ProtoArgs, reply *ProtoReply) error {
	panic("ERROR")
}

func TestProtobufCodec(t *testing.T) {
	// server
	server := rpc2.NewServer(60e9, 0, 0, NewProtobufServerCodec)
	group, _ := server.Group(codec.ServiceGroup)
	err := group.RegisterName(codec.ServiceName, new(ProtoArith))
	if err != nil {
		t.Fatal(err)
	}
	go server.ListenAndServe(codec.Network, codec.ServerAddr)
	time.Sleep(2e9)

	// client
	var args = &ProtoArgs{7, 8}
	var reply ProtoReply

	err = rpc2.
		NewDialer(codec.Network, codec.ServerAddr, NewProtobufClientCodec).
		Remote(func(client rpc2.IClient) error {
			return client.Call(codec.ServiceMethodName, args, &reply)
		})

	if err != nil {
		t.Errorf("error for Arith: %d*%d, %v \n", args.A, args.B, err)
	} else {
		t.Logf("Arith: %d*%d=%d \n", args.A, args.B, reply.C)
	}
}
