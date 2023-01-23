package computer

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

func StartServer(world World, address string) error {
	rpc.Register(&world)
	rpc.HandleHTTP()
	// address ":1993"
	port, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("connection error %v", err)
	}
	go http.Serve(port, nil)
	return nil
}

func RetrieveRemoteState(world World, address string) (int8, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return -1, fmt.Errorf("connection error %v", err)
	}
	var reply int8

	client.Call("World.GetRemoteState", !world.right, &reply)
	client.Close()
	return reply, nil
}
