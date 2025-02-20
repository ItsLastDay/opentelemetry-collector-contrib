// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chronyreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver"

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver/internal/chrony"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver/internal/metadata"
)

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	metadata.MetricsSettings                `mapstructure:"metrics"`
	// Endpoint is the published address or unix socket
	// that allows clients to connect to:
	// The allowed format is:
	//   unix:///path/to/chronyd/unix.sock
	//   udp://localhost:323
	//
	// The default value is unix:///var/run/chrony/chronyd.sock
	Endpoint string `mapstructure:"endpoint"`
	// Timeout controls the max time allowed to read data from chronyd
	Timeout time.Duration `mapstructure:"timeout"`
}

var (
	_ config.Receiver = (*Config)(nil)

	errInvalidValue = errors.New("invalid value")
)

func newDefaultCongfig() config.Receiver {
	return &Config{
		ScraperControllerSettings: scraperhelper.NewDefaultScraperControllerSettings(typeStr),
		MetricsSettings:           metadata.DefaultMetricsSettings(),

		Endpoint: "unix:///var/run/chrony/chronyd.sock",
		Timeout:  10 * time.Second,
	}
}

func (c *Config) Validate() error {
	if c.Timeout < 1 {
		return fmt.Errorf("must have a positive timeout: %w", errInvalidValue)
	}
	_, _, err := chrony.SplitNetworkEndpoint(c.Endpoint)
	return err
}
