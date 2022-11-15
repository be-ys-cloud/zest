package services

import (
	"errors"
	"github.com/be-ys-cloud/zest/internal/database"
	"github.com/be-ys-cloud/zest/internal/routines"
	"github.com/be-ys-cloud/zest/internal/utils"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/be-ys-cloud/zest/internal/utils/helpers"
	cp "github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetData(url string) (dataToWrite []byte, err error) {
	fileFullPath := url
	proxyMode := false

	if strings.HasPrefix(url, "pass/") {
		proxyMode = true
		url = strings.Replace(url, "pass/", "", 1)
		fileFullPath = strings.Replace(url, "pass/", "", 1)
	}

	url = strings.Replace(strings.Replace(url, "http/", "http://", 1), "https/", "https://", 1)

	protocol := fileFullPath[0:strings.Index(fileFullPath, "/")]

	fileFullPath = strings.Replace(fileFullPath, "http/", "", 1)
	fileFullPath = strings.Replace(fileFullPath, "https/", "", 1)

	lastSlashIndex := strings.LastIndex(fileFullPath, "/")
	if lastSlashIndex == -1 {
		lastSlashIndex = len(fileFullPath)
	}

	fileName := fileFullPath[strings.LastIndex(fileFullPath, "/")+1:]
	folderName := fileFullPath[0:lastSlashIndex]

	for configuration.UpdateInProgress {
		time.Sleep(1 * time.Second)
	}

	if _, err := os.Stat(configuration.DataStorage + fileFullPath); (url != "https" && url != "http") &&
		errors.Is(err, os.ErrNotExist) || (proxyMode && !strings.HasSuffix(fileName, ".deb")) {

		code, body, _, err := helpers.WSProvider("GET", url, nil, nil)

		if err != nil {
			logrus.Warnln("Unable to GET " + url + ". Error was : " + err.Error())
			return nil, err
		}

		if code != 200 {
			logrus.Warnln("Unable to GET " + url + ". Server returned HTTP Code : " + strconv.Itoa(code))
			return nil, err
		}

		if fileName == "InRelease" {
			resignedBody, err := utils.TransformInReleaseFile(folderName, fileFullPath, body)
			if !proxyMode {
				body = resignedBody
			}

			if err != nil {
				logrus.Warnln("Unable to sign InRelease file " + fileFullPath + ". Error was : " + err.Error())
				return nil, err
			}
		}

		if fileName == "Release" {
			err = utils.TransformReleaseFile(folderName, fileFullPath, body)
			if err != nil {
				logrus.Warnln("Unable to sign Release file " + fileFullPath + ". Error was : " + err.Error())
				return nil, err
			}
		}

		if strings.Contains(string(body), "<html>") {
			return nil, errors.New("html")
		}

		dataToWrite = body

		// Store file to server
		if err = os.MkdirAll(configuration.DataStorage+folderName, 0777); err != nil {
			logrus.Warnln("Unable to create " + configuration.DataStorage + folderName)
			return dataToWrite, err
		}

		if err = os.WriteFile(configuration.DataStorage+fileFullPath, body, 0777); err != nil {
			logrus.Warnln("Unable to write " + configuration.DataStorage + fileFullPath)
			return dataToWrite, err
		}

		if b, _ := regexp.MatchString("(.*)/dists/(.*)/(i18n/(.*)|dep11/(.*)|Contents(.*).gz)", url); b {
			return dataToWrite, nil
		}

		if strings.Contains(url, "/by-hash/") {
			findingReleaseFile := strings.TrimSuffix(fileFullPath, fileName)
			found := false

			for findingReleaseFile != "" && !found {
				if _, err = os.Stat(configuration.DataStorage + findingReleaseFile + "/Release"); err != nil {
					folderTMpParts := strings.Split(findingReleaseFile, "/")
					findingReleaseFile = strings.TrimSuffix(findingReleaseFile, "/"+folderTMpParts[len(folderTMpParts)-1])
				} else {
					found = true
				}
			}

			if !found {
				logrus.Warnln("Unable to find Release file for file hash " + url)
			} else {
				hashLineInRelease, err := exec.Command("bash", "-c", "grep -E '^ "+fileName+" (.*)$' "+configuration.DataStorage+findingReleaseFile+"/Release").Output()
				if err != nil {
					logrus.Warnln("Unable to find hash " + fileName + " in " + configuration.DataStorage + findingReleaseFile + "/Release")
				} else {
					hashParts := strings.Split(strings.ReplaceAll(string(hashLineInRelease), "\n", ""), " ")
					if err = cp.Copy(configuration.DataStorage+fileFullPath, configuration.DataStorage+findingReleaseFile+"/"+hashParts[len(hashParts)-1]); err != nil {
						logrus.Warnln("Unable to copy " + configuration.DataStorage + fileFullPath + " to " + configuration.DataStorage + findingReleaseFile + "/" + hashParts[len(hashParts)-1])
					} else {
						fileNameTmp := strings.Split(hashParts[len(hashParts)-1], "/")
						fileName = fileNameTmp[len(fileNameTmp)-1]
						fileFullPath = findingReleaseFile + "/" + hashParts[len(hashParts)-1]
					}
				}
			}
		}

		// Deflate file if we retrieved an archive
		format := strings.Split(fileName, ".")
		if helpers.IncludesString(format[len(format)-1], configuration.SupportedExtensions) {
			go routines.DeflateFile(fileFullPath, format[len(format)-1])
		}

		// Find and save file index for futur uses
		if match, _ := regexp.MatchString("(.*).(deb|dsc)", fileName); match {
			go routines.GetFileIndex(fileFullPath)
		}
	} else {
		if fileFullPath == "https" || fileFullPath == "http" {
			fileFullPath = ""
		}

		if a, err := os.Stat(configuration.DataStorage + fileFullPath); err == nil && a.IsDir() {
			files, err := ioutil.ReadDir(configuration.DataStorage + fileFullPath)
			if err != nil {
				return nil, err
			}
			dataToWrite = []byte("<html><h3>Index of " + fileFullPath + "</h3><ul>")
			dataToWrite = append(dataToWrite, []byte("<li><a href=\"/"+protocol+"/"+folderName+"\">Parent Directory</a></li>")...)

			for _, file := range files {
				path := ""
				if file.IsDir() {
					path = "/"
				}
				dataToWrite = append(dataToWrite, []byte("<li><a href=\"/"+protocol+"/"+fileFullPath+"/"+file.Name()+"\">"+file.Name()+path+"</a></li>")...)
			}
			dataToWrite = append(dataToWrite, []byte("</ul></html>")...)
		} else {
			dataToWrite, _ = os.ReadFile(configuration.DataStorage + fileFullPath)
			_ = database.UpdateFile(fileFullPath)
		}
	}

	return dataToWrite, nil
}
