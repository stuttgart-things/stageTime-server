package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
	revisionrun "github.com/stuttgart-things/stageTime-server/revisionrun"
)

var (
	redisServer   = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	ctx           = context.Background()
)

func main() {

	// INITALIZE REDIS
	redisClient := sthingsCli.CreateRedisClient(redisServer+":"+redisPort, redisPassword)
	redisJSONHandler := rejson.NewReJSONHandler()
	redisJSONHandler.SetGoRedisClient(redisClient)

	revisionRun := GetRevisionRunFromRedis("st-0-execute-ansible-rke2-cluster-1807283c5a", redisJSONHandler)
	fmt.Println(revisionRun)

}

func GetRevisionRunFromRedis(pipelineRunName string, redisJSONHandler *rejson.Handler) (revisionRun *revisionrun.CreateRevisionRunRequest) {

	revisionRunsJSON, err := redis.Bytes(redisJSONHandler.JSONGet(pipelineRunName, "."))
	if err != nil {
		log.Fatalf("Failed to JSONGet")
		return
	}

	fmt.Println(string(revisionRunsJSON))

	err = json.Unmarshal(revisionRunsJSON, &revisionRun)

	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal ", pipelineRunName)
		return
	}

	fmt.Printf("PipelineRun read from redis: %#v\n", revisionRun)

	return
}
