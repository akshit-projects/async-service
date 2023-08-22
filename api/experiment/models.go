package experiment_apis

type ExperimentData struct {
	Name string `json:"name"`
	Id string `json:"id"`
	Variants []string `json:"variants"`
	Type string `json:"type"`
	LayerId string `json:"layerId"`
}

type GetExperimentsResponse struct {
	Experiments []ExperimentData `json:"experiments"`
}


type OverrideVariantRequest struct {
	LayerId string
	ExperimentId string
	Variant string
	UserId string
}

// Response from e13n service
type ExperimentServiceResponse struct {
	Tag string `json:"tag"`
	Name string `json:"name"`
	Id string
	variants []string
}