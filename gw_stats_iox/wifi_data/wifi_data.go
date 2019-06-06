package wifi_data

import (
	"../hdm_api"
	"regexp"
	"fmt"
	"gopkg.in/ini.v1"

	"../iox_conf"
)

type Wifi_struct struct {
	ID                 string       `json:",omitempty"`
	Type               string       `json:",omitempty"`
	Connected          string       `json:",omitempty"`
	BSSID              string       `json:",omitempty"`
	SSID               string       `json:",omitempty"`
	SNR                string       `json:",omitempty"`
	Encryption         string       `json:",omitempty"`
	RSSI               string       `json:",omitempty"`
	Technology         string       `json:",omitempty"`
	Authentication     string       `json:",omitempty"`
	//TxRate             string
	//RxRate             string
	TXbytes            string       `json:",omitempty"`
	RXbytes            string       `json:",omitempty"`
}

type Wifi_dict_data struct {
	WifiInterface Wifi_struct
}



func wifi_data_find(str string, dt string) (string) {
	orig_data := str
	reg := dt
	re := regexp.MustCompile(reg)
	var Data string
	preData := len(re.FindAllStringSubmatch(orig_data, -1))
	if (preData < 1){
		Data = ""
	}else {
		Data = re.FindAllStringSubmatch(orig_data, -1)[0][1]
	}

	return Data
}

func wifi_ip_find(str string, dt string) (string) {
	orig_data := str
	reg := dt
	re := regexp.MustCompile(reg)
	var Data string
	preData := len(re.FindAllStringSubmatch(orig_data, -1))
	if (preData < 1){
		Data = ""
	}else {
		Data = re.FindAllStringSubmatch(orig_data, -1)[0][1]
	}

	return Data
}

func Wifi_all(dt string) (*Wifi_struct){
	// Cellular information variable definitions
	var ID string
	var Type string
	var Connected string
	var BSSID string
	var SSID string
	var SNR string
	var Encryption string
	var RSSI string
	var Technology string
	var Authentication string
	//var TxRate string
	//var RxRate string
	var TXbytes string
	var RXbytes string

	var str = dt

	ID = "Wlan-GigabitEthernet0"
	Type = "WiFi"

	//Data Cleaning for Cell Data

	conn := wifi_data_find(str, `State\s+:\s+(\S+)`)
	if (conn == "Assoc") {
		Connected = "Connected"
	} else {
		Connected = "Disconnected"
	}


	BSSID = wifi_data_find(str, `Address\s+: (\S+)`)

	SSID = wifi_data_find(str, `SSID\s+: (\S+)`)

	SNR = wifi_data_find(str, `Signal to Noise\s+: (\S+)`)

	Encryption = wifi_data_find(str, `Encryption\s+:\s+(\S+)`)

	Authentication = wifi_data_find(str, `Key Mgmt type\s+: (\S+ +\S+)`)

	RSSI = wifi_data_find(str, `Signal Strength\s+: (\S+)`)

	Technology = "802.11n"

	//TxRate = wifi_data_find(str, `Mobile\s+Network\s+Code\s+\(MNC\)\s+=\s+(\S+)`)

	//RxRate = wifi_data_find(str, `Current\s+Roaming\s+Status\s+=\s+(\S+)`)

	TXbytes = wifi_data_find(str, `Bytes\s+Output\s+:\s+(\d+)`)

	RXbytes = wifi_data_find(str, `Bytes\s+Input\s+:\s+(\d+)`)

	w_data := Wifi_struct{}

	// Save data to struct
	w_data.ID = ID
	w_data.Type = Type
	w_data.Connected = Connected
	w_data.BSSID = BSSID
	w_data.SSID = SSID
	w_data.SNR = SNR
	w_data.Encryption = Encryption
	w_data.RSSI = RSSI
	w_data.Technology = Technology
	w_data.Authentication = Authentication
	//w_data.RxRate = RxRate
	//w_data.TxRate = TxRate
	w_data.TXbytes = TXbytes
	w_data.RXbytes = RXbytes

	//fmt.Println(c_data)

	return &w_data


}

func Wifi_data() (*Wifi_dict_data){
	final_wifi_data := Wifi_dict_data{}

	cfg, err := iox_conf.Conf_file()

	if err != nil {
		cfg, err = ini.Load("package_config.ini")
		if err != nil {
			panic(err)
		}
	}

	ap_hostname := cfg.Section("ap_info").Key("host").String() + "." + cfg.Section("ir_router_info").Key("domain").String()
	ip_cmd := hdm_api.Show_cmd_ssh("sh cdp nei wlan-gigabitEthernet 0 detail | inc IP address")

	fmt.Println(ip_cmd)

	fmt.Println(ap_hostname)

	if (ip_cmd == "Error" ) {
		return &final_wifi_data
	}else {

		ip_out := wifi_ip_find(ip_cmd, `IP address: (\S+)`)

		fmt.Println(ip_out)

		cmd_out1 := hdm_api.Show_cmd_ssh_ap(ip_out, "show dot11 associations all-client interface Dot11Radio 1")
		//cmd_out2 := hdm_api.Show_cmd_ssh_ap(ip_out, "show controllers dot11Radio 1 radio-stats")
		cmd_out2 := hdm_api.Show_cmd_ssh_ap(ip_out, "show controllers dot11Radio 1 radio-stats")

		cmd_out0 := cmd_out1 + "\r\n\r\n" + cmd_out2

		if (cmd_out1 == "Error" ) {
			return &final_wifi_data
		}else {

			dt_out := *Wifi_all(cmd_out0)
			//fmt.Println(dt_out)

			final_wifi_data.WifiInterface = dt_out

			//dt_out := Cell_all(cmd_out)
			//fmt.Println(dt_out)
			return &final_wifi_data
		}
	}


}
