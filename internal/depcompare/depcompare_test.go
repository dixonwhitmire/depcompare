package depcompare

import "testing"

func TestLoadErrors(t *testing.T) {
	tests := map[string]struct {
		depType string
		depPath string
	}{
		"invalid-depType": {depType: "invalid", depPath: "testdata/gradle-dep-a.txt"},
		"invalid-depPath": {depType: GradleBuildDependency, depPath: "testdata/invalid-path"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Load(tc.depType, tc.depPath)
			if err == nil {
				t.Error("expected error did not occur")
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := map[string]struct {
		depType string
		depPath string
		want    map[string]string
	}{
		"gradle-dep-a.txt": {
			depType: GradleTextDependency,
			depPath: "testdata/gradle-dep-a.txt",
			want: map[string]string{
				"org.apache.commons:commons-lang3":                  "3.12.0",
				"org.apache.commons:commons-collections4":           "4.4",
				"org.springframework.boot:spring-boot-starter-web":  "",
				"org.postgresql:postgresql":                         "",
				"org.springframework.boot:spring-boot-starter-test": "",
			},
		},
		"gradle-dep-b.txt": {
			depType: GradleTextDependency,
			depPath: "testdata/gradle-dep-b.txt",
			want: map[string]string{
				"org.apache.commons:commons-lang3":        "3.12.0",
				"org.apache.commons:commons-collections4": "4.4",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Load(tc.depType, tc.depPath)
			if err != nil {
				t.Errorf("unexpected error occurred: %v", err)
			}
			// compare dependency maps
			for k := range tc.want {
				if v, ok := got[k]; !ok {
					t.Errorf("missing key %s", k)
				} else if v != tc.want[k] {
					t.Errorf("got %s, want %s", v, tc.want[k])
				}
			}
		})
	}
}
