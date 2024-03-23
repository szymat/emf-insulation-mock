package insulationmock_test

import (
	"testing"

	"github.com/szymat/emf-insulation-mock/pkg/insulationmock"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	c := insulationmock.NewConfig("scriptsPath", "routesPath", "beforeScript", "afterScript", "port")

	if c.ScriptsPath != "scriptsPath" {
		t.Errorf("NewConfig() = %v, want %v", c.ScriptsPath, "scriptsPath")
	}
	if c.RoutesPath != "routesPath" {
		t.Errorf("NewConfig() = %v, want %v", c.RoutesPath, "routesPath")
	}

	if c.BeforeScript != "beforeScript" {
		t.Errorf("NewConfig() = %v, want %v", c.BeforeScript, "beforeScript")
	}

	if c.AfterScript != "afterScript" {
		t.Errorf("NewConfig() = %v, want %v", c.AfterScript, "afterScript")
	}

	if c.Port != "port" {
		t.Errorf("NewConfig() = %v, want %v", c.Port, "port")
	}
}
