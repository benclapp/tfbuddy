package rtfp

import (
	"errors"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestEnumerateKnownResourceTypes(t *testing.T) {
	assert.Equal(t,
		[]string{"aws_cloudfront_distribution"},
		EnumerateKnownResourceTypes(),
	)
}

func TestFilterKnownResources(t *testing.T) {
	p := &tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{Type: "known"},
			{Type: "unknown"},
		},
	}
	expected := []*tfjson.ResourceChange{{Type: "known"}}
	actual, err := filterKnownResources(p, []string{"known"})

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func TestFilterKnownResourcesError(t *testing.T) {
	p := &tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{Type: "a"},
			{Type: "b"},
		},
	}
	_, err := filterKnownResources(p, []string{"c"})

	assert.Equal(t, errors.New("No resources found to diff."), err)
}

func TestGetResourceDifferForType(t *testing.T) {
	expected := newDiffAwsCloudfrontDistribution()
	actual, err := getResourceDifferForType("aws_cloudfront_distribution")

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func TestGetResourceDifferForTypeError(t *testing.T) {
	a, err := getResourceDifferForType("non_existent_type")

	assert.Equal(t, nil, a)
	assert.Equal(t, errors.New("No differ found for resource type: non_existent_type"), err)
}
