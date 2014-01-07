package mailgun

type Complaint struct {
	Count     int
	CreatedAt string
	Address   string
}

func (m *mailgunImpl) GetComplaints(limit, skip int) (int, []interface{}, error) {
	return -1, nil, nil
}
