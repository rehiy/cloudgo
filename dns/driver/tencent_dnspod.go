package drivers

import (
	"strconv"

	"github.com/rehiy/cloudgo/dns"
	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/provider/tencent"

	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

type TecentDnspodDriver struct {
	client *tencent.Client
	dnspod *dnspod.Client
	rq     *provider.ReqeustParam
}

func NewTencentDnspodDriver(rq *provider.ReqeustParam) *TecentDnspodDriver {

	client := tencent.NewClient(rq)
	dp, _ := client.Dnspod()

	return &TecentDnspodDriver{client, dp, rq}

}

func (p *TecentDnspodDriver) ListZones() ([]*dns.Zone, error) {

	resp, err := p.dnspod.DescribeDomainList(&dnspod.DescribeDomainListRequest{})

	if err != nil {
		return nil, err
	}

	zones := make([]*dns.Zone, 0)

	for _, domain := range resp.Response.DomainList {
		dnsServers := make([]string, 0)
		for _, dnsServer := range domain.EffectiveDNS {
			dnsServers = append(dnsServers, *dnsServer)
		}
		zones = append(zones, &dns.Zone{
			Id:          strconv.Itoa(int(*domain.DomainId)),
			Domain:      *domain.Name,
			PunyCode:    *domain.Punycode,
			DnsServers:  dnsServers,
			MinTTL:      int(*domain.TTL),
			Description: *domain.Remark,
		})
	}

	return zones, nil

}

func (p *TecentDnspodDriver) DetailZone(zone *dns.Zone) (*dns.Zone, error) {

	resp, err := p.dnspod.DescribeDomain(&dnspod.DescribeDomainRequest{
		Domain: &zone.Domain,
	})

	if err != nil {
		return nil, err
	}

	dnsServers := make([]string, 0)
	for _, dnsServer := range resp.Response.DomainInfo.DnspodNsList {
		dnsServers = append(dnsServers, *dnsServer)
	}

	data := &dns.Zone{
		Id:          strconv.Itoa(int(*resp.Response.DomainInfo.DomainId)),
		Domain:      *resp.Response.DomainInfo.Domain,
		PunyCode:    *resp.Response.DomainInfo.Punycode,
		DnsServers:  dnsServers,
		MinTTL:      int(*resp.Response.DomainInfo.TTL),
		Description: *resp.Response.DomainInfo.Remark,
	}

	return data, nil

}

func (p *TecentDnspodDriver) CreateZone(zone *dns.Zone) (*dns.Zone, error) {

	reps, err := p.dnspod.CreateDomain(&dnspod.CreateDomainRequest{
		Domain: &zone.Domain,
	})

	if err != nil {
		return nil, err
	}

	dnsServers := make([]string, 0)
	for _, dnsServer := range reps.Response.DomainInfo.GradeNsList {
		dnsServers = append(dnsServers, *dnsServer)
	}

	data := &dns.Zone{
		Id:         strconv.Itoa(int(*reps.Response.DomainInfo.Id)),
		Domain:     *reps.Response.DomainInfo.Domain,
		PunyCode:   *reps.Response.DomainInfo.Punycode,
		DnsServers: []string{},
	}

	return data, nil

}

func (p *TecentDnspodDriver) UpdateZone(zone *dns.Zone) (*dns.Zone, error) {

	_, err := p.dnspod.ModifyDomainRemark(&dnspod.ModifyDomainRemarkRequest{
		Domain: &zone.Domain,
		Remark: &zone.Description,
	})

	if err != nil {
		return nil, err
	}

	return zone, nil

}

func (p *TecentDnspodDriver) DeleteZone(zone *dns.Zone) error {

	_, err := p.dnspod.DescribeDomain(&dnspod.DescribeDomainRequest{
		Domain: &zone.Domain,
	})

	return err

}

func (p *TecentDnspodDriver) ListRecords(zone *dns.Zone) ([]*dns.Record, error) {

	resp, err := p.dnspod.DescribeRecordList(&dnspod.DescribeRecordListRequest{
		Domain: &zone.Domain,
	})

	if err != nil {
		return nil, err
	}

	records := make([]*dns.Record, 0)

	for _, record := range resp.Response.RecordList {
		recordType := dns.RecordType(*record.Type)

		records = append(records, &dns.Record{
			Id:       strconv.Itoa(int(*record.RecordId)),
			Name:     *record.Name,
			Type:     recordType,
			Value:    *record.Value,
			TTL:      int(*record.TTL),
			Priority: int(*record.MX),
		})
	}

	return records, nil

}

func (p *TecentDnspodDriver) DetailRecord(zone *dns.Zone, record *dns.Record) (*dns.Record, error) {

	id, _ := strconv.Atoi(record.Id)
	recordId := uint64(id)

	resp, err := p.dnspod.DescribeRecord(&dnspod.DescribeRecordRequest{
		Domain:   &zone.Domain,
		RecordId: &recordId,
	})

	if err != nil {
		return nil, err
	}

	recordType := dns.RecordType(*resp.Response.RecordInfo.RecordType)

	data := &dns.Record{
		Id:       strconv.Itoa(int(*resp.Response.RecordInfo.Id)),
		Name:     *resp.Response.RecordInfo.SubDomain,
		Type:     recordType,
		Value:    *resp.Response.RecordInfo.Value,
		TTL:      int(*resp.Response.RecordInfo.TTL),
		Priority: int(*resp.Response.RecordInfo.MX),
	}

	return data, nil

}

func (p *TecentDnspodDriver) CreateRecord(zone *dns.Zone, record *dns.Record) (*dns.Record, error) {

	ttl := uint64(record.TTL)
	priority := uint64(record.Priority)

	resp, err := p.dnspod.CreateRecord(&dnspod.CreateRecordRequest{
		Domain:     &zone.Domain,
		SubDomain:  &record.Name,
		RecordType: (*string)(&record.Type),
		RecordLine: &record.Line,
		Value:      &record.Value,
		TTL:        &ttl,
		MX:         &priority,
	})

	if err != nil {
		return nil, err
	}

	data := &dns.Record{
		Id: strconv.Itoa(int(*resp.Response.RecordId)),
	}

	return data, nil

}

func (p *TecentDnspodDriver) UpdateRecord(zone *dns.Zone, record *dns.Record) (*dns.Record, error) {

	id, _ := strconv.Atoi(record.Id)
	recordId := uint64(id)

	ttl := uint64(record.TTL)
	priority := uint64(record.Priority)

	_, err := p.dnspod.ModifyRecord(&dnspod.ModifyRecordRequest{
		Domain:     &zone.Domain,
		RecordId:   &recordId,
		SubDomain:  &record.Name,
		RecordType: (*string)(&record.Type),
		RecordLine: &record.Line,
		Value:      &record.Value,
		TTL:        &ttl,
		MX:         &priority,
	})

	if err != nil {
		return nil, err
	}

	return record, nil

}

func (p *TecentDnspodDriver) DeleteRecord(zone *dns.Zone, record *dns.Record) error {

	id, _ := strconv.Atoi(record.Id)
	recordId := uint64(id)

	_, err := p.dnspod.DeleteRecord(&dnspod.DeleteRecordRequest{
		Domain:   &zone.Domain,
		RecordId: &recordId,
	})

	return err

}
