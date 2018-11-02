package torrent

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
)

type ClosableBuf struct {
	bytes.Buffer
}

func (buf *ClosableBuf) Close() error { return nil }

type Loader func(name string) io.WriteCloser

func (client Client) LoadFiles(ID int64, loader Loader) error {
	var torrents, errGetList = client.GetTorrentList()
	if errGetList != nil {
		return errGetList
	}
	var torrent, ok = torrents.GetByID(ID)
	if !ok {
		return fmt.Errorf("unable to find torrent with ID=%d", ID)
	}
	for _, file := range torrent.Files.Completed() {
		var filePath = path.Join(client.config.DownloadDir, file.Name)
		var f, errOpen = os.Open(filePath)
		if errOpen != nil {
			return errOpen
		}
		var target = loader(file.Name)
		if _, err := io.Copy(target, f); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}
