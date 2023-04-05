package main

import "testing"

const (
	firstGradleTextPath   = "testdata/gradle-dep-a.txt"
	secondGradleTextPath  = "testdata/gradle-dep-b.txt"
	firstGradleBuildPath  = "testdata/build-a.gradle"
	secondGradleBuildPath = "testdata/build-a.gradle"
)

func TestValidateNoErrors(t *testing.T) {

	tests := map[string]struct {
		depType    string
		mode       string
		firstFile  string
		secondFile string
	}{
		"gradleb-intersect": {
			depType: "gradleb", mode: "intersect", firstFile: firstGradleBuildPath, secondFile: secondGradleBuildPath,
		},
		"gradleb-diff": {
			depType: "gradleb", mode: "diff", firstFile: firstGradleBuildPath, secondFile: secondGradleBuildPath,
		},
		"gradlet-intersect": {
			depType: "gradleb", mode: "intersect", firstFile: firstGradleTextPath, secondFile: secondGradleTextPath,
		},
		"gradlet-diff": {
			depType: "gradleb", mode: "diff", firstFile: firstGradleTextPath, secondFile: secondGradleTextPath,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			args := []string{tc.firstFile, tc.secondFile}
			err := validate(tc.depType, tc.mode, args)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateWithErrors(t *testing.T) {

	tests := map[string]struct {
		depType    string
		mode       string
		firstFile  string
		secondFile string
	}{
		"invalid-dependency-type": {
			depType: "invalid", mode: "diff", firstFile: firstGradleTextPath, secondFile: secondGradleTextPath,
		},
		"invalid-mode": {
			depType: "gradlea", mode: "invalid", firstFile: firstGradleTextPath, secondFile: secondGradleTextPath,
		},
		"first-path-invalid": {
			depType: "gradleb", mode: "diff", firstFile: "/tmp/invalid/file.txt", secondFile: secondGradleTextPath,
		},
		"second-path-invalid": {
			depType: "gradleb", mode: "diff", firstFile: firstGradleTextPath, secondFile: "/tmp/invalid/file.txt",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			args := []string{tc.firstFile, tc.secondFile}
			err := validate(tc.depType, tc.mode, args)
			if err == nil {
				t.Error("expecting error")
			}
		})
	}
}

func TestIsValidFlag(t *testing.T) {
	tests := map[string]struct {
		flagValue   string
		validValues []string
		want        bool
	}{
		"deptype-gradleb": {flagValue: "gradleb", validValues: validDepTypeFlags, want: true},
		"deptype-gradlet": {flagValue: "gradlet", validValues: validDepTypeFlags, want: true},
		"mode-intersect":  {flagValue: "intersect", validValues: validModeFlags, want: true},
		"mode-diff":       {flagValue: "diff", validValues: validModeFlags, want: true},
		"invalid-type":    {flagValue: "invalid", validValues: validDepTypeFlags, want: false},
		"invalid-diff":    {flagValue: "invalid", validValues: validModeFlags, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := isValidFlag(tc.flagValue, tc.validValues)
			if tc.want != got {
				t.Errorf("want = %v, got = %v", tc.want, got)
			}
		})
	}
}
