/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

type RevisionRunStatus struct {
	RevisionRun       string
	CountStages       int
	CountPipelineRuns int
	LastUpdated       string
	Status            string
}

type StageStatus struct {
	RevisionRun       string
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
