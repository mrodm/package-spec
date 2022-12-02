package spec

import (
	"io/fs"
	"log"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	spec "github.com/elastic/package-spec/v2"
	"github.com/elastic/package-spec/v2/code/go/internal/jsonschema"
	"github.com/elastic/package-spec/v2/code/go/internal/loader"
)

// Spec represents a package specification
type Spec struct {
	version semver.Version
	fs      fs.FS
}

// NewSpec creates a new Spec for the given version
func NewSpec(version semver.Version) (*Spec, error) {
	specVersion, err := spec.CheckVersion(version)
	if err != nil {
		return nil, errors.Wrapf(err, "could not load specification for version [%s]", version.String())
	}
	if specVersion.Prerelease() != "" {
		log.Printf("Warning: package using an unreleased version of the spec (%s)", specVersion)
	}

	s := Spec{
		version,
		spec.FS(),
	}

	return &s, nil
}

func (s Spec) RenderJsonSchema(itemPath, pkgType string) (*jsonschema.RenderedJSONSchema, error) {
	rootSpec, err := loader.LoadSpec(s.fs, s.version, pkgType)
	if err != nil {
		return nil, err
	}

	return jsonschema.JSONSchema(rootSpec, itemPath)
}

func (s Spec) RenderAllJsonSchema(pkgType string) ([]jsonschema.RenderedJSONSchema, error) {
	rootSpec, err := loader.LoadSpec(s.fs, s.version, pkgType)
	if err != nil {
		return nil, err
	}

	return jsonschema.AllJSONSchemas(rootSpec)
}
