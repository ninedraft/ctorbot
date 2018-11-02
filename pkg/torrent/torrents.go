package torrent

import (
	"misc/ctorbot/pkg/utils/it"
	"sort"
)

type Torrents []Torrent

func TorrentsFromIter(n int, factory func(i int) Torrent) Torrents {
	var torrents = make(Torrents, 0, n)
	for i := range it.It(n) {
		torrents = append(torrents, factory(i))
	}
	return torrents
}

func (torrents Torrents) Len() int {
	return len(torrents)
}

func (torrents Torrents) New() Torrents {
	return make(Torrents, 0, torrents.Len())
}

func (torrents Torrents) Copy() Torrents {
	return append(torrents.New(), torrents...)
}

func (torrents Torrents) SortByLess(less func(a, b Torrent) bool) Torrents {
	var sorted = torrents.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})
	return sorted
}

func (torrents Torrents) SortByKey(key func(torr Torrent) int) Torrents {
	return torrents.SortByLess(func(a, b Torrent) bool {
		return key(a) < key(b)
	})
}

func (torrents Torrents) SortByNames() Torrents {
	return torrents.SortByLess(func(a, b Torrent) bool {
		return a.Name < b.Name
	})
}

func (torrents Torrents) NewestFirst() Torrents {
	return torrents.SortByLess(func(a, b Torrent) bool {
		return a.Created.Before(b.Created)
	})
}

func (torrents Torrents) Filter(pred func(torr Torrent) bool) Torrents {
	var filtered = torrents.New()
	for _, torrent := range torrents {
		if pred(torrent) {
			filtered = append(filtered, torrent)
		}
	}
	return filtered
}

func (torrents Torrents) Map(op func(torr Torrent)) {
	for _, torrent := range torrents {
		op(torrent)
	}
}

func (torrents Torrents) Strings() []string {
	var strs = make([]string, 0, torrents.Len())
	torrents.Map(func(torrent Torrent) {
		strs = append(strs, torrent.String())
	})
	return strs
}

func (torrents Torrents) Names() []string {
	var names = make([]string, 0, torrents.Len())
	torrents.Map(func(torrent Torrent) {
		names = append(names, torrent.Name)
	})
	return names
}

func (torrents Torrents) Active() Torrents {
	return torrents.Filter(func(torrent Torrent) bool {
		return torrent.IsActive
	})
}

func (torrents Torrents) IDs() []int64 {
	var IDs = make([]int64, 0, torrents.Len())
	torrents.Map(func(torrent Torrent) {
		IDs = append(IDs, torrent.ID)
	})
	return IDs
}

func (torrents Torrents) Head() (Torrent, bool) {
	if torrents.Len() == 0 {
		return Torrent{}, false
	}
	return torrents[0], true
}

func (torrents Torrents) GetByID(ID int64) (Torrent, bool) {
	for _, torrent := range torrents {
		if torrent.ID == ID {
			return torrent, true
		}
	}
	return Torrent{}, false
}
