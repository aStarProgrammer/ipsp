// IPSP project main.go
package main

import (
	"ipsp/Configuration"
	"ipsp/Utils"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ReadSiteLastGuid(sitePath string) (string, error) {
	currentPath, errCurrentPath := Configuration.GetCurrentPath()

	if errCurrentPath != nil {
		return "", errors.New("current Path is empty")
	}

	var dataFilePath = filepath.Join(currentPath, "data")

	if Utils.PathIsExist(dataFilePath) == false {
		return "", errors.New("Data File Path " + dataFilePath + " not exist")
	}

	data, errReadFile := ioutil.ReadFile(dataFilePath)

	if errReadFile != nil {
		return "", errors.New("Data Read Fail ")
	}

	var sData = string(data)
	var siteLastGuids = strings.Split(sData, "|")

	if len(siteLastGuids) > 0 {
		for _, lastGuid := range siteLastGuids {
			if strings.Contains(lastGuid, ";") == true {
				var siteSettings = strings.Split(lastGuid, ";")
				if len(siteSettings) == 2 {
					var tempSitePath = siteSettings[0]
					var guid = siteSettings[1]

					if sitePath == tempSitePath {
						return guid, nil
					}
				}
			}
		}
	}

	return "", nil
}

func SaveSiteLastGuid(sitePath, oldGuid, guid string) (bool, error) {
	currentPath, errCurrentPath := Configuration.GetCurrentPath()

	if errCurrentPath != nil {
		return false, errors.New("current Path is empty")
	}

	var dataFilePath = filepath.Join(currentPath, "data")

	if Utils.PathIsExist(dataFilePath) == false {
		return false, errors.New("Data File Path " + dataFilePath + " not exist")
	}

	data, errReadFile := ioutil.ReadFile(dataFilePath)

	if errReadFile != nil {
		return false, errors.New("Data Read Fail ")
	}

	var sData = string(data)

	var oldInfo = sitePath + ";" + oldGuid
	var newInfo = sitePath + ";" + guid

	if oldGuid == "" {
		sData = sData + newInfo
	} else {
		sData = strings.Replace(sData, oldInfo, newInfo, -1)
	}
	//create，文件存在则会覆盖原始内容（其实就相当于清空），不存在则创建
	fp, error := os.Create(dataFilePath)
	if error != nil {
		return false, error
	}
	//延迟调用，关闭文件
	defer fp.Close()

	_, errWriteFile := fp.WriteString(sData)

	if errWriteFile != nil {
		return false, errors.New("Write data File Fail")
	}

	return true, nil
}

func Monitor(sitePath string) {
	fmt.Println("---------------------")
	fmt.Println(Utils.CurrentTime())
	lastGuid, errLastGuid := ReadSiteLastGuid(sitePath)

	if errLastGuid != nil {
		fmt.Println("Cannot read last Guid")
		return
	}

	newGuid, errPublishSite := PublishSite(sitePath)

	if errPublishSite != nil {
		fmt.Println("Cannot Publish site")
		return
	}

	if newGuid != lastGuid {
		//Site Changed, publish ipns again
		fmt.Println("Site Content changed, will publish it to ipns again ,new id " + newGuid)
		_, errPublishIPNS := PublishSite2IPNS(newGuid)

		if errPublishIPNS != nil {
			fmt.Println("Cannot publish 2 ipns")
			return
		}

		_, errSave := SaveSiteLastGuid(sitePath, lastGuid, newGuid)

		if errSave != nil {
			fmt.Println("Cannot save new Guid")
			return
		}
	}
	fmt.Println("---------------------")
}

func PublishSite(sitePath string) (string, error) {
	currentPath, errCurrentPath := Configuration.GetCurrentPath()
	if errCurrentPath != nil {
		return "", errors.New("Cannot get current path")
	}
	var outputFilePath = filepath.Join(currentPath, "output.txt")
	var strIPFSPublishCmd = `ipfs add -r ` + sitePath + ` > ` + outputFilePath

	//fmt.Println(strIPFSPublishCmd)
	sysType := runtime.GOOS

	var publishCmd *exec.Cmd

	if "windows" == sysType {
		publishCmd = exec.Command("cmd", "/c", strIPFSPublishCmd)
	} else if "linux" == sysType || "darwin" == sysType {
		publishCmd = exec.Command("bash", "-c", strIPFSPublishCmd)
	} else { //Not support other platforms now
		var errMsg string
		errMsg = "Complie Markdown, not supported platform " + sysType
		return "", errors.New(errMsg)
	}

	var out bytes.Buffer
	var error bytes.Buffer
	publishCmd.Stdout = &out
	publishCmd.Stderr = &error

	fmt.Println("Publishing site " + sitePath + " to ipfs ")
	err := publishCmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + error.String())
	}

	outputFileContent, errReadOutputFile := ioutil.ReadFile(outputFilePath)
	if errReadOutputFile != nil {
		return "", errors.New("Cannot read ipsc publish output File")
	}

	var sOutputFile = string(outputFileContent)
	var sTemp = strings.Split(sOutputFile, " ")
	if len(sTemp) == 0 {
		return "", errors.New("Output File content error")
	}

	var index = len(sTemp) - 2

	if index < 0 {
		return "", errors.New("Output File Content error, not enought data")
	}

	var ret = sTemp[len(sTemp)-2]
	fmt.Println("Publishing site done , QmID of site folder is " + ret)
	return ret, nil
}

func PublishSite2IPNS(siteGuid string) (bool, error) {
	currentPath, errCurrentPath := Configuration.GetCurrentPath()
	if errCurrentPath != nil {
		return false, errors.New("Cannot get current path")
	}

	var outputFilePath = filepath.Join(currentPath, "outputIPNS.txt")

	var strIPFSPublishCmd = `ipfs name publish ` + siteGuid + ` > ` + outputFilePath

	//fmt.Println(strIPFSPublishCmd)
	sysType := runtime.GOOS

	var publishCmd *exec.Cmd

	if "windows" == sysType {
		publishCmd = exec.Command("cmd", "/c", strIPFSPublishCmd)
	} else if "linux" == sysType || "darwin" == sysType {
		publishCmd = exec.Command("bash", "-c", strIPFSPublishCmd)
	} else { //Not support other platforms now
		var errMsg string
		errMsg = "Complie Markdown, not supported platform " + sysType
		return false, errors.New(errMsg)
	}

	fmt.Println("Publishing to ipns, it will take a while, please wait ...")
	_, errPandoc := publishCmd.Output()
	if errPandoc != nil {
		fmt.Println(errPandoc.Error())
		return false, errPandoc
	}

	outputFileContent, errReadOutputFile := ioutil.ReadFile(outputFilePath)
	if errReadOutputFile != nil {
		return false, errors.New("Cannot read ipsc name publish output File")
	}

	var sOutputFile = string(outputFileContent)
	var checkString2 = "Published to "

	if strings.Contains(sOutputFile, checkString2) {
		fmt.Println("Publish done , result " + sOutputFile)
		return true, nil
	}

	return false, errors.New("Name Publish failed")

}

func Run() {
	var cp CommandParser
	bParse := cp.ParseCommand()
	if bParse == true {
		var ch chan int
		//定时任务
		fmt.Println(Utils.CurrentTime())
		fmt.Println("Start Monitor, Wait for " + strconv.FormatInt(cp.MonitorInternal, 10) + " seconds to start the first monitor")
		ticker := time.NewTicker(time.Second * time.Duration(cp.MonitorInternal))
		go func() {
			for range ticker.C {
				Monitor(cp.SiteFolderPath)
			}
			ch <- 1
		}()
		<-ch
	}
}

func checkIFSPFolder() bool {
	currentPath, errCurrentPath := Configuration.GetCurrentPath()
	if errCurrentPath != nil {
		fmt.Println("Cannot get path of ipfs executable file ")
		return false
	}

	if strings.Contains(currentPath, " ") {
		fmt.Println("IPSP cannot put in a folder whose path has spaces")
		return false
	}

	return true
}

func main() {
	if checkIFSPFolder() == false {
		return
	}

	Run()
}
