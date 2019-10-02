package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/cjreeder/via_networking_script/via"
	"github.com/fatih/color"
)

type ViaList struct {
	vianame    string
	ipaddress  string
	subnetmask string
	gateway    string
	dns        string
}

func ReadCsv(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// read file into a variable to be able to usue later
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func SetNetwork(vianame string, ipaddress string, subnetmask string, gateway string, dns string) (string, error) {
	defer color.Unset()
	color.Set(color.FgYellow)

	address := vianame

	var command via.Command
	command.Command = "IpSetting"
	command.Param1 = ipaddress
	command.Param2 = subnetmask
	command.Param3 = gateway
	command.Param4 = dns
	command.Param5 = vianame

	fmt.Printf("Setting IP Info for %s", vianame)
	resp, err := via.SendCommand(command, address)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error in setting IP on %s", address))
	}
	return resp, nil
}

func main() {
	lines, err := ReadCsv("/home/creeder/Desktop/test_via_replacement_addresses.csv")
	if err != nil {
		panic(err)
	}

	// loop through the lines and turn it into an object
	for _, line := range lines {
		data := ViaList{
			vianame:    line[0],
			ipaddress:  line[1],
			subnetmask: line[2],
			gateway:    line[3],
			dns:        line[4],
		}
		fmt.Printf("Changing over %v", data.vianame)
		ret, err := SetNetwork(data.vianame, data.ipaddress, data.subnetmask, data.gateway, data.dns)
		if err != nil {
			fmt.Printf("%v returned an error: %v", data.vianame, err)
		}
		fmt.Printf("Chang over completed successfully with %v", ret)
	}
}
