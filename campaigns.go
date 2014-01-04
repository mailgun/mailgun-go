package mailgun

type Campaign struct {
	Id                string
	Name              string
	CreatedAt         string
	DeliveredCount    int
	ClickedCount      int
	OpenedCount       int
	SubmittedCount    int
	UnsubscribedCount int
	BouncedCount      int
	ComplainedCount   int
	DroppedCount      int
}
