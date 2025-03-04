package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: new paths or path operations
func TestApiAdded_DetectsNewPathsAndNewOperations(t *testing.T) {
	s1, err := open("../data/new_endpoints/base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/new_endpoints/revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIAddedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)

	require.Equal(t, "endpoint-added", errs[0].Id)
	require.Equal(t, "POST", errs[0].Operation)
	require.Equal(t, "/api/test2", errs[0].Path)

	require.Equal(t, "endpoint-added", errs[1].Id)
	require.Equal(t, "GET", errs[1].Operation)
	require.Equal(t, "/api/test3", errs[1].Path)
}

// CL: new paths or path operations
func TestApiAdded_DetectsModifiedPathsWithPathParam(t *testing.T) {
	s1, err := open("../data/new_endpoints/base_with_path_param.yaml")
	require.NoError(t, err)

	s2, err := open("../data/new_endpoints/revision_with_path_param.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIAddedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)

	require.Equal(t, "endpoint-added", errs[0].Id)
	require.Equal(t, "POST", errs[0].Operation)
	require.Equal(t, "/api/test/{id}", errs[0].Path)
}
