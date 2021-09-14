package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	api "stream-cam-api/stream-camera/api"
	"time"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
)

const (
	port      = "localhost:50051"
	bufferImg = 30
)

type server struct {
	api.UnimplementedStreamCameServiceServer
	cams *api.CameReq
}

func (svr *server) Streaming(stream api.StreamCameService_StreamingServer) error {

	var lastTimeFrame int64 = 0
	fmt.Println("client stream connect !! ")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("End.. Recv")

			stream.SendAndClose(&api.CameRsp{
				CameId: req.CameId,
			})
			return nil
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		mainErr := stream.Context().Err()
		if mainErr != nil {
			fmt.Println(mainErr)
			return mainErr
		}

		if req.TimeFrame >= lastTimeFrame {
			lastTimeFrame = req.TimeFrame
			svr.cams = req
		}

	}

	// return nil
}

func (svr *server) View(in *api.VeiwReq, stream api.StreamCameService_ViewServer) error {

	log.Printf("Client connect view id : %v", in.CameId)
	var lastTimeFrame int64 = 0
	for {
		item := svr.cams
		if item == nil {
			return errors.New("fail cam id not found")
		}

		if item.TimeFrame > lastTimeFrame {

			sendErr := stream.Send(&api.VeiwRsp{
				CameId:    item.CameId,
				Img:       item.Img,
				TimeFrame: item.TimeFrame,
			})

			if sendErr != nil {
				log.Println(sendErr)
				return sendErr
			}
			lastTimeFrame = item.TimeFrame
		}

		time.Sleep(10 * time.Millisecond)
	}

}

func Run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterStreamCameServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
