package main

import "fmt"
import "os"
import "strings"
import "encoding/json"
import "path/filepath"
import "io/ioutil"

type LinesOnlyPosition struct {
	Begin int `json:"begin,omitempty"`
	End   int `json:"end,omitempty"`
}

type LineColumnPosition struct {
	Begin *LineColumn `json:"begin,omitempty"`
	End   *LineColumn `json:"end,omitempty"`
}

type LineColumn struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

type Location struct {
	Path      string              `json:"path"`
	Lines     *LinesOnlyPosition  `json:"lines,omitempty"`
	Positions *LineColumnPosition `json:"positions,omitempty"`
}

type Issue struct {
	Type              string    `json:"type"`
	Check             string    `json:"check_name"`
	Description       string    `json:"description"`
	RemediationPoints int32     `json:"remediation_points"`
	Location          *Location `json:"location"`
	Categories        []string  `json:"categories"`
}

type Warning struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Location    *Location `json:"location"`
}

var HCL2FileExtensions = []string{".hcl", ".tf", ".nomad"}

func HCL2FileWalk(rootPath string, includePaths []string) (fileList []string, err error) {
	walkFunc := func(path string, f os.FileInfo, err error) error {
		if suffixInArr(path, HCL2FileExtensions) && prefixInArr(path, includePaths) {
			fileList = append(fileList, path)
			return nil
		}
		return err
	}

	err = filepath.Walk(rootPath, walkFunc)
	return fileList, err
}

func LoadConfig() (config map[string]interface{}, err error) {
	var parsedConfig map[string]interface{}
	var configPath string
	var ok bool

	if configPath, ok = os.LookupEnv("CC_CONFIG"); !ok {
		configPath = "/config.json"
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &parsedConfig)

		if err != nil {
			return nil, err
		}
	}

	return parsedConfig, nil
}

func IncludePaths(rootPath string, config map[string]interface{}) []string {
	if iArr, ok := config["include_paths"].([]interface{}); ok {
		paths := make([]string, len(iArr))
		for i, iVal := range iArr {
			if strVal, ok := iVal.(string); ok {
				paths[i] = filepath.Join(rootPath, strVal)
			} else {
				fmt.Fprintf(os.Stderr, "include_paths should be an array of strings, but an invalid value was encountered (%s) in %s\n", iVal, iArr)
				os.Exit(1)
			}
		}
		return paths
	}
	return []string{rootPath} // will be a prefix of any path
}

func PrintIssue(issue *Issue) (err error) {
	jsonOutput, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	jsonOutput = append(jsonOutput, 0)
	os.Stdout.Write(jsonOutput)

	return nil
}

func PrintWarning(warning *Warning) (err error) {
	warning.Type = "warning"
	jsonOutput, err := json.Marshal(warning)
	if err != nil {
		return err
	}

	jsonOutput = append(jsonOutput, 0)
	os.Stdout.Write(jsonOutput)

	return nil
}
func suffixInArr(str string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}
func prefixInArr(str string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}
