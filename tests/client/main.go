package main

import (
	revisionrun "github.com/stuttgart-things/stageTime-server/revisionrun"

	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := revisionrun.NewStatusesClient(conn)

	status, err := client.GetStatus(context.Background(), &revisionrun.StatusGetRequest{RevisionRunId: "ad3121246532123"})
	if err != nil {
		log.Fatalf("FAILED TO GET STATUS: %v", err)
	}
	log.Printf("STATUS: %v", status)

}
