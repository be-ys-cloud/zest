package services

import (
	"github.com/be-ys-cloud/zest/internal/database"
	"github.com/be-ys-cloud/zest/internal/structures"
	"github.com/be-ys-cloud/zest/internal/utils/configuration"
	"strconv"
	"strings"
)

func GetMetrics() (string, error) {

	stats, err := database.GetStats()
	if err != nil {
		return "", err
	}

	metrics := []structures.Metric{
		{
			Name:  "zest_deflate_goroutines",
			Help:  "indicates the number of running goroutines that are deflating compressed archives",
			Value: strconv.FormatInt(int64(configuration.UnpackRoutines.RunningCount()), 10),
		},
		{
			Name:  "zest_debfile_goroutines",
			Help:  "indicates the number of running goroutines that are searching for debfiles in indexes",
			Value: strconv.FormatInt(int64(configuration.Routines.RunningCount()), 10),
		},
		{
			Name:  "zest_max_goroutines",
			Help:  "indicates the maximum number of concurrent routines",
			Value: configuration.MaxRoutines,
		},
		{
			Name:  "zest_saved_bytes",
			Help:  "indicates the number of bytes served from mirror since first boot",
			Value: strconv.Itoa(stats.SavedBytes),
		},
		{
			Name:  "zest_file_number",
			Help:  "indicates the number of files in cache",
			Value: strconv.Itoa(stats.NbFiles),
		},
		{
			Name:  "zest_total_size_on_disk",
			Help:  "indicates the total size used on disk",
			Value: strconv.Itoa(stats.TotalSizeOnDisk),
		},
	}

	result := strings.Builder{}

	for _, metric := range metrics {
		result.Write([]byte("# HELP " + metric.Name + " " + metric.Help + "\n"))
		result.Write([]byte("# TYPE " + metric.Name + " gauge\n"))
		result.Write([]byte(metric.Name + "{} " + metric.Value + "\n"))
	}

	return result.String(), nil

}
