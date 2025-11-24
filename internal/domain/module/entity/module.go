package entity

import versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"

type Module struct {
	Name    string
	Version versionEntity.Version
	Direct  bool
}

func NewModule(name, versionStr string, direct bool) (Module, error) {
	version, err := versionEntity.NewVersion(versionStr)

	if err != nil {
		return Module{}, err
	}

	return Module{
		Name:    name,
		Version: version,
		Direct:  direct,
	}, nil
}
