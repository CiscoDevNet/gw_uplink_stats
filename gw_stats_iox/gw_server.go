package main

import (
	"./cell_data_zero"
	"./cell_data_one"
	"./gps_data"
	"./version_data"
	"./act_int_data"
	"./wifi_data"
	"./iox_conf"
	"fmt"
	"encoding/json"
	"time"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"bytes"
)

var all_data CiscoGW
var gw_json  string
var cell0 cell_data_zero.Cellular_data
var zero_val_cell0 cell_data_zero.Cellular_data
var cell1 cell_data_one.Cellular_data
var zero_val_cell1 cell_data_one.Cellular_data
var wifi_int wifi_data.Wifi_struct
var zero_val_wifi_int wifi_data.Wifi_struct

type CiscoGW struct {
	ActiveInterface     string        `json:",omitempty"`
	ID                  string        `json:",omitempty"`
	FWVersion           string        `json:",omitempty"`
	HWVersion           string        `json:",omitempty"`
	Interface           []interface{} `json:",omitempty"`
	Location            gps_data.Gps  `json:",omitempty"`
	Manufacturer        string        `json:",omitempty"`
	Model               string        `json:",omitempty"`

}



func JSONMarshalIndent(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "    ")

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

/*func gps() {
	for {
		gps := *gps_data.Gps_data()

		all_data.Location = gps.Location

		time.Sleep(500 * time.Millisecond)

	}

}*/

func gps() {
	gps := *gps_data.Gps_data()

	all_data.Location = gps.Location

	time.Sleep(500 * time.Millisecond)

}

/*func cellular0() {
	for {
		pre_cell0 := cell_data_zero.Cell_data().CellularInterface
		if pre_cell0.Connected != "Disconnected" {
			cell0 = cell_data_zero.Cellular_data(pre_cell0)
		}
		*//*pre_cell1 := cell_data_one.Cell_data().CellularInterface
		if pre_cell1.Connected != "Disconnected" {
			cell1 = cell_data_one.Cellular_data(pre_cell1)
		}*//*
		time.Sleep(500 * time.Millisecond)
	}
}*/

func cellular0() {
	pre_cell0 := cell_data_zero.Cell_data().CellularInterface
	if pre_cell0.ID != "" {
		cell0 = cell_data_zero.Cellular_data(pre_cell0)
	}
	/*pre_cell1 := cell_data_one.Cell_data().CellularInterface
	if pre_cell1.Connected != "Disconnected" {
		cell1 = cell_data_one.Cellular_data(pre_cell1)
	}*/
	time.Sleep(500 * time.Millisecond)
}

/*func cellular1() {
	for {
		pre_cell1 := cell_data_one.Cell_data().CellularInterface
		if pre_cell1.Connected != "Disconnected" {
			cell1 = cell_data_one.Cellular_data(pre_cell1)
		}
		time.Sleep(500 * time.Millisecond)
	}
}*/

func cellular1() {
	pre_cell1 := cell_data_one.Cell_data().CellularInterface
	if pre_cell1.ID != "" {
		cell1 = cell_data_one.Cellular_data(pre_cell1)
	}
	time.Sleep(500 * time.Millisecond)
}

func wifi() {
	for {
		pre_wifi := wifi_data.Wifi_data().WifiInterface
		if pre_wifi.ID != "" {
			wifi_int = wifi_data.Wifi_struct(pre_wifi)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

/*
func act_int() {
	for {
		act := act_int_data.Active_int()
		all_data.ActiveInterface = act.ActiveInterface
		time.Sleep(500 * time.Millisecond)
	}
}

*/
func act_int() {
	act := act_int_data.Active_int()
	all_data.ActiveInterface = act.ActiveInterface
	time.Sleep(500 * time.Millisecond)
}

func version() {
	vers := version_data.Version_data()


	all_data.Manufacturer = "Cisco"
	all_data.Model = "IR829"

	all_data.ID = vers.DeviceID
	all_data.FWVersion = vers.FWVersion
	all_data.HWVersion = vers.HWVersion
}

func interfaces() {
	for {
		if cell0 == zero_val_cell0 {
			if cell1 == zero_val_cell1 {
				if wifi_int == zero_val_wifi_int {
					all_data.Interface = append([]interface{}{})
				}else {
					all_data.Interface = append([]interface{}{wifi_int})
				}
			}else if wifi_int == zero_val_wifi_int {
				all_data.Interface = append([]interface{}{cell1})
			}else {
				all_data.Interface = append([]interface{}{cell1, wifi_int})
			}
		}else if cell1 == zero_val_cell1 {
			if cell0 == zero_val_cell0 {
				if wifi_int == zero_val_wifi_int {
					all_data.Interface = append([]interface{}{})
				}else {
					all_data.Interface = append([]interface{}{wifi_int})
				}
			}else if wifi_int == zero_val_wifi_int {
				all_data.Interface = append([]interface{}{cell0})
			}else {
				all_data.Interface = append([]interface{}{cell0, wifi_int})
			}
		}else if wifi_int == zero_val_wifi_int {
			if cell1 == zero_val_cell1 {
				if cell0 == zero_val_cell0 {
					all_data.Interface = append([]interface{}{})
				}else {
					all_data.Interface = append([]interface{}{cell0})
				}
			}else if cell0 == zero_val_cell0 {
				all_data.Interface = append([]interface{}{cell1})
			}else {
				all_data.Interface = append([]interface{}{cell0, cell1})
			}
		}else {
			all_data.Interface = append([]interface{}{cell0, cell1, wifi_int})
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func long_running() {
	for {
		act_int()
		cellular0()
		cellular1()
		gps()
	}
}

func json_proc() {
	for {
		read_data := all_data
		cisco_json, err := JSONMarshalIndent(read_data, true)
		if err != nil {
			fmt.Println(err)
			return
		}
		gw_json = string(cisco_json)
        }
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
    	fmt.Fprint(w, gw_json)
}

func main() {
    cfg, err := iox_conf.Conf_file()
    if err != nil {
		fmt.Println(err)
	}
	wgb_interface := cfg.Section("wgb_interface").Key("wgb_enabled").String()
	if (wgb_interface == "true") {
	    go wifi()
	}

	//go cellular0()
	//go cellular1()
	//go gps()
	go long_running()
	go interfaces()
	//go act_int()
	go version()
	go json_proc()


	router := httprouter.New()
    	router.GET("/api/gw_stats", Index)

	log.Fatal(http.ListenAndServe(":8080", router))

}
