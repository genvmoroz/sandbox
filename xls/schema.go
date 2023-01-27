package main

// schema

type Report struct {
	Date   string  `json:"date"`
	Fields []Field `json:"fields"`
	Error  string  `json:"error,omitempty"`
}

type Field struct {
	Raw   string `json:"raw"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Units string `json:"units"`
}
