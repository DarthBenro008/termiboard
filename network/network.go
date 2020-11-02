package network

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"termiboard/utils"
)

func GetPublicIPAddress() {
	URL := "https://api.ipify.org" //Using Third Party Service to Ping
	resp, err := http.Get(URL)     //Get the JSON Response
	if err != nil {
		utils.StandardPrinter(utils.ErrorRedColor, "Could not get response from "+URL)
		utils.StandardPrinter(utils.WarningYellowColor, "Check your Internet Connection")
		utils.PrintVerbosePanic(err) //Exit upon error, below code must not be executed
	}
	defer resp.Body.Close()              //The client must close the response body when finished with it
	IP, err := ioutil.ReadAll(resp.Body) //Reading all data from a io.Reader until EOF
	if err != nil {
		utils.StandardPrinter(utils.ErrorRedColor, "Could not find IPv4 Address")
		utils.PrintVerbosePanic(err) //Exit upon error, below code must not be executed
	}
	utils.ResultPrinter("Public IPv4 Address: ", string(IP)) //Cast []bByte to String and return
}

func GetLocalIPAddress() {
	networkInterfaces, err := net.Interfaces() //Get the name of all Network Interface Cards
	if err != nil {                            //Error Handling
		utils.StandardPrinter(utils.ErrorRedColor, "Unfortunately it is not possible to get your local IP")
		utils.PrintVerbosePanic(err) //Exit upon error, below code must not be executed
	}
	utils.ResultPrinter("Local IP Address:", "")
	for _, networkInterface := range networkInterfaces { //Iterate over each Network Interface Card
		utils.StandardPrinter(utils.BoldWhite, "\t["+string(networkInterface.Name)+"] "+networkInterface.HardwareAddr.String()) // Print Name and Hardware Address
		interfaceAddresses, err := networkInterface.Addrs()                                                                     //Get the Internet Addresses of a Network Interface
		if err == nil {
			for _, interfaceAddress := range interfaceAddresses {
				if strings.Contains(interfaceAddress.String(), ".") { //Check for IPv4 Address
					utils.StandardPrinter(utils.None, "\t\t"+"IPv4"+" -> "+interfaceAddress.String())
				} else if strings.Contains(interfaceAddress.String(), ":") { //Check for IPv6 Address
					utils.StandardPrinter(utils.None, "\t\t"+"IPv6"+" -> "+interfaceAddress.String())
				}
			}
		}
	}
}

// Returns a list of unavailable ports from the received list.
func PortChecker(ports []string) []string {
	var unavailable []string
	for _, port := range ports {
		if isPortUnavailable(port) {
			unavailable = append(unavailable, port)
		}
	}

	return unavailable
}

func isPortUnavailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if ln != nil {
		defer ln.Close()
	}

	return err != nil
}
