package main

import (
	"ipsp/Utils"
	"flag"
	"fmt"
	"os"
	"strings"
)

type CommandParser struct {
	SiteFolderPath  string
	MonitorInternal int64
}

func (cpp *CommandParser) ParseCommand() bool {
	//Set All Arguments
	flag.StringVar(&cpp.SiteFolderPath, "SiteFolder", "", "Site Folder that contains the site needs to publish")
	flag.Int64Var(&cpp.MonitorInternal, "MonitorInternal", 600, "Monitor Internal, second as unit, defualt 600 seconds")

	//Parse
	flag.Parse()

	//Trim all String properties
	cpp.SiteFolderPath = strings.TrimSpace(cpp.SiteFolderPath)

	if cpp.SiteFolderPath == "" || Utils.PathIsExist(cpp.SiteFolderPath) == false {
		fmt.Fprintln(os.Stderr, "Cannot find site")
		return false
	}

	if strings.Contains(cpp.SiteFolderPath, " ") {
		fmt.Fprintln(os.Stderr, "Site Folder Path Canot contains space")
		return false
	}

	if cpp.MonitorInternal < 0 {
		fmt.Println("Monitor Internal minus, will use 600 ")
		cpp.MonitorInternal = 600
	}

	return true
}
