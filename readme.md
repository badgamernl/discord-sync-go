# discord-sync-go
This repository is made for explosivegaming to sync the Discord player roles to the custom Factorio [explosivegaming-main](https://github.com/badgamernl/explosivegaming-main) scenario made by [Cooldude2606](https://github.com/Cooldude2606).

# Installation
* Create a directory where you want to store this script along with configuration. (either copy-paste the files or clone from github):
```
$ cd /opt/
$ git clone https://github.com/badgamernl/discord-sync-go.git
$ 
```
* Install the go packages:
```
$ go get github.com/bwmarrin/discordgo
$ go get github.com/gtaylor/factorio-rcon
```
* Build the project:
```
$ go build .
```
* Rename config.example.json to config.json and modify the values within according to your setup.

# Usage
Runs sync with the default config `./config.json`.

To change the config location add `-c /config/file/config.json`

## Dev
```
$ go run main.go
```
## Prod
```
$ ./discord-sync-go
```

# Config
The config consists of all the information for the Discord bot to login and the RCON to connect to the headless server.
```json
{
  "discord_bot_token": "TheBotToken",
  "discord_guild_id": "TheDiscordGuildIDiWantToFetchMembersFrom",
  "rcon_ip": "localhost",
  "rcon_port": "12345",
  "rcon_password": "myrconpassword",
  "roles": [
    {
      "discord": ["Owners"],
      "factorio": "Owner"
    },
    {
      "discord": ["Developers","Factorio-Team"],
      "factorio": "Developer"
    },
    {
      "discord": ["Admins"],
      "factorio": "Admin"
    },
    {
      "discord": ["Factorio-Moderators"],
      "factorio": "Mod"
    },
    {
      "discord": ["Members","Server owners"],
      "factorio": "Member"
    }
  ]
}
```
* To get the `discord_guild_id` go into [developer](https://support.discordapp.com/hc/en-us/articles/206346498-Where-can-I-find-my-server-ID-) mode and right click on a channel.

* To get the `discord_bot_token` Go to the [Discord developer tools](https://discordapp.com/developers/applications/me) and create a bot or use a already existing one. Scroll down to bot then `click to reveal token`
## Roles
The `discord` property contains the roles that the players will need to have to get the `Member` role ingame
```json
{
  "discord": ["Members","Server owners"],
  "factorio": "Member"
}
```

# Output
Program output:
```
2018/01/26 15:15:50 Config loaded: 335µs
2018/01/26 15:15:50 Config Loaded:
  Token: TOKEN
  GuildID: GUILDID
  Roles: %ROLES
2018/01/26 15:15:50 New Discord bot created: 296µs
2018/01/26 15:15:51 Discord bot is running and ready: 762.27ms - BOT NAME
2018/01/26 15:15:51 Member get: 214.197ms
2018/01/26 15:15:51 Member parse: 6.711ms
2018/01/26 15:15:51 Member to command: 50µs
2018/01/26 15:15:51 Rcon ready: 25.633ms
2018/01/26 15:15:51 Response: &{Size:10 ID:-1013797929 Type:0 Body:}
2018/01/26 15:15:51 Rcon command send: 26.92ms
2018/01/26 15:15:52 Discord bot Closed & process end: 2.037596s
```

---

Small example command output:
```
/interface Ranking._base_preset{
  ["cydes"]="Mod",
  ["BADgamerNL"]="Owner",
  ["Cooldude2606"]="Developer",
  ["Klonan"]="Developer",
  ["mark9064"]="Admin"
}
```

# Credits
Packages I used in this project.
* [discordgo](http://github.com/bwmarrin/discordgo)
* [factorio-rcon](http://github.com/gtaylor/factorio-rcon)
* [explosivegaming-main](https://github.com/badgamernl/explosivegaming-main) - [Cooldude2606](https://github.com/Cooldude2606)