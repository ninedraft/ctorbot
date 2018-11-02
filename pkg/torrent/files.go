package torrent

import (
	"misc/ctorbot/pkg/utils/it"

	"github.com/hekmon/transmissionrpc"
)

type Files []File

func FilesFromIter(n int, factory func(i int) File) Files {
	var files = make(Files, 0, n)
	for i := range it.It(n) {
		files = append(files, factory(i))
	}
	return files
}

func FilesFromTrans(transFiles []*transmissionrpc.TorrentFile, transStats []*transmissionrpc.TorrentFileStat) Files {
	return FilesFromIter(len(transFiles), func(i int) File {
		var transFile = transFiles[i]
		var stat = transStats[i]
		return FileFromTrans(transFile, stat)
	})
}

func (files Files) Len() int {
	return len(files)
}

func (files Files) New() Files {
	return make(Files, 0, files.Len())
}

func (files Files) Copy() Files {
	return append(files.New(), files...)
}

func (files Files) Filter(pred func(file File) bool) Files {
	var filtered = make(Files, 0, files.Len())
	for _, file := range files {
		if pred(file) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func (files Files) Completed() Files {
	return files.Filter(func(file File) bool {
		return file.IsCompleted()
	})
}

func (files Files) Names() []string {
	var names = make([]string, 0, files.Len())
	for _, file := range files {
		names = append(names, file.Name)
	}
	return names
}

func (files Files) Strings() []string {
	var strs = make([]string, 0, files.Len())
	for _, file := range files {
		strs = append(strs, file.String())
	}
	return strs
}
