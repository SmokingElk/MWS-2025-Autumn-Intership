package entity

import (
	"strconv"
	"strings"

	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/errors"
)

type Version struct {
	major int
	minor int
	patch int
}

func NewVersion(versionStr string) (Version, error) {
	cleanVersion := strings.TrimPrefix(strings.Split(versionStr, "-")[0], "v")

	parts := strings.Split(cleanVersion, ".")

	if len(parts) != 3 {
		return Version{}, errors.ErrBadVersionStr
	}

	partsInt := make([]int, 0, len(parts))

	for _, part := range parts {
		partInt, err := strconv.Atoi(part)

		if err != nil {
			return Version{}, errors.ErrBadVersionStr
		}

		partsInt = append(partsInt, partInt)
	}

	return Version{
		major: partsInt[0],
		minor: partsInt[0],
		patch: partsInt[0],
	}, nil
}

func (v Version) Less(other Version) bool {
	if v.major < other.major {
		return true
	}
	if v.major > other.major {
		return false
	}

	if v.minor < other.minor {
		return true
	}
	if v.minor > other.minor {
		return false
	}

	return v.patch < other.patch
}
