package sapcontrol

import "time"

// SapcontrolOutput contains output of sapcontrol command.
type SapcontrolOutput struct {
	ExecTimePoint time.Time
	Function      string
	Status        string
	Content       []string
}

// SupportedSapControlFunctions contains the supported functions of sapcontrol.
var SupportedSapControlFunctions = []string{
	"GetAlertTree",
	"GetProcessList",
}

// SAPControlTimeLayout describes time format which is used by sapcontrol tool.
const SAPControlTimeLayout = "02.01.2006 15:04:05MST"

// TimeLayout describes time format which is used in SAP CCMS.
const TimeLayout = "2006 01 02 15:04:05MST"

// AlertColors contains all valid COLORs of SAP CCMS.
var AlertColors = []string{
	"GREEN",
	"YELLOW",
	"RED",
	"GRAY",
}

// AlertNode represents an AlertNode/record from SAP GetAlertTree.
type AlertNode struct {
	ID             int
	Name           string
	Parent         int
	ActualValue    string
	Description    string
	Time           time.Time
	AnalyseTool    string
	VisibleLevel   string
	HighAlertValue string
	AlDescription  string
	AlTime         time.Time
	Tid            string
}

// AlertNodes represents a slice of AlertNode.
type AlertNodes []AlertNode

// AlertNodeRecordSize is the size of an valid record from SAP GetAlertTree.
const AlertNodeRecordSize = 11

// AlertNodeFixLandmarkField is offset to field with fixed landmark.
const AlertNodeFixLandmarkField = 3

// Process represents a process from SAP GetProcessList.
type Process struct {
	Name        string
	Description string
	DispStatus  string
	TextStatus  string
	StartTime   time.Time
	ElapsedTime time.Duration
	Pid         int
}

// ProcessList represents a slice of Process.
type ProcessList []Process

// ProcessListRecordSize is the size of an valid record from SAP GetAlertTree.
const ProcessListRecordSize = 7

// ProcessingError represents a warning which occured during processing of sapcontrol output.
type ProcessingError struct {
	Action     string
	LineNumber int
	Message    string
}

// ProcessingErrors represents a slice of ProcessingError.
type ProcessingErrors []ProcessingError
