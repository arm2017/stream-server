package client

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	api "stream-cam-api/stream-camera/api"
	"time"

	"gocv.io/x/gocv"
	"google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc"
)

const (
	// address  = "localhost:50051"
	address  = "0.tcp.ap.ngrok.io:17368"
	deviceID = 0
)

func Run() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewStreamCameServiceClient(conn)

	//webcame
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	// webcam.Set(gocv.VideoCaptureFrameHeight, 1080)
	// webcam.Set(gocv.VideoCaptureFrameWidth, 1980)
	defer webcam.Close()
	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	stream, err := c.Streaming(context.Background(), grpc.UseCompressor(gzip.Name))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {

		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}
		// fmt.Printf(" X : %v, Y : %v\n", img.Cols(), img.Rows())

		timenow := time.Now()
		timef := timenow.UnixMilli()
		timetxt := timenow.Format(time.RFC3339Nano)
		gocv.PutText(&img, timetxt, image.Point{
			X: 10,
			Y: 50,
		}, gocv.FontHersheyComplex, 0.7, color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 1,
		}, 2)
		jpg, jpgerr := gocv.IMEncode(gocv.JPEGFileExt, img)
		defer jpg.Close()

		if jpgerr != nil {
			fmt.Println("jpg encode error")
			continue
		}
		jpbbyte := jpg.GetBytes()
		sendErr := stream.Send(&api.CameReq{
			CameId:    "1234",
			TimeFrame: timef,
			Img:       jpbbyte,
		})
		if sendErr != nil {
			log.Fatalln(sendErr)
		}
		// fmt.Printf("Send... : %v , Size : %v\n", timef, len(jpbbyte))
		time.Sleep(100 * time.Microsecond)
	}
}
