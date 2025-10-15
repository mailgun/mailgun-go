package mtypes

type ListIPWarmupsResponse struct {
	Paging Paging     `json:"paging"`
	Items  []IPWarmup `json:"items"`
}

// IPWarmup - is a response for ListIPWarmups.
type IPWarmup struct {
	IP               string      `json:"ip"`
	SentWithinStage  string      `json:"sent_within_stage"`
	Throttle         int         `json:"throttle"`
	StageNumber      int         `json:"stage_number"`
	StageStartVolume int         `json:"stage_start_volume"`
	StageStartTime   RFC3339Time `json:"stage_start_time"`
	StageVolumeLimit int         `json:"stage_volume_limit"`
}

type IPWarmupDetailsResponse struct {
	Details IPWarmupDetails `json:"details"`
}

type IPWarmupDetails struct {
	SentWithinStage   string                 `json:"sent_within_stage"`
	Throttle          int                    `json:"throttle"`
	StageNumber       int                    `json:"stage_number"`
	StageStartVolume  int                    `json:"stage_start_volume"`
	StageStartTime    RFC3339Time            `json:"stage_start_time"`
	StageVolumeLimit  int                    `json:"stage_volume_limit"`
	StageStartedAt    RFC3339Time            `json:"stage_started_at"`
	HourStartedAt     RFC3339Time            `json:"hour_started_at"`
	PlanStartedAt     RFC3339Time            `json:"plan_started_at"`
	PlanLastUpdatedAt RFC3339Time            `json:"plan_last_updated_at"`
	TotalStages       int                    `json:"total_stages"`
	StageHistory      []IPWarmupStageHistory `json:"stage_history"`
}

type IPWarmupStageHistory struct {
	FirstUpdatedAt RFC3339Time `json:"first_updated_at"`
	CompletedAt    RFC3339Time `json:"completed_at"`
	Limit          int         `json:"limit"`
}
