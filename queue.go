package main

import "errors"

type SongQueue struct {
	queue    []string
	currSize int
	maxSize  int
}

func NewSongQueue(maxSize int) *SongQueue {
	currSize := 0

	return &SongQueue{
		queue:    make([]string, currSize, maxSize),
		currSize: currSize,
		maxSize:  maxSize,
	}
}

func (sq *SongQueue) Add(url string) error {
	if sq.currSize+1 > sq.maxSize {
		return errors.New("song queue is full")
	}
	sq.queue = append(sq.queue, url)
	sq.currSize += 1

	return nil
}

func (sq *SongQueue) GetNext() (string, error) {
	if sq.currSize == 0 {
		return "", errors.New("there is no song in queue")
	}

	nextSong := sq.queue[0]
	sq.queue = sq.queue[:1]
	sq.currSize -= 1

	return nextSong, nil
}

func (sq *SongQueue) Size() int {
	return sq.currSize
}
