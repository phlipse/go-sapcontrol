package sapcontrol

import (
	"regexp"
	"strconv"
)

var (
	// compiled regex should be global to compile it only once on start up
	reParseValueUnit = regexp.MustCompile(`^(-|-?\d+(.\d+)?)( ([^ ]+))?$`)
)

// stringInSlice returns bool if given string is in given slice.
func stringInSlice(s string, slice []string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}

// ParseValueUnit returns value and unit from given string. If no value/unit pair could be parsed, returned unit equals given string.
func ParseValueUnit(s string) (interface{}, string) {
	t := reParseValueUnit.FindStringSubmatch(s)
	if len(t) == 0 {
		return 0, s
	}

	// convert value "-" to "0"
	if t[1] == "-" {
		t[1] = "0"
	}

	// check conversion to int, it is mostly int ;)
	var v interface{}
	v, err := strconv.ParseInt(t[1], 10, 64)
	if err != nil {
		// check conversion to float
		v, err = strconv.ParseFloat(t[1], 64)
		if err != nil {
			return 0, s
		}
	}

	return v, t[4]
}
