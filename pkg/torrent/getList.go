package torrent

func (client Client) GetTorrentList() (Torrents, error) {
	var transTorrents, errGetTorrentList = client.tr.TorrentGetAll()
	if errGetTorrentList != nil {
		return nil, errGetTorrentList
	}
	var torrents = TorrentsFromIter(len(transTorrents),
		func(i int) Torrent {
			return TorrentFromTransmission(transTorrents[i])
		})
	return torrents, nil
}
