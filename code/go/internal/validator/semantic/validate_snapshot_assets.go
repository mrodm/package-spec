// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package semantic

import (
	"fmt"
	"log"
	"path"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	ve "github.com/elastic/package-spec/code/go/internal/errors"
	"github.com/elastic/package-spec/code/go/internal/fspath"
	"github.com/elastic/package-spec/code/go/internal/pkgpath"
	"github.com/elastic/package-spec/code/go/internal/validator/common"
)

const elasticPrereleaseTag = "SNAPSHOT"

var assetsToCheck = []string{
	"dashboard",
	"visualization",
	"lens",
	"map",
}

// ValidateSnapshotVersionsInAssets validates additional restrictions on the Elastic stack versions used to generate assets.
func ValidateSnapshotVersionsInAssets(fsys fspath.FS) ve.ValidationErrors {
	warningsAsErrors := common.IsDefinedWarningsAsErrors()
	var errs ve.ValidationErrors

	manifestVersion, err := readManifestVersion(fsys)
	if err != nil {
		return ve.ValidationErrors{err}
	}

	allowSnapshot, err := readAllowSnapshotManifest(fsys)
	if err != nil {
		return ve.ValidationErrors{err}
	}

	// technical previous and prerelease package versions are allowed to contain
	// assets built from Elastic stack SNAPSHOT versions
	packageVersion, err := semver.NewVersion(manifestVersion)
	if packageVersion.Major() == 0 || packageVersion.Prerelease() != "" {
		// no retrictions, it can contain -SNAPSHOT
		return nil
	}

	// stable versions allowed to contain assets built from Elastic stack SNAPSHOT versions
	// if allow_snapshot is set
	if allowSnapshot {
		return nil
	}

	// stable package versions are not allowed to contain assets built from Elastic stack SNAPSHOT versions
	for _, asset := range assetsToCheck {
		filePaths := path.Join("kibana", asset, "*.json")
		objectFiles, err := pkgpath.Files(fsys, filePaths)
		if err != nil {
			continue
		}

		for _, objectFile := range objectFiles {
			filePath := objectFile.Path()

			assetVersion, err := readMigrationVersionField(objectFile)
			if err != nil {
				errs = append(errs, errors.Wrap(err, "can't get elastic stack version of asset"))
				continue
			}

			snapshot, err := usingSnapshotVersion(assetVersion)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			if snapshot {
				message := fmt.Sprintf("Warning: snapshot version found in %s %s: %s", asset, filePath, assetVersion)
				if warningsAsErrors {
					errs = append(errs, errors.New(message))
				} else {
					log.Printf(message)
				}
			}
		}
	}

	return errs
}

// readMigrationVersionField return the version in migrationVersion from an asset file
func readMigrationVersionField(objectFile pkgpath.File) (string, error) {
	// there are some assets that the field under migrationVersion do not match with the asset type field
	// "migrationVersion": {
	//     "search": "7.9.3"
	// },
	// "references": [],
	// "type": "ml-module"

	versionReference, err := objectFile.Values(`$.migrationVersion.*`)
	if err != nil {
		return "", err
	}
	versions, err := toStringSlice(versionReference)
	if err != nil {
		return "", errors.Errorf("conversion error to string %s %s", versionReference, objectFile.Path())
	}
	if len(versions) > 1 {
		return "", errors.Errorf("unexpected number of versions in migrationVersion field: %s", versions)
	}
	if len(versions) == 0 {
		// some assets do not have migrationVersion field, no error raised
		return "", nil
	}
	return versions[0], nil
}

// usingSnapshotVersionVersion returns a boolean indicating if version is from a Snapshot version
func usingSnapshotVersion(version string) (bool, error) {
	if version == "" {
		// some assets do not have migrationVersion field, allowed
		return false, nil
	}

	semVersion, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}
	return semVersion.Prerelease() == elasticPrereleaseTag, nil
}

func readAllowSnapshotManifest(fsys fspath.FS) (bool, error) {
	manifestPath := "manifest.yml"
	f, err := pkgpath.Files(fsys, manifestPath)
	if err != nil {
		return false, errors.Wrap(err, "can't locate manifest file")
	}

	if len(f) != 1 {
		return false, errors.New("single manifest file expected")
	}

	val, err := f[0].Values("$.allow_snapshot")
	if err != nil {
		return false, nil
	}

	bVal, ok := val.(bool)
	if !ok {
		return false, errors.New("allow_snapshot unexpected value")
	}
	return bVal, nil
}
