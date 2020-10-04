package pkg

import (
    "testing"
)

// TestCreateDeployment tests that the SeldonDeployment struct is properly formed.
func TestCreateDeployment(t *testing.T) {
    sd, _ := CreateDeployment("../seldon-crd.json")

    tests := []struct {
        name     string
        expected string
        actual   string
    }{
        {"Name", "seldon-deployment-example", sd.Name},
        {"Namespace", "", sd.Namespace},
    }

    for _, tt := range tests {
        if tt.actual != tt.expected {
            t.Errorf("Fail, expected resource name of %v, got %v", tt.expected, tt.actual)
        }
    }

}

// TestCreateDeploymentNotJSON tests that non json manifests are rejected.
func TestCreateDeploymentNotJSON(t *testing.T) {
    _, err := CreateDeployment("../seldon-crd.yaml")
    if err == nil {
        t.Errorf("Fail, expected err, got nil")
    }

}
