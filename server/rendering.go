/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	revisionrun "github.com/stuttgart-things/stageTime-server/revisionrun"
)

var (
	now               = time.Now()
	pipelineNamespace = os.Getenv("PIPELINE_WORKSPACE")
)

type PipelineRun struct {
	Name                 string
	RevisionRunAuthor    string
	RevisionRunRepoName  string
	RevisionRunRepoUrl   string
	RevisionRunCommitId  string
	RevisionRunCreation  string
	RevisionRunDate      string
	CanFail              bool
	ResolverParams       map[string]string
	Namespace            string
	PipelineRef          string
	TimeoutPipeline      string
	Params               map[string]string
	ListParams           map[string][]string
	Workspaces           []Workspace
	NamePrefix           string
	NameSuffix           string
	Stage                string
	TaskRunTemplate      TaskRunTemplate
	VolumeClaimTemplates []VolumeClaimTemplate
}

type Workspace struct {
	Name                   string
	WorkspaceKind          string
	WorkspaceRef           string
	WorkspaceKindShortName string
}

type VolumeClaimTemplate struct {
	Name             string
	StorageClassName string
	AccessModes      string
	Storage          string
}

type TaskRunTemplate struct {
	fsGroup int
}

const PipelineRunTemplate = `
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: {{ .NamePrefix }}-{{ .Stage }}-{{ .Name }}-{{ .NameSuffix }}
  namespace: {{ .Namespace }}
  annotations:
    canfail: "{{ .CanFail }}"
  labels:
    stagetime/commit: "{{ .RevisionRunCommitId }}"
    stagetime/repo: {{ .RevisionRunRepoName }}
    stagetime/author: {{ .RevisionRunAuthor }}
    stagetime/stage: "{{ .Stage }}"
    stagetime/date: "{{ .RevisionRunDate }}"
spec:{{ if .TaskRunTemplate }}
  taskRunTemplate:
    podTemplate:
      securityContext:
        fsGroup: 65532{{ end }}
  pipelineRef:
    {{ if .PipelineRef }}name: {{ .PipelineRef }}{{ else }}resolver: git{{ end }}
    params:{{ range $name, $value := .ResolverParams }}
      - name: {{ $name }}
        value: {{ $value }}{{ end }}
  timeouts:
    pipeline: "{{ .TimeoutPipeline }}"
    tasks: "{{ .TimeoutPipeline }}"
  params:{{ range $name, $value := .Params }}
    - name: {{ $name }}
      value: {{ $value }}{{ end }}{{ if .ListParams }}{{ range $name, $values := .ListParams }}
    - name: {{ $name }}
	  value: {{ range $values }}
        - {{ . }}{{ end }}{{ end }}{{ end }}
  workspaces:{{ range .Workspaces }}
    - name: {{ .Name }}
      {{ .WorkspaceKind }}:
        {{ .WorkspaceKindShortName }}: {{ .WorkspaceRef }}{{ end }}{{ if .VolumeClaimTemplates }}{{ range .VolumeClaimTemplates }}
    - name: {{ .Name }}
      volumeClaimTemplate:
        spec:
          storageClassName: {{ .StorageClassName }}
          accessModes:
            - {{ .AccessModes }}
          resources:
            requests:
              storage: {{ .Storage }}{{ end }}{{ end }}
`

const RevisionRunTemplate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
data:
  revisionRun: {{ .RevisionRun }}
  repository: {{ .Repository }}
  revision: {{ .Repository }}
`

type VariableDelimiter struct {
	begin        string `mapstructure:"begin"`
	end          string `mapstructure:"end"`
	regexPattern string `mapstructure:"regex-pattern"`
}

var Patterns = map[string]VariableDelimiter{
	"curly":  VariableDelimiter{"{{", "}}", `\{\{(.*?)\}\}`},
	"square": VariableDelimiter{"[[", "]]", `\[\[(.*?)\]\]`},
}

func RenderPipelineRuns(gRPCRequest *revisionrun.CreateRevisionRunRequest) (renderedPipelineruns map[int][]string) {

	// GET CURRENT TIME
	dt := time.Now()

	// INIT PR MAP
	renderedPipelineruns = make(map[int][]string)

	// LOOP OVER PR MAP
	for _, pipelinerun := range gRPCRequest.Pipelineruns {

		// DECLARE VARIABLES
		listPipelineParams := make(map[string][]string)
		resolverParams := make(map[string]string)
		pipelineParams := make(map[string]string)
		var pipelineWorkspaces []Workspace
		var pipelinevolumeClaimTemplates []VolumeClaimTemplate

		// SET RESOLVER PARAMS
		resolverValues := strings.Split(pipelinerun.ResolverParams, ",")
		for _, v := range resolverValues {
			values := strings.Split(v, "=")
			resolverParams[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
		}

		// SET PARAMS
		paramValues := strings.Split(pipelinerun.Params, ",")
		for _, v := range paramValues {
			values := strings.Split(v, "=")
			pipelineParams[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
		}

		// SET LIST PARAMETERS IF GIVEN
		if pipelinerun.Listparams != "" {
			for _, v := range strings.Split(pipelinerun.Listparams, ",") {

				keyValues := strings.Split(v, "=")
				var values []string

				for _, v := range strings.Split(strings.TrimSpace(keyValues[1]), ";") {
					values = append(values, v)
				}
				listPipelineParams[strings.TrimSpace(keyValues[0])] = values
			}
		}

		// SET WORKSPACES IF GIVEN
		if len(pipelinerun.Workspaces) > 0 {
			workspaces := strings.Split(pipelinerun.Workspaces, ",")
			for _, v := range workspaces {
				values := strings.Split(v, "=")
				workspaces := strings.Split(values[1], ";")
				pipelineWorkspaces = append(pipelineWorkspaces, Workspace{strings.TrimSpace(values[0]), strings.TrimSpace(workspaces[0]), strings.TrimSpace(workspaces[1]), strings.TrimSpace(workspaces[2])})
			}
		}

		// SET VOLUMECLAIMTEMPLATES IF GIVEN
		if len(pipelinerun.VolumeClaimTemplates) > 0 {
			volumeClaimTemplates := strings.Split(pipelinerun.VolumeClaimTemplates, ",")
			for _, v := range volumeClaimTemplates {
				values := strings.Split(v, "=")
				claims := strings.Split(values[1], ";")
				pipelinevolumeClaimTemplates = append(pipelinevolumeClaimTemplates, VolumeClaimTemplate{strings.TrimSpace(values[0]), strings.TrimSpace(claims[0]), strings.TrimSpace(claims[1]), strings.TrimSpace(claims[2])})
			}
		}

		pr := PipelineRun{
			Name:                pipelinerun.Name,
			RevisionRunAuthor:   gRPCRequest.Author,
			RevisionRunCreation: gRPCRequest.PushedAt,
			RevisionRunCommitId: gRPCRequest.CommitId,
			RevisionRunRepoUrl:  gRPCRequest.RepoUrl,
			RevisionRunRepoName: gRPCRequest.RepoName,
			RevisionRunDate:     now.Format("2006-01-02 15:04:05"),
			CanFail:             pipelinerun.Canfail,
			Namespace:           pipelineNamespace,
			// PipelineRef:         pipelinerun.Name,
			ResolverParams:       resolverParams,
			TimeoutPipeline:      "0h30m0s",
			Params:               pipelineParams,
			ListParams:           listPipelineParams,
			Stage:                fmt.Sprintf("%v", pipelinerun.Stage),
			NamePrefix:           "st",
			NameSuffix:           dt.Format("020405") + gRPCRequest.CommitId[0:4],
			Workspaces:           pipelineWorkspaces,
			VolumeClaimTemplates: pipelinevolumeClaimTemplates,
		}

		// RENDER REVISIONRUN
		renderedPipelineRun, _ := RenderPipelineRun(PipelineRunTemplate, pr)

		// ADD RENDERED PRS TO REVISIONRUN
		renderedPipelineruns[int(pipelinerun.Stage)] = append(renderedPipelineruns[int(pipelinerun.Stage)], renderedPipelineRun)
	}

	return
}

func RenderPipelineRun(PipelineRunTemplate string, pr PipelineRun) (string, error) {

	var buf bytes.Buffer
	tmpl, err := template.New("pipelineRun").Option("missingkey=error").Parse(PipelineRunTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(&buf, pr)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

	return buf.String(), nil
}
