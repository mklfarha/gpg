package monitoring

import "fmt"

func (i *Implementation) emitDatadogMetric(req EmitRequest) error {
	name := fmt.Sprintf("%s.%s.%s.%s", req.Layer, req.EntityIdentifier, req.ActionIdentifier, req.Type)
	if req.LayerSubtype != "" {
		name = fmt.Sprintf("%s.%s.%s.%s.%s", req.Layer, req.LayerSubtype, req.EntityIdentifier, req.ActionIdentifier, req.Type)
	}
	if req.EntityIdentifier == "" {
		name = fmt.Sprintf("%s.%s.%s", req.Layer, req.ActionIdentifier, req.Type)
		if req.LayerSubtype != "" {
			name = fmt.Sprintf("%s.%s.%s.%s", req.Layer, req.LayerSubtype, req.ActionIdentifier, req.Type)
		}
	}
	fmt.Printf("emitting metric: %v \n", name)

	tags := []string{
		fmt.Sprintf("action:%s", req.ActionIdentifier),
		fmt.Sprintf("entity:%s", req.EntityIdentifier),
		fmt.Sprintf("layer:%s", string(req.Layer)),
		fmt.Sprintf("type:%s", string(req.Type)),
	}
	if req.LayerSubtype != "" {
		tags = append(tags, fmt.Sprintf("layer-subtype:%s", string(req.LayerSubtype)))
	}
	return i.datadogClient.Incr(name, tags, 1)
}
