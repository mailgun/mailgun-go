package mailgun

type Domain struct {
	CreatedAt    string
	SMTPLogin    string
	Name         string
	SMTPPassword string
	Wildcard     bool
	SpamAction   bool
}

type DomainDns struct {

}

func (m *mailgunImpl) GetDomains(limit, skip int) (int, []Domain, error) {
	return -1, nil, nil
}

func (m *mailgunImpl) GetSingleDomain(domain string) (Domain, DomainDns, error) {
	return Domain{}, DomainDns{}, nil
}

func (m *mailgunImpl) CreateDomain(name string, smtpPassword string, spamAction bool, wildcard bool) error {
	return nil
}

func (m *mailgunImpl) DeleteDomain(name string) error {
	return nil
}
