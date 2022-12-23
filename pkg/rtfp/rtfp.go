package rtfp

import (
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-json/sanitize"
)

var (
	TerraformResourceDiffers = []TerraformResourceDiff{newDiffAwsCloudfrontDistribution()}
)

type TerraformResourceDiff interface {
	getTerraformResourceType() string
	diffTerraformResource(rc *tfjson.ResourceChange) (string, error)
}

func PrettyTerraformPlan(planJson []byte) (string, error) {
	plan := &tfjson.Plan{}
	err := plan.UnmarshalJSON(planJson)
	if err != nil {
		return "", err
	}

	// sanitise!
	plan, err = sanitize.SanitizePlan(plan)
	if err != nil {
		return "", err
	}

	resourceChanges, err := filterKnownResources(plan, EnumerateKnownResourceTypes())
	if err != nil {
		return "", err
	}

	var str string
	for _, resourceChange := range resourceChanges {
		if !resourceChange.Change.Actions.Update() {
			continue
		}

		differ, err := getResourceDifferForType(resourceChange.Type)
		if err != nil {
			return "", err
		}

		s, err := differ.diffTerraformResource(resourceChange)
		if err != nil {
			return "", err
		}
		if s == "" {
			s = "No changes detected. This is typically because only sensitive values have changed."
		}

		str += fmt.Sprintf(`
Plan for resource: %s
%s
		`, resourceChange.Address, s)
	}

	return str, nil
}

func EnumerateKnownResourceTypes() (types []string) {
	for _, d := range TerraformResourceDiffers {
		types = append(types, d.getTerraformResourceType())
	}
	return types
}

func filterKnownResources(p *tfjson.Plan, types []string) (rc []*tfjson.ResourceChange, err error) {
	for _, change := range p.ResourceChanges {
		if contains(types, change.Type) {
			rc = append(rc, change)
		}
	}

	if len(rc) == 0 {
		return nil, errors.New("No resources found to diff.")
	}
	return rc, nil
}

func getResourceDifferForType(t string) (TerraformResourceDiff, error) {
	for _, d := range TerraformResourceDiffers {
		if d.getTerraformResourceType() == t {
			return d, nil
		}
	}
	return nil, errors.New("No differ found for resource type: " + t)
}
