package mailgun

var mockAcceptedIncomingCount uint64 = 10

var expectedResponse = MetricsResponse{
	Start:      "Mon, 15 Apr 2024 00:00:00 +0000",
	Dimensions: []string{"time"},
	Items: []MetricsItem{
		{
			Dimensions: []MetricsDimension{
				{
					Dimension:    "time",
					Value:        "Mon, 15 Apr 2024 00:00:00 +0000",
					DisplayValue: "Mon, 15 Apr 2024 00:00:00 +0000",
				},
			},
			Metrics: Metrics{
				AcceptedIncomingCount: &mockAcceptedIncomingCount,
				ClickedRate:           "0.8300",
			},
		},
	},
}
