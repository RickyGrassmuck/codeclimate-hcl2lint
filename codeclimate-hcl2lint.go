package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

var parser = hclparse.NewParser()

func main() {
	rootPath := "/code/"

	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(1)
	}

	analysisFiles, err := HCL2FileWalk(rootPath, IncludePaths(rootPath, config))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %s", err)
		os.Exit(1)
	}

	for _, path := range analysisFiles {
		processFile(path, nil)
	}
}

func parseHCL2Error(diag *hcl.Diagnostic) Issue {
	var locationParse []string

	firstParse := strings.Split(diag.Subject.String(), ":")
	fmt.Println(firstParse)

	if len(firstParse) == 2 {
		// location.Path = firstParse[0]

		locationParse = strings.Split(firstParse[1], ",")
		positionParse := strings.Split(locationParse[1], "-")
		line, _ := strconv.Atoi(locationParse[0])
		startColumn, _ := strconv.Atoi(positionParse[0])
		endColumn, _ := strconv.Atoi(positionParse[1])

		location := &Location{
			Path: firstParse[0],
			Lines: &LinesOnlyPosition{
				Begin: line,
				End:   line,
			},
			Positions: &LineColumnPosition{
				Begin: &LineColumn{
					Line:   line,
					Column: startColumn,
				},
				End: &LineColumn{
					Line:   line,
					Column: endColumn,
				},
			},
		}

		issue := Issue{
			Type:              "issue",
			Check:             codeClimateCheckName(diag),
			Description:       diag.Detail,
			RemediationPoints: 50000,
			Categories:        []string{"Syntax"},
			Location:          location,
		}
		return issue
	}
	return Issue{}
}

func processFile(fn string, in *os.File) {
	var err error
	if in == nil {
		in, err = os.Open(fn)
		if err != nil {
			fmt.Errorf("failed to open %s: %s", fn, err)
			os.Exit(1)
		}
	}

	inSrc, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Errorf("failed to read %s: %s", fn, err)
		os.Exit(1)
	}

	_, problems := parser.ParseHCL(inSrc, fn)
	for _, problem := range problems {
		var issue Issue = parseHCL2Error(problem)
		PrintIssue(&issue)
	}
}

func codeClimateCheckName(l *hcl.Diagnostic) string {
	var sev string
	if l.Severity == hcl.DiagError {
		sev = "error"
	} else {
		sev = "warn"
	}
	return "HCL2Lint/" + strings.Title(sev) + "/" + l.Summary
}
