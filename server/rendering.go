/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	revisionrun "github.com/stuttgart-things/stageTime-server/revisionrun"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	pipelineNamespace = os.Getenv("PIPELINE_WORKSPACE")
)

type PipelineRun struct {
	Name                 string
	RevisionRunAuthor    string
	RevisionRunRepoName  string
	RevisionRunRepoUrl   string
	RevisionRunCommitId  string
	RevisionRunCreation  string
	Namespace            string
	PipelineRef          string
	ServiceAccount       string
	Timeout              string
	Params               map[string]string
	ListParams           map[string][]string
	Workspaces           []Workspace
	NamePrefix           string
	NameSuffix           string
	Stage                string
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

const PipelineRunTemplate = `
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  name: "{{ .NamePrefix }}-{{ .Stage }}-{{ .Name }}-{{ .NameSuffix }}"
  namespace: {{ .Namespace }}
  labels:
    argocd.argoproj.io/instance: tekton-runs
    stagetime/commit: "{{ .RevisionRunCommitId }}"
    stagetime/repo: {{ .RevisionRunRepoName }}
    stagetime/author: {{ .RevisionRunAuthor }}
    stagetime/stage: "{{ .Stage }}"
    tekton.dev/pipeline: {{ .PipelineRef }}
spec:
  serviceAccountName: {{ .ServiceAccount }}
  timeout: {{ .Timeout }}
  pipelineRef:
    name: {{ .PipelineRef }}
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

//   stages: {{ range .Stages }}
//     - "{{ . }}"{{ end }}
//   pipelineRuns: {{ range .PipelineRuns }}
//     - "{{ . }}"{{ end }}

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

		listPipelineParams := make(map[string][]string)
		pipelineParams := make(map[string]string)
		var pipelineWorkspaces []Workspace

		paramValues := strings.Split(pipelinerun.Params, ",")
		for _, v := range paramValues {
			values := strings.Split(v, "=")
			pipelineParams[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
		}

		for _, v := range strings.Split(pipelinerun.Listparams, ",") {

			keyValues := strings.Split(v, "=")
			var values []string

			for _, v := range strings.Split(strings.TrimSpace(keyValues[1]), ";") {
				values = append(values, v)
			}
			listPipelineParams[strings.TrimSpace(keyValues[0])] = values
		}

		workspaces := strings.Split(pipelinerun.Workspaces, ",")

		for _, v := range workspaces {
			values := strings.Split(v, "=")
			workspaces := strings.Split(values[1], ";")
			pipelineWorkspaces = append(pipelineWorkspaces, Workspace{strings.TrimSpace(values[0]), strings.TrimSpace(workspaces[0]), strings.TrimSpace(workspaces[1]), strings.TrimSpace(workspaces[2])})
		}

		pr := PipelineRun{
			Name:                pipelinerun.Name,
			RevisionRunAuthor:   gRPCRequest.Author,
			RevisionRunCreation: gRPCRequest.PushedAt,
			RevisionRunCommitId: gRPCRequest.CommitId,
			RevisionRunRepoUrl:  gRPCRequest.RepoUrl,
			RevisionRunRepoName: gRPCRequest.RepoName,
			Namespace:           pipelineNamespace,
			PipelineRef:         pipelinerun.Name,
			ServiceAccount:      "default",
			Timeout:             "1h",
			Params:              pipelineParams,
			ListParams:          listPipelineParams,
			Stage:               fmt.Sprintf("%v", pipelinerun.Stage),
			NamePrefix:          "st",
			NameSuffix:          dt.Format("020405") + gRPCRequest.CommitId[0:4],
			Workspaces:          pipelineWorkspaces,
		}

		// RENDERING
		var buf bytes.Buffer
		tmpl, err := template.New("pipelinerun").Parse(PipelineRunTemplate)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(&buf, pr)
		if err != nil {
			log.Fatalf("execution: %s", err)
		}

		// ADD RENDERED PRS TO REVISIONRUN
		renderedPipelineruns[int(pipelinerun.Stage)] = append(renderedPipelineruns[int(pipelinerun.Stage)], buf.String())
	}

	return
}

// TEST DATA - TO BE REPLACED
func RenderRevisionRunCR() (renderedCR []byte) {

	cr := make(map[string]interface{})
	cr["Name"] = "44c6fec0098"
	cr["Namespace"] = "tekton"
	cr["Repository"] = "stuttgart-things"
	cr["RevisionRun"] = "44c6fec0098-123"
	cr["Stages"] = []string{"0", "1", "2"}
	cr["PipelineRuns"] = []string{"0-2321", "1-312", "2-321312"}

	renderedCR, _ = sthingsBase.RenderTemplateInline(RevisionRunTemplate, "missingkey=error", "{{", "}}", cr)

	return renderedCR
}
