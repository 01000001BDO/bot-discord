package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	"layeh.com/gopus"
)

const (
	dh         string        = "/dh"
	ai         string        = "/ai"
	v          string        = "/ait-akinator"
	valoPing   string        = "/dh valo-ping"
	morphineCmd string 		 = "/dh lmorphine"
	RED        string        = "\033[31m"
	YELLOW     string        = "\033[33m"
	BLUE       string        = "\033[34m"
	GREEN      string        = "\033[32m"
	channels   int           = 2
	frameRate  int           = 48000
	frameSize  int           = 960
	NoActivity VoiceActivity = iota
	MusicPlaying
	TTSPlaying
)
var (
	firstTime = make(map[string]bool)
	proba     = map[string][]string{
		"wach chl7(a) ?":                                       {"mrkapl4n", "KARIM", "Kernel.rs", "aka_bousta", "h_a_n_a_n", "Lynna"},
		"chl7(a) o sakn 3la bra ?":                             {"Kernel.rs", "KARIM"},
		"pizza tafarnout ?":                                    {"KARIM"},
		"chl7(a) o sakn finzegane ?":                           {"mrkapl4n"},
		"chl7(a) o sakn fagadir?":                              {"aka_bousta", "Lynna", "h_a_n_a_n"},
		"faux compte ?":                                        {"h_a_n_a_n"},
		"la denya la akhira la 7awlawala9owatailbilah":         {"aka_bousta"},
		"skayri(a) ?":                                          {"aka_bousta"},
		"khona(tna) fhmator(a)? :":                             {"Riquelme2.0"},
		"khdam(a) b twitter ?":                                 {"Kernel.rs", "mrkapl4n", "aka_bousta", "Lynna"},
		"influencer(a) ftwitter ?":                             {"Lynna"},
		"tykhra ola ttkhra ftwitter ?":                         {"aka_bousta"},
		"3ndo twitter pro ?":                                   {"mrkapl4n", "Kernel.rs"},
		"awal 7aja tydirha mli tyfi9 hiya tytl 3la linkedin ?": {"mrkapl4n"},
		"tysm3 ola ttsm3 lmorphine":                            {"mrkapl4n", "Kernel.rs", "Lynna"},
		"7M7 ?":                                                {"Riquelme2.0"},
		"tysm3 ola ttsm3 l Oudadn":                             {"KARIM"},
		"Tymot ola ttmot 3la seddam heussin ?":                 {"shaxmax"},
		"m3awd flbac ?":                                        {"aka_bousta", "Bacharnaciri"},
		"m3awd flbac mra w7da ?":                               {"Bacharnaciri"},
		"m3awd bzf flbac ?":                                    {"aka_bousta"},
		"tyl3b minceraft ?":                                    {"aka_bousta", "KARIM", "shaxmax", "XDXG", "! Madara⭐✨"},
		"tyl3b the finals ?":                                   {"wolfmen", "! Madara⭐✨"},
		"tyl3b mta ?":                                          {"aka_bousta"},
		"ty9ra ola tt9ra for fun hh":                           {"Lynna"},
		"khona ola khtna fhmator(a) ?":                         {"Riquelme2.0"},
		"tyl3b b bog3bob + neon ?":                             {"Bacharnaciri"},
		"tyl3B ola ttl3B league of legends 7achak ?":           {"Lynna", "XDXG"},
		"tyl3b valo ? madawich 3la hatba 7it ghi tykhra":       {"Kernel.rs", "mrkapl4n"},
		"tytfj fseriyat bzf ?":                                 {"mrk4plan"},
		"ta7t liya souris ?":                                   {"Kernel.rs"},
		"dildo ?":                                              {"XDXG"},
		"3aych bra lmghrib ?":                                  {"XDXG", "Kernel.rs", "KARIM", "shaxmax"},
		"fr3 lina krna b compte epic ?":                        {"Bacharnaciri"},
		"bronze f valo ?":                                      {"mrkapl4n"},
		"tydir azar ?":                                         {"! Madara⭐✨", "aka_bousta"},
		"azarat lmla7 ?":                                       {"! Madara⭐✨"},
		"zefzafi v2 ?":                                         {"! Madara⭐✨"},
		"mtykhrjch mn  geforece ?":                             {"wolfmen"},
		"tygol hbibi bzf ?":                                    {"mrkapl4n"},
		"9ari ola baghi i9ra fcmc ?":                           {"mrkapl4n", "aka_bousta", "wolfmen"},
		"9ra fcmc ?":                                           {"mrkapl4n"},
		"baghi i9ra fcmc ?":                                    {"aka_bousta", "wolfmen"},
		"kayn amal i9bloh fcmc ?":                              {"wolfmen"},
		"7afd ga3 flags dyal dowal ?":                          {"XDXG"},
		"Khdam fglovo ?":                                       {"! Madara⭐✨"},
		"AI ghaydilih khdmto ?":                                {"mrkapl4n", "Kernel.rs"},
	}
	LastReq      time.Time
	activeGames  = make(map[string]*GameState)
	players      = make(map[string]*MusicPlayer)
	yt           = youtube.Client{}
    voiceManager = &VoiceStateManager{
        guildStates: make(map[string]*GuildVoiceState),
    }
	valorantServers = []ValorantServer{
		{Name: "EU Frankfurt 1", IP: "35.198.119.251", Location: "Frankfurt, Germany"},
		{Name: "EU Frankfurt 2", IP: "35.198.119.252", Location: "Frankfurt, Germany"},
		{Name: "EU Paris", IP: "35.198.119.253", Location: "Paris, France"},
		{Name: "EU London", IP: "35.198.119.254", Location: "London, UK"},
		{Name: "EU Stockholm", IP: "35.198.119.255", Location: "Stockholm, Sweden"},
		{Name: "EU Warsaw", IP: "35.198.119.250", Location: "Warsaw, Poland"},
	}
)

type VoiceActivity int
type GameState struct {
	CurrQ  string
	RemCan []string
	AskedQ []string
	Score  map[string]int
}

type Song struct {
	URL      string
	Title    string
	Duration string
}

type MusicPlayer struct {
	queue     []Song
	isPlaying bool
	voiceConn *discordgo.VoiceConnection
	stopChan  chan bool
	mu        sync.Mutex
}

type GuildVoiceState struct {
	currentActivity VoiceActivity
	mu              sync.Mutex
}

type VoiceStateManager struct {
	guildStates map[string]*GuildVoiceState
	mu          sync.RWMutex
}

type ValorantServer struct {
    Name     string
    IP       string
    Location string
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

func (v *VoiceStateManager) GetGuildState(guildID string) *GuildVoiceState {
    v.mu.Lock()
    defer v.mu.Unlock()

    state, exists := v.guildStates[guildID]
    if !exists {
        state = &GuildVoiceState{
            currentActivity: NoActivity,
        }
        v.guildStates[guildID] = state
    }
    return state
}
func (v *VoiceStateManager) SetActivity(guildID string, activity VoiceActivity) bool {
    state := v.GetGuildState(guildID)
    state.mu.Lock()
    defer state.mu.Unlock()
    if state.currentActivity != NoActivity && state.currentActivity != activity {
        return false
    }

    state.currentActivity = activity
    return true
}

func (v *VoiceStateManager) ClearActivity(guildID string, activity VoiceActivity) {
    state := v.GetGuildState(guildID)
    state.mu.Lock()
    defer state.mu.Unlock()
    if state.currentActivity == activity {
        state.currentActivity = NoActivity
    }
}

func (v *VoiceStateManager) GetCurrentActivity(guildID string) VoiceActivity {
    state := v.GetGuildState(guildID)
    state.mu.Lock()
    defer state.mu.Unlock()
    return state.currentActivity
}

func loadData() {
	f, err := os.Open("whitelist.json")
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error opening whitelist.json: %v", err))
		firstTime = make(map[string]bool)
		return
	}
	defer f.Close()
	r := json.NewDecoder(f)
	if err := r.Decode(&firstTime); err != nil {
		logMsg("ERROR", fmt.Sprintf("Error decoding whitelist.json: %v", err))
		firstTime = make(map[string]bool)
		return
	}
}

func saveToJson() {
	f, err := os.Create("whitelist.json")
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error creating whitelist.json: %v", err))
		return
	}
	defer f.Close()
	w := json.NewEncoder(f)
	if err := w.Encode(firstTime); err != nil {
		logMsg("ERROR", fmt.Sprintf("Error encoding whitelist.json: %v", err))
		return
	}
}

func joinVoiceChannel(s *discordgo.Session, guildID, channelID string) (*discordgo.VoiceConnection, error) {
	return s.ChannelVoiceJoin(guildID, channelID, false, true)
}

func getOrCreatePlayer(guildID string) *MusicPlayer {
	if player, exists := players[guildID]; exists {
		return player
	}

	player := &MusicPlayer{
		queue:     make([]Song, 0),
		isPlaying: false,
		stopChan:  make(chan bool),
	}
	players[guildID] = player
	return player
}

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
    if !strings.HasPrefix(m.Content, valoPing) {
        return
    }
    embed := &discordgo.MessageEmbed{
        Title:       "🔍 Omoro ta7dot ",
        Description: "tsna wa7d chwiya , l omor ta7dot",
        Color:       0x00FF00,
    }
    msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
    if err != nil {
        logMsg("ERROR", fmt.Sprintf("Error sending initial message: %v", err))
        return
    }

    type PingResult struct {
        Name     string  `json:"name"`
        Location string  `json:"location"`
        Ping     float64 `json:"ping"`
    }
    var results []PingResult

    for _, server := range valorantServers {
        cmd := exec.Command("ping", "-c", "3", server.IP)
        output, err := cmd.CombinedOutput()
        var pingTime float64
        if err != nil {
            pingTime = 999
        } else {
            outputStr := string(output)
            re := regexp.MustCompile(`time=(\d+\.?\d*)`)
            matches := re.FindAllStringSubmatch(outputStr, -1)
            
            if len(matches) > 0 {
                var total float64
                count := 0
                for _, match := range matches {
                    if len(match) > 1 {
                        if val, err := strconv.ParseFloat(match[1], 64); err == nil {
                            total += val
                            count++
                        }
                    }
                }
                if count > 0 {
                    pingTime = total / float64(count)
                }
            }
        }

        results = append(results, PingResult{
            Name:     server.Name,
            Location: server.Location,
            Ping:     math.Round(pingTime*100) / 100,
        })
    }

    embed = &discordgo.MessageEmbed{
        Title:       "🌐 Valorant EU Server Pings",
        Description: "Lpingat lmla7:",
        Color:       0x00FF00,
        Fields:      make([]*discordgo.MessageEmbedField, 0),
    }

    for _, result := range results {
        pingStatus := "🟢" 
        if result.Ping > 100 {
            pingStatus = "🔴" 
        } else if result.Ping > 50 {
            pingStatus = "🟡" 
        }

        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name:   fmt.Sprintf("%s %s", pingStatus, result.Name),
            Value:  fmt.Sprintf("```\nLocation: %s\nPing: %.2fms\n```", result.Location, result.Ping),
            Inline: true,
        })
    }
    _, err = s.ChannelMessageEditEmbed(m.ChannelID, msg.ID, embed)
    if err != nil {
        logMsg("ERROR", fmt.Sprintf("Error updating ping message: %v", err))
    }
}



func handleMusic(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a command: play, skip, stop, queue")
		return
	}
	vs, err := findUserVoiceState(s, m.GuildID, m.Author.ID)
	if err != nil {
		embed := &discordgo.MessageEmbed{
			Title:       "Chi 7aja trat !!!",
			Description: "Khoya wach 7mar ? khsk tkon fchi room wach baghi ndkhl lik fkrk ?",
			Color:       0xFF0000,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	player := getOrCreatePlayer(m.GuildID)

	switch args[1] {
	case "play":
		if len(args) < 3 {
			embed := &discordgo.MessageEmbed{
				Title:       "Chi 7aja trat !!!",
				Description: "Hbibi lien dyal mzika",
				Color:       0xFF0000,
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
		handlePlay(s, m, vs.ChannelID, args[2], player)

	case "skip":
		handleSkip(s, m, player)

	case "stop":
		handleStop(s, m, player)

	case "queue":
		handleQueue(s, m, player)
	}
}

func handlePlay(s *discordgo.Session, m *discordgo.MessageCreate, voiceChannelID string, url string, player *MusicPlayer) {
    if !voiceManager.SetActivity(m.GuildID, MusicPlaying) {
        currentActivity := voiceManager.GetCurrentActivity(m.GuildID)
        var message string
        if currentActivity == TTSPlaying {
            message = "bot tydwi, tsna hta ysali"
        } else {
            message = "bot mkhdm mzika, tsna hta ysali"
        }
        
        embed := &discordgo.MessageEmbed{
            Title: "3a9o bika",
            Description: message,
            Color: 0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }

    player.mu.Lock()
    defer player.mu.Unlock()
    cmd := exec.Command("yt-dlp", "-j", 
    "--no-check-certificates",
    "--ignore-errors", 
    "--no-playlist",
    "--extractor-args", "youtube:player_client=android",
    url)
    output, err := cmd.CombinedOutput() 
    if err != nil {
        voiceManager.ClearActivity(m.GuildID, MusicPlaying)
        embed := &discordgo.MessageEmbed{
            Title: "Chi 7aja trat !!!",
            Description: fmt.Sprintf("Error ma9drtch njib video: %v\nOutput: %s", err, string(output)), 
            Color: 0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }

    var videoInfo struct {
        Title    string `json:"title"`
        Duration int    `json:"duration"`
    }
    if err := json.Unmarshal(output, &videoInfo); err != nil {
        return
    }

    song := Song{
        URL:      url,
        Title:    videoInfo.Title,
        Duration: fmt.Sprintf("%d:%02d", videoInfo.Duration/60, videoInfo.Duration%60),
    }
    player.queue = append(player.queue, song)
    
    embed := &discordgo.MessageEmbed{
        Title: "Jdid fl Queue",
        Description: fmt.Sprintf("🎵 **%s**\n⏱️ Duration: %s", song.Title, song.Duration),
        Color: 0x00FF00,
        Footer: &discordgo.MessageEmbedFooter{
            Text: "Added by " + m.Author.Username,
            IconURL: m.Author.AvatarURL(""),
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)

    if !player.isPlaying {
        go startPlaying(s, m.GuildID, voiceChannelID, player)
    }
}
func startPlaying(s *discordgo.Session, guildID string, voiceChannelID string, player *MusicPlayer) {
    defer voiceManager.ClearActivity(guildID, MusicPlaying)
    
    player.mu.Lock()
    if player.isPlaying {
        player.mu.Unlock()
        return
    }
    player.isPlaying = true
    player.mu.Unlock()

	for {
		player.mu.Lock()
		if len(player.queue) == 0 {
			player.isPlaying = false
			if player.voiceConn != nil {
				player.voiceConn.Disconnect()
				player.voiceConn = nil
			}
			embed := &discordgo.MessageEmbed{
				Title:       "Queue Finished",
				Description: "📭 Queue khawya, zid chi 7aja akhora",
				Color:       0x00FF00,
			}
			s.ChannelMessageSendEmbed(voiceChannelID, embed)
			player.mu.Unlock()
			return
		}

		currentSong := player.queue[0]
		player.queue = player.queue[1:]
		player.mu.Unlock()
		if player.voiceConn == nil {
			vc, err := joinVoiceChannel(s, guildID, voiceChannelID)
			if err != nil {
				continue
			}
			player.voiceConn = vc
		}
		embed := &discordgo.MessageEmbed{
			Title:       "Sm3 Sm3 🎵",
			Description: fmt.Sprintf("**%s**\n⏱️ Duration: %s", currentSong.Title, currentSong.Duration),
			Color:       0x00FF00,
		}
		s.ChannelMessageSendEmbed(voiceChannelID, embed)

		video, err := yt.GetVideo(currentSong.URL)
		if err != nil {
			continue
		}

		formats := video.Formats.WithAudioChannels()
		stream, _, err := yt.GetStream(video, &formats[0])
		if err != nil {
			continue
		}

		ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
		ffmpeg.Stdin = stream
		stdout, err := ffmpeg.StdoutPipe()
		if err != nil {
			continue
		}

		err = ffmpeg.Start()
		if err != nil {
			continue
		}
		encoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
		if err != nil {
			continue
		}
		finished := make(chan bool)
		interrupt := make(chan bool)
		go func() {
			defer func() {
				ffmpeg.Process.Kill()
				close(finished)
			}()

			buffer := make([]int16, frameSize*channels)
			for {
				select {
				case <-player.stopChan:
					return
				case <-interrupt:
					return
				default:
					err := binary.Read(stdout, binary.LittleEndian, &buffer)
					if err != nil {
						if err == io.EOF {
							return
						}
						fmt.Println("Error reading from ffmpeg stdout:", err)
						return
					}
					opus, err := encoder.Encode(buffer, frameSize, frameSize*2)
					if err != nil {
						fmt.Println("Error encoding to Opus:", err)
						return
					}
					select {
					case player.voiceConn.OpusSend <- opus:
					case <-player.stopChan:
						return
					case <-interrupt:
						return
					}
				}
			}
		}()
		select {
		case <-finished:
			time.Sleep(500 * time.Millisecond)
		case <-player.stopChan:
			close(interrupt)
			player.stopChan = make(chan bool)
		}
	}
}

func handleSkip(s *discordgo.Session, m *discordgo.MessageCreate, player *MusicPlayer) {
	player.mu.Lock()
	defer player.mu.Unlock()

	if len(player.queue) == 0 && !player.isPlaying {
		embed := &discordgo.MessageEmbed{
			Title:       "Chi 7aja trat !!!",
			Description: "Queue khawya akhoya",
			Color:       0xFF0000,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	close(player.stopChan)
	player.stopChan = make(chan bool)

	embed := &discordgo.MessageEmbed{
		Title:       "Skip ✅",
		Description: "⏭️ tskipat a hbibi ",
		Color:       0x00FF00,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Skipped by " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleQueue(s *discordgo.Session, m *discordgo.MessageCreate, player *MusicPlayer) {
	player.mu.Lock()
	defer player.mu.Unlock()

	if len(player.queue) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "Queue",
			Description: "📭 Queue khawi azb",
			Color:       0xFF0000,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	var queueMsg strings.Builder
	for i, song := range player.queue {
		queueMsg.WriteString(fmt.Sprintf("%d. 🎵 **%s**\n⏱️ Duration: %s\n\n", i+1, song.Title, song.Duration))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Queue ✅",
		Description: queueMsg.String(),
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Total Songs",
				Value:  fmt.Sprintf("%d", len(player.queue)),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Requested by " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleStop(s *discordgo.Session, m *discordgo.MessageCreate, player *MusicPlayer) {
    player.mu.Lock()
    defer player.mu.Unlock()

    close(player.stopChan)
    player.stopChan = make(chan bool)
    player.queue = make([]Song, 0)
    if player.voiceConn != nil {
        player.voiceConn.Disconnect()
        player.voiceConn = nil
    }
    player.isPlaying = false
    
    voiceManager.ClearActivity(m.GuildID, MusicPlaying)

    embed := &discordgo.MessageEmbed{
        Title: "Stop ✅",
        Description: "⏹️ tfi lbolice jwan jay hh",
        Color: 0x00FF00,
        Footer: &discordgo.MessageEmbedFooter{
            Text: "Stopped by " + m.Author.Username,
            IconURL: m.Author.AvatarURL(""),
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func startGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	gameState := &GameState{
		RemCan: getAllCandidates(),
		Score:  make(map[string]int),
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
		Title:       "Ait Akinator",
		Description: question + "\n\n jawb b **/ait-akinator ah** or **/ait-akinator la**",
		Color:       0x00FF00,
	}
	s.ChannelMessageSendEmbed(channelID, embed)
}
func sendResult(s *discordgo.Session, channelID string, state *GameState) {
	var r string
	if len(state.RemCan) == 1 {
		r = "Howa ola machi howa: **@" + state.RemCan[0] + "**"

	} else {
		r = " AYkon ghi wa7d  min hado , skill issues (probelem fl'arbre ) : " + strings.Join(state.RemCan, ", ")
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Ezzz Ezzz",
		Description: r,
		Color:       0x00FF00,
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
func handleTTS(s *discordgo.Session, m *discordgo.MessageCreate, text string) {
    if !voiceManager.SetActivity(m.GuildID, TTSPlaying) {
        currentActivity := voiceManager.GetCurrentActivity(m.GuildID)
        var message string
        if currentActivity == MusicPlaying {
            message = "Bot mkhdm mzika, tsna hta ysali"
        } else {
            message = "Bot  tydwi , tsna hta ysali"
        }
        
        embed := &discordgo.MessageEmbed{
            Title: "3a9o bika !!!",
            Description: message,
            Color: 0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }
    defer voiceManager.ClearActivity(m.GuildID, TTSPlaying)

    vs, err := findUserVoiceState(s, m.GuildID, m.Author.ID)
    if err != nil {
        voiceManager.ClearActivity(m.GuildID, TTSPlaying)
        embed := &discordgo.MessageEmbed{
            Title: "Chi 7aja trat !!!",
            Description: "Khask tkon f voice channel bach nkhlih yhdr m3ak",
            Color: 0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }
	if _, err := exec.LookPath("espeak"); err != nil {
		embed := &discordgo.MessageEmbed{
			Title:       "Chi 7aja trat !!!",
			Description: "espeak machi installed f server",
			Color:       0xFF0000,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}
	vc, err := joinVoiceChannel(s, m.GuildID, vs.ChannelID)
	if err != nil {
		embed := &discordgo.MessageEmbed{
			Title:       "Chi 7aja trat !!!",
			Description: "Ma9dertch ndkhl l voice channel",
			Color:       0xFF0000,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}
	defer vc.Disconnect()
	tempFile, err := os.CreateTemp("", "tts-*.wav")
	if err != nil {
		logMsg("ERROR", fmt.Sprintf("Error creating temp file: %v", err))
		return
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	espeakCmd := exec.Command("espeak",
		"-v", "fr+f2",     
		"-s", "150",       
		"-p", "60",      
		"-a", "200",       	
		"-g", "10",       
		"-k", "5",         
		"-w", tempFile.Name(),
		text,
	)
	if output, err := espeakCmd.CombinedOutput(); err != nil {
		embed := &discordgo.MessageEmbed{
			Title:       "Chi 7aja trat !!!",
			Description: fmt.Sprintf("Error generating speech: %v\nOutput: %s", err, string(output)),
			Color:       0xFF0000,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}
	embed := &discordgo.MessageEmbed{
		Title:       "TTS 🗣️",
		Description: fmt.Sprintf("Speaking: %s", text),
		Color:       0x00FF00,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Requested by " + m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err := playAudioFile(vc, tempFile.Name()); err != nil {
		logMsg("ERROR", fmt.Sprintf("Error playing audio: %v", err))
	}
}

func playAudioFile(vc *discordgo.VoiceConnection, filename string) error {
	audioFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening audio file: %v", err)
	}
	defer audioFile.Close()
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpeg.Stdin = audioFile

	stdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}

	if err := ffmpeg.Start(); err != nil {
		return fmt.Errorf("error starting ffmpeg: %v", err)
	}
	defer ffmpeg.Process.Kill()
	encoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		return fmt.Errorf("error creating encoder: %v", err)
	}

	buffer := make([]int16, frameSize*channels)
	for {
		err := binary.Read(stdout, binary.LittleEndian, &buffer)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error reading from ffmpeg stdout: %v", err)
		}
		opus, err := encoder.Encode(buffer, frameSize, frameSize*2)
		if err != nil {
			return fmt.Errorf("error encoding to Opus: %v", err)
		}
		vc.OpusSend <- opus
	}
}

func handleClearChat(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
    if len(args) < 3 {
        s.ChannelMessageSend(m.ChannelID, "Khoya wach 7mar, zid chi number dyal messages li baghi tmsa7. ex: `/dh clear [number]`")
        return
    }
    num, err := strconv.Atoi(args[2])
    if err != nil {
        logMsg("ERROR", fmt.Sprintf("Error converting argument to integer: %v", err))
        s.ChannelMessageSend(m.ChannelID, "Khoya wach 7mar, khsk tdkhl wa7d number li howa positif o appartient à N (int).")
        return
    }
    if num <= 0 {
        s.ChannelMessageSend(m.ChannelID, "Khoya wach 7mar, khsk tdkhl wa7d number li howa positif o appartient à N (int).")
        return
    }

    if num > 100 {
        s.ChannelMessageSend(m.ChannelID, "Ma ymknch tmsa7 ktr mn 100 message f mra .")
        return
    }
    messages, err := s.ChannelMessages(m.ChannelID, num+1, "", "", "")
    if err != nil {
        logMsg("ERROR", fmt.Sprintf("Error retrieving messages: %v", err))
        return
    }
    var messageIDs []string
    for _, message := range messages {
        if message.ID != m.ID {
            messageIDs = append(messageIDs, message.ID)
        }
    }
    logMsg("INFO", fmt.Sprintf("Deleting %d messages", len(messageIDs)))
    err = s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
    if err != nil {
        logMsg("ERROR", fmt.Sprintf("Error deleting messages: %v", err))
        return
    }
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Msa7t %d dyal messages.", len(messageIDs)))
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

func findUserVoiceState(s *discordgo.Session, guildID, userID string) (*discordgo.VoiceState, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs, nil
		}
	}
	return nil, fmt.Errorf("user not in a voice channel")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !firstTime[m.Author.ID] {
		dm, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			logMsg("ERROR", fmt.Sprintf("Error creating DM channel: %v", err))
		} else {
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

	if strings.HasPrefix(m.Content, valoPing) {
        handlePing(s, m)
        return
    }

	args := strings.Split(m.Content, " ")
	if args[0] == ai && len(args) > 1 {
		if time.Since(LastReq) < 5*time.Second {
			embed := &discordgo.MessageEmbed{
				Title:       "Tsna a w9",
				Description: "Wach baghi ty7ni , ra tnmil mtnti7ch",
				Color:       0xFF0000,
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
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
		if args[1] == "dwi" {
			if len(args) < 3 {
				embed := &discordgo.MessageEmbed{
					Title:       "Chi 7aja trat !!!",
					Description: "3tini text bach n9rah ya had w9",
					Color:       0xFF0000,
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				return
			}
			text := strings.Join(args[2:], " ")
			handleTTS(s, m, text)
			return
		}
		if args[1] == "clear" {
			handleClearChat(s, m, args)
			return
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

	if args[0] == dh && len(args) > 1 {
		command := strings.ToLower(args[1])
		if command == "play" || command == "skip" || command == "stop" || command == "queue" {
			handleMusic(s, m, args)
			return
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
