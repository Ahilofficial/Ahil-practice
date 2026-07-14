package model

import "time"

type Log struct {
	Service   string      `json:"service"`
	Method    string      `json:"method"`
	Endpoint  string      `json:"endpoint"`
	Request   interface{} `json:"request"`
	Response  interface{} `json:"response"`
	Status    int         `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}