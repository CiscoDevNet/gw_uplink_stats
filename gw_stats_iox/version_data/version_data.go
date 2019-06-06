package version_data

import (
	"../hdm_api"
	"regexp"

)

type version_data struct {
	DeviceID     string     `json:",omitempty"`
	HWVersion    string     `json:",omitempty"`
	FWVersion    string     `json:",omitempty"`
}

func cell_sn_ver(dt string) (string) {
	var re = regexp.MustCompile(`Processor\s+board\s+ID\s+(\S+\d+)`)
	var str = dt
	var Data string
	preData := len(re.FindAllStringSubmatch(str, -1))
	if (preData < 1){
		Data = ""
	}else {
		Data = re.FindAllStringSubmatch(str, -1)[0][1]
	}

	return Data

}

func cell_hw_ver(dt string) (string) {
	var re = regexp.MustCompile(`Cisco\s+(\S+\s+\S+\s+\S+)`)
	var str = dt
	var Data string
	preData := len(re.FindAllStringSubmatch(str, -1))
	if (preData < 1){
		Data = ""
	}else {
		Data = re.FindAllStringSubmatch(str, -1)[0][1]
	}

	return Data

}

func cell_fw_ver(dt string) (string) {
	var re = regexp.MustCompile(`Version\s+(\S+),`)
	var str = dt
	var Data string
	preData := len(re.FindAllStringSubmatch(str, -1))
	if (preData < 1){
		Data = ""
	}else {
		Data = re.FindAllStringSubmatch(str, -1)[0][1]
	}

	return Data

}

func Version_data() (*version_data){
	cmd_out0 := hdm_api.Show_cmd_ssh("show version | inc Processor")
	cmd_out1 := hdm_api.Show_cmd_ssh("show version | inc Version")
	cmd_out2 := hdm_api.Show_cmd_ssh("show version | inc revision")

	var DeviceID string
	var HWVersion string
	var FWVersion string

	if (cmd_out0 == "Error") {
		DeviceID = ""
	}else {
		DeviceID = cell_sn_ver(cmd_out0)
	}

	if (cmd_out1 == "Error") {
		FWVersion = ""
	}else {
		FWVersion = cell_fw_ver(cmd_out1)
	}

	if (cmd_out2 == "Error") {
		HWVersion = ""
	}else {
		HWVersion = cell_hw_ver(cmd_out2)
	}

	ver_data := version_data{}
	ver_data.DeviceID = DeviceID
	ver_data.HWVersion = HWVersion
	ver_data.FWVersion = FWVersion

	return &ver_data
}