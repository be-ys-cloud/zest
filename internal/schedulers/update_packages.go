package schedulers

import (
	"context"
	"github.com/be-ys-cloud/zest/internal/utils"
	"github.com/be-ys-cloud/zest/internal/utils/compress"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/be-ys-cloud/zest/internal/utils/gpg"
	"github.com/be-ys-cloud/zest/internal/utils/helpers"
	cp "github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// UpdatePackagesFiles updates all repositories present in server, renewing Release and Packages files + injecting our old indexes to keep versions available.
func UpdatePackagesFiles(_ context.Context) {
	var err error

	if err = os.MkdirAll(configuration.TempStorage, 0777); err != nil {
		logrus.Warnln("Could not create temporary dir: " + configuration.TempStorage)
		return
	}

	_ = filepath.Walk(configuration.DataStorage, func(path string, info os.FileInfo, inputError error) (err error) {

		if info == nil || info.IsDir() || info.Name() != "Release" {
			return
		}

		logrus.Infoln("Updating " + path)

		if err = os.RemoveAll(configuration.TempStorage + "ongoing/"); err != nil {
			logrus.Warnln("Unable to remove ongoing folder. The process will continue, but you may find some inconsistencies.")
		}

		// We are in a folder containing InRelease ! Moving to a temp folder to do our stuff.
		if err = cp.Copy(strings.TrimSuffix(path, "Release"), configuration.TempStorage+"ongoing/"); err != nil {
			logrus.Warnln("Unable to copy Release file to working directory. Aborting.")
			return
		}

		url := strings.TrimPrefix(path, configuration.DataStorage)

		scheme := "https://"
		//Define weather we are HTTP or HTTPS
		_, _, _, err = helpers.WSProvider("GET", scheme+url, nil, nil)
		if err != nil {
			scheme = "http://"
		}

		url = scheme + url

		//Try to get Release file, as it is easier
		code, body, _, err := helpers.WSProvider("GET", url, nil, nil)
		if err != nil || code != 200 {
			code, body, _, err = helpers.WSProvider("GET", strings.TrimSuffix(url, "Release")+"InRelease", nil, nil)
			if err != nil || code != 200 {
				logrus.Warnln("Unable to get Release or InRelease file for " + path + ". Aborting update for this repository.")
				return
			} else {
				file := configuration.TempStorage + "ongoing/Release"

				if err = os.WriteFile(file, body, 0777); err != nil {
					logrus.Warnln("Unable to write Release file !")
					return
				}

				if err = exec.Command("sed", "-i", "-e", "1,3d", file).Run(); err != nil {
					logrus.Warnln("Could not create sed in file " + file)
					return
				}

				if err = exec.Command("bash", "-c", "truncate -s $(grep -b '^-----BEGIN PGP SIGNATURE-----' "+file+" | grep -o -E '^[0-9]*') "+file).Run(); err != nil {
					logrus.Warnln("Could not grep in file " + file)
					return
				}
			}
		} else {
			if err = os.WriteFile(configuration.TempStorage+"ongoing/Release", body, 0777); err != nil {
				logrus.Warnln("Unable to write Release file !")
				return
			}
		}

		hashes := map[string]map[string]string{}

		_ = filepath.Walk(configuration.TempStorage+"ongoing/", func(packagePath string, packageInfo os.FileInfo, packageError error) (err error) {

			if packageInfo == nil || packageInfo.IsDir() {
				return
			}

			if strings.Contains(packagePath, "/i18n/") || strings.Contains(packagePath, "/dep11/") || strings.Contains(packagePath, "/Contents/") {
				//Remove file and update from dist
				packageUrl := strings.TrimSuffix(url, "Release") + strings.TrimPrefix(packagePath, configuration.TempStorage+"ongoing/")
				code, body, _, err = helpers.WSProvider("GET", packageUrl, nil, nil)
				if err == nil && code == 200 {
					_ = os.WriteFile(packagePath, body, 0777)
				} else {
					if err != nil {
						logrus.Warnln("Could not GET " + packageUrl + ". Error was : " + err.Error())
					} else {
						logrus.Warnln("Could not GET " + packageUrl + ". Code was " + strconv.Itoa(code))
					}
				}
			}

			if packageInfo.Name() != "Packages" && packageInfo.Name() != "Sources" {
				return
			}

			// Remove old archives
			for _, k := range configuration.SupportedExtensions {
				if err = os.Remove(packagePath + "." + k); err != nil {
					logrus.Warnln("Unable to remove " + packagePath + "." + k)
				}
			}
			// Remove package file
			if err = moveFile(packagePath, packagePath+"_old"); err != nil {
				return
			}

			packageUrl := strings.TrimSuffix(url, "Release") + strings.TrimPrefix(packagePath, configuration.TempStorage+"ongoing/")

			//Try to update Packages files
			packagesFileUpdated := false
			code, body, _, err := helpers.WSProvider("GET", packageUrl, nil, nil)
			if err != nil || code != 200 {
				for _, k := range configuration.SupportedExtensions {
					code, body, _, err = helpers.WSProvider("GET", packageUrl+"."+k, nil, nil)
					if err == nil && code == 200 {
						if err = os.WriteFile(packagePath+"."+k, body, 0777); err != nil {
							logrus.Warnln("Error while writing file " + packagePath + "." + k)
							continue
						}
						if err = compress.Deflate(packagePath+"."+k, packagePath); err != nil {
							logrus.Warnln("Error while deflating file " + packagePath + "." + k)
							continue
						}
						if err = os.Remove(packagePath + "." + k); err != nil {
							logrus.Warnln("Unable to remove " + packagePath + "." + k + ". Update will continue.")
						}

						packagesFileUpdated = true
						break
					}
				}
			} else {
				if err = os.WriteFile(packagePath, body, 0777); err != nil {
					logrus.Warnln("Error while writing file " + packagePath)
				} else {
					packagesFileUpdated = true
				}
			}

			if !packagesFileUpdated {
				logrus.Warnln("Unable to update " + packagePath + ". Rollback in progress.")
				if err = moveFile(packagePath+"_old", packagePath); err != nil {
					logrus.Warnln("Failed to rollback " + packagePath)
				}
				return
			}

			// Remove old file
			if err = os.Remove(packagePath + "_old"); err != nil {
				logrus.Warnln("Unable to remove old index file. Update will continue.")
			}

			//If partials, add all
			if _, err := os.Stat(packagePath + "_partials/"); err == nil {
				_ = filepath.Walk(packagePath+"_partials/", func(pathPartial string, infoPartial os.FileInfo, errPartial error) (err error) {
					if infoPartial == nil || infoPartial.IsDir() {
						return
					}

					if err = exec.Command("bash", "-c", "grep '/"+strings.TrimSuffix(infoPartial.Name(), "_index")+"' "+packagePath).Run(); err != nil {
						if err = exec.Command("bash", "-c", "cat "+pathPartial+" >> "+packagePath).Run(); err != nil {
							logrus.Warnln("Unable to add partial file " + pathPartial + " to package file " + packagePath)
						}
					}
					return
				})
			}

			cleanedString := strings.TrimSuffix(strings.TrimSuffix(packagePath, "Sources"), "Packages") + "by-hash"

			if err = os.RemoveAll(cleanedString); err != nil {
				logrus.Warnln("Error while removing path " + cleanedString)
			}

			if err = os.MkdirAll(cleanedString, 0777); err != nil {
				logrus.Warnln("Error while creating path " + cleanedString)
			}

			//Create archives and compute their hashes
			hashes[packagePath] = helpers.CalculateHash(packagePath)
			for key := range hashSize {
				if err = cp.Copy(packagePath, cleanedString+"/"+key+"/"+hashes[packagePath][key]); err != nil {
					logrus.Warnln("Could not copy " + packagePath + " to " + cleanedString + "/" + key + "/" + hashes[packagePath][key])
				}
			}

			for _, k := range configuration.SupportedExtensions {
				if err = compress.Inflate(packagePath, packagePath+"."+k); err != nil {
					logrus.Warnln("Error while creating archive file " + packagePath + "." + k)
				} else {
					hashes[packagePath+"."+k] = helpers.CalculateHash(packagePath + "." + k)
					for key := range hashSize {
						if err = cp.Copy(packagePath+"."+k, cleanedString+"/"+key+"/"+hashes[packagePath+"."+k][key]); err != nil {
							logrus.Warnln("Could not copy " + packagePath + "." + k + " to " + cleanedString + "/" + key + "/" + hashes[packagePath+"."+k][key])
						}
					}
				}
			}
			return
		})

		//Replace all hashes in Release file
		err = updateReleaseFile(hashes)
		if err != nil {
			logrus.Warnln("Unable to update release file !")
		}

		// Updating GPG signatures for Release and InRelease files
		if err = os.Remove(configuration.TempStorage + "ongoing/InRelease"); err != nil {
			logrus.Warnln("Unable to remove " + configuration.TempStorage + "ongoing/InRelease file.")
		}
		if err = os.Remove(configuration.TempStorage + "ongoing/Release.gpg"); err != nil {
			logrus.Warnln("Unable to remove " + configuration.TempStorage + "ongoing/Release.gpg file.")
		}

		//Change Signed-By field
		if err = utils.ReplaceValidUntilAndSignedBy(configuration.TempStorage + "ongoing/Release"); err != nil {
			logrus.Infoln("Could not change Signed-By field. Program will continue, but some clients may reject your file.")
		}

		//Generate detached signature
		if err = gpg.DetachedSign(configuration.KeyFile, configuration.KeyPass, configuration.TempStorage+"ongoing/Release", configuration.TempStorage+"ongoing/Release.gpg"); err != nil {
			logrus.Warnln("Unable to generate Release.gpg file !")
		}

		//Generate infile signature
		if err = cp.Copy(configuration.TempStorage+"ongoing/Release", configuration.TempStorage+"ongoing/InRelease"); err != nil {
			logrus.Warnln("Error while copying Release to InRelease file.")
		}
		if err = gpg.InlineSign(configuration.KeyFile, configuration.KeyPass, configuration.TempStorage+"ongoing/InRelease"); err != nil {
			logrus.Warnln("Error while signing InRelease file.")
		}

		//Operations done, moving folder to new data folder.
		if err = os.MkdirAll(strings.TrimSuffix(configuration.TempStorage+path, "Release"), 0777); err != nil {
			logrus.Warnln("Error while creating " + strings.TrimSuffix(configuration.TempStorage+path, "Release"))
		}

		if err = cp.Copy(configuration.TempStorage+"ongoing/", strings.TrimSuffix(configuration.TempStorage+path, "Release")); err != nil {
			logrus.Warnln("Error while copying data.")
		}

		_ = os.RemoveAll(configuration.TempStorage + "ongoing")

		logrus.Infoln("Update done for " + path)

		return
	})

	//Copying all new data to production data folder
	configuration.UpdateInProgress = true

	if err = exec.Command("bash", "-c", "cp -rf "+configuration.TempStorage+configuration.DataStorage+"*"+" "+configuration.DataStorage).Run(); err != nil {
		logrus.Warnln("Error while copying files from tmp to data folder" + err.Error())
	}

	configuration.UpdateInProgress = false

	// Remove tmp folder.
	err = os.RemoveAll(configuration.TempStorage + configuration.DataStorage)
	if err != nil {
		logrus.Warnln("Error while removing temporary folder." + err.Error())
	}

	logrus.Infoln("Update done for all repositories.")
}

// moveFile moves a file from source to dest.
func moveFile(source string, dest string) (err error) {
	err = cp.Copy(source, dest)
	if err != nil {
		return
	}

	err = os.Remove(source)
	if err != nil {
		return
	}
	return
}

func updateReleaseFile(hashes map[string]map[string]string) (err error) {
	data, err := os.ReadFile(configuration.TempStorage + "ongoing/Release")
	if err != nil {
		return
	}

	for fileName := range hashes {
		hashName := strings.ReplaceAll(regexp.QuoteMeta(strings.TrimPrefix(fileName, configuration.TempStorage+"ongoing/")), "\\", "")

		for hashType := range hashes[fileName] {
			if hashType != "Size" {
				regex := regexp.MustCompile(`(?m)^ ([a-z0-9A-Z]{` + strconv.Itoa(hashSize[hashType]) + `}) +([0-9]*) ` + regexp.QuoteMeta(hashName) + `$`)
				data = regex.ReplaceAll(data, []byte(" "+hashes[fileName][hashType]+" "+hashes[fileName]["Size"]+" "+hashName))
			}
		}
	}

	if err = os.WriteFile(configuration.TempStorage+"ongoing/Release", data, 0777); err != nil {
		logrus.Warnln("Unable to update Release file ; write failed.")
	}

	return
}
