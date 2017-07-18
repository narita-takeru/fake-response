package main

import (
	"fmt"
	"github.com/narita-takeru/tcpstream"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

var resRegex = regexp.MustCompile("(?ms)^(.*)\r\n\r\n{")
var contentLengthRegex = regexp.MustCompile("Content-Length: .*")

type Spec struct {
	Ports     Ports
	Extract   string
	Endpoints map[string]string
}

type Ports struct {
	Src string
	Dst string
}

func fileToSpec(path string) Spec {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var spec Spec
	err = yaml.Unmarshal(buf, &spec)
	if err != nil {
		panic(err)
	}

	return spec
}

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Please spec file path argument.")
		return
	}

	specPath := os.Args[1]
	spec := fileToSpec(specPath)

	endpointRegex := regexp.MustCompile(spec.Extract)

	fmt.Println("Start Fake Response.")

	t := tcpstream.Thread{}
	processingEndpoints := map[int]string{}
	t.SrcToDstHook = func(id, seq int, b []byte) []byte {
		group := endpointRegex.FindSubmatch(b)
		if len(group) <= 1 {
			return b
		}

		endpoint := string(group[1])
		for expectEndpoint, _ := range spec.Endpoints {
			if endpoint == expectEndpoint {
				processingEndpoints[id] = endpoint
				break
			}
		}

		return b
	}

	t.DstToSrcHook = func(id, seq int, b []byte) []byte {
		endpoint, ok := processingEndpoints[id]
		if !ok {
			return b
		}

		group := resRegex.FindSubmatch(b)
		if len(group) <= 1 {
			return b
		}

		resHeaders := string(group[1])
		replaced := spec.Endpoints[endpoint]
		contentLen := len([]byte(replaced))

		resHeaders = contentLengthRegex.ReplaceAllString(resHeaders, fmt.Sprintf("Content-Length: %d", contentLen))

		resBody := fmt.Sprintf("%s\r\n\r\n%s\n", resHeaders, replaced)
		fmt.Println("Replaced: " + endpoint)
		return []byte(resBody)
	}

	t.Do(spec.Ports.Src, spec.Ports.Dst)
}
