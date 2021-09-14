package viewer

import (
	"context"
	"fmt"
	"io"
	"log"
	api "stream-cam-api/stream-camera/api"

	"gocv.io/x/gocv"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func Run() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewStreamCameServiceClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	stream, err := c.View(context.Background(), &api.VeiwReq{
		CameId: "1234",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	//view
	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("End.. Recv")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		mainErr := stream.Context().Err()
		if mainErr != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Recive %v , Time : %v \n", req.CameId, req.TimeFrame)
		jpg, jpgErr := gocv.IMDecode(req.Img, gocv.IMReadAnyColor)
		defer jpg.Close()

		if jpgErr != nil {
			fmt.Println("Oopsie daisy. I made a boo boo!", err)
			continue
		}
		defer jpg.Close()
		window.IMShow(jpg)
		window.WaitKey(1)
	}
}
