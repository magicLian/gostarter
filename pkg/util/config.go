package util

import (
	"os"

	"gopkg.in/ini.v1"
)

func GetEnvOrIniValue(iniFile *ini.File, section string, key string) string {
	osname := ""
	osnameBeta := ""
	if section != "" {
		osname = (section + "." + key)
		osnameBeta = (section + "_" + key)
	} else {
		osname = key
	}

	if os.Getenv(osname) != "" {
		return os.Getenv(osname)
	} else if os.Getenv(osnameBeta) != "" {
		return os.Getenv(osnameBeta)
	} else {
		return iniFile.Section(section).Key(key).MustString("")
	}
}
