package gps_data

import (
	"fmt"
	"../hdm_api"
	"../iox_conf"
	"regexp"
	"strconv"

)

type Gps struct {
	Latitude   string     `json:",omitempty"`
	Longitude  string     `json:",omitempty"`
	HDOP       string     `json:",omitempty"`
	Status     string     `json:",omitempty"`
}

type Gps_loc struct {
	Location Gps
}


func cell_data_find(str string, dt string) (string) {
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

func dms2dd(degrees string, minutes string, seconds string, direction string) (string){
	var dd float64
	var ddn string

	if (degrees == "" || minutes == "" || seconds == "" || direction == "") {
		ddn = ""
	}else {
		dd0, _ := strconv.ParseFloat(degrees, 64)
		dd1, _ := strconv.ParseFloat(minutes, 64)
		dd01 := dd1 / float64(60)
		dd2, _ := strconv.ParseFloat(seconds, 64)
		dd02 := dd2 / float64(60*60)
		dd = dd0 + dd01 + dd02
		if (direction == "W" || direction == "S") {
			dd *= -1
			ddn = strconv.FormatFloat(dd, 'f', -2, 64)
		} else {
			ddn = strconv.FormatFloat(dd, 'f', -2, 64)
		}
	}

	return ddn
}


func Cell_gps(dt string) (*Gps_loc) {
	// degreee variable definitions
	var lat_deg string
	var lat_min string
	var lat_sec string
	var lat_dir string

	var long_deg string
	var long_min string
	var long_sec string
	var long_dir string

	var str = dt

	// HDOP and Status Variable declarations

	var HDOP string
	var Status string

	//Data Cleaning for Location Data

	lat_deg = cell_data_find(str, `Latitude:\s+(\d+)\s+Deg`)

	lat_min = cell_data_find(str, `Latitude:\s+\d+\s+Deg\s+(\d+)\s+Min`)

	lat_sec = cell_data_find(str, `Latitude:\s+\d+\s+Deg\s+\d+\s+Min\s+(\S+)\s+Sec`)

	lat_dir = cell_data_find(str, `Latitude:\s+\d+\s+Deg\s+\d+\s+Min\s+\S+\s+Sec\s+(\S)`)

	long_deg = cell_data_find(str, `Longitude:\s+(\d+)\s+Deg`)

	long_min = cell_data_find(str, `Longitude:\s+\d+\s+Deg\s+(\d+)\s+Min`)

	long_sec = cell_data_find(str, `Longitude:\s+\d+\s+Deg\s+\d+\s+Min\s+(\S+)\s+Sec`)

	long_dir = cell_data_find(str, `Longitude:\s+\d+\s+Deg\s+\d+\s+Min\s+\S+\s+Sec\s+(\S)`)

	Latitude := dms2dd(lat_deg, lat_min, lat_sec, lat_dir)
	Longitude := dms2dd(long_deg, long_min, long_sec, long_dir)

	//Data Cleaning for GPS Status and HDOP
	//GPS\s+auto\s+tracking\s+status:\s+(\S+)
	HDOP = cell_data_find(str, `HDOP:\s+(\S+),`)

	Status = cell_data_find(str, `GPS\s+Mode\s+Used:\s+(\S+)`)

	gps_data := Gps{}

	if Status == "standalone" {
		gps_data.Status = "Valid"
	}else {
		gps_data.Status = "Invalid"
	}
	gps_data.HDOP = HDOP
	gps_data.Latitude = Latitude
	gps_data.Longitude = Longitude

	cell_gps_data := Gps_loc{}
	cell_gps_data.Location = gps_data

	return &cell_gps_data


}

func Gps_data() (*Gps_loc) {
	var cell_command string
    cfg, err := iox_conf.Conf_file()
    if err != nil {
		fmt.Println(err)
	}
	dual_lte := cfg.Section("dual_lte").Key("lte_enabled").String()
	if (dual_lte == "true") {
	    cell_command = "show cell 0/0 gps detail"
	} else {
	    cell_command = "show cell 0 gps detail"
	}
	cmd_out := hdm_api.Show_cmd_ssh(cell_command)
	gps_out := Gps_loc{}
	if (cmd_out == "Error") {
		return &gps_out
	}else {

		//fmt.Println(Cell_gps(cmd_out))

		dt_out := *Cell_gps(cmd_out)
		gps_out = dt_out

		return &gps_out
	}
}