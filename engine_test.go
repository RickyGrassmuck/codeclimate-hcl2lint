package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testIssue = Issue{
	Type:              "issue",
	Check:             "HCL2Lint/Error/Missing key/value separator",
	Description:       "Expected an equals sign (\"=\") to mark the beginning of the attribute value. If you intended to given an attribute name containing periods or spaces, write the name in quotes to create a string literal.",
	RemediationPoints: 50000,
	Categories:        []string{"Syntax"},
	Location: &Location{
		Path: "/code/base-images.pkr.hcl",
		Lines: &LinesOnlyPosition{
			Begin: 15,
			End:   15,
		},
		Positions: &LineColumnPosition{
			Begin: &LineColumn{
				Line:   15,
				Column: 23,
			},
			End: &LineColumn{
				Line:   15,
				Column: 32,
			},
		},
	},
}

var testWarning = Warning{
	Type:        "warning",
	Description: "This is a test warning",
	Location: &Location{
		Path: "/code/base-images.pkr.hcl",
		Lines: &LinesOnlyPosition{
			Begin: 15,
			End:   15,
		},
		Positions: &LineColumnPosition{
			Begin: &LineColumn{
				Line:   15,
				Column: 23,
			},
			End: &LineColumn{
				Line:   15,
				Column: 32,
			},
		},
	},
}

func TestHCL2FileWalk(t *testing.T) {
	fileList, err := HCL2FileWalk("./", []string{"test/data/"})
	testFileList := []string{"test/data/testBad.hcl", "test/data/testGood.hcl"}
	assert.NoError(t, err, "The test should not error")
	assert.Equal(t, testFileList, fileList)
}

func TestLoadConfig(t *testing.T) {
	validConf := map[string]interface{}{
		"include_paths": []interface{}{"/code", "test/data/"},
	}
	os.Setenv("CC_CONFIG", "test_config.json")
	config, err := LoadConfig()
	assert.NoError(t, err, "It should not error")
	assert.Equal(t, validConf, config, "It should be a valid config")
}

func TestPrintIssue(t *testing.T) {
	var printNoError = PrintIssue(&testIssue)
	assert.NoError(t, printNoError, "This should not error")
}
func TestWarning(t *testing.T) {
	var printNoError = PrintWarning(&testWarning)
	assert.NoError(t, printNoError, "This should not error")
}
func TestSuffixInArr(t *testing.T) {
	var suffixArray = []string{".hcl", ".tf", ".nomad"}
	var testSuffixPass = suffixInArr(".hcl", suffixArray)
	var testSuffixFail = suffixInArr(".nothcl", suffixArray)

	assert.True(t, testSuffixPass, "It should find the suffix")
	assert.False(t, testSuffixFail, "It should not find the suffix")
}
func TestPrefixInArr(t *testing.T) {
	var prefixArray = []string{"/code/", "/more-code/"}
	var testPrefixPass = prefixInArr("/code/", prefixArray)
	var testPrefixFail = prefixInArr("/not-code/", prefixArray)

	assert.True(t, testPrefixPass, "It should find the prefix")
	assert.False(t, testPrefixFail, "It should not find the prefix")
}
