package main

import "strings"

type ReportData struct {
	date   string
	values map[string]string
}

func NewReportData(values map[string]string, date string) ReportData {
	if values == nil {
		values = make(map[string]string)
	}
	return ReportData{
		date:   date,
		values: values,
	}
}

func (r *ReportData) SetValue(key, value string) {
	if len(strings.TrimSpace(key)) == 0 {
		return
	}

	if r.values == nil {
		r.values = make(map[string]string)
	}
	r.values[key] = value
}

func (r *ReportData) GetValue(key string) (string, bool) {
	v, ok := r.values[key]
	return v, ok
}

func (r *ReportData) GetDate() string {
	return r.date
}
