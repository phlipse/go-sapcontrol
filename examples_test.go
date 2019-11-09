package sapcontrol_test

import (
	"fmt"
	"io"

	"github.com/phlipse/go-sapcontrol"
)

func ExampleGetFilesystemUsage(r io.Reader) {
	// read in AlertNodes from io.Reader with given offset
	sap, err := sapcontrol.Read(r)
	if err != nil {
		panic(err)
	}
	nodes, _ := sap.GetAlertNodes()

	// get ID of 'OperatingSystem' child node
	osID := nodes.GetChildNodeByName("OperatingSystem", 0).ID
	// get ID of 'Filesystems' child node
	fsID := nodes.GetChildNodeByName("Filesystems", osID).ID
	// get all mounts below 'Filesystems' child node
	mounts := nodes.GetNodesByParentID(fsID)

	for _, m := range mounts {
		fmt.Printf("%s: %s used\n", m.Name, nodes.GetChildNodeByName("Percentage_Used", m.ID).Description)
	}
}

func ExampleGetNumberOfCPUs(r io.Reader) {
	// read in AlertNodes from io.Reader with given offset
	sap, err := sapcontrol.Read(r)
	if err != nil {
		panic(err)
	}
	nodes, _ := sap.GetAlertNodes()

	// get ID of 'OperatingSystem' child node
	osID := nodes.GetChildNodeByName("OperatingSystem", 0).ID
	// get ID of 'CPU' child node
	cpuID := nodes.GetChildNodeByName("CPU", osID).ID

	fmt.Printf("number of CPUs: %s\n", nodes.GetChildNodeByName("Number of CPUs", cpuID).Description)
}

func ExampleParseValueUnit() {
	s := "162 %/s"

	v, u := sapcontrol.ParseValueUnit(s)
	if u == s {
		fmt.Println("[ERROR] could not parse")
	}

	fmt.Printf("%d - %s\n", v, u)
	// Output:
	// 162 - %/s
}
