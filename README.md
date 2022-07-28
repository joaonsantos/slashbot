# slashbot

This repository contains code for a golang discord music bot implementation that uses slash commands.

# why does this exist?

Discord has changed their [discord apps api](https://discord.com/blog/slash-commands-are-here) starting in 2022.

The [previous bot implementation I did](https://github.com/cloud-discord/cloudfeed) was written using the old way.
It also didn't support playlists.

Also, this seemed like a good idea to practice go channels and concurrency anyway.

# features
- streams audio from youtube video to discord voice channel
- basic playlist support (adding item, moving to next item after play)
- ability to skip audio stream, moves to next item on playlist

# running

## locally

Make sure [ffmpeg](https://ffmpeg.org/) is installed.
```bash
$ go run . -guild <guild_id> -channel <channel_id> -token <bot_token>
```

## docker (detached)
```bash
$ docker build -t ffmpeg-alpine -f dockerfile.ffmpeg .
$ docker build -t bot . && docker run --rm -d bot -guild <guild_id> -channel <channel_id> -token <bot_token>
```
You can get the guild and channel id values easily by [activating developer mode](https://apps.uk/discord-developer-mode/).

For the bot token, you need to create an account on https://discord.com/developers and create a bot.

# limitations
- the bot was designed to only handle one connection to a server, 
  if it is connected to multiple servers <em>bad things will happen™️</em>

# todo
- stap da warudo! (clean the playlist, stop audio and disconnect)

# kudos
Big kudos to:
- https://github.com/bwmarrin/discordgo
- https://github.com/jonas747/dca
- https://github.com/kkdai/youtube
