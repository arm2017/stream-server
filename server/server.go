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
	bufferImg = 15
)

type server struct {
	api.UnimplementedStreamCameServiceServer
	cams map[string]api.CameReq
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

		// _, ok := svr.cams[req.CameId]
		// if !ok {
		// 	fmt.Println("client stream id : " + req.CameId)
		// }

		if req.TimeFrame >= lastTimeFrame {
			lastTimeFrame = req.TimeFrame
			svr.cams[req.CameId] = api.CameReq{
				CameId:    req.CameId,
				Img:       req.Img,
				TimeFrame: req.TimeFrame,
			}

		}

	}

	// return nil
}

func (svr *server) View(in *api.VeiwReq, stream api.StreamCameService_ViewServer) error {

	log.Printf("Client connect view id : %v", in.CameId)
	var lastTimeFrame int64 = 0
	for {
		item, ok := svr.cams[in.CameId]
		if !ok {
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
			time.Sleep(50 * time.Millisecond)
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
	api.RegisterStreamCameServiceServer(s, &server{
		cams: make(map[string]api.CameReq),
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
