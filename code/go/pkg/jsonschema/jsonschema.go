// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package jsonschema

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	"github.com/elastic/package-spec/v2/code/go/internal/validator"
)

func AllJSONSchemas(version, pkgType string) error {
	fmt.Printf("All json schemas (%s - %s)", pkgType, version)
	specVersion, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	spec, err := validator.NewSpec(*specVersion)
	if err != nil {
		return err
	}

	rendered, err := spec.AllJSONSchema(pkgType)

	for _, itemSpec := range rendered {
		fmt.Printf("Name: %s\n", itemSpec.Name)
		fmt.Printf("Content:\n%s\n", itemSpec.JSONSchema)
	}
	return nil
}

func JSONSchema(itemPath, version, pkgType string) ([]byte, error) {
	fmt.Printf("jsonschema for %s (%s - %s)", itemPath, pkgType, version)
	specVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	spec, err := validator.NewSpec(*specVersion)
	if err != nil {
		return nil, errors.Wrap(err, "invalid package spec version")
	}

	fmt.Printf("spec %+v", spec)
	rendered, err := spec.JSONSchema(itemPath, pkgType)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to render jsonschema for %s", itemPath)
	}

	fmt.Printf("Name: %s\n", itemPath)
	fmt.Printf("Content:\n%s\n", rendered.JSONSchema)
	return rendered.JSONSchema, nil
}
