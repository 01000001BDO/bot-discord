package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"encoding/json"
	"time"
)

const (
	dh 		string 	= 	"/dh"
	ai 		string 	= 	"/ai"
	dhPlay 	string 	= 	"/dh-play" 
	RED 	string  = 	"\033[31m"
	YELLOW 	string  = 	"\033[33m"
	BLUE 	string 	= 	"\033[34m"
	GREEN 	string 	= 	"\033[32m"
)


func logMsg(lvl, m string) {
	t := time.Now()
	color := ""
	switch lvl {
	case "ERROR":
		color = RED
	case "OK":
		color = GREEN
	case "WARN":
		color = YELLOW
	default:
		color = BLUE
	}
	log.Printf("%s%s[%s] %s\033[0m\n", color, t, lvl, m)
}


var firstTime = make(map[string]bool) 



func askGem(s *discordgo.Session, m *discordgo.MessageCreate, q string) string {
	s.ChannelTyping(m.ChannelID)
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error creating Gemini client: %v", err))
		return "Error creating Gemini client"
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(q))
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error generating content: %v", err))
		return "Error generating response"
	}

	var response strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				partJSON, _ := json.Marshal(part)
				response.WriteString(string(partJSON))
			}
		}
	}

	if response.Len() == 0 {
		return "No response generated"
	}

	r := response.String()
	r = strings.ReplaceAll(r, "\"", "")
	r = strings.ReplaceAll(r, "\\n", "\n")
	r = strings.ReplaceAll(r, "\\t", "\t")	
	return r
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !firstTime[m.Author.ID] {
		dm, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			logMsg("ERROR", fmt.Sprintf("Error creating DM channel: %v", err))
		}else {
			embed := &discordgo.MessageEmbed{
				Title:       "Contribute with Us!",
				Description: "Feel free to contribute <3 :\n\n[GitHub Repo Link](https://github.com/01000001BDO/bot-discord)",
				Color:       0x00FF00, 
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Thank you for your support!",
					IconURL: "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png", 
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png",
				},
			}
			_, err = s.ChannelMessageSendEmbed(dm.ID, embed)
			if err != nil {
				logMsg("ERROR", fmt.Sprintf("Error sending DM: %v", err))
			} else {
				logMsg("OK", fmt.Sprintf("Sent first-time embed message to user: %s", m.Author.Username))
				firstTime[m.Author.ID] = true 
			}
		}
	}

	args := strings.Split(m.Content, " ")
	if args[0] == ai && len(args) > 1 {
		prompt := strings.Join(args[1:], " ")
		response := askGem(s, m, prompt)
		embed := &discordgo.MessageEmbed{
			Title: "Gemini Response",
			Color: 0x00FF00,
		}
		if response == "No response generated" || response == "Error generating response" || response == "Error creating Gemini client" {
			embed.Description = "Error generating response"
			embed.Color = 0xFF0000 
		} else {
			embed.Description = response
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		logMsg("OK", fmt.Sprintf("bot jawbt l prompt: %s", prompt))
		return
	}

	if args[0] == dh {
		if args[1] == "latence" {
			latency := s.HeartbeatLatency().Milliseconds()
			author := &discordgo.MessageEmbedAuthor{
				Name:    m.Author.Username,
				IconURL: m.Author.AvatarURL(""),
			}
			embed := &discordgo.MessageEmbed{
				Title:       "Latence",
				Description: fmt.Sprintf("Latence: **%d ms**", latency),
				Author:      author,
			}
			if latency < 100 {
				embed.Color = 0x00FF00
			} else {
				embed.Color = 0xFFFF00 
				embed.Description = fmt.Sprintf("Latence: **%d ms**\n Mlagi m3a krk", latency)
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			logMsg("OK", fmt.Sprintf("jawb  l  /dh latence command b : %d ms", latency))
		}
		if args[1] == "memes" {
			logMsg("WARN", "Memes command is not yet implemented.")
		}


		if args[1] == "ls" {
			channels, err := s.GuildChannels(m.GuildID)
			if err != nil {
				logMsg("ERROR", fmt.Sprintf("Error fetching channels: %v", err))
				return
			}
			categoryMap := make(map[string][]string)
			var Allcategory []string
			for _, channel := range channels {
				if channel.Type == discordgo.ChannelTypeGuildCategory {
					continue
				}
				if channel.ParentID != "" {
					categoryMap[channel.ParentID] = append(categoryMap[channel.ParentID], channel.Name)
				} else {
					Allcategory = append(Allcategory, channel.Name)
				}
			}
			var description strings.Builder
			for categoryID, channels := range categoryMap {
				category, err := s.Channel(categoryID)
				if err != nil {
					logMsg("ERROR", fmt.Sprintf("Error fetching category: %v", err))
					continue
				}
		
				description.WriteString(fmt.Sprintf("**%s**:\n", category.Name))
				for _, channel := range channels {
					description.WriteString(fmt.Sprintf("- %s\n", channel))
				}
				description.WriteString("\n")
			}
			if len(Allcategory) > 0 {
				description.WriteString("**No Category**:\n")
				for _, channel := range Allcategory {
					description.WriteString(fmt.Sprintf("- %s\n", channel))
				}
			}
			embed := &discordgo.MessageEmbed{
				Title:       "Server Channels and Categories",
				Description: description.String(),
				Color:       0x00FF00,
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			logMsg("OK", "Listed channels and categories")
		}

		if args[1] == "pwd" {
			channel, err := s.Channel(m.ChannelID)
			if err != nil {
				logMsg("ERROR", fmt.Sprintf("Error fetching channel: %v", err))
				return
			}
			guild, err := s.Guild(m.GuildID)
			if err != nil {
				logMsg("ERROR", fmt.Sprintf("Error fetching guild: %v", err))
				return
			}
			path := fmt.Sprintf("./%s/%s/%s", m.Author.Username, guild.Name, channel.Name)
			embed := &discordgo.MessageEmbed{
				Title:       "Current Path",
				Description: path,
				Color:       0x00FF00, 
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.Author.Username,
					IconURL: m.Author.AvatarURL(""),
				},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			logMsg("OK", fmt.Sprintf("Output path: %s", path))
		}		
	}
}

func main() {
	godotenv.Load()
	if godotenv.Load() != nil {
		logMsg("ERROR", "Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	gemini := os.Getenv("GEMINI_API_KEY")
	if token == "" {
		logMsg("ERROR", "No token provided. Please set the TOKEN environment variable.")
		return
	}
	if gemini == "" {
		logMsg("ERROR", "No Gemini API key provided. Please set the GEMINI_API_KEY environment variable.")
		return
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error creating Discord session: %v", err))
		return
	}

	s.AddHandler(messageCreate)
	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = s.Open()
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error opening connection: %v", err))
		return
	}
	defer s.Close()

	logMsg("OK", "I'm alive !")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	logMsg("OK", "Bot shutting down...")
}
