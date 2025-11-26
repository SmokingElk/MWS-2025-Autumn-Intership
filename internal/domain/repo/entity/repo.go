package entity

import (
	"strings"

	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"
)

type Repo struct {
	Name         string
	Version      versionEntity.Version
	Dependencies map[string]moduleEntity.Module
}

func NewRepo(name, versionStr string) (Repo, error) {
	if parts := strings.Split(versionStr, "."); len(parts) < 3 {
		versionStr += ".0"
	}

	version, err := versionEntity.NewVersion(versionStr)

	if err != nil {
		return Repo{}, err
	}

	return Repo{
		Name:         name,
		Version:      version,
		Dependencies: make(map[string]moduleEntity.Module),
	}, nil
}

func (r *Repo) AddDependency(dependency moduleEntity.Module) {
	r.Dependencies[dependency.Name] = dependency
}
