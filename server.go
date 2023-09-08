/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	server "github.com/stuttgart-things/stageTime-server/server"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"google.golang.org/grpc/reflection"

	"github.com/stuttgart-things/stageTime-server/internal"

	revisionrun "github.com/stuttgart-things/stageTime-server/revisionrun"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

var (
	serverPort  = port
	logfilePath = "stageTime-server.log"
	log         = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
)

type Server struct {
	revisionrun.UnimplementedStageTimeApplicationServiceServer
}

func NewServer() Server {
	return Server{}
}

func (s Server) CreateRevisionRun(ctx context.Context, gRPCRequest *revisionrun.CreateRevisionRunRequest) (*revisionrun.Response, error) {

	log := sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)

	server.RenderPipelineRuns(gRPCRequest)

	receivedRevisionRun := bytes.Buffer{}

	mars := jsonpb.Marshaler{OrigName: true, EmitDefaults: true}
	if err := mars.Marshal(&receivedRevisionRun, gRPCRequest); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "server create revisionrun: marshal: %v", err)
	}

	log.Println("REQUEST:", receivedRevisionRun.String())

	if err := json.Unmarshal([]byte(receivedRevisionRun.Bytes()), &gRPCRequest); err != nil {
		log.Fatal(err)
	}

	// STATUS OUTPUT GRPC DATA
	fmt.Println(gRPCRequest.Author + " created RevisionRun " + gRPCRequest.CommitId + " at " + gRPCRequest.PushedAt)
	fmt.Println("Repository:", gRPCRequest.RepoName)
	fmt.Println("RepositoryUrl:", gRPCRequest.RepoUrl)
	fmt.Println("PipelineRuns:", len(gRPCRequest.Pipelineruns))

	// TEST RENDERING
	renderedPipelineruns := server.RenderPipelineRuns(gRPCRequest)
	fmt.Println(renderedPipelineruns)
	log.Info("all pipelineRuns can be rendered")

	// SEND STATS TO REDIS
	server.SendStatsToRedis(renderedPipelineruns)

	// TEST LOOPING
	for i := 0; i < (len(renderedPipelineruns)); i++ {

		for j, pr := range renderedPipelineruns[i] {

			fmt.Println(j)
			fmt.Println(pr)

			resourceName, _ := sthingsBase.GetRegexSubMatch(pr, `name: "(.*?)"`)
			prIdentifier := strings.Split(resourceName, "-")

			fmt.Println("PR", i)
			fmt.Println("RESOURCE-NAME", resourceName)
			fmt.Println("IDENTIFIER", prIdentifier)

		}
	}

	// SEND PIPELINERUN TO REDIS MessageQueue
	server.SendPipelineRunToMessageQueue(renderedPipelineruns)
	log.Info("revisionRun was stored in MessageQueue")

	return &revisionrun.Response{
		Result: revisionrun.Response_SUCCESS,
		Success: &revisionrun.Response_Success{
			Data: []byte("good job - revisionRun was stored in MessageQueue"),
		},
	}, nil
}

func main() {

	// PRINT BANNER + VERSION INFO
	internal.PrintBanner()

	if os.Getenv("SERVER_PORT") != "" {
		serverPort = ":" + os.Getenv("SERVER_PORT")
	}

	log.Info("gRPC server running on port " + serverPort)

	listener, err := net.Listen("tcp", "0.0.0.0"+serverPort)
	if err != nil {
		log.Fatalln(err)
	}

	log.Info("stageTime-server running at ", listener.Addr(), serverPort)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	stageTimeServer := NewServer()

	revisionrun.RegisterStageTimeApplicationServiceServer(grpcServer, stageTimeServer)

	log.Fatalln(grpcServer.Serve(listener))
}
