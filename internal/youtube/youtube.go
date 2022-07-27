package youtube

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
)

const DefaultIntervalTimeout = 500 * time.Millisecond

type YoutubeSession struct {
	client          youtube.Client
	encodingSession *dca.EncodeSession
	intervalTimeout time.Duration
}

func NewSession(intervalTimeout time.Duration) *YoutubeSession {
	return &YoutubeSession{
		client:          youtube.Client{},
		intervalTimeout: intervalTimeout,
	}
}

func (yt *YoutubeSession) StreamYoutubeVideo(vc *discordgo.VoiceConnection, url string) error {
	// Change these accordingly
	options := dca.StdEncodeOptions
	options.BufferedFrames = 100
	options.FrameDuration = 20
	options.CompressionLevel = 5
	options.Bitrate = 96
	options.Application = "lowdelay"

	video, err := yt.client.GetVideo(url)
	if err != nil {
		return err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := yt.client.GetStream(video, &formats[0])
	if err != nil {
		return err
	}

	encodingSession, err := dca.EncodeMem(stream, options)
	if err != nil {
		return err
	}
	defer encodingSession.Cleanup()

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(yt.intervalTimeout)
	yt.encodingSession = encodingSession
	vc.Speaking(true)

	defer func(vc *discordgo.VoiceConnection) {
		vc.Speaking(false)
		yt.encodingSession = nil
		// Sleep for a specificed amount of time before ending.
		time.Sleep(yt.intervalTimeout)
	}(vc)

	stopChan := make(chan error)
	dca.NewStream(encodingSession, vc, stopChan)

	stats := encodingSession.Stats()
	log.Printf("started streaming: %s", video.Title)
	log.Printf("stream stats: %+v", stats)

	if err := <-stopChan; err != nil && err != io.EOF {
		return err
	}
	stats = encodingSession.Stats()
	log.Printf("streaming completed: %s", video.Title)
	log.Printf("stream stats: %+v", stats)

	return nil
}

func (yt *YoutubeSession) StopStream() error {
	if yt.encodingSession == nil {
		return errors.New("no audio currently playing")
	}
	return yt.encodingSession.Stop()
}
