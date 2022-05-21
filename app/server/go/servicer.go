package main

import (
	"context"
	"github.com/asatale/envoy-rate-limit/app/server/go/hello_world"
	"github.com/rs/zerolog/log"
)

type HelloWorldServer struct {
}

func (s *HelloWorldServer) SayHello(ctx context.Context, req *hello_world.HelloRequest) (*hello_world.HelloReply, error) {
	deadline, _ := ctx.Deadline()
	log.Debug().Msgf("Handling RPC: SayHello. Deadline: %v", deadline)
	return &hello_world.HelloReply{ClientName: req.GetClientName(), SeqNum: req.GetSeqNum()}, nil
}

func (s *HelloWorldServer) LotsOfGreetings(stream hello_world.Greeter_LotsOfGreetingsServer) error {
	var clntName string
	var seqNum int32

	deadline, _ := stream.Context().Deadline()
	log.Debug().Msgf("Handling RPC: LotsOfGreetings. Deadline: %v", deadline)

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Debug().Msgf("LotsOfGreetings: Received error %v", err)
			break
		}
		clntName = req.GetClientName()
		seqNum = req.GetSeqNum()
	}

	rsp := &hello_world.HelloReply{ClientName: clntName, SeqNum: seqNum}
	err := stream.SendAndClose(rsp)
	return err
}

func (s *HelloWorldServer) LotsOfReplies(req *hello_world.HelloRequest, stream hello_world.Greeter_LotsOfRepliesServer) error {
	clntName := req.GetClientName()
	seqNum := req.GetSeqNum()

	deadline, _ := stream.Context().Deadline()
	log.Debug().Msgf("Handling RPC: LotsOfReplies. Deadline: %v", deadline)

	for i := 0; i < 10; i++ {
		err := stream.Send(&hello_world.HelloReply{
			ClientName: clntName,
			SeqNum:     seqNum,
		})
		if err != nil {
			log.Debug().Msgf("LotsOfReplies: Received error %v", err)
			break
		}
	}
	return nil
}

func (s *HelloWorldServer) BidiHello(stream hello_world.Greeter_BidiHelloServer) error {

	deadline, _ := stream.Context().Deadline()
	log.Debug().Msgf("Handling RPC: BidiHello. Deadline: %v", deadline)

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Debug().Msgf("BidiHello: Received error %v", err)
			break
		}
		log.Debug().Msgf("BidiHello: Received msg %v", req)

		stream.Send(&hello_world.HelloReply{
			ClientName: req.GetClientName(),
			SeqNum:     req.GetSeqNum(),
		})
	}
	return nil
}
