package checker_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

func getReqPropFile(file string) string {
	return fmt.Sprintf("../data/required-properties/%s", file)
}

// BC: new required property in request header is breaking
func TestBreaking_NewRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: new optional property in request header is not breaking
func TestBreaking_NewNonRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request header to required is breaking
func TestBreaking_PropertyRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: changing an existing property in request header to optional is not breaking
func TestBreaking_PropertyRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in response body to optional is breaking
func TestBreaking_RespBodyRequiredPropertyDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: changing an existing property in response body to required is not breaking
func TestBreaking_RespBodyRequiredPropertyEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request body to optional is not breaking
func TestBreaking_ReqBodyRequiredPropertyDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request body to required is breaking
func TestBreaking_ReqBodyRequiredPropertyEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: adding a new required property in request body is breaking
func TestBreaking_ReqBodyNewRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: deleting a required property in request is not breaking
func TestBreaking_ReqBodyDeleteRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: adding a new required property in response body is not breaking
func TestBreaking_RespBodyNewRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required property in response body is breaking
func TestBreaking_RespBodyDeleteRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: adding a new required property under AllOf in response body is not breaking
func TestBreaking_RespBodyNewAllOfRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required property under AllOf in response body is breaking
func TestBreaking_RespBodyDeleteAllOfRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: adding a new required property under AllOf in response body is not breaking but when multiple inline (without $ref) schemas under AllOf are modified simultaneously, we detect is as breaking
// explanation: when multiple inline (without $ref) schemas under AllOf are modified we can't correlate schemas across base and revision
// as a result we can't determine that the change was "a new required property" and the change appears as breaking
func TestBreaking_RespBodyNewAllOfMultiRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-multi-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-multi-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: adding a new required read-only property in request body is not breaking
func TestBreaking_ReadOnlyNewRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-new-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-new-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing read-only property in request body to required is not breaking
func TestBreaking_ReadOnlyPropertyRequiredEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required write-only property in response body is not breaking
func TestBreaking_WriteOnlyDeleteRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a non-required non-write-only property in response body is not breaking
func TestBreaking_WriteOnlyDeleteNonRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-partial-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-partial-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing write-only property in response body to optional is not breaking
func TestBreaking_WriteOnlyPropertyRequiredDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing required property in response body to write-only is not breaking
func TestBreaking_RequiredPropertyWriteOnlyEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing required property in response body to not-write-only is breaking
func TestBreaking_RequiredPropertyWriteOnlyDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}