package main

type Request struct {
	HttpVersion string            `json:"Version"`
	Path        string            `json:"Path"`
	Method      string            `json:"Method"`
	Body        string            `json:"Body"`
	Headers     map[string]string `json:"Headers"`
}

/////////////////////////////////////
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
