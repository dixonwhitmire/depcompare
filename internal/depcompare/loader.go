package depcompare

import "fmt"

// DependencyLoader loads dependency information from file.
type DependencyLoader interface {
	// Load loads dependency information from file returning a map of dependency names and versions.
	Load(filePath string) (map[string]string, error)
}

// LoadDependencies uses the DependencyLoader to load dependencies from file.
func LoadDependencies(loader DependencyLoader, dependencyPath string) (map[string]string, error) {
	deps, err := loader.Load(dependencyPath)

	if err != nil {
		return nil, fmt.Errorf("depcompare.LoadDependencies: error loading dependency file: %w", err)
	}
	return deps, nil
}

func NewLoader(loaderType string) {

}
