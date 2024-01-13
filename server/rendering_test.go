/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package server

import (
	"fmt"
	"testing"
)

var (
	listPipelineParams = make(map[string][]string)
	pipelineParams     = make(map[string]string)
	pipelineWorkspaces []Workspace

	testPr = PipelineRun{
		Name:                "hello",
		RevisionRunAuthor:   "hello",
		RevisionRunCreation: "hello",
		RevisionRunCommitId: "hello",
		RevisionRunRepoUrl:  "hello",
		RevisionRunRepoName: "hello",
		CanFail:             false,
		Namespace:           "tektoncd",
		PipelineRef:         "whatever",
		TimeoutPipeline:     "0h12m0s",
		Stage:               "0",
		NamePrefix:          "st",
		NameSuffix:          "083423",
		Workspaces:          pipelineWorkspaces,
		Params:              pipelineParams,
		ListParams:          listPipelineParams,
	}
)

func TestRenderPipelineRuns(t *testing.T) {

	renderedPipelineRun, _ := RenderPipelineRun(PipelineRunTemplate, testPr)
	fmt.Println(renderedPipelineRun)
	// INCOMING gRPC REQUEST:
	// {"repo_name":"stuttgart-things","pushed_at":"2023-02-20T22:40:36Z","author":"patrick","repo_url":"https://codehub.sva.de/Lab/stuttgart-things/stuttgart-things.git","commit_id":"3c5ac44c6fec00989c7e27b36630a82cdfd26e3b0","pipelineruns":[{"name":"package-publish-helmchart","canfail":false,"stage":0,"params":"git-revision=master, git-repo-url='git@codehub.sva.de:Lab/stuttgart-things/stuttgart-things.git', gitWorkspaceSubdirectory=/helm/sthings-cluster, helm-chart-path=gitops/apps, helm-chart-name=sthings-cluster, helm-chart-tag=0.2.1, registry=scr.tiab.labda.sva.de, working-image=scr.tiab.labda.sva.de/sthings-k8s-workspace/sthings-k8s-workspace:281126-1644","listparams":"git-revision=master;hello;whatever","workspaces":"ssh-credentials=secret;codehub-ssh;secretName, source=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"},{"name":"build-kaniko-image","canfail":true,"stage":0,"params":"context=/kaniko/decksman, dockerfile=./Dockerfile, git-revision=main, gitRepoUrl='git@codehub.sva.de:Lab/stuttgart-things/dev/decksman.git', gitWorkspaceSubdirectory=/kaniko/decksman, image=scr.tiab.labda.sva.de/decksman/decksman, registry=scr.tiab.labda.sva.de, tag=0.8.66","listparams":"git-revision=master;hello;whatever","workspaces":"ssh-credentials=secret;codehub-ssh;secretName, shared-workspace=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"},{"name":"package-publish-helmchart","canfail":false,"stage":1,"params":"git-revision=master, git-repo-url='git@codehub.sva.de:Lab/stuttgart-things/stuttgart-things.git', gitWorkspaceSubdirectory=/helm/sthings-tekton, helm-chart-path=gitops/apps, helm-chart-name=sthings-tekton, helm-chart-tag=0.4.4, registry=scr.tiab.labda.sva.de, working-image=scr.tiab.labda.sva.de/sthings-k8s-workspace/sthings-k8s-workspace:281126-1644","listparams":"git-revision=master;hello;whatever","workspaces":"ssh-credentials=secret;codehub-ssh;secretName, source=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"}]}

}
