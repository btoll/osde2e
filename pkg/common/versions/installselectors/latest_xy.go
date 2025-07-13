package installselectors

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	viper "github.com/openshift/osde2e/pkg/common/concurrentviper"
	"github.com/openshift/osde2e/pkg/common/config"
	"github.com/openshift/osde2e/pkg/common/spi"
)

func init() {
	registerSelector(latestXYVersion{})
}

type latestXYVersion struct{}

func (m latestXYVersion) ShouldUse() bool {
	return viper.GetString(config.Cluster.InstallLatestXY) != ""
}

func (m latestXYVersion) Priority() int {
	return 60
}

func (m latestXYVersion) SelectVersion(versionList *spi.VersionList) (*semver.Version, string, error) {
	latestXY := viper.GetString(config.Cluster.InstallLatestXY)
	versions := versionList.AvailableVersions()

	semVersion, err := semver.NewVersion(latestXY)
	if err != nil {
		return nil, m.String(), fmt.Errorf("error parsing semantic version for %s", latestXY)
	}

	for _, version := range versions {
		if version.Version().Major() == semVersion.Major() && version.Version().Minor() == semVersion.Minor() {
			return version.Version(), m.String(), nil
		}
	}

	return nil, m.String(), fmt.Errorf("unable to locate latest version for %s", latestXY)
}

func (m latestXYVersion) String() string {
	return "latest X.Y version available"
}
