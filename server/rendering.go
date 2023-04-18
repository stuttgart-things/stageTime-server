/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package server

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"text/template"
	"time"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	revisionrun "github.com/stuttgart-things/yacht-application-server/revisionrun"
)

var (
	pipelineNamespace = os.Getenv("PIPELINE_WORKSPACE")
)

type PipelineRun struct {
	Name                string
	RevisionRunAuthor   string
	RevisionRunRepoName string
	RevisionRunRepoUrl  string
	RevisionRunCommitId string
	RevisionRunCreation string
	Namespace           string
	PipelineRef         string
	ServiceAccount      string
	Timeout             string
	Params              map[string]string
	Workspaces          []Workspace
	NamePrefix          string
	NameSuffix          string
	Stage               string
}

type Workspace struct {
	Name                   string
	WorkspaceKind          string
	WorkspaceRef           string
	WorkspaceKindShortName string
}

const PipelineRunTemplate = `
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  name: "{{ .NamePrefix }}-{{ .Stage }}-{{ .Name }}-{{ .NameSuffix }}"
  namespace: {{ .Namespace }}
  labels:
    argocd.argoproj.io/instance: tekton-runs
    yacht/commit: "{{ .RevisionRunCommitId }}"
    yacht/repo: {{ .RevisionRunRepoName }}
    yacht/author: {{ .RevisionRunAuthor }}
    tekton.dev/pipeline: {{ .Name }}
spec:
  serviceAccountName: {{ .ServiceAccount }}
  timeout: {{ .Timeout }}
  pipelineRef:
    name: {{ .PipelineRef }}
  podTemplate:
    hostAliases:
      - ip: 10.10.210.114
        hostnames:
          - codehub.sva.de
  params:{{ range $name, $value := .Params }}
  - name: {{ $name }}
    value: {{ $value }}{{ end }}
  workspaces:{{ range .Workspaces }}
  - name: {{ .Name }}
    {{ .WorkspaceKind }}:
      {{ .WorkspaceKindShortName }}: {{ .WorkspaceRef }}{{ end }}
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

func RenderPipelineRuns(req *revisionrun.CreateRevisionRunRequest) (renderedPipelineruns map[int][]string) {

	dt := time.Now()

	renderedPipelineruns = make(map[int][]string)

	for _, pipelinerun := range req.Pipelineruns {

		pipelineParams := make(map[string]string)
		var pipelineWorkspaces []Workspace

		fmt.Println(pipelinerun.Name)
		fmt.Println(pipelinerun.Stage)

		paramValues := strings.Split(pipelinerun.Params, ",")

		for i, v := range paramValues {
			values := strings.Split(v, "=")

			pipelineParams[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
			fmt.Println(i)
			fmt.Println(strings.TrimSpace(values[0]))
			fmt.Println(strings.TrimSpace(values[1]))

		}

		workspaces := strings.Split(pipelinerun.Workspaces, ",")

		for _, v := range workspaces {
			values := strings.Split(v, "=")
			workspaces := strings.Split(values[1], ";")

			pipelineWorkspaces = append(pipelineWorkspaces, Workspace{strings.TrimSpace(values[0]), strings.TrimSpace(workspaces[0]), strings.TrimSpace(workspaces[1]), strings.TrimSpace(workspaces[2])})

		}

		fmt.Println(pipelineWorkspaces)

		pr := PipelineRun{
			Name:                pipelinerun.Name,
			RevisionRunAuthor:   req.Author,
			RevisionRunCreation: req.PushedAt,
			RevisionRunCommitId: req.CommitId,
			RevisionRunRepoUrl:  req.RepoUrl,
			RevisionRunRepoName: req.RepoName,
			Namespace:           pipelineNamespace,
			PipelineRef:         pipelinerun.Name,
			ServiceAccount:      "default",
			Timeout:             "1h",
			Params:              pipelineParams,
			Stage:               fmt.Sprintf("%f", math.RoundToEven(pipelinerun.Stage)),
			Workspaces:          pipelineWorkspaces,
			NamePrefix:          "y",
			NameSuffix:          dt.Format("020405") + req.CommitId[0:4],
		}

		tmpl, err := template.New("pipelinerun").Parse(PipelineRunTemplate)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer

		err = tmpl.Execute(&buf, pr)

		if err != nil {
			log.Fatalf("execution: %s", err)
		}

		fmt.Println(buf.String())
		renderedPipelineruns[int(pipelinerun.Stage)] = append(renderedPipelineruns[int(pipelinerun.Stage)], buf.String())

	}

	return
}

func RenderOutputData(template, delimiter string, templateKeyValues map[string]string) {

	// convert string to interface map
	templateValueData := make(map[string]interface{})
	for k, v := range templateKeyValues {
		templateValueData[k] = v
	}

	// render template
	renderedTemplate, err := sthingsBase.RenderTemplateInline(template, "missingkey=zero", Patterns[delimiter].begin, Patterns[delimiter].end, templateValueData)

	if err != nil {
		log.Fatal(err)
	}

	renderedData := strings.ReplaceAll(string(renderedTemplate), "&#34;", " ")

	fmt.Println(renderedData)

}
