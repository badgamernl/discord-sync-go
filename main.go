package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/gtaylor/factorio-rcon"
)

// ConfigJSON -
type ConfigJSON struct {
	Token    string `json:"discord_bot_token"`
	Guild    string `json:"discord_guild_id"`
	IP       string `json:"rcon_ip"`
	Port     string `json:"rcon_port"`
	Password string `json:"rcon_password"`
	Roles    []Role `json:"roles"`
}

// Role -
type Role struct {
	DiscordRoles []string `json:"discord"`
	Factorio     string   `json:"factorio"`
}

// Player -
type Player struct {
	Name string `json:"username"`
	Role string `json:"role"`
}

// Variables used for command line parameters
var (
	ConfigFile string
	Config     ConfigJSON
)

func init() {
	flag.StringVar(&ConfigFile, "c", "./config.json", "Config File location")
	flag.Parse()
	start := time.Now()
	Config = loadConfiguration(ConfigFile)
	log.Printf("Config loaded: %s", time.Since(start))
}

func main() {
	start := time.Now()
	// Print the config
	log.Printf("Config Loaded:\n  Token: %v\n  GuildID: %v\n  Roles: %v\n", Config.Token, Config.Guild, strconv.Itoa(len(Config.Roles)))
	// Create a new Discord session using the provided bot token.
	GoBot, err := discordgo.New("Bot " + Config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	log.Printf("New Discord bot created: %s", time.Since(start))
	// Open a websocket connection to Discord and begin listening.
	err = GoBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	log.Printf("Discord bot is running and ready: %s - %v", time.Since(start), GoBot.State.User.Username)

	// Parse members
	members, err := members(GoBot)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Make command from players
	command := generateCommand(members)

	startRcon := time.Now()
	r, err := rcon.Dial(Config.IP + ":" + Config.Port)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	err = r.Authenticate(Config.Password)
	if err != nil {
		panic(err)
	}
	log.Printf("Rcon ready: %s", time.Since(startRcon))

	startRconCommand := time.Now()
	response, err := r.Execute(command)
	if err != nil {
		panic(err)
	}
	log.Printf("Response: %+v\n", response)
	log.Printf("Rcon command send: %s", time.Since(startRconCommand))

	GoBot.Close()
	log.Printf("Discord bot Closed & process end: %s", time.Since(start))
}

func loadConfiguration(file string) ConfigJSON {
	var config ConfigJSON
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func generateCommand(members []*Player) string {
	start := time.Now()
	var players []string
	for _, member := range members {
		players = append(players, `["`+member.Name+`"]="`+member.Role+`"`)
	}
	command := `/interface Ranking._base_preset{` + strings.Join(players, ",") + `}`
	log.Printf("Member to command: %s", time.Since(start))
	return command
}

func members(s *discordgo.Session) ([]*Player, error) {
	start := time.Now()
	players := []*Player{}
	members, err := s.GuildMembers(Config.Guild, "", 1000)
	if err != nil {
		return players, err
	}
	log.Printf("Member get: %s", time.Since(start))
	start = time.Now()
	for _, member := range members {
		players = memberCheck(s, member, players)
	}
	log.Printf("Member parse: %s", time.Since(start))
	return players, nil
}

func memberCheck(s *discordgo.Session, member *discordgo.Member, players []*Player) []*Player {
	if member.User.Bot == true {
		return players
	}
	username, err := memberName(member)
	if err != nil {
		return players
	}
	for _, role := range Config.Roles {
		for _, discordRole := range role.DiscordRoles {
			if memberHasRole(s, member, discordRole) {
				player := new(Player)
				player.Name = username
				player.Role = role.Factorio
				players = append(players, player)
				return players
			}
		}
	}
	return players
}

func memberHasRole(s *discordgo.Session, member *discordgo.Member, roleName string) bool {
	for _, roleID := range member.Roles {
		role, err := s.State.Role(Config.Guild, roleID)
		if err != nil {
			continue
		}
		if role.Name == roleName {
			return true
		}
	}
	return false
}

func memberName(member *discordgo.Member) (username string, err error) {
	validName := regexp.MustCompile(`^([a-zA-Z0-9]|\_|\.|\-)+$`)
	if member.Nick != "" {
		if validName.MatchString(member.Nick) {
			return member.Nick, nil
		}
	}
	if validName.MatchString(member.User.Username) {
		return member.User.Username, nil
	}
	return "", errors.New("Username not valid")
}
