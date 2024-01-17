/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"strings"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

func ValidateStorePipelineRuns(pipelienRun string) (bool, map[string]string) {

	// REWRITE PIPELINERUN FOR VALIDATION
	tmpPipelienRun := strings.Split(pipelienRun, "timeouts")
	paramsAsDefaults := strings.ReplaceAll("timeouts"+tmpPipelienRun[1], "value:", "default:")
	rewrittenPipelineRun := tmpPipelienRun[0] + paramsAsDefaults

	// VALIDATE PIIPELINERUN
	validPipelineRun, pipelineRun, _ := sthingsK8s.ConvertYAMLtoPipelineRun(rewrittenPipelineRun)

	if validPipelineRun {

		prInformation := pipelineRun.Labels

		prInformation["name"] = pipelineRun.Name
		prInformation["identifier"] = pipelineRun.Name
		prInformation["revision-id"] = strings.Split(pipelineRun.Name, "-")[0]

		if pipelineRun.Annotations["canfail"] != "" {
			prInformation["canFail"] = pipelineRun.Annotations["canfail"]
		} else {
			prInformation["canFail"] = "false"
		}

		return true, prInformation

	} else {
		return false, nil
	}

}
