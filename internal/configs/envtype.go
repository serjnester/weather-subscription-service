package configs

import (
	"errors"
	"strings"
)

type EnvType string

const (
	EnvTypeDev   EnvType = "dev"
	EnvTypeStage EnvType = "stage"
	EnvTypeProd  EnvType = "prod"
)

var knownEnvTypes = map[EnvType]bool{
	EnvTypeDev:   true,
	EnvTypeStage: true,
	EnvTypeProd:  true,
}

var ErrUnknownEnvType = errors.New("unknown env type")

func (e *EnvType) Decode(v string) error {
	if *e = EnvType(strings.ToLower(v)); !knownEnvTypes[*e] {
		return ErrUnknownEnvType
	}

	return nil
}

func (e *EnvType) IsProd() bool {
	return e != nil && *e == EnvTypeProd
}
