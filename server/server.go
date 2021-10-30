package server

import (
	"context"
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
	port      = ":50051"
	bufferImg = 30
)

type server struct {
	api.UnimplementedStreamCameServiceServer
	cams   *api.CameReq
	movehw MoveHW
}

type MoveHW struct {
	hwId        string
	online      bool
	lastConnect int64
	cmd         api.MoveRsp
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
			// fmt.Println("send...")
			lastTimeFrame = item.TimeFrame
		}

		cerr := stream.Context().Err()
		if cerr != nil {
			log.Println(cerr)
			return cerr
		}

		time.Sleep(10 * time.Millisecond)
	}

}

func (svr *server) MoveRegister(in *api.MoveRegisterReq, stream api.StreamCameService_MoveRegisterServer) error {

	if in.HwId == "" {
		return errors.New("hwid is null")
	}

	svr.movehw.hwId = in.HwId
	svr.movehw.lastConnect = in.RegisTime
	svr.movehw.online = true
	svr.movehw.cmd = api.MoveRsp{
		Direction: "",
		TimeMove:  0,
	}
	lasttime := int64(0)
	fmt.Printf("MoveRegister : %v\n", in.HwId)

	for {
		if svr.movehw.cmd.TimeMove > lasttime {
			lasttime = svr.movehw.cmd.TimeMove

			err := stream.Send(&api.MoveRsp{
				Direction: svr.movehw.cmd.Direction,
				TimeMove:  svr.movehw.cmd.TimeMove,
			})

			if err != nil {
				fmt.Println(err)
				svr.movehw.online = false
				return err
			}

		}

		time.Sleep(10 * time.Millisecond)

		cterr := stream.Context().Err()
		if cterr != nil {
			fmt.Println(cterr)
			svr.movehw.online = false
			break
		}

	}

	return nil

}

func (svr *server) Move(ct context.Context, in *api.MoveReq) (*api.MoveRsp, error) {

	fmt.Printf("IN : %v\n", in.HwId)

	if in.HwId == "" {
		return nil, errors.New("id not match : " + in.HwId)
	}

	if !svr.movehw.online {
		return nil, errors.New("move hw is offline")
	}

	svr.movehw.cmd.Direction = in.Direction
	svr.movehw.cmd.TimeMove = in.TimeMove

	return &api.MoveRsp{
		Direction: in.Direction,
		TimeMove:  in.TimeMove,
	}, nil
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
