package routines

import (
	"github.com/be-ys-cloud/zest/internal/utils/compress"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/sirupsen/logrus"
	"strings"
)

func DeflateFile(fileName string, format string) {
	configuration.UnpackRoutines.Wait()

	if err := compress.Deflate(configuration.DataStorage+fileName, configuration.DataStorage+strings.TrimSuffix(fileName, "."+format)); err != nil {
		logrus.Warnln("Unable to deflate " + fileName)
	}

	configuration.UnpackRoutines.Done()
}
