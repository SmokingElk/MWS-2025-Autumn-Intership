package errors

import "errors"

var (
	ErrRepoNotFound = errors.New("repo not found")
	ErrNotGoRepo    = errors.New("not go repo")
	ErrBadGomod     = errors.New("bad gomod")
	ErrBadUrl       = errors.New("bad url")
)
