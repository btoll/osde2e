package installselectors

import (
	"fmt"
	"log"
	"strings"

	"github.com/Masterminds/semver/v3"
	viper "github.com/openshift/osde2e/pkg/common/concurrentviper"
	"github.com/openshift/osde2e/pkg/common/config"
	"github.com/openshift/osde2e/pkg/common/spi"
	"github.com/openshift/osde2e/pkg/common/versions/common"
)

func init() {
	registerSelector(specificNightlies{})
}

// SpecificNightlies attempts to parse a config option as semver and use the major.minor to look for nightlies
type specificNightlies struct{}

func (m specificNightlies) ShouldUse() bool {
	log.Printf("specific nightly value: %q", viper.GetString(config.Cluster.InstallSpecificNightly))
	return viper.GetString(config.Cluster.InstallSpecificNightly) != ""
}

func (m specificNightlies) Priority() int {
	return 60
}

func (m specificNightlies) SelectVersion(versionList *spi.VersionList) (*semver.Version, string, error) {
	specificNightly := viper.GetString(config.Cluster.InstallSpecificNightly)
	if specificNightly == "" {
		return nil, m.String(), fmt.Errorf("no version to match nightly found")
	}

	versionsWithoutDefault := removeDefaultVersion(versionList.AvailableVersions())
	common.SortVersions(versionsWithoutDefault)

	versionToMatch := semver.MustParse(specificNightly)

	if versionToMatch == nil {
		return nil, m.String(), fmt.Errorf("error parsing semver version for %s", specificNightly)
	}

	for i := len(versionsWithoutDefault) - 1; i > -1; i-- {
		if strings.Contains(versionsWithoutDefault[i].Version().Original(), "nightly") && versionsWithoutDefault[i].Version().Major() == versionToMatch.Major() && versionsWithoutDefault[i].Version().Minor() == versionToMatch.Minor() {
			// Since we're going through a list in reverse-order, the first X.Y that matches should be the latest!
			return versionsWithoutDefault[i].Version(), m.String(), nil
		}
	}

	return nil, m.String(), fmt.Errorf("no valid nightly found for version %s", specificNightly)
}

func (m specificNightlies) String() string {
	return "specific nightly"
}
