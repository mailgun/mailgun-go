package mailgun

type MetricsResponse struct {
	Start      string            `json:"start"`
	End        string            `json:"end"`
	Resolution Resolution        `json:"resolution"`
	Duration   string            `json:"duration"`
	Dimensions []string          `json:"dimensions"`
	Pagination MetricsPagination `json:"pagination"`
	Items      []MetricsItem     `json:"items"`
	Aggregates MetricsAggregates `json:"aggregates"`
}

type MetricsItem struct {
	Dimensions []MetricsDimension `json:"dimensions"`
	Metrics    Metrics            `json:"metrics"`
}

type MetricsAggregates struct {
	Metrics Metrics `json:"metrics"`
}

type Metrics struct {
	AcceptedIncomingCount         *uint64 `json:"accepted_incoming_count"`
	AcceptedOutgoingCount         *uint64 `json:"accepted_outgoing_count"`
	AcceptedCount                 *uint64 `json:"accepted_count"`
	DeliveredSMTPCount            *uint64 `json:"delivered_smtp_count"`
	DeliveredHTTPCount            *uint64 `json:"delivered_http_count"`
	DeliveredOptimizedCount       *uint64 `json:"delivered_optimized_count"`
	DeliveredCount                *uint64 `json:"delivered_count"`
	StoredCount                   *uint64 `json:"stored_count"`
	ProcessedCount                *uint64 `json:"processed_count"`
	SentCount                     *uint64 `json:"sent_count"`
	OpenedCount                   *uint64 `json:"opened_count"`
	ClickedCount                  *uint64 `json:"clicked_count"`
	UniqueOpenedCount             *uint64 `json:"unique_opened_count"`
	UniqueClickedCount            *uint64 `json:"unique_clicked_count"`
	UnsubscribedCount             *uint64 `json:"unsubscribed_count"`
	ComplainedCount               *uint64 `json:"complained_count"`
	FailedCount                   *uint64 `json:"failed_count"`
	TemporaryFailedCount          *uint64 `json:"temporary_failed_count"`
	PermanentFailedCount          *uint64 `json:"permanent_failed_count"`
	ESPBlockCount                 *uint64 `json:"esp_block_count"`
	WebhookCount                  *uint64 `json:"webhook_count"`
	PermanentFailedOptimizedCount *uint64 `json:"permanent_failed_optimized_count"`
	PermanentFailedOldCount       *uint64 `json:"permanent_failed_old_count"`
	BouncedCount                  *uint64 `json:"bounced_count"`
	HardBouncesCount              *uint64 `json:"hard_bounces_count"`
	SoftBouncesCount              *uint64 `json:"soft_bounces_count"`
	DelayedBounceCount            *uint64 `json:"delayed_bounce_count"`
	SuppressedBouncesCount        *uint64 `json:"suppressed_bounces_count"`
	SuppressedUnsubscribedCount   *uint64 `json:"suppressed_unsubscribed_count"`
	SuppressedComplaintsCount     *uint64 `json:"suppressed_complaints_count"`
	DeliveredFirstAttemptCount    *uint64 `json:"delivered_first_attempt_count"`
	DelayedFirstAttemptCount      *uint64 `json:"delayed_first_attempt_count"`
	DeliveredSubsequentCount      *uint64 `json:"delivered_subsequent_count"`
	DeliveredTwoPlusAttemptsCount *uint64 `json:"delivered_two_plus_attempts_count"`

	DeliveredRate     string `json:"delivered_rate"`
	OpenedRate        string `json:"opened_rate"`
	ClickedRate       string `json:"clicked_rate"`
	UniqueOpenedRate  string `json:"unique_opened_rate"`
	UniqueClickedRate string `json:"unique_clicked_rate"`
	UnsubscribedRate  string `json:"unsubscribed_rate"`
	ComplainedRate    string `json:"complained_rate"`
	BounceRate        string `json:"bounce_rate"`
	FailRate          string `json:"fail_rate"`
	PermanentFailRate string `json:"permanent_fail_rate"`
	TemporaryFailRate string `json:"temporary_fail_rate"`
	DelayedRate       string `json:"delayed_rate"`

	// usage metrics
	EmailValidationCount        *uint64 `json:"email_validation_count"`
	EmailValidationPublicCount  *uint64 `json:"email_validation_public_count"`
	EmailValidationValidCount   *uint64 `json:"email_validation_valid_count"`
	EmailValidationSingleCount  *uint64 `json:"email_validation_single_count"`
	EmailValidationBulkCount    *uint64 `json:"email_validation_bulk_count"`
	EmailValidationListCount    *uint64 `json:"email_validation_list_count"`
	EmailValidationMailgunCount *uint64 `json:"email_validation_mailgun_count"`
	EmailValidationMailjetCount *uint64 `json:"email_validation_mailjet_count"`
	EmailPreviewCount           *uint64 `json:"email_preview_count"`
	EmailPreviewFailedCount     *uint64 `json:"email_preview_failed_count"`
	LinkValidationCount         *uint64 `json:"link_validation_count"`
	LinkValidationFailedCount   *uint64 `json:"link_validation_failed_count"`
	SeedTestCount               *uint64 `json:"seed_test_count"`
}

type MetricsDimension struct {
	// The dimension
	Dimension string `json:"dimension"`
	// The dimension value
	Value string `json:"value"`
	// The dimension value in displayable form
	DisplayValue string `json:"display_value"`
}
