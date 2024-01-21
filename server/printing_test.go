/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package server

import (
	"testing"
)

func TestPrintTable(t *testing.T) {

	testPr := PipelineRun{
		Name:                 "simulate-stagetime-pipelinerun-25",
		RevisionRunAuthor:    "patrick-hermann-sva",
		RevisionRunCreation:  "2024-01-16",
		RevisionRunCommitId:  "43232e232",
		RevisionRunRepoUrl:   "https://github.com/stuttgart-things/stuttgart-things.git",
		RevisionRunRepoName:  "stuttgart-things",
		CanFail:              false,
		Namespace:            "tektoncd",
		ResolverParams:       resolverParams,
		TimeoutPipeline:      "0h12m0s",
		Stage:                "0",
		NamePrefix:           "st",
		NameSuffix:           "083423",
		VolumeClaimTemplates: volumeClaimTemplates,
		Params:               pipelineParams,
		ListParams:           listPipelineParams,
	}

	PrintTable(testPr)
}

// func TestGetPipelineRunStatus(t *testing.T) {

// 	redisClient = goredis.NewClient(&goredis.Options{Addr: redisUrl, Password: redisPassword, DB: 0})

// 	hello := GetPipelineRunStatus(redisClient, "st-0-simulate-stagetime-2057253c5a-status")
// 	fmt.Println(hello)
// }
