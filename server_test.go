/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package main

import (
	"fmt"
	"testing"
)

var (
	pipelineRun = `
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  annotations:
    canfail: "true"
  labels:
    stagetime/author: patrick-hermann-sva
    stagetime/commit: 3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0
    stagetime/repo: stuttgart-things
    stagetime/stage: "1"
  name: st-1-simulate-stagetime-1713293c5a
  namespace: stagetime-tekton
spec:
  params:
  - name: gitRepoUrl
    value: https://github.com/stuttgart-things/stageTime-server.git
  - name: gitRevision
    value: main
  - name: gitWorkspaceSubdirectory
    value: stageTime
  - name: scriptPath
    value: tests/prime.sh
  - name: scriptTimeout
    value: 25s
  pipelineRef:
    params:
    - name: pathInRepo
      value: stageTime/pipelines/simulate-stagetime-pipelineruns.yaml
    - name: revision
      value: main
    - name: url
      value: https://github.com/stuttgart-things/stuttgart-things.git
    resolver: git
  taskRunTemplate:
    podTemplate:
      securityContext:
        fsGroup: 65532
  timeouts:
    pipeline: 0h30m0s
    tasks: 0h30m0s
  workspaces:
  - name: source
    volumeClaimTemplate:
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 20Mi
        storageClassName: openebs-hostpath
`
	pipelineRuns         []string
	renderedPipelineruns = make(map[int][]string)
)

func TestCreateRevisionRun(t *testing.T) {

	// SET PIPELINERUN ON SLICE
	pipelineRuns = append(pipelineRuns, pipelineRun)
	renderedPipelineruns[0] = pipelineRuns

	// LOOP OVER PRS
	for i := 0; i < (len(renderedPipelineruns)); i++ {

		for _, pr := range renderedPipelineruns[i] {

			fmt.Println(pr)

		}
	}

}
