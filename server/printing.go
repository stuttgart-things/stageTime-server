/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nitishm/go-rejson/v4"

	goredis "github.com/redis/go-redis/v9"

	"github.com/jedib0t/go-pretty/v6/table"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	redisUrl = os.Getenv("REDIS_SERVER") + ":" + os.Getenv("REDIS_PORT")
	// redisPassword    = os.Getenv("REDIS_PASSWORD")
	redisClient      = goredis.NewClient(&goredis.Options{Addr: redisUrl, Password: redisPassword, DB: 0})
	redisJSONHandler = rejson.NewReJSONHandler()
)

type RevisionRunStatus struct {
	RevisionRun       string
	CountStages       int
	CountPipelineRuns int
	LastUpdated       string
	Status            string
}

type StageStatus struct {
	StageID           string
	CountPipelineRuns int
	LastUpdated       string
	Status            string
}

type PipelineRunStatus struct {
	Stage           int
	PipelineRunName string
	CanFail         bool
	LastUpdated     string
	Status          string
}

func PrintTable(printObject interface{}) {

	tw := table.NewWriter()
	header := sthingsCli.CreateTableHeader(printObject)
	tw.AppendHeader(header)
	tw.AppendRow(sthingsCli.CreateTableRows(printObject))
	tw.AppendSeparator()
	tw.SetStyle(table.StyleColoredBright)
	tw.SetOutputMirror(os.Stdout)
	tw.Render()
}

func SetStage(stages map[string]int, stage string) (updatedValue int) {
	existingValue, ok := stages[stage]

	if ok {
		updatedValue = existingValue + 1
	} else {
		updatedValue = 1
	}

	return
}

func GetPipelineRunStatus(jsonKey string, redisJSONHandler *rejson.Handler) PipelineRunStatus {

	pipelineRunStatusJson := sthingsCli.GetRedisJSON(redisJSONHandler, jsonKey)
	pipelineRunStatus := PipelineRunStatus{}

	err := json.Unmarshal(pipelineRunStatusJson, &pipelineRunStatus)
	if err != nil {
		fmt.Println("FAILED TO JSON UNMARSHAL")
	}

	return pipelineRunStatus
}

func GetStageStatus(jsonKey string, redisJSONHandler *rejson.Handler) StageStatus {

	stageStatusJson := sthingsCli.GetRedisJSON(redisJSONHandler, jsonKey)
	stageStatus := StageStatus{}

	err := json.Unmarshal(stageStatusJson, &stageStatus)
	if err != nil {
		fmt.Println("FAILED TO JSON UNMARSHAL")
	}

	return stageStatus
}
