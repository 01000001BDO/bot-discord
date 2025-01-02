package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"sort"
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
	swlCmd     string        = "/dh swl"
	nmapCmd    string        = "/dh scan "
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
    g1 = "👊"
    g2 = "✌️"
    g3 = "✋"
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
	Response = []string{
		"9witi 3liya  blas2ila",
		"fhmtk wakha dwiti mn trmtk",
		"miybi",
		"kayna had lhdri" ,
		"ah",
		"la",
		"owo",
		"salam ,  3jbni so2al dyalk o jatni lfikra ntsa7bo",
		"khoya sir gha t7wa ",
		"dwiti bzf",
		"mnkhdmch",
		"gha ttbz m3a krk",
		"oMMMMMMMar",
		"NONAYMOROZA",
	}
    playerScores = make(map[string]*PlayerScore) 
    scoresFile = "game_scores.json"
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


type ScanResult struct {
    URL           string
    IP            []string
    Hostname      string
    Technologies  []string
    Headers       map[string]string
    TLSInfo       *tls.ConnectionState
    OpenPorts     []int
    DNS           DNSInfo
    WebTech       WebTechnologies
}

type DNSInfo struct {
    MXRecords     []string
    TXTRecords    []string
    NSRecords     []string
    CNAMERecords  []string
}

type WebTechnologies struct {
    Frameworks    []string
    Libraries     []string
    BuildTools    []string
    UILibraries   []string
    Analytics     []string
    Deployment    []string
}

type SecurityScanner struct {
    client        *http.Client
    mutex         sync.Mutex
    techPatterns  map[string]map[string][]string
}

type PlayerScore struct {
    UserID    string `json:"user_id"`
    Username  string `json:"username"`
    Wins      int    `json:"wins"`
    Losses    int    `json:"losses"`
    Draws     int    `json:"draws"`
    LastPlayed time.Time `json:"last_played"`
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
func loadScores() error {
    data, err := os.ReadFile(scoresFile)
    if err != nil {
        if os.IsNotExist(err) {
            return nil 
        }
        return err
    }

    return json.Unmarshal(data, &playerScores)
}

func saveScores() error {
    data, err := json.MarshalIndent(playerScores, "", "    ")
    if err != nil {
        return err
    }

    return os.WriteFile(scoresFile, data, 0644)
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

func NewSecurityScanner() *SecurityScanner {
    return &SecurityScanner{
        client: &http.Client{
            Timeout: 15 * time.Second,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    InsecureSkipVerify: true, 
                },
            },
        },
        techPatterns: initTechPatterns(),
    }
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
    embed := &discordgo.MessageEmbed{
        Title:       "🤖 Commands dyal L'bot",
        Description: "List dyal Commands li kaynin fl bot:",
        Color:       0x00FF00,
        Fields: []*discordgo.MessageEmbedField{
            {
                Name: "🎵 Music Commands",
                Value: "• `/dh play [url]` - play lchi track\n" +
                    "• `/dh playlist [url]` - play  playlist (fiha mochkil lakan jhd n9adha ) \n" +
                    "• `/dh skip` - Bach tskipi track\n" +
                    "• `/dh stop` - Bach tw9f music\n" +
                    "• `/dh queue` - Bach tchouf playlist",
                Inline: false,
            },
            {
                Name: "🎮 Game Commands",
                Value: "• `/ait-akinator` -  ait Akinator \n" +
                    "• `/ait-akinator (ah/la)` - Jawb 3la les questions\n" +
                    "• `/dh game` - L3b 7ajara wara9 mi9ass m3a lbot\n" +
                    "   - 👊 = 7ajara\n" +
                    "  - ✌️ = wra9a\n" +
                    "  - ✋ = m9ass\n" +
                    "  - 🔄 = rematch\n" +
                    "  - ❌ = exit game\n" +
                    "  - 📊 = leaderboard",
                Inline: false,
            },
            {
				Name: "🎯 Valorant Command",
				Value: "• `/dh valo-ping` - Tchouf ping dyal servers\n" +
					   "• `/dh valo-ping [ip]` - Tchouf ping dyal IP dyalk ",
				Inline: false,
			},
            {
                Name: "🤖 AI Commands",
                Value: "• `/ai [prompt]` - ask  Gemini",
                Inline: false,
            },
            {
                Name: "🗣️ Voice Commands",
                Value: "• `/dh dwi [text]` - Bot ghayi9ra text li ktbti",
                Inline: false,
            },
            {
                Name: "⚙️ Utility Commands",
                Value: "• `/dh clear [number]` - Msa7 messages\n" +
                    "• `/dh latence` - Tchouf latency dyal bot\n" +
                    "• `/dh ls` - Tchouf channels li kaynin\n" +
                    "• `/dh pwd` - Tchouf current path",
                Inline: false,
            },
			{
				Name: "🎱 Magic 8-Ball",
				Value: "• `/dh swl [question]` - swl lbot ijawbk ",
				Inline: false,
			},
			{
				Name: "🔍 Security Scanner",
				Value: "• `/dh scan [url]` - Scan website (gathering information public : IP/Tech/ports ..)\n" +
					"• Example: `/dh scan diamondhands.com`\n" +
					"**NOTE:** Matkhdmch l scan 3la sites li machi dyalk.",
				Inline: false,
			},
        },
        Footer: &discordgo.MessageEmbedFooter{
            Text: "Written in Go  by  @aka_bousta",
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleMagic8Ball(s *discordgo.Session, m *discordgo.MessageCreate) {
    question := strings.TrimPrefix(m.Content, swlCmd)
    question = strings.TrimSpace(question)

    if question == "" {
        embed := &discordgo.MessageEmbed{
            Title:       "Chi 7aja trat !!!",
            Description: "Gha ttbz m3a krk , Khassek tktb chi so2al.",
            Color:       0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }

    rand.Seed(time.Now().UnixNano())
    response := Response[rand.Intn(len(Response))]
    embed := &discordgo.MessageEmbed{
        Title:       "🎱 " + question,
        Description: "**" + response + "**",
        Color:       0x9B59B6, 
        Footer: &discordgo.MessageEmbedFooter{
            Text:    "Asked by " + m.Author.Username,
            IconURL: m.Author.AvatarURL(""),
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
}


func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
    args := strings.Split(m.Content, " ")
    isCustomIP := len(args) > 2 && args[2] != ""
    
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
    if isCustomIP {
        customIP := args[2]
        cmd := exec.Command("ping", "-c", "3", customIP)
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
            Name:     "Custom IP",
            Location: "User Specified",
            Ping:     math.Round(pingTime*100) / 100,
        })
    } else {
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
    }
    title := "🌐 Valorant EU Server Pings"
    if isCustomIP {
        title = "🌐 Custom IP Ping Results"
    }

    embed = &discordgo.MessageEmbed{
        Title:       title,
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


func initTechPatterns() map[string]map[string][]string {
    return map[string]map[string][]string{
        "Frameworks": {
            "Next.js": {
                `"__NEXT_DATA__"`,
                `/_next/static`,
                `next/dist/pages/_app`,
            },
            "React": {
                `react.development.js`,
                `react.production.min.js`,
                `__REACT_DEVTOOLS_GLOBAL_HOOK__`,
            },
            "Vue.js": {
                `vue.js`,
                `vue.min.js`,
                `__vue__`,
            },
            "Angular": {
                `ng-version`,
                `angular.js`,
                `angular.min.js`,
            },
        },
        "BuildTools": {
            "Webpack": {
                `webpackJsonp`,
                `__webpack_require__`,
            },
            "Vite": {
                `@vite/client`,
                `vite/dist`,
            },
            "Parcel": {
                `parcelRequire`,
            },
        },
        "UILibraries": {
            "Tailwind CSS": {
                `tailwind`,
                `tw-`,
            },
            "Material-UI": {
                `MuiButton`,
                `MuiTypography`,
            },
            "Bootstrap": {
                `bootstrap.min.css`,
                `bootstrap.bundle.js`,
            },
        },
        "Analytics": {
            "Google Analytics": {
                `google-analytics.com`,
                `gtag`,
                `ga.js`,
            },
            "Plausible": {
                `plausible.io`,
            },
        },
		"CloudServices": {
			"Vercel": {
				`vercel.app`,
				`vercel-analytics`,
				`vercel.com`,
			},
			"Netlify": {
				`netlify.app`,
				`netlify.com`,
				`netlify-headers`,
			},
			"Cloudflare": {
				`cloudflare`,
				`__cf_bm`,
				`cf-ray`,
			},
			"AWS": {
				`amazonaws.com`,
				`aws-amplify`,
				`x-amz-`,
			},
			"GitHub Pages": {
				`github.io`,
				`githubusercontent`,
			},
			"Firebase": {
				`firebaseapp.com`,
				`firebase-hosting`,
			},
		},
    }
}

func (s *SecurityScanner) ScanTarget(target string) (*ScanResult, error) {
    if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
        target = "https://" + target
    }

    result := &ScanResult{
        URL:     target,
        Headers: make(map[string]string),
        WebTech: WebTechnologies{},
    }

    host := strings.TrimPrefix(strings.TrimPrefix(target, "https://"), "http://")
    host = strings.Split(host, "/")[0]
    ips, err := net.LookupIP(host)
    if err == nil {
        for _, ip := range ips {
            result.IP = append(result.IP, ip.String())
        }
    }
    if len(result.IP) > 0 {
        hostnames, err := net.LookupAddr(result.IP[0])
        if err == nil && len(hostnames) > 0 {
            result.Hostname = hostnames[0]
        }
    }
    result.DNS = s.getDNSInfo(host)
    if len(result.IP) > 0 {
        result.OpenPorts = s.scanPorts(result.IP[0])
    }
    req, err := http.NewRequest("GET", target, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Set("User-Agent", "SecurityScanner/1.0")
    resp, err := s.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading body: %v", err)
    }
    for name, values := range resp.Header {
        result.Headers[name] = strings.Join(values, ", ")
    }
    result.Technologies = s.detectTechnologies(resp.Header, body)
	cloudServices := s.detectCloudServices(target, resp.Header)
	result.Technologies = append(result.Technologies, cloudServices...)
    if resp.TLS != nil {
        result.TLSInfo = resp.TLS
    }

    return result, nil
}


func (s *SecurityScanner) detectTechnologies(headers http.Header, body []byte) []string {
    var technologies []string
    bodyStr := string(body)
    if server := headers.Get("Server"); server != "" {
        technologies = append(technologies, "Server: "+server)
    }
    if powered := headers.Get("X-Powered-By"); powered != "" {
        technologies = append(technologies, "Powered By: "+powered)
    }
    for category, techs := range s.techPatterns {
        for techName, patterns := range techs {
            for _, pattern := range patterns {
                if strings.Contains(bodyStr, pattern) {
                    technologies = append(technologies, fmt.Sprintf("%s: %s", category, techName))
                    break
                }
            }
        }
    }
    scriptRegex := regexp.MustCompile(`<script[^>]*src=['"](.*?)['"][^>]*>`)
    scripts := scriptRegex.FindAllStringSubmatch(bodyStr, -1)
    for _, script := range scripts {
        if len(script) > 1 {
            technologies = append(technologies, "Script: "+script[1])
        }
    }

    return unique(technologies)
}

func (s *SecurityScanner) detectCloudServices(url string, headers http.Header) []string {
    var services []string
    for cloudName, patterns := range s.techPatterns["CloudServices"] {
        for _, pattern := range patterns {
            if strings.Contains(url, pattern) {
                services = append(services, fmt.Sprintf("Cloud Platform: %s", cloudName))
                break
            }
        }
    }
    cloudHeaders := map[string]string{
        "x-vercel-id":        "Vercel",
        "x-netlify":          "Netlify",
        "cf-ray":             "Cloudflare",
        "x-github-request":   "GitHub Pages",
        "x-firebase-hosting": "Firebase",
    }

    for header, platform := range cloudHeaders {
        if value := headers.Get(header); value != "" {
            services = append(services, "Cloud Platform: "+platform)
        }
    }
    return unique(services)
}

func (s *SecurityScanner) scanPorts(ip string) []int {
    commonPorts := []int{80, 443, 8080, 8443, 3000, 4000, 5000}
    var openPorts []int
    var wg sync.WaitGroup
    var mutex sync.Mutex

    for _, port := range commonPorts {
        wg.Add(1)
        go func(p int) {
            defer wg.Done()
            address := fmt.Sprintf("%s:%d", ip, p)
            conn, err := net.DialTimeout("tcp", address, 2*time.Second)
            if err == nil {
                mutex.Lock()
                openPorts = append(openPorts, p)
                mutex.Unlock()
                conn.Close()
            }
        }(port)
    }
    wg.Wait()
    return openPorts
}

func (s *SecurityScanner) getDNSInfo(domain string) DNSInfo {
    var info DNSInfo

    if mxRecords, err := net.LookupMX(domain); err == nil {
        for _, mx := range mxRecords {
            info.MXRecords = append(info.MXRecords, mx.Host)
        }
    }

    if txtRecords, err := net.LookupTXT(domain); err == nil {
        info.TXTRecords = txtRecords
    }

    if nsRecords, err := net.LookupNS(domain); err == nil {
        for _, ns := range nsRecords {
            info.NSRecords = append(info.NSRecords, ns.Host)
        }
    }

    if cname, err := net.LookupCNAME(domain); err == nil {
        info.CNAMERecords = append(info.CNAMERecords, cname)
    }

    return info
}

func unique(slice []string) []string {
    keys := make(map[string]bool)
    var list []string
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}
func handleEnhancedScan(s *discordgo.Session, m *discordgo.MessageCreate, target string) {
    scanner := NewSecurityScanner()
    
    embed := &discordgo.MessageEmbed{
        Title:       "🔍 Scanning Target",
        Description: fmt.Sprintf("Scanning %s for detailed information...", target),
        Color:       0xFFFF00,
    }
    msg, _ := s.ChannelMessageSendEmbed(m.ChannelID, embed)

    result, err := scanner.ScanTarget(target)
    if err != nil {
        errorEmbed := &discordgo.MessageEmbed{
            Title:       "❌ Scan Faila ahbibi hh (skill issue ola probleme flink )",
            Description: fmt.Sprintf("Error scanning target: %v", err),
            Color:       0xFF0000,
        }
        s.ChannelMessageEditEmbed(m.ChannelID, msg.ID, errorEmbed)
        return
    }

    var description strings.Builder
    description.WriteString("**🌐 Network Information**\n")
    if len(result.IP) > 0 {
        description.WriteString(fmt.Sprintf("• IP Addresses: %s\n", strings.Join(result.IP, ", ")))
    }
    if result.Hostname != "" {
        description.WriteString(fmt.Sprintf("• Hostname: %s\n", result.Hostname))
    }
    if len(result.OpenPorts) > 0 {
        description.WriteString("\n**🔌 Open Ports**\n")
        for _, port := range result.OpenPorts {
            description.WriteString(fmt.Sprintf("• %d\n", port))
        }
    }
    if len(result.Technologies) > 0 {
        description.WriteString("\n**💻 Technologies Detected**\n")
        for _, tech := range result.Technologies {
            description.WriteString(fmt.Sprintf("• %s\n", tech))
        }
    }
    if len(result.DNS.MXRecords) > 0 || len(result.DNS.NSRecords) > 0 {
        description.WriteString("\n**📡 DNS Information**\n")
        if len(result.DNS.MXRecords) > 0 {
            description.WriteString(fmt.Sprintf("• MX Records: %s\n", strings.Join(result.DNS.MXRecords, ", ")))
        }
        if len(result.DNS.NSRecords) > 0 {
            description.WriteString(fmt.Sprintf("• NS Records: %s\n", strings.Join(result.DNS.NSRecords, ", ")))
        }
    }
    if result.TLSInfo != nil {
        description.WriteString("\n**🔒 SSL/TLS Information**\n")
        description.WriteString(fmt.Sprintf("• Version: %s\n", getTLSVersion(result.TLSInfo.Version)))
        description.WriteString(fmt.Sprintf("• Cipher Suite: %v\n", result.TLSInfo.CipherSuite))
    }

    resultEmbed := &discordgo.MessageEmbed{
        Title:       "🎯 Scan Results , 3ich ahbibi",
        Description: description.String(),
        Color:       0x00FF00,
        Footer: &discordgo.MessageEmbedFooter{
            Text: "Mtn7mlch lms2oliya dyal scan  , une fois drti scan all informations are gone  ",
        },
    }
    s.ChannelMessageEditEmbed(m.ChannelID, msg.ID, resultEmbed)
}

func getTLSVersion(version uint16) string {
    versions := map[uint16]string{
        tls.VersionTLS10: "TLS 1.0",
        tls.VersionTLS11: "TLS 1.1",
        tls.VersionTLS12: "TLS 1.2",
        tls.VersionTLS13: "TLS 1.3",
    }
    if v, ok := versions[version]; ok {
        return v
    }
    return "Unknown"
}

func handleGameReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
    if r.UserID == s.State.User.ID {
        return
    }
    if _, exists := playerScores[r.UserID]; !exists {
        user, _ := s.User(r.UserID)
        username := "Unknown"
        if user != nil {
            username = user.Username
        }
        playerScores[r.UserID] = &PlayerScore{
            UserID: r.UserID,
            Username: username,
            Wins: 0,
            Losses: 0,
            Draws: 0,
        }
    }
    playerScore := playerScores[r.UserID]
    if r.Emoji.Name == "📊" {
        showLeaderboard(s, r.ChannelID)
        return
    }
    if r.Emoji.Name == "🔄" {
        embed := &discordgo.MessageEmbed{
            Title: "🎮 7jar Wra9 M9ass",
            Description: "Mchina laysm7 lina?\n\n" +
                "👊 = 7ajara\n" +
                "✌️ = wra9a\n" +
                "✋ = m9ass\n\n" +
                "Click 3la emoji bach tl3b!\n\n" +
                "🔄 = rematch\n" +
                "❌ = exit game\n" +
                "📊 = leaderboard",
            Color: 0x00FF00,
        }
       
        newMsg, _ := s.ChannelMessageSendEmbed(r.ChannelID, embed)
        s.MessageReactionAdd(newMsg.ChannelID, newMsg.ID, "👊")
        s.MessageReactionAdd(newMsg.ChannelID, newMsg.ID, "✌️")
        s.MessageReactionAdd(newMsg.ChannelID, newMsg.ID, "✋")
        s.MessageReactionAdd(newMsg.ChannelID, newMsg.ID, "🔄")
        s.MessageReactionAdd(newMsg.ChannelID, newMsg.ID, "❌")
        s.MessageReactionAdd(newMsg.ChannelID, newMsg.ID, "📊")
        return
    } else if r.Emoji.Name == "❌" {
        showPlayerScore(s, r.ChannelID, r.UserID)
        s.ChannelMessageSendEmbed(r.ChannelID, &discordgo.MessageEmbed{
            Title: "👋 yawdi yawdi z3ma rijal",
            Description: "lay9tlna m3a rjal",
            Color: 0xFF0000,
        })
        return
    }

    choices := map[string]string{
        "👊": "7ajara",
        "✌️": "wra9",
        "✋": "m9ass",
    }
 
    if _, isGameMove := choices[r.Emoji.Name]; !isGameMove {
        return
    }
 
    botEmojis := []string{"👊", "✌️", "✋"}
    botChoice := botEmojis[rand.Intn(len(botEmojis))]
    userChoice := r.Emoji.Name
   
    var result string
    var color int
    if userChoice == botChoice {
        result = "wa 3rgan 😐"
        color = 0xFFFF00
        playerScore.Draws++
    } else if (userChoice == "👊" && botChoice == "✌️") ||
        (userChoice == "✌️" && botChoice == "✋") ||
        (userChoice == "✋" && botChoice == "👊") {
        result = "9lawi 3La ji3an , atzid ola atghryha b7al l9hyba hh?"
        color = 0x00FF00
        playerScore.Wins++
    } else {
        result = "olyaaaa idk fzeb rb7t, li drbna 9a3o myhmnach sda3o"
        color = 0xFF0000
        playerScore.Losses++
    }
    playerScore.LastPlayed = time.Now()
    if err := saveScores(); err != nil {
        logMsg("ERROR", fmt.Sprintf("Error saving scores: %v", err))
    }
 
    embed := &discordgo.MessageEmbed{
        Title: "🎮 Score finale",
        Fields: []*discordgo.MessageEmbedField{
            {
                Name: "Nta:",
                Value: choices[userChoice],
                Inline: true,
            },
            {
                Name: "Ana:",
                Value: choices[botChoice],
                Inline: true,
            },
            {
                Name: "Score dyalk:",
                Value: fmt.Sprintf("Rbe7ti: %d\nKhserti: %d\nT3adol: %d\nWin Rate: %.1f%%",
                    playerScore.Wins, playerScore.Losses, playerScore.Draws,
                    calculateWinRate(playerScore)),
                Inline: false,
            },
        },
        Description: result + "\n\n🔄 = rematch\n❌ = exit game\n📊 = leaderboard",
        Color: color,
        Footer: &discordgo.MessageEmbedFooter{
            Text: "Zayd ola nayd",
        },
    } 
    resultMsg, _ := s.ChannelMessageSendEmbed(r.ChannelID, embed)
    s.MessageReactionAdd(resultMsg.ChannelID, resultMsg.ID, "🔄")
    s.MessageReactionAdd(resultMsg.ChannelID, resultMsg.ID, "❌")
    s.MessageReactionAdd(resultMsg.ChannelID, resultMsg.ID, "📊")
}

 func showPlayerScore(s *discordgo.Session, channelID string, userID string) {
    if score, exists := playerScores[userID]; exists {
        user, _ := s.User(userID)
        username := "Unknown"
        if user != nil {
            username = user.Username
        }
        embed := &discordgo.MessageEmbed{
            Title: "📊 Score dyal " + username,
            Description: fmt.Sprintf(
                "**Rbe7ti:** %d\n**Khserti:** %d\n**T3adol:** %d\n**Win Rate:** %.1f%%",
                score.Wins, score.Losses, score.Draws,
                calculateWinRate(score),
            ),
            Color: 0x00FF00,
        }
        s.ChannelMessageSendEmbed(channelID, embed)
    }
}
func calculateWinRate(score *PlayerScore) float64 {
    total := score.Wins + score.Losses + score.Draws
    if total == 0 {
        return 0
    }
    return float64(score.Wins) / float64(total) * 100
}

func showLeaderboard(s *discordgo.Session, channelID string) {
    var players []struct {
        UserID   string
        Username string
        Score    *PlayerScore
        WinRate  float64
    }
    for userID, score := range playerScores {
        user, _ := s.User(userID)
        username := "Unknown"
        if user != nil {
            username = user.Username
        }
        players = append(players, struct {
            UserID   string
            Username string
            Score    *PlayerScore
            WinRate  float64
        }{
            UserID:   userID,
            Username: username,
            Score:    score,
            WinRate:  calculateWinRate(score),
        })
    }
    sort.Slice(players, func(i, j int) bool {
        return players[i].WinRate > players[j].WinRate
    })
    var description strings.Builder
    description.WriteString("🏆 Top Players:\n\n")
    
    for i, player := range players {
        if i >= 10 {
            break
        }
        description.WriteString(fmt.Sprintf(
            "**%d.** %s\n├ Rbe7: %d | Khser: %d | T3adol: %d\n└ Win Rate: %.1f%%\n\n",
            i+1, player.Username,
            player.Score.Wins, player.Score.Losses, player.Score.Draws,
            player.WinRate,
        ))
    }

    embed := &discordgo.MessageEmbed{
        Title:       "🏆 Leaderboard",
        Description: description.String(),
        Color:       0xFFD700,
    }
    s.ChannelMessageSendEmbed(channelID, embed)
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
		url := args[2]
        if strings.Contains(url, "playlist?list=") {
            handlePlaylist(s, m, vs.ChannelID, url, player)
        } else {
            handlePlay(s, m, vs.ChannelID, url, player)
        }

    case "playlist":
        if len(args) < 3 {
            embed := &discordgo.MessageEmbed{
                Title:       "Chi 7aja trat !!!",
                Description: "Hbibi lien dyal playlist",
                Color:       0xFF0000,
            }
            s.ChannelMessageSendEmbed(m.ChannelID, embed)
            return
        }
        handlePlaylist(s, m, vs.ChannelID, args[2], player)

    case "skip":
        handleSkip(s, m, player)

    case "stop":
        handleStop(s, m, player)

    case "queue":
        handleQueue(s, m, player)
    }
}

// !TODO : mkhdamch had l9lawi : (khsni nfixi probeleme dyal link ) 
func handlePlaylist(s *discordgo.Session, m *discordgo.MessageCreate, voiceChannelID string, url string, player *MusicPlayer) {
    if !voiceManager.SetActivity(m.GuildID, MusicPlaying) {
        currentActivity := voiceManager.GetCurrentActivity(m.GuildID)
        message := "Bot mkhdm mzika, tsna hta ysali"
        if currentActivity == TTSPlaying {
            message = "Bot tydwi, tsna hta ysali"
        }
        embed := &discordgo.MessageEmbed{
            Title:       "3a9o bika",
            Description: message,
            Color:       0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }
    loadingEmbed := &discordgo.MessageEmbed{
        Title:       "Loading Playlist 🎵",
        Description: "Tsna chwiya, kanjib l playlist...",
        Color:       0xFFFF00,
    }
    msg, err := s.ChannelMessageSendEmbed(m.ChannelID, loadingEmbed)
    if err != nil {
        fmt.Println("Error sending loading message:", err)
        return
    }
    client := youtube.Client{}
    playlist, err := client.GetPlaylist(url)
    if err != nil {
        voiceManager.ClearActivity(m.GuildID, MusicPlaying)
        errorEmbed := &discordgo.MessageEmbed{
            Title:       "Chi 7aja trat !!!",
            Description: fmt.Sprintf("Error ma9drtch njib playlist: %v", err),
            Color:       0xFF0000,
        }
        s.ChannelMessageEditEmbed(m.ChannelID, msg.ID, errorEmbed)
        return
    }

    player.mu.Lock()
    defer player.mu.Unlock()
    addedSongs := 0
    for _, entry := range playlist.Videos {
        video, err := client.GetVideo(fmt.Sprintf("https://www.youtube.com/watch?v=%s", entry.ID))
        if err != nil {
            continue
        }
        var format youtube.Format
        for _, f := range video.Formats {
            if f.AudioChannels > 0 {
                format = f
                break
            }
        }

        if format == (youtube.Format{}) {
            continue
        }

        song := Song{
            URL:      format.URL,
            Title:    video.Title,
            Duration: video.Duration.String(),
        }

        player.queue = append(player.queue, song)
        addedSongs++
    }

    resultEmbed := &discordgo.MessageEmbed{
        Title:       "Playlist Added ✅",
        Description: fmt.Sprintf("Zadt **%d** songs mn playlist\n**%s**", addedSongs, playlist.Title),
        Color:       0x00FF00,
        Footer: &discordgo.MessageEmbedFooter{
            Text:    "Added by " + m.Author.Username,
            IconURL: m.Author.AvatarURL(""),
        },
    }
    s.ChannelMessageEditEmbed(m.ChannelID, msg.ID, resultEmbed)
    if !player.isPlaying {
        go startPlaying(s, m.GuildID, voiceChannelID, player)
    }
}

func handlePlay(s *discordgo.Session, m *discordgo.MessageCreate, voiceChannelID string, url string, player *MusicPlayer) {
    if !voiceManager.SetActivity(m.GuildID, MusicPlaying) {
        currentActivity := voiceManager.GetCurrentActivity(m.GuildID)
        message := "Bot mkhdm mzika, tsna hta ysali"
        if currentActivity == TTSPlaying {
            message = "Bot tydwi, tsna hta ysali"
        }

        embed := &discordgo.MessageEmbed{
            Title:       "3a9o bika",
            Description: message,
            Color:       0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }

    player.mu.Lock()
    defer player.mu.Unlock()

    client := youtube.Client{}
    video, err := client.GetVideo(url)
    if err != nil {
        voiceManager.ClearActivity(m.GuildID, MusicPlaying)
        embed := &discordgo.MessageEmbed{
            Title:       "Chi 7aja trat !!!",
            Description: fmt.Sprintf("Error ma9drtch njib video: %v", err),
            Color:       0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        return
    }

    var format youtube.Format
    for _, f := range video.Formats {
        if f.AudioChannels > 0 {
            format = f
            break
        }
    }

    song := Song{
        URL:      format.URL,
        Title:    video.Title,
        Duration: video.Duration.String(),
    }

    player.queue = append(player.queue, song)

    embed := &discordgo.MessageEmbed{
        Title:       "Tzadt fl Queue",
        Description: fmt.Sprintf("🎵 **%s**", song.Title),
        Color:       0x00FF00,
        Footer: &discordgo.MessageEmbedFooter{
            Text:    "Added by " + m.Author.Username,
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
            player.mu.Unlock()
            return
        }

        currentSong := player.queue[0]
        player.queue = player.queue[1:]
        player.mu.Unlock()

        if player.voiceConn == nil {
            vc, err := joinVoiceChannel(s, guildID, voiceChannelID)
            if err != nil {
                fmt.Println("Error joining voice channel:", err)
                continue
            }
            player.voiceConn = vc
        }

        ffmpeg := exec.Command("ffmpeg", "-i", currentSong.URL, 
            "-f", "s16le", 
            "-ar", "48000", 
            "-ac", "2",
            "-af", "volume=0.5",
            "pipe:1")
            
        ffmpeg.Stderr = nil
        stdout, err := ffmpeg.StdoutPipe()
        if err != nil {
            fmt.Println("Error creating stdout pipe:", err)
            continue
        }

        err = ffmpeg.Start()
        if err != nil {
            fmt.Println("Error starting ffmpeg:", err)
            continue
        }

        encoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
        if err != nil {
            fmt.Println("Error creating opus encoder:", err)
            continue
        }

        embed := &discordgo.MessageEmbed{
            Title:       "Sm3 Sm3 🎵",
            Description: fmt.Sprintf("**%s**", currentSong.Title),
            Color:       0x00FF00,
        }
        s.ChannelMessageSendEmbed(voiceChannelID, embed)
        audioPCM := make([]int16, frameSize*channels)
        for {
            err := binary.Read(stdout, binary.LittleEndian, &audioPCM)
            if err != nil {
                if err != io.EOF {
                    fmt.Println("Error reading from ffmpeg stdout:", err)
                }
                break
            }

			opusData, err := encoder.Encode(audioPCM, frameSize, frameSize*2)
			if err != nil {
				fmt.Println("Error encoding to opus:", err)
				continue
			}
			
			select {
			case player.voiceConn.OpusSend <- opusData:
            case <-player.stopChan:
                ffmpeg.Process.Kill()
                return
            }
        }

        ffmpeg.Wait()
        time.Sleep(200 * time.Millisecond)
    }
}


func handleSkip(s *discordgo.Session, m *discordgo.MessageCreate, player *MusicPlayer) {
    player.mu.Lock()
    if len(player.queue) == 0 && !player.isPlaying {
        embed := &discordgo.MessageEmbed{
            Title:       "Chi 7aja trat !!!",
            Description: "Queue khawya akhoya",
            Color:       0xFF0000,
        }
        s.ChannelMessageSendEmbed(m.ChannelID, embed)
        player.mu.Unlock()
        return
    }
    if player.isPlaying {
        close(player.stopChan)
        player.stopChan = make(chan bool)
        player.isPlaying = false
    }

    embed := &discordgo.MessageEmbed{
        Title:       "Skip ✅",
        Description: "⏭️ tskipat a hbibi",
        Color:       0x00FF00,
        Footer: &discordgo.MessageEmbedFooter{
            Text:    "Skipped by " + m.Author.Username,
            IconURL: m.Author.AvatarURL(""),
        },
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed)
    if len(player.queue) > 0 && !player.isPlaying {
        vs, err := findUserVoiceState(s, m.GuildID, m.Author.ID)
        if err == nil {
            go startPlaying(s, m.GuildID, vs.ChannelID, player)
        }
    }
    player.mu.Unlock()
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
        Title:       "Stop ✅",
        Description: "⏹️ tfi lbolice jwan jay hh",
        Color:       0x00FF00,
        Footer: &discordgo.MessageEmbedFooter{
            Text:    "Stopped by " + m.Author.Username,
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
    switch strings.ToLower(strings.TrimSpace(q)) {
    case "salam", "slm", "slt":
        return "9lawi 3la hdra dyal facebook , chri m3ak lbitcoin "
    case "fin saken":
        return "fkrk , ada7k ana saken fwa7d pc 7achak hh"
    case "chkoun nta" : 
        return "li 7wak hh"
    case "katdwi ghi bdarija": 
        return "tndwi bkolchi walakin wlad l97ab frdo 3lina had darija"
    case "ach tyban lik flmodawana jdida":
        return " ra mb9itch ga3 9ad njawb bsbaha akhoya"
    case "kolchi mzyan ?":
        return "swlni ki jatni lwalida 3ndk b3da hh"
    case "ach tyban lik f m6":   
        return "fiha fasl 7adari "    
    }
    s.ChannelTyping(m.ChannelID)
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
    if err != nil {
        logMsg("ERROR", fmt.Sprintf("Error creating Gemini client: %v", err))
        return "kayn chi mochkil (ima skill issues ola chi 7aja khra)"     
    }
    defer client.Close()
    darija_prompt := `Respond naturally in Moroccan Darija (using Latin script/ darija ) to the following question. 
    Give only a single, direct response without explanations or translations. Keep it casual and conversational: ` + q
    model := client.GenerativeModel("gemini-1.5-flash")
    resp, err := model.GenerateContent(ctx, genai.Text(darija_prompt))
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
	if strings.HasPrefix(m.Content, swlCmd) {
		handleMagic8Ball(s, m)
		return
	}

	if strings.HasPrefix(m.Content, nmapCmd) {
		args := strings.Split(m.Content, " ")
		if len(args) < 3 {
			embed := &discordgo.MessageEmbed{
				Title:       "Error",
				Description: "Please provide a website to scan. Usage: `/dh scan example.com`",
				Color:       0xFF0000,
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
		handleEnhancedScan(s, m, args[2])
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
		if args[1] == "help" {
			handleHelp(s, m)
			return
		}
        if args[1] == "game" {
            embed := &discordgo.MessageEmbed{
                Title:       "🎮 7ajara Wara9a Mi9ass",
                Description: "Mchina laysm7 lina?\n\n" +
                    g1 + " = 7ajara\n" +
                    g2 + " = wra9a\n" +
                    g3 + " = mi9ass\n\n" +
                    "Click 3la emoji bach tl3b!\n\n" +
                    "🔄 = reamtch \n" + 
                    "❌ = exit game",
                Color:       0x00FF00,
                Footer: &discordgo.MessageEmbedFooter{
                    Text:    "Gwima m3a hbibna :  @" + m.Author.Username,
                    IconURL: m.Author.AvatarURL(""),
                },
            }
           
            msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
            if err != nil {
                logMsg("ERROR", fmt.Sprintf("Error sending game message: %v", err))
                return
            }
            s.MessageReactionAdd(msg.ChannelID, msg.ID, g1)
            s.MessageReactionAdd(msg.ChannelID, msg.ID, g2)
            s.MessageReactionAdd(msg.ChannelID, msg.ID, g3)
            s.MessageReactionAdd(msg.ChannelID, msg.ID, "🔄")
            s.MessageReactionAdd(msg.ChannelID, msg.ID, "❌")
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
    loadScores()
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
    s.AddHandler(handleGameReaction)
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
