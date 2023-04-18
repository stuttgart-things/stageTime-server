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

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	server "github.com/stuttgart-things/yacht-application-server/server"

	"google.golang.org/grpc/reflection"

	"github.com/fatih/color"
	revisionrun "github.com/stuttgart-things/yacht-application-server/revisionrun"
	goVersion "go.hein.dev/go-version"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

var (
	shortened   = false
	version     = "unset"
	date        = "unknown"
	commit      = "unknown"
	output      = "yaml"
	serverPort  = port
	logfilePath = "yas.log"
	log         = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
)

// https://patorjk.com/software/taag/#p=testall&h=3&f=Graffiti&t=YAS
// 3D-ASCII
const banner = `
___    ___ ________  ________
|\  \  /  /|\   __  \|\   ____\
\ \  \/  / | \  \|\  \ \  \___|_
 \ \    / / \ \   __  \ \_____  \
  \/  /  /   \ \  \ \  \|____|\  \
__/  / /      \ \__\ \__\____\_\  \
|\___/ /        \|__|\|__|\_________\
\|___|/                  \|_________|


`

type Server struct {
	revisionrun.UnimplementedYachtApplicationServiceServer
}

func NewServer() Server {
	return Server{}
}

func (s Server) CreateRevisionRun(ctx context.Context, req *revisionrun.CreateRevisionRunRequest) (*revisionrun.Response, error) {

	log := sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)

	server.RenderPipelineRuns(req)

	receivedRevisionRun := bytes.Buffer{}

	mars := jsonpb.Marshaler{OrigName: true, EmitDefaults: true}
	if err := mars.Marshal(&receivedRevisionRun, req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "server create revisionrun: marshal: %v", err)
	}

	log.Println("REQUEST:", receivedRevisionRun.String())

	if err := json.Unmarshal([]byte(receivedRevisionRun.Bytes()), &req); err != nil {
		log.Fatal(err)
	}

	// STATUS OUTPUT GRPC DATA
	fmt.Println(req.Author + " created RevisionRun " + req.CommitId + " at " + req.PushedAt)
	fmt.Println("Repository:", req.RepoName)
	fmt.Println("RepositoryUrl:", req.RepoUrl)
	fmt.Println("PipelineRuns:", len(req.Pipelineruns))

	// TEST RENDERING
	renderedPipelineruns := server.RenderPipelineRuns(req)
	fmt.Println(renderedPipelineruns)
	log.Info("all pipelineRuns can be rendered")

	// TEST RENDERING
	server.SendStatsToRedis(renderedPipelineruns)

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

	if os.Getenv("SERVER_PORT") != "" {
		serverPort = ":" + os.Getenv("SERVER_PORT")
	}

	log.Info("gRPC server running on port " + serverPort)

	// Output banner + version output
	color.Cyan(banner)

	color.Cyan("YACHT APPLICATION SERVER")
	resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
	color.Cyan(resp + "\n")

	listener, err := net.Listen("tcp", "0.0.0.0"+serverPort)
	if err != nil {
		log.Fatalln(err)
	}

	log.Info("YAS running at ", listener.Addr(), serverPort)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	yachtApplicationServer := NewServer()

	revisionrun.RegisterYachtApplicationServiceServer(grpcServer, yachtApplicationServer)

	log.Fatalln(grpcServer.Serve(listener))
}
