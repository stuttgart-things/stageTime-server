/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"reflect"

	"github.com/stretchr/testify/assert"

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
)

func TestValidateStorePipelineRuns(t *testing.T) {

	assert := assert.New(t)

	expectedPrInformation := map[string]string{"identifier": "st", "name": "st-1-simulate-stagetime-1713293c5a", "revision-id": "st", "stagetime/author": "patrick-hermann-sva", "stagetime/commit": "3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0", "stagetime/repo": "stuttgart-things", "stagetime/stage": "1"}

	valid, prInformation := ValidateStorePipelineRuns(pipelineRun)

	// map[identifier:st name:st-1-simulate-stagetime-1713293c5a revision-id:st stagetime/author:patrick-hermann-sva stagetime/commit:3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0 stagetime/repo:stuttgart-things stagetime/stage:1]
	if !reflect.DeepEqual(expectedPrInformation, prInformation) {
		t.Errorf("error")
	}

	assert.Equal(true, valid)

}
