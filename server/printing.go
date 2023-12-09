/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

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
