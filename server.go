package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/zhanchengsong/grpc-chat/protobuf"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedChatServer
	mu                sync.Mutex
	writeChatStream   map[string]pb.Chat_ReceiveChatAndPresenceServer
	receiveChatStream map[string]pb.Chat_SendChatAndPresenceServer
	// Temp implementation
	queuedMessages map[string][]pb.ChatAndPresenceMessage
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// Handle request to receive chat and presence call when some client connects
func (s *server) ReceiveChatAndPresence(in *pb.StartReceivingChatsRequest, chatStream pb.Chat_ReceiveChatAndPresenceServer) error {
	// register the user and the stream
	s.mu.Lock()
	s.writeChatStream[in.GetUserId()] = chatStream
	s.mu.Unlock()
	// check if any pending messages are present
	s.mu.Lock()
	queued, ok := s.queuedMessages[in.GetUserId()]
	if ok {
		for _, msg := range queued {
			chatStream.Send(&msg)
		}
	}
	s.mu.Unlock()
	return nil
}


// Handle receive messages from the client
func (s *server) SendChatAndPresence(rec_stream pb.Chat_SendChatAndPresenceServer) error {
	for 
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// make the map
	var chat_server = server{}
	chat_server.receiveChatStream = make(map[string]pb.Chat_SendChatAndPresenceServer)
	chat_server.writeChatStream = make(map[string]pb.Chat_ReceiveChatAndPresenceServer)
	pb.RegisterChatServer(s, &chat_server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
