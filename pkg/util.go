package pkg

import (
	"errors"
	"fmt"
	"github.com/anchore/grype/grype/version"
	syftPkg "github.com/anchore/syft/syft/pkg"
	"log"
)

func MatchVuln() error {
	searchVersion := "6.2.4"
	pkgType := syftPkg.BinaryPkg

	verObj, err := version.NewVersion(searchVersion, version.FormatFromPkgType(pkgType))
	if err != nil {
		log.Printf("err:%v", verObj)
		return err
	}
	log.Printf("versionObj: %+v \n", verObj)

	versionConstraint := ">= 6.0.0, < 6.0.17 || >= 6.2.0, < 6.2.9 || >= 7.0.0, < 7.0.8"
	constraint, err := version.GetConstraint(versionConstraint, version.UnknownFormat)
	if err != nil {
		log.Printf("get constraint err:%v", err)
		return err
	}
	log.Printf("constraint:%+v\n", constraint)

	isPackageVulnerable, err := constraint.Satisfied(verObj)
	if err != nil {
		var e *version.NonFatalConstraintError
		if errors.As(err, &e) {
			log.Printf("satisfied err:%v\n", e)
		} else {
			return fmt.Errorf("failed to check constraint=%q version=%q: %w", constraint, verObj, err)
		}
	}
	log.Printf("isPackageVulnerable:%v \n", isPackageVulnerable)
	return nil
}
