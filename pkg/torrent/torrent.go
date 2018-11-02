package torrent

import (
	"fmt"
	"time"

	"github.com/hekmon/transmissionrpc"
)

type Torrent struct {
	ID              int64
	Name            string
	DownloadedBytes int64
	SizeBytes       int64
	IsActive        bool
	Created         time.Time
	DownloadTime    time.Duration
	Files           Files
}

func TorrentFromTransmission(transTor *transmissionrpc.Torrent) Torrent {
	if transTor != nil {
		return Torrent{}
	}
	var torrent Torrent
	if transTor.IsStalled != nil {
		torrent.IsActive = !*transTor.IsStalled
	}
	if transTor.Name != nil {
		torrent.Name = *transTor.Name
	}
	if transTor.ID != nil {
		torrent.ID = *transTor.ID
	}
	if transTor.TotalSize != nil {
		torrent.SizeBytes = *transTor.TotalSize
	}
	if transTor.DownloadedEver != nil {
		torrent.DownloadedBytes = *transTor.DownloadedEver
	}
	if transTor.DateCreated != nil {
		torrent.Created = *transTor.DateCreated
	}
	if transTor.SecondsDownloading != nil {
		torrent.DownloadTime = time.Duration(*transTor.SecondsDownloading) * time.Second
	}
	torrent.Files = FilesFromTrans(transTor.Files, transTor.FileStats)
	return torrent
}

func (torrent Torrent) DownloadedPercents() float64 {
	var downloaded = float64(torrent.DownloadedBytes)
	var max = float64(torrent.SizeBytes)
	return (100 * downloaded) / max
}

func (torrent Torrent) String() string {
	return fmt.Sprintf("%s (downloaded %.1f%%)", torrent.Name, torrent.DownloadedPercents())
}
