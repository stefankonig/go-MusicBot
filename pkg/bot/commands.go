package bot

import (
	"fmt"
	"strings"

	"github.com/svenwiltink/go-musicbot/pkg/music"
	"strconv"
	"time"
)

type Command struct {
	Name       string
	MasterOnly bool
	Function   func(bot *MusicBot, message Message)
}

var helpCommand = Command{
	Name: "help",
	Function: func(bot *MusicBot, message Message) {
		helpString := "Available commands: "
		for _, command := range bot.commands {
			helpString += command.Name + " "
		}

		bot.ReplyToMessage(message, helpString)
	},
}

func sanitizeSong(song string) string {
	song = strings.TrimSpace(song)
	song = strings.TrimLeft(song, "<")
	song = strings.TrimRight(song, ">")
	return song
}

var addCommand = Command{
	Name: "add",
	Function: func(bot *MusicBot, message Message) {
		words := strings.SplitN(message.Message, " ", 3)
		if len(words) <= 2 {
			bot.ReplyToMessage(message, "No song provided")
			return
		}

		song := music.Song{
			Path: sanitizeSong(words[2]),
		}

		song, err := bot.musicPlayer.AddSong(song)
		if err != nil {
			bot.ReplyToMessage(message, err.Error())
			return
		}

		if message.IsPrivate {
			bot.BroadcastMessage(fmt.Sprintf("%s: %s added by %s", song.Artist, song.Name, message.Sender.Name))
		}
		bot.ReplyToMessage(message, fmt.Sprintf("%s: %s added", song.Artist, song.Name))

	},
}

var searchCommand = Command{
	Name: "search",
	Function: func(bot *MusicBot, message Message) {
		words := strings.SplitN(message.Message, " ", 3)
		if len(words) <= 2 {
			bot.ReplyToMessage(message, "No song provided")
			return
		}

		songs, err := bot.musicPlayer.Search(strings.TrimSpace(words[2]))

		if err != nil {
			bot.ReplyToMessage(message, fmt.Sprintf("error: %v", err))
		}

		if len(songs) == 0 {
			bot.ReplyToMessage(message, "No song found")
			return
		}

		var builder strings.Builder

		for number, song := range songs {
			builder.WriteString(fmt.Sprintf("%d  %s - %s (%s)\n", number+1, song.Artist, song.Name, song.Duration))
		}

		bot.ReplyToMessage(message, builder.String())

	},
}

var searchAddCommand = Command{
	Name: "search-add",
	Function: func(bot *MusicBot, message Message) {
		words := strings.SplitN(message.Message, " ", 3)
		if len(words) <= 2 {
			bot.ReplyToMessage(message, "No song provided")
			return
		}

		songs, err := bot.musicPlayer.Search(strings.TrimSpace(words[2]))

		if err != nil {
			bot.ReplyToMessage(message, fmt.Sprintf("error: %v", err))
		}

		if len(songs) == 0 {
			bot.ReplyToMessage(message, "No song found")
			return
		}

		song := songs[0]
		song, err = bot.musicPlayer.AddSong(song)

		if err != nil {
			bot.ReplyToMessage(message, fmt.Sprintf("Error: %v", err))
			return
		}

		if message.IsPrivate {
			bot.BroadcastMessage(fmt.Sprintf("%s: %s added by %s", song.Artist, song.Name, message.Sender.Name))
		}

		bot.ReplyToMessage(message, fmt.Sprintf("%s: %s added", song.Artist, song.Name))
	},
}

var nextCommand = Command{
	Name: "next",
	Function: func(bot *MusicBot, message Message) {
		err := bot.musicPlayer.Next()
		if err != nil {
			bot.ReplyToMessage(message, fmt.Sprintf("Could not skip song: %v", err))
		} else {
			if message.IsPrivate {
				bot.BroadcastMessage(fmt.Sprintf("%s skipped the song", message.Sender.Name))
			}
			bot.ReplyToMessage(message, "Skipping song")
		}
	},
}

var pausedCommand = Command{
	Name: "pause",
	Function: func(bot *MusicBot, message Message) {
		err := bot.musicPlayer.Pause()
		if err != nil {
			bot.ReplyToMessage(message, fmt.Sprintf("Error: %v", err))
			return
		}

		if message.IsPrivate {
			bot.BroadcastMessage(fmt.Sprintf("%s stopped the music", message.Sender.Name))
		}

		bot.ReplyToMessage(message, "Music paused")

	},
}

var playCommand = Command{
	Name: "play",
	Function: func(bot *MusicBot, message Message) {
		err := bot.musicPlayer.Play()
		if err != nil {
			bot.ReplyToMessage(message, fmt.Sprintf("Error: %v", err))
			return
		}

		if message.IsPrivate {
			bot.BroadcastMessage(fmt.Sprintf("%s resumed the music", message.Sender.Name))
		}

		bot.ReplyToMessage(message, "Music resumed")

	},
}

var currentCommand = Command{
	Name: "current",
	Function: func(bot *MusicBot, message Message) {
		song, durationLeft := bot.musicPlayer.GetCurrentSong()
		if song == nil {
			bot.ReplyToMessage(message, "Nothing currently playing")
			return
		}

		if song.SongType == music.SongTypeSong {
			bot.ReplyToMessage(
				message,
				fmt.Sprintf("Current song: %s %s. %s remaining (%s)", song.Artist, song.Name, durationLeft.String(), song.Duration.Round(time.Second).String()))
		} else {
			bot.ReplyToMessage(
				message,
				fmt.Sprintf("Current song: %s %s. This is a livestream, use the next command to skip", song.Artist, song.Name))
		}
	},
}

var queueCommand = Command{
	Name: "queue",
	Function: func(bot *MusicBot, message Message) {
		queue := bot.GetMusicPlayer().GetQueue()

		queueLength := queue.GetLength()
		nextSongs, _ := queue.GetNextN(5)
		duration := queue.GetTotalDuration()

		bot.ReplyToMessage(message, fmt.Sprintf("%d songs in the queue. Total duration %s", queueLength, duration.String()))

		for index, song := range nextSongs {
			bot.ReplyToMessage(message, fmt.Sprintf("#%d, %s: %s (%s)\n", index+1, song.Artist, song.Name, song.Duration.String()))
		}

		if queueLength > 5 {
			bot.ReplyToMessage(message, fmt.Sprintf("and %d more", queueLength-5))
		}

	},
}

var flushCommand = Command{
	Name: "flush",
	Function: func(bot *MusicBot, message Message) {
		bot.musicPlayer.GetQueue().Flush()

		if message.IsPrivate {
			bot.BroadcastMessage(fmt.Sprintf("%s flushed the queue", message.Sender.Name))
		}

		bot.ReplyToMessage(message, "Queue flushed")
	},
}

var shuffleCommand = Command{
	Name: "shuffle",
	Function: func(bot *MusicBot, message Message) {
		bot.musicPlayer.GetQueue().Shuffle()

		if message.IsPrivate {
			bot.BroadcastMessage(fmt.Sprintf("%s shuffled the queue", message.Sender.Name))
		}

		bot.ReplyToMessage(message, "Queue shuffled")
	},
}

var whiteListCommand = Command{
	Name:       "whitelist",
	MasterOnly: true,
	Function: func(bot *MusicBot, message Message) {
		words := strings.SplitN(message.Message, " ", 4)
		if len(words) <= 3 {
			bot.ReplyToMessage(message, "whitelist <add|remove> <name>")
			return
		}

		name := strings.TrimSpace(words[3])
		if len(name) == 0 {
			bot.ReplyToMessage(message, "whitelist <add|remove> <name>")
			return
		}

		if words[2] == "add" {
			err := bot.whitelist.Add(name)
			if err == nil {
				bot.ReplyToMessage(message, fmt.Sprintf("added %s to the whitelist", name))
			} else {
				bot.ReplyToMessage(message, fmt.Sprintf("error: %v", err))
			}
		} else if words[2] == "remove" {
			err := bot.whitelist.Remove(name)
			if err == nil {
				bot.ReplyToMessage(message, fmt.Sprintf("removed %s from the whitelist", name))
			} else {
				bot.ReplyToMessage(message, fmt.Sprintf("error: %v", err))
			}
		} else {
			bot.ReplyToMessage(message, "whitelist <add|remove> <name>")
			return
		}
	},
}

var volCommand = Command{
	Name: "vol",
	Function: func(bot *MusicBot, message Message) {
		words := strings.SplitN(message.Message, " ", 3)

		if len(words) == 2 {
			volume, err := bot.musicPlayer.GetVolume()

			if err != nil {
				bot.ReplyToMessage(message, fmt.Sprintf("unable to get volume: %v", err))
				return
			}

			bot.ReplyToMessage(message, fmt.Sprintf("Current volume: %d", volume))
			return
		}

		// init vars here so we can use them after the switch statement
		volumeString := strings.TrimSpace(words[2])
		var volume int
		var err error

		switch volumeString {
		case "++":
			{
				volume, err = bot.musicPlayer.IncreaseVolume(10)
				if err != nil {
					bot.ReplyToMessage(message, fmt.Sprintf("unable to increase volume: %s", err))
					return
				}
			}
		case "--":
			{
				volume, err = bot.musicPlayer.DecreaseVolume(10)
				if err != nil {
					bot.ReplyToMessage(message, fmt.Sprintf("unable to decrease volume: %s", err))
					return
				}
			}
		default:
			{
				volume, err = strconv.Atoi(strings.TrimSpace(volumeString))

				if err != nil {
					bot.ReplyToMessage(message, fmt.Sprintf("%s is not a valid number", volumeString))
					return
				}

				if volume >= 0 && volume <= 100 {
					bot.musicPlayer.SetVolume(volume)
				} else {
					bot.ReplyToMessage(message, fmt.Sprintf("%s is not a valid volume", volumeString))
					return
				}
			}
		}

		if message.IsPrivate {
			bot.BroadcastMessage(fmt.Sprintf("Volume set to %d by %s", volume, message.Sender.Name))
		}

		bot.ReplyToMessage(message, fmt.Sprintf("Volume set to %d", volume))
	},
}

var aboutCommand = Command{
	Name: "about",
	Function: func(bot *MusicBot, message Message) {
		bot.ReplyToMessage(message, "go-MusicBot by Sven Wiltink: https://github.com/svenwiltink/go-MusicBot")
		bot.ReplyToMessage(message, fmt.Sprintf("Version: %s", Version))
		bot.ReplyToMessage(message, fmt.Sprintf("Go: %s", GoVersion))
		bot.ReplyToMessage(message, fmt.Sprintf("Build date: %s", BuildDate))
	},
}
