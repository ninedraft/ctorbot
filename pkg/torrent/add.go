package torrent

import (
	"github.com/hekmon/transmissionrpc"
)

func (client Client) AddMagnet(magnet string) (Torrent, error) {
	var addConfig = &transmissionrpc.TorrentAddPayload{
		Filename:    &magnet,
		DownloadDir: &client.config.DownloadDir,
	}
	var torrent, err = client.tr.TorrentAdd(addConfig)
	if err != nil {
		return Torrent{}, err
	}
	return TorrentFromTransmission(torrent), nil
}
