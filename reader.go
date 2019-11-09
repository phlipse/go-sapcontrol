package sapcontrol

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// systems timezone
var tz string

func init() {
	// get systems timezone
	tz, _ = time.Now().In(time.Local).Zone()
}

// Read reads in output of sapcontrol and stores in SapcontrolOutput struct.
func Read(r io.Reader) (SapcontrolOutput, error) {
	// compile regexs
	// do it inside of function body because Read() will be called normally once per execution
	reWord := regexp.MustCompile(`^\s*(\w+)\s*$`)
	reTime := regexp.MustCompile(`^\s*([0-9.]+\s[0-9:]+)\s*$`)

	var output SapcontrolOutput

	// create new scanner from io.Reader
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "":
			// skip empty lines
			continue
		case strings.HasPrefix(line, "Unknown webmethod:"):
			// unknown function call
			return SapcontrolOutput{}, fmt.Errorf(scanner.Text())
		case reTime.MatchString(line):
			// append local timezone
			t := reTime.FindString(line) + tz
			output.ExecTimePoint, _ = time.Parse(SAPControlTimeLayout, t)
		case reWord.MatchString(line):
			word := reWord.FindString(line)

			if stringInSlice(word, SupportedSapControlFunctions) {
				output.Function = word
			} else if word == "OK" {
				output.Status = "OK"
			} else {
				// not sure what to do...
				continue
			}
		default:
			output.Content = append(output.Content, line)
		}
	}

	// check for errors during scan
	err := scanner.Err()

	return output, err
}

// GetProcssList extracts ProcessList from sapcontrol output.
func (s SapcontrolOutput) GetProcssList() (ProcessList, ProcessingErrors) {
	var pErr ProcessingErrors

	// check if function call was "GetProcessList"
	if s.Function != "GetProcessList" {
		pErr = append(pErr, ProcessingError{
			Action:     "skip",
			LineNumber: -1,
			Message:    fmt.Sprintf("all records were skipped due to wrong function: got %s want %s", s.Function, "GetProcessList"),
		})
		return nil, pErr
	}

	var list ProcessList

	// start at 1 because of header
	for i := 1; i < len(s.Content); i++ {
		// get fields of record
		subs := strings.Split(s.Content[i], ",")

		// each record needs to have at least 11 fields
		if len(subs) < ProcessListRecordSize {
			pErr = append(pErr, ProcessingError{
				Action:     "skip",
				LineNumber: i,
				Message:    fmt.Sprintf("record in line %d was skipped due to too less fields: got %d want %d", i, len(subs), ProcessListRecordSize),
			})
			continue
		}

		// prepare fields
		name := strings.TrimSpace(subs[0])
		description := strings.TrimSpace(subs[1])
		dispStatus := strings.TrimSpace(subs[2])
		textStatus := strings.TrimSpace(subs[3])
		startTime, err := time.Parse(TimeLayout, strings.TrimSpace(subs[4])+tz)
		if err != nil {
			pErr = append(pErr, ProcessingError{
				Action:     "ignore",
				LineNumber: i,
				Message:    fmt.Sprintf("starttime field of record in line %d was ignored due to parsing error: %v", i, err),
			})
		}
		// elapsedtime is a duration
		duration := strings.TrimSpace(subs[5])
		d := strings.Split(duration, ":")
		if len(d) == 3 {
			// insert h, m and s characters to get it parsed
			duration = fmt.Sprintf("%sh%sm%ss", d[0], d[1], d[2])
		}
		elapsedTime, err := time.ParseDuration(duration)
		if err != nil {
			pErr = append(pErr, ProcessingError{
				Action:     "ignore",
				LineNumber: i,
				Message:    fmt.Sprintf("elapsedtime field of record in line %d was ignored due to parsing error: %v", i, err),
			})
		}
		pid, err := strconv.Atoi(strings.TrimSpace(subs[6]))
		if err != nil {
			pErr = append(pErr, ProcessingError{
				Action:     "ignore",
				LineNumber: i,
				Message:    fmt.Sprintf("pid field of record in line %d was ignored due to parsing error: %v", i, err),
			})
		}

		// build up Process from prepared fields
		process := Process{
			Name:        name,
			Description: description,
			DispStatus:  dispStatus,
			TextStatus:  textStatus,
			StartTime:   startTime,
			ElapsedTime: elapsedTime,
			Pid:         pid,
		}

		// add Process to slice
		list = append(list, process)
	}

	return list, pErr
}

// GetAlertNodes extracts AlertNodes from sapcontrol output.
func (s SapcontrolOutput) GetAlertNodes() (AlertNodes, ProcessingErrors) {
	var pErr ProcessingErrors

	// check if function call was "GetAlertTree"
	if s.Function != "GetAlertTree" {
		pErr = append(pErr, ProcessingError{
			Action:     "skip",
			LineNumber: -1,
			Message:    fmt.Sprintf("all records were skipped due to wrong function: got %s want %s", s.Function, "GetAlertTree"),
		})
		return nil, pErr
	}

	var nodes AlertNodes

	// lineNumber is not same as i because of random line breaks!
	lineNumber := -1
	// start at 1 because of header
	for i := 1; i < len(s.Content); i++ {
		// initially start at -1 and increment at beginning because of some continue calls, defer is slow
		lineNumber++

		// check if line ends with ; else read next line and join them
		// some csv fields contain new lines, this is bad
		// check only up to 10 times to prevent endless loops
		var line string
		for j := 0; j < 10; j++ {
			line += s.Content[i]
			if strings.HasSuffix(line, ";") {
				break
			}
			i++
		}

		// we can not use encoding/csv reader because the provided records are often invalid
		// there are often field separators inside fields
		// so do it manually
		subs := strings.Split(line, ",")

		// each record needs to have at least 11 fields
		if len(subs) < AlertNodeRecordSize {
			pErr = append(pErr, ProcessingError{
				Action:     "skip",
				LineNumber: i,
				Message:    fmt.Sprintf("record in line %d was skipped due to too less fields: got %d want %d", i, len(subs), AlertNodeRecordSize),
			})
			continue
		}

		// check if FixLandmarkField is AlertColor because it is our fixed landmark
		if !stringInSlice(strings.TrimSpace(subs[AlertNodeFixLandmarkField-1]), AlertColors) {
			pErr = append(pErr, ProcessingError{
				Action:     "skip",
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("record in line %d was skipped due to not finding fixed landmark at expected offset %d", lineNumber, AlertNodeFixLandmarkField),
			})
			continue
		}
		// search second COLOR because it is our dynamic landmark
		landmark := 0
		for i := 3; i < len(subs); i++ {
			if stringInSlice(strings.TrimSpace(subs[i]), AlertColors) {
				landmark = i
				break
			}
		}
		if landmark == 0 {
			pErr = append(pErr, ProcessingError{
				Action:     "skip",
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("record in line %d was skipped due to not finding dynamic landmark", lineNumber),
			})
			continue
		}

		// prepare fields
		name := strings.TrimSpace(subs[0])
		parent, err := strconv.Atoi(strings.TrimSpace(subs[1]))
		if err != nil {
			pErr = append(pErr, ProcessingError{
				Action:     "ignore",
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("parent field of record in line %d was ignored due to parsing error: %v", lineNumber, err),
			})
		}
		actualValue := strings.TrimSpace(subs[2])
		description := strings.TrimSpace(strings.Join(subs[3:landmark-3], ","))
		timeParsed, err := time.Parse(TimeLayout, strings.TrimSpace(subs[landmark-3])+tz)
		// only check if not empty string was parsed
		if err != nil && len(strings.TrimSpace(subs[landmark-3])) > 0 {
			pErr = append(pErr, ProcessingError{
				Action:     "ignore",
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("Time field of record in line %d was ignored due to parsing error: %v", lineNumber, err),
			})
		}
		analyseTool := strings.TrimSpace(subs[landmark-2])
		visibleLevel := strings.TrimSpace(subs[landmark-1])
		highAlertValue := strings.TrimSpace(subs[landmark])
		alDescription := strings.TrimSpace(strings.Join(subs[landmark+1:len(subs)-2], ","))
		alTimeParsed, err := time.Parse(TimeLayout, strings.TrimSpace(subs[len(subs)-2])+tz)
		if err != nil && len(strings.TrimSpace(subs[len(subs)-2])) > 0 {
			pErr = append(pErr, ProcessingError{
				Action:     "ignore",
				LineNumber: lineNumber,
				Message:    fmt.Sprintf("AlTime field of record in line %d was ignored due to parsing error: %v", lineNumber, err),
			})
		}
		tid := strings.TrimSpace(subs[len(subs)-1])

		// build up AlertNode from prepared fields
		node := AlertNode{
			ID:             lineNumber,
			Name:           name,
			Parent:         parent,
			ActualValue:    actualValue,
			Description:    description,
			Time:           timeParsed,
			AnalyseTool:    analyseTool,
			VisibleLevel:   visibleLevel,
			HighAlertValue: highAlertValue,
			AlDescription:  alDescription,
			AlTime:         alTimeParsed,
			Tid:            tid,
		}

		// add AlertNode to slice
		nodes = append(nodes, node)
	}

	return nodes, pErr
}
