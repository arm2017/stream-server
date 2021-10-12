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
	address = "172.20.10.14:50051"
	// address    = "0.tcp.ap.ngrok.io:18664"
	deviceID   = 0
	jpgQuality = 40
)

func streamCameToServer(streamClient *api.StreamCameServiceClient) {
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

	stream, err := (*streamClient).Streaming(context.Background(), grpc.UseCompressor(gzip.Name))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	loop := 0
	img := gocv.NewMat()
	defer img.Close()

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
			R: 59,
			G: 57,
			B: 244,
			A: 1,
		}, 2)

		jpg, jpgerr := gocv.IMEncodeWithParams(gocv.JPEGFileExt, img, []int{gocv.IMWriteJpegQuality, jpgQuality})
		if jpgerr != nil {
			fmt.Println("jpg encode error")
			continue
		}
		sendErr := stream.Send(&api.CameReq{
			CameId:    "1234",
			TimeFrame: timef,
			Img:       jpg.GetBytes(),
		})
		if sendErr != nil {
			log.Fatalln(sendErr)
		}
		loop = loop + 1
		// fmt.Printf("Send... : %v , Loop : %v\n", timef, (loop))
		time.Sleep(100 * time.Microsecond)
	}
}

func Run() {
	// Set up a connection to the server.
	log.Printf("connect: %v\n", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("connect: %v success.\n", address)
	defer conn.Close()
	c := api.NewStreamCameServiceClient(conn)
	//stream to server
	streamCameToServer(&c)
	fmt.Println("end")
}
