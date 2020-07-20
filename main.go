package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(fname string) []byte {
	data, err := ioutil.ReadFile(fname)
	check(err)
	return data
}

func readYaml(fname string) interface{} {
	var data interface{}
	err := yaml.Unmarshal(readFile(fname), &data)
	check(err)
	return data
}

func getPipeline(pipeline string) map[interface{}]interface{} {
	idx := readYaml("index.yml").(map[interface{}]interface{})
	if val, ok := idx["pipeline"]; ok {
		if vall, okk := val.(map[interface{}]interface{})[pipeline]; okk {
			return readYaml("pipelines/" + vall.(string) + ".yml").(map[interface{}]interface{})
		}
	}
	return make(map[interface{}]interface{})
}

func getTemplate(pipeline map[interface{}]interface{}) []interface{} {
	return readYaml("templates/" + pipeline["template"].(string) + ".yml").([]interface{})
}

func renderPipeline(pipeline map[interface{}]interface{}, tpl []interface{}) []interface{} {
	var empty interface{}
	result := make([]interface{}, 0)
	for _, module := range tpl {
		tmpl := template.Must(template.ParseFiles("modules/" + module.(string) + ".tmpl"))
		buf := &bytes.Buffer{}
		if val, ok := pipeline[module]; ok {
			tmpl.Execute(buf, val)
		} else {
			tmpl.Execute(buf, empty)
		}
		data := make(map[interface{}]interface{})
		yaml.Unmarshal([]byte(buf.String()), &data)
		str, err := yaml.Marshal(&data)
		check(err)
		result = append(result, string(str))
	}
	return result
}

func run(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(204)
		}
	}()
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	check(err)
	pipelineName := data["repo"].(map[string]interface{})["slug"]
	fmt.Println(pipelineName)
	pipeline := getPipeline(pipelineName.(string))
	fmt.Println(pipeline)
	tpl := getTemplate(pipeline)
	result := renderPipeline(pipeline, tpl)
	str := ""
	for _, v := range result {
		str += "---\n" + v.(string)
	}
	content := make(map[string]string)
	content["data"] = str
	js, err := json.Marshal(content)
	check(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/", run)
	http.ListenAndServe(":8080", nil)
}
