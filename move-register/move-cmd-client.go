package moveregister

import (
	"context"
	"fmt"
	"log"
	api "stream-cam-api/stream-camera/api"
	"time"

	"google.golang.org/grpc"
)

func RunTest() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("connect: %v\n", address)
	defer conn.Close()
	c := api.NewStreamCameServiceClient(conn)
	//stream to server
	rsp, err := c.Move(context.Background(), &api.MoveReq{
		HwId:      "raspi-01",
		Direction: "W",
		TimeMove:  time.Now().UnixMilli(),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Move : %v\n", rsp.Direction)

	fmt.Println("end")
}
