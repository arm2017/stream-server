package moveregister

import (
	"context"
	"fmt"
	"log"
	gpioapi "stream-cam-api/rpi-gpio"
	api "stream-cam-api/stream-camera/api"
	"time"

	"google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc"
)

const (
	address = "172.20.10.14:50051"
	// address = "0.tcp.ap.ngrok.io:18664"
)

var bord *gpioapi.GpioBoard

func MoveRegisterToServer(streamClient *api.StreamCameServiceClient) {

	stream, err := (*streamClient).MoveRegister(context.Background(), &api.MoveRegisterReq{
		HwId:      "raspi-01",
		RegisTime: time.Now().UnixMilli(),
	}, grpc.UseCompressor(gzip.Name))

	if err != nil {
		log.Fatalf("MoveRegister : %v", err)
	}

	for {

		move, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Move : %v, Time : %v\n", move.Direction, move.TimeMove)

		time.Sleep(10 * time.Microsecond)
	}
}

func Run() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("connect: %v\n", address)
	defer conn.Close()

	//gpio
	bord = gpioapi.GpioSetup()
	fmt.Println(bord)

	c := api.NewStreamCameServiceClient(conn)
	//stream to server
	MoveRegisterToServer(&c)
	fmt.Println("end")
}
