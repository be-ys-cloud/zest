package structures

import "time"

type File struct {
	Id             int       `json:"id" sql:"id"`
	FileName       string    `json:"fileName" sql:"fileName"`
	LastDownloaded time.Time `json:"lastDownloaded" sql:"lastDownloaded"`
	NbDownload     int       `json:"nbDownload" sql:"nbDownload"`
	FileSize       int       `json:"fileSize" sql:"fileSize"`
	Version        string
}

type FileLocation struct {
	Id          int    `json:"id" sql:"id"`
	FileId      int    `json:"fileId" sql:"fileId"`
	FilePackage string `json:"filePackage" sql:"filePackage"`
}
