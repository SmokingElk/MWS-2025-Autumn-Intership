package entity

import (
	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"
)

type Repo struct {
	Name         string
	Version      versionEntity.Version
	Dependencies map[string]moduleEntity.Module
}

func NewRepo(name, versionStr string) (Repo, error) {
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
