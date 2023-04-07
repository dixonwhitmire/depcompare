package depcompare

import (
	"testing"
)

var gradleBaseDeps = map[string]string{
	"org.apache.commons:commons-lang3":                  "3.12.0",
	"org.apache.commons:commons-collections4":           "4.4",
	"org.springframework.boot:spring-boot-starter-web":  "",
	"org.postgresql:postgresql":                         "",
	"org.springframework.boot:spring-boot-starter-test": "",
}

var gradleDeps = map[string]string{
	"org.apache.commons:commons-lang3":        "3.12.0",
	"org.apache.commons:commons-collections4": "4.4",
	"org.apache.commons:commons-csv":          "1.10.0",
}

func TestLoadErrors(t *testing.T) {
	tests := map[string]struct {
		depType string
		depPath string
	}{
		"invalid-depType": {depType: "invalid", depPath: "testdata/gradle-base.txt"},
		"invalid-depPath": {depType: gradleBuildDependency, depPath: "testdata/invalid-path"},
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
		"gradle-dep.txt": {
			depType: gradleTextDependency,
			depPath: "testdata/gradle-dep.txt",
			want:    gradleDeps,
		},
		"gradle-base.txt": {
			depType: gradleTextDependency,
			depPath: "testdata/gradle-base.txt",
			want:    gradleBaseDeps,
		},
		"build-dep.gradle": {
			depType: gradleBuildDependency,
			depPath: "testdata/build-dep.gradle",
			want:    gradleDeps,
		},
		"build-base.gradle": {
			depType: gradleBuildDependency,
			depPath: "testdata/build-base.gradle",
			want:    gradleBaseDeps,
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

func TestCompare(t *testing.T) {
	want := map[string][]string{
		IntersectKey: {
			"org.apache.commons:commons-collections4",
			"org.apache.commons:commons-lang3",
		},
		BaseOnlyKey: {
			"org.postgresql:postgresql",
			"org.springframework.boot:spring-boot-starter-test",
			"org.springframework.boot:spring-boot-starter-web",
		},
		DepOnlyKey: {
			"org.apache.commons:commons-csv",
		},
	}

	result := Compare(gradleDeps, gradleBaseDeps)

	for _, k := range []string{IntersectKey, BaseOnlyKey, DepOnlyKey} {
		got, ok := result[k]
		if !ok {
			t.Fatalf("missing %s result", k)
		}
		if len(got) != len(want[k]) {
			t.Fatalf("testing %s got length %d, want length %d", k, len(got), len(want[k]))
		}
		for i := 0; i < len(got); i++ {
			if got[i] != want[k][i] {
				t.Errorf("got = %s, want = %s", got[i], want[k][i])
			}
		}
	}
}
