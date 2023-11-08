package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	//"google.golang.org/grpc/credentials"
	revisionrun "github.com/stuttgart-things/stageTime-server/revisionrun"
	"google.golang.org/grpc/credentials"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	address = "localhost:50051" // local go
	// address = "stagetime.cd43.sthings-pve.labul.sva.de:443"
	testFilePath = "tests/prs.json"
	// testFilePath = "tests/ansible.json"
)

type Client struct {
	yasClient revisionrun.StageTimeApplicationServiceClient
	timeout   time.Duration
}

func NewClient(conn grpc.ClientConnInterface, timeout time.Duration) Client {
	return Client{
		yasClient: revisionrun.NewStageTimeApplicationServiceClient(conn),
		timeout:   timeout,
	}
}

func (c Client) CreateRevisionRun(ctx context.Context, json io.Reader) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(c.timeout))
	defer cancel()

	req := revisionrun.CreateRevisionRunRequest{}
	if err := jsonpb.Unmarshal(json, &req); err != nil {
		return fmt.Errorf("CLIENT CREATE REVISIONRUN: UNMARSHAL: %w", err)
	}

	res, err := c.yasClient.CreateRevisionRun(ctx, &req)

	fmt.Println(res)

	if err != nil {
		if er, ok := status.FromError(err); ok {
			return fmt.Errorf("CLIENT CREATE REVISIONRUN: CODE: %s - msg: %s", er.Code(), er.Message())
		}
		return fmt.Errorf("CLIENT CREATE REVISIONRUn: %w", err)
	}

	log.Println("RESULT:", res.Result)
	log.Println("RESPONSE:", res)

	return nil
}

func main() {

	// Check env vor given server port
	if os.Getenv("STAGETIME_SERVER") != "" {
		address = os.Getenv("STAGETIME_SERVER")
	}

	if os.Getenv("STAGETIME_TEST_FILES") != "" {
		testFilePath = os.Getenv("STAGETIME_TEST_FILES")
	}

	if strings.Contains(address, "localhost") {
		ConnectInsecure(address, testFilePath)
	} else {
		ConnectSecure(address, testFilePath)
	}

}

func ConnectSecure(address, testFilePath string) {

	log.Println("CLIENT STARTED CONNECTING TO.. " + address)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithBlock())

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	json, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	yasClient := NewClient(conn, time.Second)
	err = yasClient.CreateRevisionRun(context.Background(), bytes.NewBuffer(json))

	log.Println("ERR:", err)

}

func ConnectInsecure(address, testFilePath string) {

	log.Println("client started connecting to.. " + address)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	json, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	yasClient := NewClient(conn, time.Second)
	err = yasClient.CreateRevisionRun(context.Background(), bytes.NewBuffer(json))

	log.Println("ERR:", err)

}
