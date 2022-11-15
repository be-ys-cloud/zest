package schedulers

import (
	"context"
	"github.com/be-ys-cloud/zest/internal/services"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"strconv"
)

// Cleaner starts a cleaning only if < 10% of FS is available
func Cleaner(_ context.Context) {
	var stat unix.Statfs_t

	err := unix.Statfs(configuration.DataStorage, &stat)

	if err != nil {
		logrus.Warnln("Could not stat directory to get disk usage !")
		return
	}

	minSpace := 10
	if d, err := strconv.Atoi(configuration.FreeSpaceThreshold); err == nil {
		minSpace = d
	}

	if (float64(stat.Bavail*uint64(stat.Bsize))/float64(stat.Blocks*uint64(stat.Bsize)))*100 < float64(minSpace) {
		logrus.Warnln("Free space on FS is low ! Now cleaning data to recover some space...")
		err = services.Cleanup()
		if err != nil {
			logrus.Warnln("Error while cleaning ! ", err.Error())
		}
	}
}
