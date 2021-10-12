package moveregister

import (
	"context"
	"fmt"
	"log"
	rpigpio "stream-cam-api/rpi-gpio"
	api "stream-cam-api/stream-camera/api"
	"time"

	"google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc"
)

const (
	Address = "172.20.10.14:50051"
	// address = "0.tcp.ap.ngrok.io:18664"
)

var bord *rpigpio.GpioBoard

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

		if move.Direction == "W" {
			bord.MoveW()
		} else if move.Direction == "A" {
			bord.MoveA()
		} else if move.Direction == "D" {
			bord.MoveD()
		} else if move.Direction == "S" {
			bord.MoveS()
		}

		time.Sleep(10 * time.Microsecond)
	}
}

func Run() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("connect: %v\n", Address)
	defer conn.Close()

	//gpio
	bord = rpigpio.Setup()
	fmt.Println("Gpio is OK.")
	defer bord.Close()

	c := api.NewStreamCameServiceClient(conn)
	//stream to server
	MoveRegisterToServer(&c)
	fmt.Println("end")
}
