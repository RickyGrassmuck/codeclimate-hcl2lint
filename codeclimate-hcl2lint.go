package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rigrassm/codeclimate-hcl2lint/Utils"
)

var parser = hclparse.NewParser()

func main() {
	rootPath := "/code/"

	config, err := Utils.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(1)
	}

	analysisFiles, err := Utils.HCL2FileWalk(rootPath, Utils.IncludePaths(rootPath, config))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %s", err)
		os.Exit(1)
	}

	for _, path := range analysisFiles {
		processFile(path, nil)
	}
}

func parseHCL2Error(diag *hcl.Diagnostic) Utils.Issue {
	var location Utils.Location
	var locationParse []string

	firstParse := strings.Split(diag.Subject.String(), ":")

	if len(firstParse) == 2 {
		location.Path = firstParse[0]

		locationParse = strings.Split(firstParse[1], ",")

		location.Lines.Begin, _ = strconv.Atoi(locationParse[0])
		location.Lines.End, _ = strconv.Atoi(locationParse[0])
		positionParse := strings.Split(locationParse[1], "-")

		beginColumn, _ := strconv.Atoi(positionParse[0])
		endColumn, _ := strconv.Atoi(positionParse[1])

		location.Positions.Begin = &Utils.LineColumn{
			Line:   location.Lines.Begin,
			Column: beginColumn,
		}
		location.Positions.End = &Utils.LineColumn{
			Line:   location.Lines.Begin,
			Column: endColumn,
		}
		issue := Utils.Issue{
			Type:              "issue",
			Check:             codeClimateCheckName(diag),
			Description:       diag.Detail,
			RemediationPoints: 50000,
			Categories:        []string{"Syntax"},
			Location:          &location,
		}
		return issue
	}
	return Utils.Issue{}
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
		var issue Utils.Issue = parseHCL2Error(problem)
		Utils.PrintIssue(&issue)
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
