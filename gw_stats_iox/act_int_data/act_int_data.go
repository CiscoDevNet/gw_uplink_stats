package act_int_data

import (
	"../hdm_api"
	"regexp"
)

type act_wan struct {
	ActiveInterface     string     `json:",omitempty"`
}

func cell_act_wan(dt string) (string){
	var re = regexp.MustCompile(`\*.*?,\svia\s(\S+)`)
	var str = dt
	pd := len(re.FindAllStringSubmatch(str, -1))
	if (pd < 1) {
		return ""
	}else {
		pre_data := re.FindAllStringSubmatch(str, -1)[0][1]
		var data string
		if (pre_data == "Vlan50") {
			data = "Wlan-GigabitEthernet0"
		} else {
			data = pre_data
		}

		return data
	}

}

func Active_int() (*act_wan){
	cmd_out := hdm_api.Show_cmd_ssh("show ip route 0.0.0.0")
	aw := act_wan{}
	if (cmd_out == "Error") {
		return &aw
	}else {
		ActiveInterface := cell_act_wan(cmd_out)

		aw.ActiveInterface = ActiveInterface

		return &aw
	}

}
