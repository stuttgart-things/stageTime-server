/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"fmt"
	"os"
	"strings"

	redis "github.com/go-redis/redis/v7"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
)

func SendPipelineRunToMessageQueue(stageID string) {

	streamValues := map[string]interface{}{
		"stage": stageID,
	}

	sthingsCli.EnqueueDataInRedisStreams(redisAddress+":"+redisPort, redisPassword, redisStream, streamValues)
	fmt.Println("STREAM", redisStream)
	fmt.Println("VALUES", streamValues)

}

func SendStatsToRedis(renderedPipelineruns map[int][]string) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress + ":" + redisPort,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	for i := 0; i < (len(renderedPipelineruns)); i++ {

		for j, pr := range renderedPipelineruns[i] {

			fmt.Println(j)
			fmt.Println(pr)

			resourceName, _ := sthingsBase.GetRegexSubMatch(pr, `name: "(.*?)"`)
			prIdentifier := strings.Split(resourceName, "-")
			err := redisClient.Set("prIdentifier", prIdentifier[(len(prIdentifier)-1)], 0).Err()
			if err != nil {
				panic(err)
			}

			err = redisClient.Set("countPipelineRuns", len(renderedPipelineruns), 0).Err()
			if err != nil {
				panic(err)
			}

		}
	}

}
