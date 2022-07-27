package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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
					Description: "Add an url to the playlist. It must link to an audio stream on youtube.",
					Required:    true,
				},
			},
		},
		{
			Name:        "next-music",
			Description: "Get next music url from the queue",
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
			Playlist.Add(songURL)
			currSize := Playlist.Size()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("added %s to the playlist\n\n %d song(s) now in the playlist", songURL, currSize),
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

				for Playlist.Size() != 0 {
					songURL, err := Playlist.GetNext()
					if err != nil {
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
							Content: fmt.Sprintf("playing %s\n\n %d song(s) remaining in the playlist", songURL, currSize),
						},
					})
					StreamAudio(vc, *GuildID, *ChannelID, songURL)
				}
				vc.Disconnect()
			}(s)

		},
		"skip-music": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if err := StopStream(); err != nil {
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
