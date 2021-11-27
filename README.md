# 📚 WordFinder

## 🏛️ Current State
This repository has been created to make analysis about lyrics of songs, the goal is to provide tool which will be able to:

- ✔️   REST API for searching songs by artist and filtering out banned words
- ✔️   Search <span style="color:green">**600 songs** from **Eminem</span> in <span style="color:green">9 seconds**</span>. on Ryzen 7 5800X and 500MB/s isp
- ✔️   Find all songs by artist without banned words, could be used to find "family friendly" music without some kind of words
- ✔️   Provide list of keywords in many ways in ex. these keywords are going to be used as arguments
  
- ❌ Find occurrence of specific words and calculate in which songs the word were most used. 
- ❌ Database
- ❌ Better way of managing banned words sets (There will be endpoint for registering keywords sets, then it will be accessible by id as filter)
## 🚀 Future plans
- Swagger
- Make better errors logging - including sentry
- Caching in PostgresSQL and / or Redis
- Create API Clients (frontend or cli-client)
- Performance monitoring using Grafana

## 🔨 Build & Run
CLI version:
```bash
make cli
```

API version:
```bash
make api
```

<br><br>
<span style="color:red">**IMPORTANT**</span> For both of these binaries you will need to prepare `.env` file with credentials

```bash
touch .env
```

Paste template:
```bash
export USER_AGENTS="Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"
export GENIUS_RAPID_API_HOST=genius.p.rapidapi.com
export GENIUS_RAPID_API_KEY=[OBTAIN IT FROM RAPIDAPI.COM]
export GENIUS_HOST=genius.com
export GENIUS_API_HOST=genius.com/api

export REQUEST_TIMEOUT=10s
export MAX_CHANNEL_BUFFER_SIZE=30
export SERVER_PORT=8080
```

in place of <span style="color:orange">[OBTAIN IT FROM RAPIDAPI.COM]</span> put api token from https://rapidapi.com/brianiswu/api/genius/

## 🪧 Usage of API
The basic response struct:
```json5
{
  "data": null,
  "error": null,
}
```
^ ps. only one of these values may be equal to `null`

### GET https://localhost:8080/artists/:the_artist_name/songs

```json5
{
  "data": {
    "songs": [
      {
        "title": "Example",
        "url": "https://genius.com/example"
      },
      {
        "title": "Example1",
        "url": "https://genius.com/example-1"
      }
    ]
  },
  "error": null
}
```

### GET https://localhost:8080/artists/:the_artist_name/songs/words
`word_count` contains all words used in lyrics which are longer than 2 characters.
```json5
{
  "data": {
    "songs": [
      {
        "title": "Example",
        "url": "https://genius.com/example",
        "words_count": {
          "abc": 2,
          "cba": 1
        }
      },
      {
        "title": "Example1",
        "url": "https://genius.com/example-1",
        "words_count": {
          "qwe": 2,
          "rrr": 1,
        }
      }
    ]
  },
  "error": null
}
```


### GET https://localhost:8080/artists/:the_artist_name/songs/words?banned_words=:base64
This example shows how songs can be filtered out because of containing one of banned words

`word_count` contains all words used in lyrics which are longer than 2 characters.

`:base64` param in url is base64 string with banned words separated by commas, 


example url: `GET` https://localhost:8080/artists/eminem/songs/words?banned_words=d29yZCx3b3JkMix3b3JkMw
Will return eminem songs without base64 encoded words in lyrics
```json5
{
  "data": {
    "songs": [
      {
        "title": "Example",
        "url": "https://genius.com/example",
        "words_count": {
          "abc": 2,
          "cba": 1
        }
      },
    ]
  },
  "error": null
}
```

 💥 `./genius-cli` 💥
## 🪧 Usage of CLI

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

#### Simplest Usage 
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
