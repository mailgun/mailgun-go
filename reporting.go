package mailgun

// common things for stats and metrics

type Resolution string

const (
	ResolutionHour  = Resolution("hour")
	ResolutionDay   = Resolution("day")
	ResolutionMonth = Resolution("month")
)
