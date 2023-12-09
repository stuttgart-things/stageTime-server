/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"fmt"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
)

type PipelineRunStatus struct {
	Stage       int
	PipelineRun string
	CanFail     bool
	LastUpdated string
	Status      string
}

type StageStatus struct {
	ID                string
	Status            string
	PipelineRunStatus []PipelineRunStatus
}

type RevisionRunStatus struct {
	ID          string
	CountStages string
	Status      string
	StageStatus []StageStatus
}

func SendStageToMessageQueue(stageID string) {

	streamValues := map[string]interface{}{
		"stage": stageID,
	}

	sthingsCli.EnqueueDataInRedisStreams(redisAddress+":"+redisPort, redisPassword, redisStream, streamValues)
	fmt.Println("STREAM", redisStream)
	fmt.Println("VALUES", streamValues)

}

func CreateTable(renderedPipelineruns map[int][]string) {

	countPipelineRuns := 0
	stageNumber := ""

	for i := 0; i < (len(renderedPipelineruns)); i++ {

		for _, pr := range renderedPipelineruns[i] {
			countPipelineRuns += 1
			stageNumber, _ = sthingsBase.GetRegexSubMatch(pr, `stage: "(.*?)"`)
		}

	}

	countStage := sthingsBase.ConvertStringToInteger(stageNumber) + 1

	fmt.Println("TOTAL STAGES:")
	fmt.Println(countStage)

	fmt.Println("TOTAL PIPELINERUNS:")
	fmt.Println(countPipelineRuns)

}

// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr:     redisAddress + ":" + redisPort,
// 		Password: redisPassword, // no password set
// 		DB:       0,             // use default DB
// 	})

// 	for i := 0; i < (len(renderedPipelineruns)); i++ {

// 		for j, pr := range renderedPipelineruns[i] {

// 			fmt.Println(j)
// 			fmt.Println(pr)

// 			resourceName, _ := sthingsBase.GetRegexSubMatch(pr, `name: "(.*?)"`)
// 			prIdentifier := strings.Split(resourceName, "-")
// 			err := redisClient.Set("prIdentifier", prIdentifier[(len(prIdentifier)-1)], 0).Err()
// 			if err != nil {
// 				panic(err)
// 			}

// 			err = redisClient.Set("countPipelineRuns", len(renderedPipelineruns), 0).Err()
// 			if err != nil {
// 				panic(err)
// 			}

// 		}
// 	}

// }
