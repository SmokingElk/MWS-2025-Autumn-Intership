package entity

import versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"

type Module struct {
	Name    string
	Version versionEntity.Version
	Direct  bool
}

func NewModule(name string, version versionEntity.Version, direct bool) Module {
	return Module{
		Name:    name,
		Version: version,
		Direct:  direct,
	}
}
