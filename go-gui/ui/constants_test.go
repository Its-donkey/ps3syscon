package ui

import (
	"strings"
	"testing"
)

func TestAppVersion(t *testing.T) {
	if AppVersion == "" {
		t.Error("AppVersion should not be empty")
	}

	// Version should be in semantic versioning format (x.y.z)
	parts := strings.Split(AppVersion, ".")
	if len(parts) < 2 {
		t.Errorf("AppVersion %q should be in semantic versioning format", AppVersion)
	}
}

func TestAppAuthor(t *testing.T) {
	if AppAuthor == "" {
		t.Error("AppAuthor should not be empty")
	}
}

func TestAppEmail(t *testing.T) {
	if AppEmail == "" {
		t.Error("AppEmail should not be empty")
	}

	// Email should contain @
	if !strings.Contains(AppEmail, "@") {
		t.Errorf("AppEmail %q should contain @", AppEmail)
	}
}

func TestAppRepo(t *testing.T) {
	if AppRepo == "" {
		t.Error("AppRepo should not be empty")
	}

	// Repo URL should start with https://
	if !strings.HasPrefix(AppRepo, "https://") {
		t.Errorf("AppRepo %q should start with https://", AppRepo)
	}

	// Repo URL should contain github.com
	if !strings.Contains(AppRepo, "github.com") {
		t.Errorf("AppRepo %q should contain github.com", AppRepo)
	}
}

func TestConstantsAreNotEmpty(t *testing.T) {
	constants := map[string]string{
		"AppVersion": AppVersion,
		"AppAuthor":  AppAuthor,
		"AppEmail":   AppEmail,
		"AppRepo":    AppRepo,
	}

	for name, value := range constants {
		if value == "" {
			t.Errorf("%s should not be empty", name)
		}
	}
}
