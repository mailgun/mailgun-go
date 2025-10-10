package mtypes

type ListIPWarmupsResponse struct {
	Paging Paging     `json:"paging"`
	Items  []IPWarmup `json:"items"`
}

// IPWarmup - is a response for ListIPWarmups.
type IPWarmup struct {
	IP               string `json:"ip"`
	SentWithinStage  string `json:"sent_within_stage"`
	Throttle         int    `json:"throttle"`
	StageNumber      int    `json:"stage_number"`
	StageStartVolume int    `json:"stage_start_volume"`
	StageStartTime   string `json:"stage_start_time"`
	StageVolumeLimit int    `json:"stage_volume_limit"`
}
