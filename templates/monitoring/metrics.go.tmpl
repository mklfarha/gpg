package monitoring

import (
	"context"
	"fmt"

	"github.com/DataDog/datadog-go/v5/statsd"
	"go.uber.org/fx"
)

func (i *Implementation) emitMetric(req EmitRequest) error {
	if i.config.MetricsEnabled {
		switch i.config.MetricsProvider {
		case DATADOG_METRICS_PROVIDER:
			return i.emitDatadogMetric(req)
		}
	}
	return nil
}

func initMetricsClient(params Params, config MonitoringConfig, imp *Implementation) error {
	if config.MetricsEnabled && params.Lifecycle != nil {
		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				switch config.MetricsProvider {
				case DATADOG_METRICS_PROVIDER:
					addr := params.Provider.Get("monitoring.datadog-addr").String()
					dgc, err := statsd.New(addr, statsd.WithErrorHandler(func(err error) {
						fmt.Printf("datadog client error: %v \n", err)
					}))
					if err != nil {
						return err
					}
					imp.datadogClient = dgc
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				switch config.MetricsProvider {
				case DATADOG_METRICS_PROVIDER:
					imp.datadogClient.Flush()
					imp.datadogClient.Close()
				}
				return nil
			},
		})
	}

	return nil
}
