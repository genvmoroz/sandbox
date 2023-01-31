package main

import "time"

type RawReport struct {
	Email      EmailEnvelope      `json:"email"`
	Attachment AttachmentEnvelope `json:"attachment"`
	Envelope   Envelope           `json:"envelope"`
	Reports    []Report           `json:"reports"`
	Errors     []string           `json:"errors,omitempty"`
}

type EmailEnvelope struct {
	ReceivedTime time.Time `json:"received-time"`
	ID           string    `json:"email-id"`
	Subject      string    `json:"email-subject"`
	From         string    `json:"email-from"`
}

type AttachmentEnvelope struct {
	// name of the xlsx attached file, could contain voyage number and port name
	FileName string `json:"file-name"`
}

// Envelope is intended to be parsed from Email Subject
type Envelope struct {
	ReportType   string `json:"report-type"`
	VesselID     string `json:"vessel-id"`
	VesselIMO    string `json:"vessel-imo"`
	VoyageNumber string `json:"voyage-number"`
	Error        string `json:"error,omitempty"`
}

type Report struct {
	// Loading or Discharging or Bunkering operation.
	// Need to extract the strikeout text style to identify the exact operation.
	Operation string `json:"operation"`

	// Arrival or Departure.
	// Need to extract the strikeout text style to identify the exact activity.
	Activity string `json:"activity"`

	Port                         Port     `json:"port"`
	SteamingDistanceFromLastPort Distance `json:"steaming-distance-from-last-port"`
	TimeDifferenceFromUTC        string   `json:"time-difference-from-utc"`
	Fields                       []Field  `json:"fields"`
	Errors                       []string `json:"errors,omitempty"`
}

type Port struct {
	// Raw possible example: "Name of Port (Arrived/Departed) : Yokohama, Japan"
	Raw string `json:"raw"`

	// Arrived or Departed.
	// Questions: Do we really need it since Activity is already present above?
	//			  Is it possible Report.Activity == Arrival but Port.Activity == Departed?
	//
	// Need to extract the strikeout text style to identify the exact activity.
	Activity string `json:"activity"`

	Name    string `json:"name"`
	Country string `json:"country"`
}

type Distance struct {
	// Raw possible example: "Steaming distance from last port (NM) : 162 NM"
	Raw string `json:"raw"`

	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type Field struct {
	// Raw possible example: "EOSP : [local time: 0530 17th] [UTC:2030 16th] [FO ROB:408.5MT] [MGO ROB:208.11MT]"
	Raw string `json:"raw"`

	Key string `json:"key"`

	LocalTime string `json:"local-time"`
	UTC       string `json:"utc"`
	FoROB     ROB    `json:"fo-rob"`
	MgoROB    ROB    `json:"mgo-rob"`
}

type ROB struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}
