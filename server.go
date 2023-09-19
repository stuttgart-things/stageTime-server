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
	"time"

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
	now         = time.Now()
)

var (
	redisAddress     = os.Getenv("REDIS_SERVER")
	redisPort        = os.Getenv("REDIS_PORT")
	redisPassword    = os.Getenv("REDIS_PASSWORD")
	redisQueue       = os.Getenv("REDIS_QUEUE")
	redisClient      = sthingsCli.CreateRedisClient(redisAddress+":"+redisPort, redisPassword)
	redisJSONHandler = rejson.NewReJSONHandler()
)

type Server struct {
	revisionrun.UnimplementedStageTimeApplicationServiceServer
}

func NewServer() Server {
	return Server{}
}

func (s Server) CreateRevisionRun(ctx context.Context, gRPCRequest *revisionrun.CreateRevisionRunRequest) (*revisionrun.Response, error) {

	// CREATE REDIS CLIENT / JSON HANDLER
	redisJSONHandler.SetGoRedisClient(redisClient)

	receivedRevisionRun := bytes.Buffer{}

	mars := jsonpb.Marshaler{OrigName: true, EmitDefaults: true}

	if err := mars.Marshal(&receivedRevisionRun, gRPCRequest); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "server create revisionrun: marshal: %v", err)
	}

	log.Println("INCOMING gRPC REQUEST:", receivedRevisionRun.String())

	if err := json.Unmarshal([]byte(receivedRevisionRun.Bytes()), &gRPCRequest); err != nil {
		log.Fatal(err)
	}

	// STATUS OUTPUT GRPC DATA
	log.Info(gRPCRequest.Author + " created RevisionRun " + gRPCRequest.CommitId + " at " + gRPCRequest.PushedAt)
	log.Info("REPOSITORY: ", gRPCRequest.RepoName)
	log.Info("REPOSITORYURL: ", gRPCRequest.RepoUrl)
	log.Info("PIPELINERUNS: ", len(gRPCRequest.Pipelineruns))

	// TEST RENDERING
	renderedPipelineruns := server.RenderPipelineRuns(gRPCRequest)
	log.Info("ALL PIPELINERUNS CAN BE RENDERED")

	// SEND STATS TO REDIS
	server.SendStatsToRedis(renderedPipelineruns)

	// LOOP OVER REVISIONRUN
	for i := 0; i < (len(renderedPipelineruns)); i++ {

		for _, pr := range renderedPipelineruns[i] {

			resourceName, _ := sthingsBase.GetRegexSubMatch(pr, `name: "(.*?)"`)
			revisionRunID, _ := sthingsBase.GetRegexSubMatch(pr, `commit: "(.*?)"`)
			stage, _ := sthingsBase.GetRegexSubMatch(pr, `stage: "(.*?)"`)

			prIdentifier := strings.Split(resourceName, "-")

			fmt.Println("PR", i)
			fmt.Println("RESOURCE-NAME", resourceName)
			fmt.Println("IDENTIFIER", prIdentifier)
			fmt.Println("REVISIONRUN-ID", revisionRunID)
			fmt.Println("STAGE", stage)

			// SET STAGES ON LIST
			sthingsCli.AddValueToRedisSet(redisClient, now.Format(time.RFC3339)+"-"+revisionRunID+"-"+"stages", stage)
			sthingsCli.AddValueToRedisSet(redisClient, now.Format(time.RFC3339)+"-"+revisionRunID, resourceName)
			log.Info("REVISIONRUN NAME "+resourceName+" STORED ON ", now.Format(time.RFC3339)+"-"+revisionRunID)

			sthingsCli.AddValueToRedisSet(redisClient, now.Format(time.RFC3339)+"-"+revisionRunID+"-"+stage, resourceName)
			log.Info("REVISIONRUN NAME "+resourceName+" STORED ON ", now.Format(time.RFC3339)+"-"+revisionRunID+"-"+stage)

			// CONVERT PR TO JSON + ADD TO REDIS
			prJSON := sthingsCli.ConvertYAMLToJSON(pr)
			fmt.Println(string(prJSON))
			sthingsCli.SetRedisJSON(redisJSONHandler, prJSON, resourceName)
		}
	}

	fmt.Println("REVISONRUN PRINTED")
	cr := server.RenderRevisionRunCR()
	fmt.Println(string(cr))
	crJSON := sthingsCli.ConvertYAMLToJSON(string(cr))
	fmt.Println(crJSON)

	stageID := "stageTime-" + gRPCRequest.CommitId[0:4]
	fmt.Println("COMMIT ID: ", stageID)
	sthingsCli.SetRedisJSON(redisJSONHandler, crJSON, stageID)

	// SEND PIPELINERUN TO REDIS MessageQueue
	server.SendPipelineRunToMessageQueue(stageID)
	log.Info("REVISIONRUN WAS STORED IN MESSAGEQUEUE")

	return &revisionrun.Response{
		Result: revisionrun.Response_SUCCESS,
		Success: &revisionrun.Response_Success{
			Data: []byte("GOOD JOB - REVISIONRUN WAS STORED IN MESSAGEQUEUE"),
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
