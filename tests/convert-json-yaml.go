package main

import (
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"
)

func main() {
	j := []byte(`{"name": "John", "age": 30}`)
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Println(string(y))

	f := string(y)
	fmt.Println(f)

	g2 := `apiVersion: tekton.dev/v1beta1
	kind: PipelineRun
	metadata:
	  name: hello
	  namespace: ansible
	`
	g := "apiVersion: tekton.dev/v1beta1\nkind: PipelineRun\nmetadata:\n  name: hello\n  namespace: ansible\nspec:"
	fmt.Println(g)

	j2, err := yaml.YAMLToJSON([]byte(strings.Replace(g2, "\t", "\n", 5)))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(j2))
	/* Output:
	{"age":30,"name":"John"}
	*/
}
