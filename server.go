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

	rejson "github.com/nitishm/go-rejson/v4"
	"github.com/stuttgart-things/stageTime-server/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
	"google.golang.org/grpc/reflection"

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

var (
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisQueue    = os.Getenv("REDIS_QUEUE")
)

type Server struct {
	revisionrun.UnimplementedStageTimeApplicationServiceServer
}

func NewServer() Server {
	return Server{}
}

func (s Server) CreateRevisionRun(ctx context.Context, gRPCRequest *revisionrun.CreateRevisionRunRequest) (*revisionrun.Response, error) {

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

	// LOOP OVER REVISIONRUN

	for i := 0; i < (len(renderedPipelineruns)); i++ {

		for j, pr := range renderedPipelineruns[i] {

			fmt.Println(j)
			fmt.Println(pr)

			resourceName, _ := sthingsBase.GetRegexSubMatch(pr, `name: "(.*?)"`)
			revisionRunID, _ := sthingsBase.GetRegexSubMatch(pr, `commit: "(.*?)"`)
			stage, _ := sthingsBase.GetRegexSubMatch(pr, `stage: "(.*?)"`)

			prIdentifier := strings.Split(resourceName, "-")

			fmt.Println("PR", i)
			fmt.Println("RESOURCE-NAME", resourceName)
			fmt.Println("IDENTIFIER", prIdentifier)
			fmt.Println("REVISIONRUN-ID", revisionRunID)
			fmt.Println("STAGE", stage)

			// CREATE REDIS CLIENT / JSON HANDLER
			redisClient := sthingsCli.CreateRedisClient(redisAddress+":"+redisPort, redisPassword)
			redisJSONHandler := rejson.NewReJSONHandler()
			redisJSONHandler.SetGoRedisClient(redisClient)

			// SET PR ON LIST
			sthingsCli.AddValueToRedisSet(redisClient, revisionRunID, resourceName)

			// CONVERT PR TO JSON + ADD TO REDIS
			prJSON := sthingsCli.ConvertYAMLToJSON(pr)
			fmt.Println(string(prJSON))
			sthingsCli.SetRedisJSON(redisJSONHandler, prJSON, resourceName)
		}
	}

	// SEND PIPELINERUN TO REDIS MessageQueue
	streamValues := map[string]interface{}{
		"stage": "stage0",
	}

	server.SendPipelineRunToMessageQueue(streamValues)
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
	log.Info("redis server " + redisAddress)
	log.Info("redis port " + redisPort)
	log.Info("redis queue " + redisQueue)

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
