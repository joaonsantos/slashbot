package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/cloud-discord/slashbot/internal/playlist"
	"github.com/cloud-discord/slashbot/internal/youtube"
)

const (
	PlaylistMaxSize = 10
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	ChannelID      = flag.String("channel", "", "Test Channel ID")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	VoiceBuffer    = make([][]byte, 0)
	Playlist       = playlist.New(PlaylistMaxSize)
	YtSession      = youtube.NewSession()
	StopStreamChan chan error
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open discord websocket connection, needed for voice
	if err := s.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Registering commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer s.Close()

	// Listen for kill signals from os
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-sig

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
	log.Println("Gracefully shutting down.")
}

func StreamAudio(vc *discordgo.VoiceConnection, guildID, channelID, songURL string) error {
	if err := YtSession.StreamYoutubeVideo(vc, songURL); err != nil {
		return err
	}
	return nil
}

func StopStream() error {
	return YtSession.StopStream()
}
