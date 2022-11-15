package routines

import (
	"github.com/be-ys-cloud/zest/internal/database"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetFileIndex(fileName string) {
	for configuration.UnpackRoutines.RunningCount() != 0 {
		time.Sleep(2 * time.Second)
	}

	configuration.Routines.Wait()

	fileParts := strings.Split(fileName, "/")
	count := 1
	folderName2 := fileName

	// From the most specific to the most generic path : trying to find the dists/ folder.
	for count <= len(fileParts) {
		folderName2 = strings.TrimSuffix(folderName2, "/"+fileParts[len(fileParts)-count])
		count += 1
		if _, err := os.Stat(configuration.DataStorage + folderName2 + "/dists"); err == nil {
			break
		}
	}
	if count != len(fileParts) {
		poolFileName := strings.TrimPrefix(fileName, folderName2+"/")

		var releaseList []string

		_ = filepath.Walk(configuration.DataStorage+folderName2+"/dists", func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}

			if err != nil {
				logrus.Warnln("Unable to walk on " + configuration.DataStorage + folderName2 + "/dists")
				logrus.Warnln("Error was : " + err.Error())
				return err
			}

			if !info.IsDir() && (info.Name() == "Packages" || info.Name() == "Sources") {
				var data []byte
				if info.Name() == "Packages" {
					data, err = exec.Command("bash", "-c", "grep -zoP '([^\\n]+\\n)*Filename: "+regexp.QuoteMeta(poolFileName)+"\\n([^\\n]+[\\n]{1})*' "+path).Output()
				} else {
					data, err = exec.Command("bash", "-c", "grep -zoP '([^\\n]+\\n)*(.*)"+regexp.QuoteMeta(fileParts[len(fileParts)-1])+"\\n([^\\n]+[\\n]{1})*' "+path).Output()
				}

				if err == nil {
					releaseList = append(releaseList, strings.TrimPrefix(path, configuration.DataStorage))
					err = os.MkdirAll(path+"_partials/", 0777)
					if err != nil {
						logrus.Warnln("Unable to create folder " + path + "_partials/")
						logrus.Warnln("Error was : " + err.Error())
						return err
					}

					err = os.WriteFile(path+"_partials/"+fileName[strings.LastIndex(fileName, "/")+1:]+"_index", data[:len(data)-1], 0777)
					if err != nil {
						logrus.Warnln("Unable to write file " + path + "_partials/" + fileName[strings.LastIndex(fileName, "/")+1:] + "_index")
						logrus.Warnln("Error was : " + err.Error())
						return err
					}
				}
			}
			return nil
		})

		_ = database.InsertFile(fileName, releaseList)

	} else {
		logrus.Warnln("Unable to find folder for " + fileName)
	}

	configuration.Routines.Done()
}
