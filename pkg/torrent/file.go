package torrent

import (
	"fmt"

	"github.com/hekmon/transmissionrpc"
)

type File struct {
	BytesCompleted int64  `json:"bytesCompleted"`
	Length         int64  `json:"length"`
	Name           string `json:"name"`
	Priority       int64  `json:"priority"`
}

func FileFromTrans(trans *transmissionrpc.TorrentFile, stat *transmissionrpc.TorrentFileStat) File {
	if trans == nil {
		return File{}
	}
	var file = File{
		Name:           trans.Name,
		BytesCompleted: trans.BytesCompleted,
		Length:         trans.Length,
	}
	if stat != nil {
		file.Priority = stat.Priority
	}
	return file
}

func (file File) String() string {
	var completed = float64(file.BytesCompleted)
	var length = float64(file.Length)
	return fmt.Sprintf("%s (downloaded %.1f%%)", file.Name, 100*completed/length)
}

func (file File) IsCompleted() bool { return file.BytesCompleted == file.Length }
