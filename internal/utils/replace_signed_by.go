package utils

import (
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// ReplaceValidUntilAndSignedBy Signed-By field should be ignored if Valid-Until is not present, so there is no problem to combine these two updates.
func ReplaceValidUntilAndSignedBy(filePath string) (err error) {

	//Replace Valid-Until
	validUntilData, err := exec.Command("bash", "-c", "grep -E '^Valid-Until: (.*)$' "+filePath).Output()
	if err != nil {
		return err
	}

	oldValidUntil := strings.ReplaceAll(strings.TrimPrefix(strings.ReplaceAll(regexp.QuoteMeta(string(validUntilData)), "/", "\\/"), " "), "\n", "")
	var cstSh, _ = time.LoadLocation("UTC")
	newValidUntil := strings.ReplaceAll(regexp.QuoteMeta("Valid-Until: "+time.Now().Add(time.Hour*24*7).In(cstSh).Format(time.RFC1123Z)), "/", "\\/")

	err = exec.Command("bash", "-c", "sed -i 's/"+oldValidUntil+"/"+newValidUntil+"/g' "+filePath).Run()
	if err != nil {
		return err
	}

	//Replace Signed-By
	signedData, err := exec.Command("bash", "-c", "grep -E '^Signed-By: (.*)$' "+filePath).Output()
	if err != nil {
		return err
	}

	oldSignedData := strings.ReplaceAll(strings.TrimPrefix(strings.ReplaceAll(regexp.QuoteMeta(string(signedData)), "/", "\\/"), " "), "\n", "")
	newSignedData := strings.ReplaceAll(regexp.QuoteMeta("Signed-By: "+configuration.KeyName), "/", "\\/")

	err = exec.Command("bash", "-c", "sed -i 's/"+oldSignedData+"/"+newSignedData+"/g' "+filePath).Run()
	if err != nil {
		return err
	}

	return
}
