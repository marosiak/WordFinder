# 📚 WordFinder

## 🏛️ Current State
This repository has been created to make analysis about lyrics of songs, the goal is to provide tool which will be able to:

- ✔️   Search <span style="color:green">**600 songs** from **Eminem</span> in <span style="color:green">3 seconds**</span>. on Ryzen 7 5800X and 500MB/s isp
- ✔️   Find all songs by artist without banned words, could be used to find "family friendly" music without some kind of words
- ✔️   Provide list of keywords in many ways in ex. these keywords are going to be used as arguments
  
- ❌ Find occurrence of specific words and calculate in which songs the word were most used. 

## 🚀 Future plans
- Caching in PostgresSQL and / or Redis
- Creating API for clients so it will be hosted

## 🔨 Build & Run
```bash
./build.sh
```
In order to access Genius you will need `.env` file with credentials

```bash
touch .env
```

Paste template:
```bash
USER_AGENTS="Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"
GENIUS_RAPID_API_HOST=genius.p.rapidapi.com
GENIUS_RAPID_API_KEY=[OBTAIN IT FROM RAPIDAPI.COM]
GENIUS_HOST=genius.com
GENIUS_API_HOST=genius.com/api

REQUEST_TIMEOUT=10s
MAX_CHANNEL_BUFFER_SIZE=30
```

in place of <span style="color:orange">[OBTAIN IT FROM RAPIDAPI.COM]</span> put api token from https://rapidapi.com/brianiswu/api/genius/
   
 💥 `./genius-cli` 💥
## 🪧 Usage

This is generic --help view, you may use it with other commands for more details for ex.

`genius-cli --help`
`genius-cli command-name --help`

```text
NAME:
   genius-cli - genius-cli --help

USAGE:
   genius-cli [global options] command [command options] [arguments...]

COMMANDS:
   songs-by-artist-without-banned-words  Will return list of songs which does not contains any of --keywords or --keyword
   help, h                               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### 🚫🍆 genius-cli songs-by-artist-without-banned-words --help
**yes, this name sucks, give me better one pls**.
This command will output list of songs without provided words 

####Simplest Usage 
```bash
genius-cli songs-by-artist-without-banned-words --keywords-file="swears.txt"
```
The `swears.txt` file should contain words separated by new lines or commas(",")
```bash
NAME:
   genius-cli songs-by-artist-without-banned-words - Will return list of songs which does not contains any of --keywords or --keyword

USAGE:
   genius-cli songs-by-artist-without-banned-words [command options] [arguments...]

OPTIONS:
   --query value, -q value                  --query="the_name"
   --keyword value, --kwd value             --keyword="the_keyword"
   --keywords value, --kwds value           --keywords="the_keyword","another_keyword"
   --keywords-file value, --kwds-f value    --keywords-file="keywords.txt"
   --keywords-files value, --kwds-fs value  --keywords-files="swears.txt,drugs.txt"
   --help, -h                               show help (default: false)
```