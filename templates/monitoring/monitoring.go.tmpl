package monitoring

import (
	"fmt"
)

func New(params Params) (*Implementation, error) {
	config := MonitoringConfig{}
	err := params.Provider.Get("monitoring").Populate(&config)
	if err != nil {
		return nil, err
	}

	imp := &Implementation{
		logger:   params.Logger,
		provider: params.Provider,
		config:   config,
	}

	err = initMetricsClient(params, config, imp)
	if err != nil {
		return nil, err
	}

	return imp, nil
}

func (i *Implementation) Emit(req EmitRequest) {
	// emit log
	i.emitLog(req)

	// emit metric and return
	err := i.emitMetric(req)
	if err != nil {
		fmt.Printf("ERROR emiting metrics: %v", err)
	}
}
