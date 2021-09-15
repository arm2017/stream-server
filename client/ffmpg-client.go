package client

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	api "stream-cam-api/stream-camera/api"

	"gocv.io/x/gocv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

func captureImg() {
	// ffmpeg -f v4l2 -i /dev/video0 -vframes 1 test.jpeg
	cmd := exec.Command("ffmpeg", "-f", "v4l2", "-i", "/dev/video0", "-vframes", "1", "test.jpeg", "-y")
	err := cmd.Wait()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stdout, _ := cmd.Output()

	// Print the output
	fmt.Println(string(stdout))
}

func Run2() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("connect: %v\n", address)
	defer conn.Close()
	c := api.NewStreamCameServiceClient(conn)

	stream, err := c.Streaming(context.Background(), grpc.UseCompressor(gzip.Name))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// loop := 0
	img := gocv.NewMat()
	defer img.Close()

	captureImg()

	stream.CloseSend()

	// for {

	// 	timenow := time.Now()
	// 	timef := timenow.UnixMilli()
	// 	// timetxt := timenow.Format(time.RFC3339Nano)

	// 	sendErr := stream.Send(&api.CameReq{
	// 		CameId:    "1234",
	// 		TimeFrame: timef,
	// 		Img:       nil,
	// 	})
	// 	if sendErr != nil {
	// 		log.Fatalln(sendErr)
	// 	}
	// 	loop = loop + 1
	// 	fmt.Printf("Send... : %v , Loop : %v\n", timef, (loop))
	// 	time.Sleep(100 * time.Microsecond)
	// }
}
