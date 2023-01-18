package v2

import "time"

type RawReport struct {
	Email      EmailEnvelope      `json:"email"`
	Attachment AttachmentEnvelope `json:"attachment"`
	Envelope   Envelope           `json:"envelope"`
	Content    Content            `json:"content"`
	Remarks    []Remark           `json:"remarks"`
	Errors     []string           `json:"errors,omitempty"`
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
	ReportType string `json:"report-type"`
	VesselUUID string `json:"vessel-uuid"`
	// Should we have VesselIMO field here since we can obtain VesselUUID instead
	VesselIMO string `json:"vessel-imo"`
	// We must define where we should take the value from,
	// since the EMAIL has three places for this value
	VoyageNumber string `json:"voyage-number"`
	Error        string `json:"error,omitempty"`
}

type Content struct {
	Heads []Head       `json:"heads"`
	Data  []DataColumn `json:"data"`
}

type Head struct {
	RowID int `json:"row-id"`
	Key   Key `json:"keys"`
}

type Key struct {
	Raw     string `json:"raw"`
	Name    string `json:"name"`
	Details string `json:"details"`
	Error   string `json:"error,omitempty"`
}

type DataColumn struct {
	ReportDate string   `json:"report-date"`
	Records    []Record `json:"records"`
	Error      string   `json:"error,omitempty"`
}

type Record struct {
	RowID  int    `json:"row-id"`
	Values string `json:"values"`
	Error  string `json:"error,omitempty"`
}

type Remark struct {
	RawLine string `json:"raw-line"`
}
