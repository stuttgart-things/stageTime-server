/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package server

import (
	"fmt"
	"strings"
	"testing"
)

var (
	listPipelineParams   = make(map[string][]string)
	pipelineParams       = make(map[string]string)
	resolverParams       = make(map[string]string)
	volumeClaimTemplates []VolumeClaimTemplate
	workspaces           []Workspace
)

func TestRenderPipelineRuns(t *testing.T) {
	resolverParams["url"] = "https://github.com/stuttgart-things/stuttgart-things.git"
	resolverParams["revision"] = "main"
	resolverParams["pathInRepo"] = "stageTime/pipelines/simulate-stagetime-pipelineruns.yaml"

	pipelineParams["gitRepoUrl"] = "https://github.com/stuttgart-things/stageTime-server.git"
	pipelineParams["gitRevision"] = "main"
	pipelineParams["gitWorkspaceSubdirectory"] = "stageTime"
	pipelineParams["scriptPath"] = "tests/prime.sh"
	pipelineParams["scriptTimeout"] = "15s"

	volumeClaimTemplates = append(volumeClaimTemplates, VolumeClaimTemplate{"source", "openebs-hostpath", "ReadWriteOnce", "20Mi"})
	workspaces = append(workspaces, Workspace{"dockerconfig", "secret", "scr-labda", "secretName"})
	fsGroup := TaskRunTemplate{65532}

	testPr := PipelineRun{
		Name:                "simulate-stagetime-pipelinerun-25",
		RevisionRunAuthor:   "patrick-hermann-sva",
		RevisionRunCreation: "2024-01-16",
		RevisionRunCommitId: "43232e232",
		RevisionRunRepoUrl:  "https://github.com/stuttgart-things/stuttgart-things.git",
		RevisionRunRepoName: "stuttgart-things",
		CanFail:             false,
		Namespace:           "tektoncd",
		// PipelineRef:         "simulate-stagetime-pipelineruns",
		ResolverParams:       resolverParams,
		TimeoutPipeline:      "0h12m0s",
		Stage:                "0",
		NamePrefix:           "st",
		NameSuffix:           "083423",
		VolumeClaimTemplates: volumeClaimTemplates,
		TaskRunTemplate:      fsGroup,
		// Workspaces:           workspaces,
		Params:     pipelineParams,
		ListParams: listPipelineParams,
	}

	renderedPipelineRun, _ := RenderPipelineRun(PipelineRunTemplate, testPr)

	fmt.Println(renderedPipelineRun)

	split := strings.Split(renderedPipelineRun, "timeouts")

	fmt.Println(split[0])
	fmt.Println(split[1])

	paramsAsDefaults := strings.ReplaceAll("timeouts"+split[1], "value:", "default:")

	renderedPipelineRun = split[0] + paramsAsDefaults

	fmt.Println(renderedPipelineRun)
	// // Using the Replace Function
	// testresults := strings.ReplaceAll(renderedPipelineRun, "value:", "default:")
	// // Display the ReplaceAll Output
	// fmt.Println(testresults)

	// // INCOMING gRPC REQUEST:
	// // {"repo_name":"stuttgart-things","pushed_at":"2023-02-20T22:40:36Z","author":"patrick","repo_url":"https://codehub.sva.de/Lab/stuttgart-things/stuttgart-things.git","commit_id":"3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0","pipelineruns":[{"name":"package-publish-helmchart","canfail":false,"stage":0,"params":"git-revision=master, git-repo-url='git@codehub.sva.de:Lab/stuttgart-things/stuttgart-things.git', gitWorkspaceSubdirectory=/helm/sthings-cluster, helm-chart-path=gitops/apps, helm-chart-name=sthings-cluster, helm-chart-tag=0.2.1, registry=scr.tiab.labda.sva.de, working-image=scr.tiab.labda.sva.de/sthings-k8s-workspace/sthings-k8s-workspace:281126-1644","listparams":"git-revision=master;hello;whatever","workspaces":"ssh-credentials=secret;codehub-ssh;secretName, source=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"},{"name":"build-kaniko-image","canfail":true,"stage":0,"params":"context=/kaniko/decksman, dockerfile=./Dockerfile, git-revision=main, gitRepoUrl='git@codehub.sva.de:Lab/stuttgart-things/dev/decksman.git', gitWorkspaceSubdirectory=/kaniko/decksman, image=scr.tiab.labda.sva.de/decksman/decksman, registry=scr.tiab.labda.sva.de, tag=0.8.66","listparams":"git-revision=master;hello;whatever","workspaces":"ssh-credentials=secret;codehub-ssh;secretName, shared-workspace=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"},{"name":"package-publish-helmchart","canfail":false,"stage":1,"params":"git-revision=master, git-repo-url='git@codehub.sva.de:Lab/stuttgart-things/stuttgart-things.git', gitWorkspaceSubdirectory=/helm/sthings-tekton, helm-chart-path=gitops/apps, helm-chart-name=sthings-tekton, helm-chart-tag=0.4.4, registry=scr.tiab.labda.sva.de, working-image=scr.tiab.labda.sva.de/sthings-k8s-workspace/sthings-k8s-workspace:281126-1644","listparams":"git-revision=master;hello;whatever","workspaces":"ssh-credentials=secret;codehub-ssh;secretName, source=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"}]}

}
