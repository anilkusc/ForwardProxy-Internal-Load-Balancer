package main

import (
	"time"
)

type Request struct {
	HttpVersion string            `json:"Version"`
	Host        string            `json:"Host"`
	Method      string            `json:"Method"`
	Body        string            `json:"Body"`
	Headers     map[string]string `json:"Headers"`
	Date        time.Time         `json:"Date"`
}

type Response struct {
	HttpVersion string            `json:"Version"`
	Status      string            `json:"Status"`
	Body        string            `json:"Body"`
	Headers     map[string]string `json:"Headers"`
}

type Log struct {
	LogRequest  Request  `json:"Request"`
	LogResponse Response `json:"Response"`
}
