package main

import (
	"github.com/kataras/iris"
)

type DomainRoot struct {
	Domain struct {
		Id            int    `json:"id"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		Ttl           string `json:"ttl"`
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

type RecordRoot struct {
	Record struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"record"`
}

type PDnsWriteZone struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Masters     []string `json:"masters"`
	NameServers []string `json:"nameservers"`
}

type PDnsReadZone struct {
	Account        string   `json:"account"`
	DnsSec         bool     `json:"dnssec"`
	Id             string   `json:"id"`
	Kind           string   `json:"kind"`
	LastCheck      int      `json:"last_check"`
	Masters        []string `json:"masters"`
	Name           string   `json:"name"`
	NotifiedSerial int      `json:"notified_serial"`
	Serial         int      `json:"serial"`
	SoaEdit        string   `json:"soa_edit"`
	SoaEditApi     string   `json:"soa_edit_api"`
	Url            string   `json:"url"`
	RRSets         []struct {
		Comments []string `json:"comments"`
		Name     string   `json:"name"`
		Records  []struct {
			Content  string `json:"content"`
			Disabled bool   `json:"disabled"`
		} `json:"records"`
		Ttl  int    `json:"ttl"`
		Type string `json:"type"`
	} `json:"rrsets"`
}

type PDnsRRSets struct {
	RRSets []struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		Ttl        int    `json:"ttl"`
		ChangeType string `json:"changetype"`
		Records    []struct {
			Content  string `json:"content"`
			Disabled bool   `json:"disabled"`
		} `json:"records"`
	} `json:"rrsets"`
}

func resultWithCreated(c *iris.Context, r interface{}) {
	c.ReadJSON(&r)
	c.JSON(iris.StatusCreated, r)
}

func createRecords(c *iris.Context) {
	c.Param("domain_id")
	record := RecordRoot{}
	record.Record.Id = 1
	resultWithCreated(c, &record)
}

func createZone(c *iris.Context) {
	domain := DomainRoot{}
	domain.Domain.Id = 1
	resultWithCreated(c, &domain)
}

func notify(c *iris.Context) {
	// TODO: notify
	c.SetStatusCode(iris.StatusNoContent)
}

func main() {
	iris.Post("/domains/:domain_id/records.json", createRecords)
	iris.Post("/domains.json", createZone)
	iris.Post("/bind9/export.json", notify)
	iris.Post("/bind9/schedule_export.json", notify)

	iris.Listen(":8080")
}
