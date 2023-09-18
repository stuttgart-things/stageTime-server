/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"fmt"
	"os"
	"strings"

	sthingsCli "github.com/stuttgart-things/sthingsCli"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	redis "github.com/go-redis/redis/v7"
)

var (
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisQueue    = os.Getenv("REDIS_QUEUE")
)

func SendPipelineRunToMessageQueue(streamValues map[string]interface{}) {

	sthingsCli.EnqueueDataInRedisStreams(redisAddress+":"+redisPort, redisPassword, redisQueue, streamValues)

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
