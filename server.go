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
	serverPort        = port
	logfilePath       = "stageTime-server.log"
	log               = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	now               = time.Now()
	stage             string
	stageNumber       string
	revisionRunID     string
	countPipelineRuns = 0
	pipelineRunStatus []server.PipelineRunStatus
)

var (
	redisAddress     = os.Getenv("REDIS_SERVER")
	redisPort        = os.Getenv("REDIS_PORT")
	redisPassword    = os.Getenv("REDIS_PASSWORD")
	redisStream      = os.Getenv("REDIS_STREAM")
	redisClient      = sthingsCli.CreateRedisClient(redisAddress+":"+redisPort, redisPassword)
	redisJSONHandler = rejson.NewReJSONHandler()
)

type Server struct {
	revisionrun.UnimplementedStageTimeApplicationServiceServer
}

func NewServer() Server {
	return Server{}
}

func main() {

	// PRINT BANNER + VERSION INFO
	internal.PrintBanner()

	if os.Getenv("SERVER_PORT") != "" {
		serverPort = ":" + os.Getenv("SERVER_PORT")
	}

	log.Info("GRPC SERVER RUNNING ON PORT: " + serverPort)
	log.Info("REDIS SERVER: " + redisAddress)
	log.Info("REDIS PORT: " + redisPort)
	log.Info("REDIS QUEUE: " + redisStream)

	listener, err := net.Listen("tcp", "0.0.0.0"+serverPort)
	if err != nil {
		log.Fatalln(err)
	}

	log.Info("STAGETIME-SERVER RUNNING AT ", listener.Addr(), serverPort)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	stageTimeServer := NewServer()

	revisionrun.RegisterStageTimeApplicationServiceServer(grpcServer, stageTimeServer)

	log.Fatalln(grpcServer.Serve(listener))
}

func (s Server) CreateRevisionRun(ctx context.Context, gRPCRequest *revisionrun.CreateRevisionRunRequest) (*revisionrun.Response, error) {

	prInformation := make(map[string]string)

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

	// LOOP OVER REVISIONRUN
	stages := make(map[string]int)
	for i := 0; i < (len(renderedPipelineruns)); i++ {

		for _, pr := range renderedPipelineruns[i] {

			prValid, prInformation := internal.ValidateStorePipelineRuns(pr)

			if !prValid {
				log.Error("PIPELINERUN NOT VALID - SKIPPING")
				break
			}

			log.Info("COUNT-PR: ", i)
			log.Info("RESOURCE-NAME: ", prInformation["name"])
			log.Info("IDENTIFIER: ", prInformation["name"])
			log.Info("REVISIONRUN-ID: ", prInformation["stagetime/commit"])
			log.Info("STAGE: ", prInformation["stagetime/stage"])
			log.Info("DATE: ", prInformation["stagetime/date"])

			countPipelineRuns += 1

			stage = prInformation["stagetime/stage"]
			stages[stage] = server.SetStage(stages, stage)
			revisionRunID = prInformation["stagetime/commit"]
			// SET STAGES ON LIST
			// sthingsCli.AddValueToRedisSet(redisClient, now.Format(time.RFC3339)+"-"+prInformation["stagetime/commit"]+"-"+"stages", stage)
			// sthingsCli.AddValueToRedisSet(redisClient, now.Format(time.RFC3339)+"-"+prInformation["stagetime/commit"], prInformation["name"])

			// CONVERT PR TO JSON + ADD TO REDIS
			prJSON := sthingsCli.ConvertYAMLToJSON(pr)
			fmt.Println(string(prJSON))
			sthingsCli.SetRedisJSON(redisJSONHandler, prJSON, prInformation["name"])
			log.Info("PIPELINERUN WAS ADDED TO REDIS (JSON): ", prInformation["name"])

			// CREATE ON REVISIONRUN STATUS ON REDIS + PRINT AS TABLE
			initialPrs := server.PipelineRunStatus{
				Stage:           sthingsBase.ConvertStringToInteger(prInformation["stagetime/stage"]),
				PipelineRunName: prInformation["name"],
				CanFail:         sthingsBase.ConvertStringToBoolean(prInformation["canFail"]),
				LastUpdated:     now.Format("2006-01-0215-04-05"),
				Status:          "NOT STARTED (YET)",
			}
			pipelineRunStatus = append(pipelineRunStatus, initialPrs)

			// sthingsCli.DeleteRedisSet(redisClient, prInformation["stagetime/commit"]+"-"+stage)
			sthingsCli.AddValueToRedisSet(redisClient, prInformation["stagetime/date"]+"-"+prInformation["stagetime/commit"]+"-"+prInformation["stagetime/stage"], prInformation["name"])
			log.Info("ADDED PIPELINERUN NAMES TO REDIS (SET): ", prInformation["stagetime/date"]+"-"+prInformation["stagetime/commit"]+"-"+prInformation["stagetime/stage"])
		}
	}

	countStage := sthingsBase.ConvertStringToInteger(prInformation["stage"]) + 1

	// CREATE REVISIONRUN STATUS ON REDIS + PRINT AS TABLE
	initialRrs := server.RevisionRunStatus{
		RevisionRun:       revisionRunID,
		CountStages:       countStage,
		CountPipelineRuns: countPipelineRuns,
		LastUpdated:       now.Format("2006-01-0215-04-05"),
		Status:            "CREATED W/ STAGETIME-SERVER",
	}

	sthingsCli.SetRedisJSON(redisJSONHandler, initialRrs, revisionRunID+"-status")
	log.Info("INITIAL REVISIONRUNSTATUS WAS ADDED TO REDIS (JSON): ", revisionRunID+"-status")
	server.PrintTable(initialRrs)

	// CREATE PIPELINERUN STATUS ON REDIS + PRINT AS TABLE
	for _, pr := range pipelineRunStatus {
		sthingsCli.SetRedisJSON(redisJSONHandler, pr, pr.PipelineRunName+"-status")
		log.Info("INITIAL PIPELINERUN STATUS WAS ADDED TO REDIS (JSON): ", pr.PipelineRunName+"-status")
		server.PrintTable(pr)
	}

	// CREATE STAGE STATUS ON REDIS + PRINT AS TABLE
	for stageNumber, _ := range stages {

		fmt.Println("STAGENUMBER: ", stageNumber)

		initialStageStatus := server.StageStatus{
			StageID:           now.Format("2006-01-0215-04-05") + "-" + revisionRunID + "-" + stageNumber,
			CountPipelineRuns: stages[stageNumber],
			LastUpdated:       now.Format("2006-01-0215-04-05"),
			Status:            "CREATED W/ STAGETIME-SERVER",
		}
		// "2024-01-2405-54-56-093uohkjscbod32903de-0"

		log.Info("INITIAL STATE STATUS WAS ADDED TO REDIS (JSON): ", revisionRunID+stageNumber)
		sthingsCli.SetRedisJSON(redisJSONHandler, initialStageStatus, revisionRunID+stageNumber)
		server.PrintTable(initialStageStatus)
	}

	// OUTPUT RevisionRun STATUS
	server.OutputRevisonRunStatus(renderedPipelineruns)

	// SEND STAGE TO STREAM
	server.SendStageToMessageQueue(now.Format("2006-01-0215-04-05") + "+" + revisionRunID + "+0")
	log.Info("STAGE WAS QUEUED FOR PIPELINERUN CREATION ON SERVER (STREAM): ", revisionRunID+"+0")

	// HANDLING OF REVISONRUN CR
	// stageID := "stageTime-" + gRPCRequest.CommitId[0:4]
	// fmt.Println("REVISIONRUN ID: ", stageID)
	// sthingsCli.AddValueToRedisSet(redisClient, now.Format(time.RFC3339)+"-"+stageID, stageID)
	// cr := server.RenderRevisionRunCR()
	// fmt.Println(string(cr))
	// crJSON := sthingsCli.ConvertYAMLToJSON(string(cr))
	// fmt.Println(crJSON)
	// sthingsCli.SetRedisJSON(redisJSONHandler, crJSON, stageID)

	// SEND gRPC RESPONSE
	return &revisionrun.Response{
		Result: revisionrun.Response_SUCCESS,
		Success: &revisionrun.Response_Success{
			Data: []byte("GOOD JOB - REVISIONRUN WAS CREATED"),
		},
	}, nil
}
