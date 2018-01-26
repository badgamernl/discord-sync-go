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

type ConfigJSON struct {
	Token    string `json:"discord_bot_token"`
	Guild    string `json:"discord_guild_id"`
	IP       string `json:"rcon_ip"`
	Port     string `json:"rcon_port"`
	Password string `json:"rcon_password"`
	Roles    []Role `json:"roles"`
}

type Role struct {
	DiscordRoles []string `json:"discord"`
	Factorio     string   `json:"factorio"`
}

type Player struct {
	Name string `json:"username"`
	Role string `json:"role"`
}

// Variables used for command line parameters
var (
	ConfigFile string
	Config     ConfigJSON
	GoBot      *discordgo.Session
)

func init() {
	flag.StringVar(&ConfigFile, "c", "./config.json", "Config File location")
	flag.Parse()
	Config = loadConfiguration(ConfigFile)
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
	// Open a websocket connection to Discord and begin listening.
	err = GoBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Print the bot status
	log.Println("Discord bot is running and ready:", GoBot.State.User.Username)
	// Parse members
	members, err := members(GoBot)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Make command from pla
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
	elapsedRcon := time.Since(startRcon)
	log.Printf("Rcon ready: %s", elapsedRcon)
	startRconCommand := time.Now()
	response, err := r.Execute(command)
	if err != nil {
		panic(err)
	}
	log.Printf("Response: %+v\n", response)
	elapsedRconCommand := time.Since(startRconCommand)
	log.Printf("Rcon command send: %s", elapsedRconCommand)

	GoBot.Close()
	log.Println("Discord bot Closed")
	elapsed := time.Since(start)
	log.Printf("Process time: %s", elapsed)
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
	elapsed := time.Since(start)
	log.Printf("Member to command: %s", elapsed)
	return command
}

func members(s *discordgo.Session) ([]*Player, error) {
	start := time.Now()
	players := []*Player{}
	members, err := s.GuildMembers(Config.Guild, "", 1000)
	if err != nil {
		return players, err
	}
	elapsed := time.Since(start)
	log.Printf("Member get: %s", elapsed)
	start = time.Now()
	for _, member := range members {
		players = memberCheck(s, member, players)
	}
	elapsed = time.Since(start)
	log.Printf("Member parse: %s", elapsed)
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
