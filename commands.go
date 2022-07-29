package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/joaonsantos/slashbot/internal/youtube"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help-command",
			Description: "Help command",
		},
		{
			Name:        "add-music",
			Description: "Add a music url to the queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "music-url",
					Description: "Add an url item to the playlist, it must link to a youtube video",
					Required:    true,
				},
			},
		},
		{
			Name:        "next-music",
			Description: "Get next item from the queue",
		},
		{
			Name:        "skip-music",
			Description: "Skip currently playing audio stream",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there, I am a bot built on slash commands! I can play music for you!\n" +
						"\nAdd items to the playlist with `/add-music`." +
						"\nGet the party started with `/next-music`." +
						"\nTired of the current music? Use `/skip-music`!",
				},
			})
		},
		"add-music": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			songURL := optionMap["music-url"].StringValue()
			if err := youtube.ValidateURL(songURL); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("failed to validate url: %s", err.Error()),
					},
				})
				return
			}
			Playlist.Add(songURL)
			currSize := Playlist.Size()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("added %s to the playlist\n\n %d item(s) now in the playlist", songURL, currSize),
				},
			})
		},
		"next-music": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			go func(s *discordgo.Session) {
				vc, err := s.ChannelVoiceJoin(*GuildID, *ChannelID, false, true)
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("failed to join voice channel: %s", err.Error()),
						},
					})
					return
				}
				defer vc.Disconnect()

				for Playlist.Size() != 0 {
					songURL, err := Playlist.GetNext()
					if err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("failed to get an item from the playlist: %s", err.Error()),
							},
						})
						return
					}
					currSize := Playlist.Size()
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("playing %s\n\n %d item(s) remaining in the playlist", songURL, currSize),
						},
					})
					if err := YtSession.StreamYoutubeVideo(vc, songURL); err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("audio stream failed: %s", err.Error()),
							},
						})
					}
				}
			}(s)

		},
		"skip-music": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := YtSession.StopStream(); err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("failed to get music from playlist: %s", err.Error()),
					},
				})
				return
			}

			currSize := Playlist.Size()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("stopping audio stream\n\n %d song(s) remaining in the playlist", currSize),
				},
			})
		},
	}
)
