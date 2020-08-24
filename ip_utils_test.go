package main

import (
	"testing"
)

type testSpec struct {
	input    string
	expected string
}

func TestApplyCIDRMask(t *testing.T) {
	specs := []testSpec{
		{"10.244.1.34/24", "10.244.1.0"},
		{"10.244.23.3/16", "10.244.0.0"},
	}

	for _, spec := range specs {
		output := applyCIDRMask(spec.input)
		if output != spec.expected {
			t.Errorf("Expected %s -> Got %s", spec.expected, output)
		}
	}
}

func TestReduceCIDRSpecificity(t *testing.T) {
	cidrInput := "10.244.2.0/24"
	expectedOutput := "10.244.2.0/16"
	output := reduceCIDRSpecificity(cidrInput)

	if expectedOutput != output {
		t.Errorf("Expected %s -> Got %s", expectedOutput, output)
	}
}

func TestGetCIDRMaskSize(t *testing.T) {
	specs := []testSpec{
		{"10.244.1.34/24", "24"},
		{"10.244.23.3/16", "16"},
	}

	for _, spec := range specs {
		output := getCIDRMaskSize(spec.input)
		if output != spec.expected {
			t.Errorf("Expected %s -> Got %s", spec.expected, output)
		}
	}
}
