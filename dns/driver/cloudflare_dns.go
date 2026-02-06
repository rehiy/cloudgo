package drivers

import (
	"github.com/rehiy/cloudgo/dns"
	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/provider/cloudflare"

	cf "github.com/cloudflare/cloudflare-go"
)

type CloudflareDnsDriver struct {
	client *cloudflare.Client
	api    *cf.API
	rq     *provider.ReqeustParam
}

func NewCloudflareDnsDriver(rq *provider.ReqeustParam) *CloudflareDnsDriver {

	client := cloudflare.NewClient(rq)
	api, _ := client.NewApi()

	return &CloudflareDnsDriver{client, api, rq}

}

func (p *CloudflareDnsDriver) ListZones() ([]*dns.Zone, error) {

	resp, err := p.api.ListZones(p.client.Ctx)

	if err != nil {
		return nil, err
	}

	zones := make([]*dns.Zone, 0)

	for _, zone := range resp {
		zones = append(zones, &dns.Zone{
			Id:         zone.ID,
			Domain:     zone.Name,
			CreateTime: 0,
		})
	}

	return zones, nil

}

func (p *CloudflareDnsDriver) DetailZone(zone *dns.Zone) (*dns.Zone, error) {

	resp, err := p.api.ZoneDetails(p.client.Ctx, zone.Id)

	if err != nil {
		return nil, err
	}

	dnsServers := make([]string, 0)
	for _, dnsServer := range resp.NameServers {
		dnsServers = append(dnsServers, dnsServer)
	}

	data := &dns.Zone{
		Id:         resp.ID,
		Domain:     resp.Name,
		CreateTime: 0,
		DnsServers: dnsServers,
	}

	return data, nil

}

func (p *CloudflareDnsDriver) CreateZone(zone *dns.Zone) (*dns.Zone, error) {

	account := cf.Account{ID: ""}

	resp, err := p.api.CreateZone(p.client.Ctx, zone.Domain, false, account, "full")

	if err != nil {
		return nil, err
	}

	data := &dns.Zone{
		Id:         resp.ID,
		Domain:     resp.Name,
		CreateTime: 0,
	}

	return data, nil

}

func (p *CloudflareDnsDriver) UpdateZone(zone *dns.Zone) (*dns.Zone, error) {

	_, err := p.api.ZoneDetails(p.client.Ctx, zone.Id)

	if err != nil {
		return nil, err
	}

	return zone, nil

}

func (p *CloudflareDnsDriver) DeleteZone(zone *dns.Zone) error {

	_, err := p.api.DeleteZone(p.client.Ctx, zone.Id)

	if err != nil {
		return err
	}

	return nil

}

func (p *CloudflareDnsDriver) ListRecords(zone *dns.Zone) ([]*dns.Record, error) {

	rc := &cf.ResourceContainer{
		Identifier: zone.Id,
	}

	resp, _, err := p.api.ListDNSRecords(p.client.Ctx, rc, cf.ListDNSRecordsParams{})

	if err != nil {
		return nil, err
	}

	records := make([]*dns.Record, 0)

	for _, record := range resp {
		recordType := dns.RecordType(record.Type)

		records = append(records, &dns.Record{
			Id:       record.ID,
			Name:     record.Name,
			Type:     recordType,
			Value:    record.Content,
			TTL:      record.TTL,
			Priority: int(*record.Priority),
		})
	}

	return records, nil

}

func (p *CloudflareDnsDriver) DetailRecord(zone *dns.Zone, record *dns.Record) (*dns.Record, error) {

	rc := &cf.ResourceContainer{
		Identifier: zone.Id,
	}

	resp, err := p.api.GetDNSRecord(p.client.Ctx, rc, record.Id)

	if err != nil {
		return nil, err
	}

	recordType := dns.RecordType(resp.Type)

	data := &dns.Record{
		Id:    resp.ID,
		Name:  resp.Name,
		Type:  recordType,
		Value: resp.Content,
		TTL:   resp.TTL,
	}

	return data, nil

}

func (p *CloudflareDnsDriver) CreateRecord(zone *dns.Zone, record *dns.Record) (*dns.Record, error) {

	rc := &cf.ResourceContainer{
		Identifier: zone.Id,
	}

	_, err := p.api.CreateDNSRecord(p.client.Ctx, rc, cf.CreateDNSRecordParams{
		Type:    string(record.Type),
		Name:    record.Name,
		Content: record.Value,
		TTL:     record.TTL,
	})

	if err != nil {
		return nil, err
	}

	return record, nil

}

func (p *CloudflareDnsDriver) UpdateRecord(zone *dns.Zone, record *dns.Record) (*dns.Record, error) {

	rc := &cf.ResourceContainer{
		Identifier: zone.Id,
	}

	_, err := p.api.UpdateDNSRecord(p.client.Ctx, rc, cf.UpdateDNSRecordParams{
		Type:    string(record.Type),
		Name:    record.Name,
		Content: record.Value,
		TTL:     record.TTL,
	})

	if err != nil {
		return nil, err
	}

	return record, nil

}

func (p *CloudflareDnsDriver) DeleteRecord(zone *dns.Zone, record *dns.Record) error {

	rc := &cf.ResourceContainer{
		Identifier: zone.Id,
	}

	err := p.api.DeleteDNSRecord(p.client.Ctx, rc, record.Id)

	return err

}
