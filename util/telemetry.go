package toputils

import (
	"net"
	"strconv"

	"github.com/pkg/errors"

	gometrics "github.com/circonus-labs/circonus-gometrics"
)

func (engine *Engine) SetupTelemetry(host string, port int, instanceId, displayName string) error {
	cfg := &gometrics.Config{}
	cfg.Interval = "10s"

	circonusPort := strconv.Itoa(port)
	hostPort := net.JoinHostPort(host, circonusPort)

	circonusUrl := "http://" + hostPort + "/write/" + instanceId

	cfg.CheckManager.Check.SubmissionURL = circonusUrl
	cfg.CheckManager.Check.InstanceID = instanceId
	cfg.CheckManager.Check.DisplayName = displayName

	metrics, err := gometrics.NewCirconusMetrics(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to bootstrap circonus")
	}

	engine.CirconusMetrics = metrics

	return nil
}

func (engine *Engine) PostMetricsGaugeValue(gaugeKey string, gaugeValue int64) {
	if engine.CirconusMetrics == nil {
		return
	}

	engine.CirconusMetrics.SetGaugeFunc(gaugeKey, func() int64 {
		return gaugeValue
	})
}

func (engine *Engine) PostMetricsHistogramValue(histogramKey string, histogramValue float64) {
	if engine.CirconusMetrics == nil {
		return
	}

	engine.CirconusMetrics.SetHistogramValue(histogramKey, histogramValue)
}

func (engine *Engine) PostMetricsTextValue(key string, value string) {
	if engine.CirconusMetrics == nil {
		return
	}

	engine.CirconusMetrics.SetTextFunc(key, func() string {
		return value
	})
}
