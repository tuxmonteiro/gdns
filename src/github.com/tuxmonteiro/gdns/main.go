package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
	"gopkg.in/square/go-jose.v1/json"
	_ "os/exec"
	"strconv"
)

var (
	maxIdleUpstreamConns = flag.Int("maxIdleUpstreamConns", 50, "The maximum idle connections to upstream host")
	pdnsServer           = flag.String("pdnsServer", "127.0.0.1:8000", "PowerDNS host. May include port in the form 'host:port'")
	pdnsServerToken      = flag.String("pdnsServerToken", "password", "PowerDNS token.")

	client *fasthttp.HostClient

	domainCount = 0
	domains     map[int]string
)

type domainRoot struct {
	Domain struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		TTL           string `json:"ttl"`
		Notes         string `json:"notes"`
		PrimaryNS     string `json:"primary_ns"`
		Contact       string `json:"contact"`
		Refresh       int    `json:"refresh"`
		Retry         int    `json:"retry"`
		Expire        int    `json:"expire"`
		Minimum       int    `json:"minimum"`
		AuthorityType string `json:"authority_type"`
	} `json:"domain"`
}

type recordRoot struct {
	Record struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"record"`
}

type pDNSWriteZone struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Masters     []string `json:"masters"`
	NameServers []string `json:"nameservers"`
}

type pDNSReadZone struct {
	Account        string   `json:"account"`
	DNSSec         bool     `json:"dnssec"`
	ID             string   `json:"id"`
	Kind           string   `json:"kind"`
	LastCheck      int      `json:"last_check"`
	Masters        []string `json:"masters"`
	Name           string   `json:"name"`
	NotifiedSerial int      `json:"notified_serial"`
	Serial         int      `json:"serial"`
	SoaEdit        string   `json:"soa_edit"`
	SoaEditAPI     string   `json:"soa_edit_api"`
	URL            string   `json:"url"`
	RRSets         []struct {
		Comments []string `json:"comments"`
		Name     string   `json:"name"`
		Records  []struct {
			Content  string `json:"content"`
			Disabled bool   `json:"disabled"`
		} `json:"records"`
		TTL  int    `json:"ttl"`
		Type string `json:"type"`
	} `json:"rrsets"`
}

type simpleRecord struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}

type rr struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	TTL        int            `json:"ttl"`
	ChangeType string         `json:"changetype"`
	Records    []simpleRecord `json:"records"`
}

type pDNSRRSets struct {
	RRSets []rr `json:"rrsets"`
}

func createRecords(c *iris.Context) {
	domainId, _ := strconv.Atoi(c.Param("domain_id"))
	zoneName := domains[domainId]

	record := recordRoot{}
	json.Unmarshal(c.Request.Body(), &record)

	aRecord := simpleRecord{}
	aRecord.Disabled = false
	aRecord.Content = record.Record.Content

	aRr := rr{}
	aRr.Name = record.Record.Name
	aRr.ChangeType = "REPLACE"
	aRr.Type = record.Record.Type
	aRr.TTL = 300
	aRr.Records = make([]simpleRecord, 0)
	aRr.Records = append(aRr.Records, aRecord)

	recordPdns := pDNSRRSets{}
	recordPdns.RRSets = make([]rr, 0)
	recordPdns.RRSets = append(recordPdns.RRSets, aRr)

	do(iris.MethodPatch, c, fmt.Sprintf("http://%s/api/v1/servers/localhost/zones/%s", *pdnsServer, zoneName), &recordPdns, &record)
	pdnsUpdate(zoneName)
}

func createZone(c *iris.Context) {
	domain := domainRoot{}
	json.Unmarshal(c.Request.Body(), &domain)

	domainCount++
	domain.Domain.ID = domainCount
	domains[domainCount] = domain.Domain.Name

	zone := pDNSWriteZone{}
	zone.Masters = make([]string, 0)
	zone.NameServers = make([]string, 0)
	zone.Name = fmt.Sprintf("%s.", domain.Domain.Name)
	zone.Kind = "Native"

	do(iris.MethodPost, c, fmt.Sprintf("http://%s/api/v1/servers/localhost/zones", *pdnsServer), &zone, &domain)
	rndcAddZone(domain.Domain.Name)
}

func rndcAddZone(zone string) {
	//exec.Command("/usr/sbin/rndc", "addzone", zone, "'{type slave; masters port 5353 { ::1; }; allow-notify { ::1; };};'").Output()
}

func pdnsUpdate(zone string) {
	//exec.Command("/usr/bin/pdnsutil", "increase-serial", zone).Output()
	//exec.Command("/usr/bin/pdns_control", "notify", zone).Output()
}

func do(method string, c *iris.Context, url string, data interface{}, originalData interface{}) {
	request := fasthttp.Request{}
	request.Header.SetMethod(method)
	request.Header.SetContentType("application/json")
	request.Header.Add("X-API-Key", *pdnsServerToken)
	request.SetRequestURI(url)
	b, _ := json.Marshal(data)
	request.SetBody(b)
	response := fasthttp.Response{}
	client.Do(&request, &response)
	c.JSON(response.StatusCode(), &originalData)
}

func notify(c *iris.Context) {
	// TODO: notify
	c.SetStatusCode(iris.StatusNoContent)
}

func newClient() *fasthttp.HostClient {
	return &fasthttp.HostClient{
		Addr:                          *pdnsServer,
		MaxConns:                      *maxIdleUpstreamConns,
		DisableHeaderNamesNormalizing: true,
	}
}

func main() {
	flag.Parse()

	client = newClient()
	domains = make(map[int]string)

	iris.Post("/domains/:domain_id/records.json", createRecords)
	iris.Post("/domains.json", createZone)
	iris.Post("/bind9/export.json", notify)
	iris.Post("/bind9/schedule_export.json", notify)

	iris.Listen(":8080")
}
