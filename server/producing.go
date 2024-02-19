/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nitishm/go-rejson/v4"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
	// log           = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	// logfilePath   = "stageTime-server.log"
)

func SendStageToMessageQueue(stageID string) {

	streamValues := map[string]interface{}{
		"stage": stageID,
	}

	sthingsCli.EnqueueDataInRedisStreams(redisAddress+":"+redisPort, redisPassword, redisStream, streamValues)
	fmt.Println("STREAM", redisStream)
	fmt.Println("VALUES", streamValues)
}

func OutputRevisonRunStatus(renderedPipelineruns map[int][]string) {

	t := time.Now()
	countPipelineRuns := 0
	stageNumber := ""
	revisionRunID := ""

	for i := 0; i < (len(renderedPipelineruns)); i++ {
		for _, pr := range renderedPipelineruns[i] {
			countPipelineRuns += 1
			stageNumber, _ = sthingsBase.GetRegexSubMatch(pr, `stage: "(.*?)"`)
			revisionRunID, _ = sthingsBase.GetRegexSubMatch(pr, `commit: "(.*?)"`)
		}
	}
	countStage := sthingsBase.ConvertStringToInteger(stageNumber) + 1

	rrs := RevisionRunStatus{
		RevisionRun:       revisionRunID,
		CountStages:       countStage,
		CountPipelineRuns: countPipelineRuns,
		LastUpdated:       t.Format("2006-01-02 15:04:05"),
		Status:            "REVISIONRUN CREATED W/ STAGETIME-SERVER",
	}

	tw := table.NewWriter()
	header := sthingsCli.CreateTableHeader(rrs)
	tw.AppendHeader(header)
	tw.AppendRow(sthingsCli.CreateTableRows(rrs))
	tw.AppendSeparator()
	tw.SetStyle(table.StyleColoredBright)
	tw.SetOutputMirror(os.Stdout)
	tw.Render()
}

// GET REVISIONRUN STATUS FROM REDIS
func GetRevisionRunFromRedis(redisJSONHandler *rejson.Handler, revisionRunStatusID string, print bool) (revisionRunFromRedis RevisionRunStatus) {

	revisionRunStatus := sthingsCli.GetRedisJSON(redisJSONHandler, revisionRunStatusID)

	err := json.Unmarshal(revisionRunStatus, &revisionRunFromRedis)
	if err != nil {
		log.Fatalf("FAILED TO JSON UNMARSHAL REVISIONRUN STATUS")
	}

	if print {
		PrintTable(revisionRunFromRedis)
	}

	return
}

// GET STAGE STATUS FROM REDIS
func GetStageFromRedis(redisJSONHandler *rejson.Handler, stageID string, print bool) (stageStatus StageStatus) {

	revisionRunStatus := sthingsCli.GetRedisJSON(redisJSONHandler, stageID)

	err := json.Unmarshal(revisionRunStatus, &stageStatus)
	if err != nil {
		log.Fatalf("FAILED TO JSON UNMARSHAL STAGE STATUS")
	}

	if print {
		PrintTable(stageStatus)
	}

	return
}

// SET REVISIONRUN STATUS IN REDIS
func SetRevisionRunStatusInRedis(redisJSONHandler *rejson.Handler, revisionRunStatusID, updatedMessage string, revisionRun RevisionRunStatus, print bool) {

	// UPDATE MESSAGE
	revisionRun.Status = updatedMessage

	sthingsCli.SetRedisJSON(redisJSONHandler, revisionRun, revisionRunStatusID)

	if print {
		PrintTable(revisionRun)
	}

}

// SET STAGE STATUS IN REDIS
func SetStageStatusInRedis(redisJSONHandler *rejson.Handler, stageID, updatedMessage string, stageStatus StageStatus, print bool) {

	// UPDATE MESSAGE
	stageStatus.Status = updatedMessage

	sthingsCli.SetRedisJSON(redisJSONHandler, stageStatus, stageID)

	if print {
		PrintTable(stageStatus)
	}

}
