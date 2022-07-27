package playlist

import "errors"

const DefaultPlaylistMaxSize = 10

type Playlist struct {
	queue    []string
	currSize int
	maxSize  int
}

func New(maxSize int) *Playlist {
	currSize := 0

	return &Playlist{
		queue:    make([]string, currSize, maxSize),
		currSize: currSize,
		maxSize:  maxSize,
	}
}

func (pl *Playlist) Add(url string) error {
	if pl.currSize+1 > pl.maxSize {
		return errors.New("song queue is full")
	}
	pl.queue = append(pl.queue, url)
	pl.currSize += 1

	return nil
}

func (pl *Playlist) GetNext() (string, error) {
	if pl.currSize == 0 {
		return "", errors.New("there is no song in queue")
	}

	nextSong := pl.queue[0]
	pl.queue = pl.queue[:1]
	pl.currSize -= 1

	return nextSong, nil
}

func (pl *Playlist) Size() int {
	return pl.currSize
}
