package model

import "time"

type RawReport struct {
	Email                   EmailEnvelope            `json:"email"`
	Attachment              AttachmentEnvelope       `json:"attachment"`
	Envelope                Envelope                 `json:"envelope"`
	Reports                 []Report                 `json:"reports"`
	ArrivalDepartureReports []ArrivalDepartureReport `json:"arrival-departure-reports"`
	ParsedReportBody        ParsedReportBody         `json:"parsed-report-body"`
	Remarks                 []Remark                 `json:"remarks"`
	Errors                  []string                 `json:"errors,omitempty"`
}

type EmailEnvelope struct {
	ReceivedTime time.Time `json:"received-time"`
	ID           string    `json:"email-id"`
	Subject      string    `json:"email-subject"`
	From         string    `json:"email-from"`
}

type AttachmentEnvelope struct {
	// name of the xlsx attached file, could contain voyage number.
	FileName string `json:"file-name"`

	// name of the xlsx sheet of attached file, could contain voyage number.
	//
	// Note: sheet name also contains voyage number,
	// but these two values could be different,
	// need to decide which one to use
	SheetName string `json:"sheet-name"`
}

// Envelope is intended to be parsed from Email Subject
type Envelope struct {
	ReportType   string `json:"report-type"`
	VesselIMO    string `json:"vessel-imo"`
	VesselID     string `json:"vessel-id"`
	VoyageNumber string `json:"voyage-number"`
	Error        string `json:"error,omitempty"`
}

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

type Remark struct {
	RawLine string `json:"raw-line"`
}

type Vessel struct {
	VesselID   string
	VesselName string
	VesselIMO  string
}

// ParsedReportBody represents arrival/departure report key/value pairs from email body content
type ParsedReportBody struct {
	Raw string `json:"raw"`
	// KeyValuePairs represents key/value list from email body
	// example:
	// 		1) Arrival Port Name : Yokohama, Japan
	// 		2) Time Difference from UTC on Arrival : +09:00
	// 		3) Time at S/B Eng.(EOSP) : 0530lt, 17th Jan. 2023
	// 		4) Arrival Pilot Station Position : 35-09.2N, 139-46.1E
	// 		5) Time of Arrival at Pilot Station : 0646lt, 17th Jan. 2023
	// 		6) Time of NOR tendered : N/A
	Fields  []Field  `json:"fields"`
	Remarks []Remark `json:"remarks"`
	Errors  []string `json:"errors,omitempty"`
}

type ArrivalDepartureReport struct {
	// Raw possible example: "Arrival/~Departure~ Report (for ~loading~/discharging operation)"
	// strikeout text could be marked with '~' before and after the word.
	// Or how are we going to indicate that this particular word was strikeout from the source?
	// What do you think?
	Raw string `json:"raw"`
	// Loading or Discharging or Bunkering operation.
	// Need to extract the strikeout text style to identify the exact operation.
	Operation string `json:"operation"`

	Arrival   bool `json:"arrival,omitempty"`
	Departure bool `json:"departure,omitempty"`

	Headings []Heading `json:"headings"`
	Rows     []Row     `json:"rows"`
	Errors   []string  `json:"errors,omitempty"`
}

type Heading struct {
	// Raw possible example: "Steaming distance from last port (NM) : 162 NM"
	Raw   string `json:"raw"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Row struct {
	// Raw possible example: "EOSP : [local time: 0530 17th] [UTC:2030 16th] [FO ROB:408.5MT] [MGO ROB:208.11MT]"
	Raw    string  `json:"raw"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}
