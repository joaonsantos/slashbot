package youtube

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
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
	const audioFilename = "/tmp/audio.mp4"

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
	defer stream.Close()

	if err := saveAudio(audioFilename, stream); err != nil {
		return err
	}
	yt.encodingSession, err = dca.EncodeFile(audioFilename, options)
	if err != nil {
		return err
	}

	stopChan := make(chan error)
	dca.NewStream(yt.encodingSession, vc, stopChan)

	yt.startSession(vc)
	defer yt.stopSession(vc)

	stats := yt.encodingSession.Stats()
	log.Printf("started streaming: %s", video.Title)
	log.Printf("stream stats: %+v", stats)

	if err := <-stopChan; err != nil && err != io.EOF {
		return err
	}
	stats = yt.encodingSession.Stats()
	log.Printf("streaming completed: %s", video.Title)
	log.Printf("stream stats: %+v", stats)

	return nil
}

func (yt *YoutubeSession) startSession(vc *discordgo.VoiceConnection) {
	// Sleep for a specified amount of time before playing the sound
	time.Sleep(yt.intervalTimeout)
	vc.Speaking(true)
}

func (yt *YoutubeSession) stopSession(vc *discordgo.VoiceConnection) {
	log.Println("cleaning up encoding session")
	vc.Speaking(false)
	yt.encodingSession.Cleanup()
	yt.encodingSession = nil
	// Sleep for a specificed amount of time before ending.
	time.Sleep(yt.intervalTimeout)
}

func (yt *YoutubeSession) StopStream() error {
	if yt.encodingSession == nil {
		return errors.New("no audio currently playing")
	}
	return yt.encodingSession.Stop()
}

func saveAudio(filename string, src io.ReadCloser) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		return err
	}
	return nil
}

func ValidateURL(url string) error {
	if strings.Contains(url, "playlist?list") {
		return errors.New("this command does not support adding playlists")
	}
	if !strings.HasPrefix(url, "https://www.youtube.com/watch?v=") {
		return errors.New("provide a valid link to a youtube audio stream")
	}
	_, err := youtube.ExtractVideoID(url)
	if err != nil {
		return err
	}
	return nil
}
