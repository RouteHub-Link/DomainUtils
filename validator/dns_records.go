package validator

import (
	"github.com/miekg/dns"
)

type DNSRecord struct {
	RecordType string
	Records    []string
}

func fetchDNSRecords(domain string, recordType uint16) ([]dns.RR, error) {

	client := dns.Client{}
	server := "1.1.1.1:53"

	// Query for DNS records
	msg := dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), recordType)
	msg.RecursionDesired = true

	// Send the DNS query to the server
	response, _, err := client.Exchange(&msg, server)
	if err != nil {
		return nil, err
	}

	return response.Answer, nil
}
