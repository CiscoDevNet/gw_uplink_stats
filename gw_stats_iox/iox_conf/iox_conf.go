package iox_conf

import (
	"gopkg.in/ini.v1"
)


func Conf_file() (*ini.File, error){
	conf, err := ini.Load("/data/package_config.ini")
	if err != nil {
		conf, err = ini.Load("package_config.ini")
		if err != nil {
			panic(err)
		}
	}

	return conf, err

}
