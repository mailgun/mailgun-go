package mailgun

type Complaint struct {
	Count     int
	CreatedAt string
	Address   string
}

func (m *mailgunImpl) GetComplaints(limit, skip int) (int, []interface{}, error) {
	// TODO - this is NOT complete!
	return -1, nil, nil
}
