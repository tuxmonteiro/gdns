package main

import (
	"encoding/json"
	"gopkg.in/gin-gonic/gin.v1"
	"io/ioutil"
	"net/http"
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

func readBody(c *gin.Context) []byte {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func resultWithCreated(c *gin.Context, r interface{}) {
	if err := json.Unmarshal(readBody(c), &r); err != nil {
		panic(err)
	}
	c.JSON(http.StatusCreated, r)
}

func main() {
	router := gin.Default()
	router.POST("/domains/:domain_id/records.json", func(c *gin.Context) {
		c.Param("domain_id")
		record := RecordRoot{}
		record.Record.Id = 1
		resultWithCreated(c, record)
	})
	router.POST("/domains.json", func(c *gin.Context) {
		domain := DomainRoot{}
		domain.Domain.Id = 1
		resultWithCreated(c, domain)
	})
	notify := func(c *gin.Context) {
		// TODO: notify
		c.Status(http.StatusNoContent)
		c.Done()
	}
	router.POST("/bind9/export.json", notify)
	router.POST("/bind9/schedule_export.json", notify)

	router.Run(":8080")
}
