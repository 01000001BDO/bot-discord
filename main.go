package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"math"
	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const (
	dh 		string 	= 	"/dh"
	ai 		string 	= 	"/ai"
	v 		string  = 	"/ait-akinator"
	dhPlay 	string 	= 	"/dh-play" 
	RED 	string  = 	"\033[31m"
	YELLOW 	string  = 	"\033[33m"
	BLUE 	string 	= 	"\033[34m"
	GREEN 	string 	= 	"\033[32m"
)

var (
	firstTime = make(map[string]bool) 
	proba = map[string][]string{
		"wach chl7(a) ?" : {"mrkapl4n" , "KARIM" , "Kernel.rs" , "aka_bousta" , "h_a_n_a_n" , "Lynna" },
		"chl7(a) o sakn 3la bra ?" : {"Kernel.rs"  , "KARIM" },
		"pizza tafarnout ?" : {"KARIM"},
		"chl7(a) o sakn finzegane ?" : {"mrkapl4n"},
		"chl7(a) o sakn fagadir?" : {"aka_bousta" , "Lynna" , "h_a_n_a_n"},
		"faux compte ?" : {"h_a_n_a_n"},
		"la denya la akhira la 7awlawala9owatailbilah" : {"aka_bousta"},
		"skayri(a) ?" : {"aka_bousta"},
		"khona(tna) fhmator(a)? :" : {"Riquelme2.0"},
		"khdam(a) b twitter ?" : {"Kernel.rs" , "mrkapl4n" ,"aka_bousta"  , "Lynna"},
		"influencer(a) ftwitter ?" : {"Lynna"},
		"tykhra ola ttkhra ftwitter ?" : {"aka_bousta"},
		"3ndo twitter pro ?" : {"mrkapl4n" , "Kernel.rs"},
		"awal 7aja tydirha mli tyfi9 hiya tytl 3la linkedin ?" : {"mrkapl4n"},
		"tysm3 ola ttsm3 lmorphine" : {"mrkapl4n" , "Kernel.rs" , "Lynna"},
		"7M7 ?" : {"Riquelme2.0"},
		"tysm3 ola ttsm3 l Oudadn" : {"KARIM"},
		"Tymot ola ttmot 3la seddam heussin ?" : {"shaxmax"},
		"m3awd flbac ?" : {"aka_bousta" , "Bacharnaciri"},
		"m3awd flbac mra w7da ?" : {"Bacharnaciri"},
		"m3awd bzf flbac ?" : {"aka_bousta"},
		"tyl3b minceraft ?" : {"aka_bousta" , "KARIM" , "shaxmax" , "XDXG" , "! Madara⭐✨" },
		"tyl3b the finals ?"  : {"wolfmen" , "! Madara⭐✨"},
		"tyl3b mta ?" : {"aka_bousta" },
		"ty9ra ola tt9ra for fun hh" : {"Lynna"},
		"khona ola khtna fhmator(a) ?" : {"Riquelme2.0"},
		"tyl3b b bog3bob + neon ?" : {"Bacharnaciri"},
		"tyl3B ola ttl3B league of legends 7achak ?" : {"Lynna" , "XDXG" },
		"tyl3b valo ? madawich 3la hatba 7it ghi tykhra" : {"Kernel.rs" , "mrkapl4n" },
		"tytfj fseriyat bzf ?" : {"mrk4plan"},
		"ta7t liya souris ?" : {"Kernel.rs"},
		"dildo ?" : {"XDXG"},
		"3aych bra lmghrib ?" : {"XDXG" , "Kernel.rs" , "KARIM" , "shaxmax" },
		"fr3 lina krna b compte epic ?" : {"Bacharnaciri"},
		"bronze f valo ?" : {"mrkapl4n"},
		"tydir azar ?" : {"! Madara⭐✨" , "aka_bousta" },
		"azarat lmla7 ?" : {"! Madara⭐✨" },
		"zefzafi v2 ?" : {"! Madara⭐✨"},
		"mtykhrjch mn  geforece ?"  : {"wolfmen"},
		"tygol hbibi bzf ?" : {"mrkapl4n" },
		"9ari ola baghi i9ra fcmc ?" : {"mrkapl4n" , "aka_bousta" , "wolfmen"},
		"9ra fcmc ?" : {"mrkapl4n" },
		"baghi i9ra fcmc ?" : { "aka_bousta" , "wolfmen"},
		"kayn amal i9bloh fcmc ?" : {"wolfmen"},
		"7afd ga3 flags dyal dowal ?": {"XDXG"},
		"Khdam fglovo ?" : {"! Madara⭐✨"},
		"AI ghaydilih khdmto ?" : {"mrkapl4n" , "Kernel.rs"},
	}
	 LastReq  time.Time
	 activeGames = make(map[string]*GameState)
)

type GameState struct {
    CurrQ         string
    RemCan      []string
    AskedQ      []string
    Score       map[string]int
}


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

func loadData() {
	f , err := os.Open("whitelist.json")
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error opening whitelist.json: %v", err))
		firstTime = make(map[string]bool)
		return
	}
	defer f.Close()
	r := json.NewDecoder(f)
	if err := r.Decode(&firstTime) ; err != nil {
		logMsg("ERROR", fmt.Sprintf("Error decoding whitelist.json: %v", err))
		firstTime = make(map[string]bool)
		return
	}
}

func saveToJson() {
	f , err := os.Create("whitelist.json")
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error creating whitelist.json: %v", err))
		return
	}
	defer f.Close()
	w := json.NewEncoder(f) 
	if err := w.Encode(firstTime) ; err != nil {
		logMsg("ERROR", fmt.Sprintf("Error encoding whitelist.json: %v", err))
		return
	}
}


func startGame(s *discordgo.Session, m *discordgo.MessageCreate) {
    gameState := &GameState{
        RemCan: getAllCandidates(),
        Score:              make(map[string]int),
    }
    activeGames[m.ChannelID] = gameState
    nextQuestion := selectNextQuestion(gameState)
    sendQuestion(s, m.ChannelID, nextQuestion)
}

func handleAnswer(s *discordgo.Session, m *discordgo.MessageCreate, answer string) {
	gameState := activeGames[m.ChannelID]
    if gameState == nil || gameState.CurrQ == "" {
        logMsg("ERROR", "No active game or question found")
        return
    }
    for _, candidate := range gameState.RemCan {
        for question, validAnswers := range proba {
            if question == gameState.CurrQ {
                isValid := contains(validAnswers, candidate)
                if (answer == "ah" && isValid) || (answer == "la" && !isValid) {
                    gameState.Score[candidate]++
                }
            }
        }
    }

    var newCandidates []string
    maxScore := 0
    for _, candidate := range gameState.RemCan {
        if score := gameState.Score[candidate]; score >= maxScore {
            if score > maxScore {
                newCandidates = []string{candidate}
                maxScore = score
            } else {
                newCandidates = append(newCandidates, candidate)
            }
        }
    }
    gameState.RemCan = newCandidates
    if len(gameState.RemCan) == 1 || len(gameState.AskedQ) >= 10 {
        sendResult(s, m.ChannelID, gameState)
        delete(activeGames, m.ChannelID)
        return
    }

    nextQuestion := selectNextQuestion(gameState)
    sendQuestion(s, m.ChannelID, nextQuestion)
}

func selectNextQuestion(state *GameState) string {
    if len(state.RemCan) == 0 {
        return ""
    }

    bestQuestion := ""
    bestSplit := 0.0
    for question := range proba {
        if contains(state.AskedQ, question) {
            continue
        }

        validCount := 0
        for _, candidate := range state.RemCan {
            if contains(proba[question], candidate) {
                validCount++
            }
        }

        split := math.Abs(float64(validCount) - float64(len(state.RemCan))/2)
        if bestQuestion == "" || split > bestSplit {
            bestSplit = split
            bestQuestion = question
        }
    }

    if bestQuestion != "" {
        state.CurrQ = bestQuestion
        state.AskedQ = append(state.AskedQ, bestQuestion)
    }
    return bestQuestion
}

func getAllCandidates() []string {
	arr := make(map[string]bool)
    var c []string
    for _, validAnswers := range proba {
        for _, i := range validAnswers {
            if !arr[i] {
                arr[i] = true
                c = append(c, i)
            }
        }
    }
    return c
}

func sendQuestion(s *discordgo.Session, channelID string, question string) {
	embed := &discordgo.MessageEmbed{
		Title: "Ait Akinator",
		Description: question + "\n\n jawb b **/ait-akinator ah** or **/ait-akinator la**",
		Color: 0x00FF00,
	}
	s.ChannelMessageSendEmbed(channelID, embed)
}
func sendResult(s *discordgo.Session, channelID string, state *GameState) {
	var r string
    if len(state.RemCan) == 1 {
		r = "Howa ola machi howa: **@" + state.RemCan[0]+"**"

    } else {
        r = " AYkon ghi wa7d  min hado , skill issues (probelem fl'arbre ) : " + strings.Join(state.RemCan, ", ")
    }
    embed := &discordgo.MessageEmbed{
        Title: "Ezzz Ezzz",
        Description: r,
        Color: 0x00FF00,
    }
    s.ChannelMessageSendEmbed(channelID, embed)
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}


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
				saveToJson()
			}
		}
	}

	args := strings.Split(m.Content, " ")
	if args[0] == ai && len(args) > 1 {
		if time.Since(LastReq) < 5*time.Second {
			embed := &discordgo.MessageEmbed{
				Title:       "Tsna a w9",
				Description: "Wach baghi ty7ni , ra tnmil mtnti7ch",
				Color : 0xFF0000,
			}
			_ , err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				logMsg("ERROR", fmt.Sprintf("Error sending message: %v", err))
			}
			return
		}
		LastReq = time.Now()
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
	if args[0] == v {
		if len(args) == 1 {
			startGame(s, m)
			return
		}
		answer := strings.ToLower(args[1])
		if answer == "ah" || answer == "la" {
			handleAnswer(s, m, answer)
		}
	}
}

func main() {
	loadData()
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