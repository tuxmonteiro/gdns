package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/gin-gonic/gin.v1"
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

func main() {
	router := gin.Default()
	router.POST("/domains/:domain_id/records.json", func(c *gin.Context) {
		c.Param("domain_id")
		body := new(bytes.Buffer)
		body.ReadFrom(c.Request.Body)
		record := RecordRoot{}
		record.Record.Id = 1
		if err := json.Unmarshal(body.Bytes(), &record); err != nil {
			panic(err)
		}
		c.JSON(http.StatusCreated, record)
	})
	router.POST("/domains.json", func(c *gin.Context) {
		body := new(bytes.Buffer)
		body.ReadFrom(c.Request.Body)
		domain := DomainRoot{}
		domain.Domain.Id = 1
		if err := json.Unmarshal(body.Bytes(), &domain); err != nil {
			panic(err)
		}
		c.JSON(http.StatusCreated, domain)
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
