package torrent

import (
	"github.com/hekmon/transmissionrpc"
)

func (client Client) DeleteTorrentsAndData(IDs ...int64) error {
	var rmConfig = &transmissionrpc.TorrentRemovePayload{
		DeleteLocalData: true,
		IDs:             IDs,
	}
	return client.tr.TorrentRemove(rmConfig)
}

func (client Client) DeleteTorrents(IDs ...int64) error {
	var rmConfig = &transmissionrpc.TorrentRemovePayload{
		DeleteLocalData: false,
		IDs:             IDs,
	}
	return client.tr.TorrentRemove(rmConfig)
}
