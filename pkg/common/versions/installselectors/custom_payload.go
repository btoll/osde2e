package installselectors

import (
	"math"

	"github.com/Masterminds/semver/v3"
	viper "github.com/openshift/osde2e/pkg/common/concurrentviper"
	"github.com/openshift/osde2e/pkg/common/config"
	"github.com/openshift/osde2e/pkg/common/spi"
)

func init() {
	registerSelector(customPayload{})
}

type customPayload struct{}

func (c customPayload) ShouldUse() bool {
	return viper.GetString(config.Cluster.CustomPayload) != ""
}

func (c customPayload) Priority() int {
	return math.MaxInt32
}

func (c customPayload) SelectVersion(versionList *spi.VersionList) (*semver.Version, string, error) {
	customPayload := viper.GetString(config.Cluster.CustomPayload)
	//_ = versionList.AvailableVersions()
	semVersion, _ := semver.NewVersion(customPayload)
	return semVersion, c.String(), nil
}

func (c customPayload) String() string {
	return "custom payload"
}
