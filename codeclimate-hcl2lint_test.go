package main

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

func TestCodeClimateCheckName(t *testing.T) {
	var testDiag = &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "test_summary",
	}
	expectedName := "HCL2Lint/" + strings.Title("error") + "/test_summary"
	checkName := codeClimateCheckName(testDiag)
	assert.Equal(t, expectedName, checkName, "The name should match")
}
