package computer

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

func StartServer(world World, address string) error {
	remote_err := make(chan error, 1)
	err := rpc.Register(&world)
	if err != nil {
		return fmt.Errorf("connection error %v", err)
	}
	rpc.HandleHTTP()
	// address ":1993"
	port, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("connection error %v", err)
	}
	go func() {
		remote_err <- http.Serve(port, nil)
	}()
	//go http.Serve(port, nil)
	return nil
}

func RetrieveRemoteState(world World, address string) (int8, error) {
	var err error
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return -1, fmt.Errorf("connection error %v", err)
	}
	var reply int8

	err = client.Call("World.GetRemoteState", !world.right, &reply)
	client.Close()
	if err != nil {
		return -1, fmt.Errorf("connection error %v", err)
	}
	return reply, nil
}
