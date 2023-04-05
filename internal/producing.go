/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	redis "github.com/go-redis/redis/v7"
	redisqueue "github.com/robinjoseph08/redisqueue/v2"
)

var (
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisQueue    = os.Getenv("REDIS_QUEUE")
)

func SendPipelineRunToMessageQueue(renderedPipelineruns map[int][]string) {

	pipelineRuns, err := json.Marshal(renderedPipelineruns)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress + ":" + redisPort,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	p, err := redisqueue.NewProducerWithOptions(&redisqueue.ProducerOptions{
		StreamMaxLength:      10000,
		ApproximateMaxLength: true,
		RedisClient:          redisClient,
	})

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	err = p.Enqueue(&redisqueue.Message{
		Stream: redisQueue,
		Values: map[string]interface{}{
			"revisionRun": string(pipelineRuns),
		},
	})
	if err != nil {
		panic(err)
	}

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
