package internal_test

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/internal"
	"gopkg.in/yaml.v3"
)

func cmdToArgs(cmd string) []string {
	return strings.Fields(cmd)
}

func Test_NoArgs(t *testing.T) {
	require.Equal(t, 101, internal.Run([]string{}, io.Discard, io.Discard))
}

func Test_OneArg(t *testing.T) {
	require.Equal(t, 101, internal.Run(cmdToArgs("oasdiff"), io.Discard, io.Discard))
}

func Test_NoRevision(t *testing.T) {
	require.Equal(t, 101, internal.Run(cmdToArgs("oasdiff -base base.yaml"), io.Discard, io.Discard))
}

func Test_InvalidArg(t *testing.T) {
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff -base data/openapi-test1.yaml -revision data/openapi-test1.yaml -deprecation-days 23s"), io.Discard, io.Discard))
}

func Test_BasicDiff(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -exclude-elements endpoints"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_DiffJson(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -format json -exclude-elements endpoints"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, json.Unmarshal(stdout.Bytes(), &bc))
}

func Test_DiffHtml(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -format html"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), `<h3 id="new-endpoints-none">New Endpoints: None</h3>`)
}

func Test_DiffText(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -format text"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), `### New Endpoints: None`)
}

func Test_Summary(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -summary"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), `diff: true`)
}

func Test_InvalidFile(t *testing.T) {
	require.Equal(t, 102, internal.Run(cmdToArgs("oasdiff -base no-file -revision ../data/openapi-test3.yaml"), io.Discard, io.Discard))
}

func Test_DiffInvalidFormat(t *testing.T) {
	require.Equal(t, 108, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -format xxx"), io.Discard, io.Discard))
}

func Test_BreakingChangesIncludeChecks(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/run_test/breaking_changes_include_checks_base.yaml -revision ../data/run_test/breaking_changes_include_checks_revision.yaml -check-breaking -include-checks response-non-success-status-removed,api-tag-removed -format json"), &stdout, io.Discard))
	bc := checker.BackwardCompatibilityErrors{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Len(t, bc, 2)
	for _, c := range bc {
		require.Equal(t, c.Level, checker.ERR)
	}
}

func Test_BasicBreakingChanges(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking"), io.Discard, io.Discard))
}

func Test_BreakingChangesInvalidFormat(t *testing.T) {
	require.Equal(t, 108, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -format html"), io.Discard, io.Discard))
}

func Test_BreakingChangesYaml(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -format yaml"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesJson(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -format json"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, json.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesText(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking"), &stdout, io.Discard))
	var bc interface{}
	require.Error(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Error(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesFailOnDiff(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -fail-on-diff"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarns(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -fail-on-diff -fail-on-warns"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarnsErrsOnly(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test2.yaml -revision ../data/openapi-test4.yaml -check-breaking -fail-on-diff -fail-on-warns"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnDiffNoDiff(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test1.yaml -check-breaking -fail-on-diff"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarnsNoDiff(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test1.yaml -check-breaking -fail-on-diff -fail-on-warns"), io.Discard, io.Discard))
}

func Test_BreakingChangesIgnoreErrs(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -err-ignore ../data/ignore-err-example.txt -format json"), &stdout, io.Discard))
	bc := checker.BackwardCompatibilityErrors{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Len(t, bc, 5)
}

func Test_BreakingChangesIgnoreErrsAndWarns(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -err-ignore ../data/ignore-err-example.txt -warn-ignore ../data/ignore-warn-example.txt -format json"), &stdout, io.Discard))
	bc := checker.BackwardCompatibilityErrors{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Len(t, bc, 4)
}

func Test_BreakingChangesInvalidIgnoreFile(t *testing.T) {
	require.Equal(t, 121, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -err-ignore no-file"), io.Discard, io.Discard))
}

func Test_ComposedMode(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -composed -base ../data/composed/base/*.yaml -revision ../data/composed/revision/*.yaml -exclude-elements endpoints"), &stdout, io.Discard))
	var bc interface{}
	require.NoError(t, yaml.Unmarshal(stdout.Bytes(), &bc))
	require.Equal(t, map[string]interface{}{"paths": map[string]interface{}{"deleted": []interface{}{"/api/old-test"}}}, bc)
}

func Test_Help(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -help"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), "Usage of oasdiff")
}

func Test_HelpShortcut(t *testing.T) {
	var stdout bytes.Buffer
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff -h"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), "Usage of oasdiff")
}

func Test_Version(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -version"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), "oasdiff version:")
}

func Test_StripPrefixBase(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -check-breaking -base ../data/simple.yaml -revision ../data/simple.yaml -strip-prefix-base /partner-api"), io.Discard, io.Discard))
}

func Test_DuplicatePathsFail(t *testing.T) {
	require.NotZero(t, internal.Run(cmdToArgs("oasdiff -base ../data/duplicate_endpoints/base.yaml -revision ../data/duplicate_endpoints/revision.yaml -check-breaking"), io.Discard, io.Discard))
}

func Test_DuplicatePathsOK(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/duplicate_endpoints/base.yaml -revision ../data/duplicate_endpoints/revision.yaml -check-breaking -match-path-params"), io.Discard, io.Discard))
}

func Test_Changelog(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/run_test/changelog_base.yaml -revision ../data/run_test/changelog_revision.yaml -changelog -format json"), &stdout, io.Discard))
	cl := checker.BackwardCompatibilityErrors{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &cl))
	require.Len(t, cl, 1)
}

func Test_BreakingChangesChangelogOptionalCheckersAreInfoLevel(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/run_test/changelog_include_checks_base.yaml -revision ../data/run_test/changelog_include_checks_revision.yaml -changelog -format json"), &stdout, io.Discard))
	cl := checker.BackwardCompatibilityErrors{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &cl))
	require.Len(t, cl, 2)
	for _, c := range cl {
		require.Equal(t, c.Level, checker.INFO)
	}
}
