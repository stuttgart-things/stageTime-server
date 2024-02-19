/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	redisAddress  = os.Getenv("REDIS_SERVER")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisStream   = os.Getenv("REDIS_STREAM")
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
