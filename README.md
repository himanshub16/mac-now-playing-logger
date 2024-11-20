# mac-now-playing-logger

**What it does?**

Logs and records what is currently in "Now Playing" to console and duckdb.

**Why?**

I want to eventually figure out - 
- which song was the one I looped most number of times today
- what was my most played song/artist for this week

**How does this work?**

One way would be to have some chrome extension track Youtube Music or Spotify. I use Youtube Music.
But that is too much work, and there's no API to get history.

[nowplaying-cli](https://github.com/kirtan-shah/nowplaying-cli/) is a great tool which calls Mac OS's APIs to get what is playing right now.

This gives metadata as - `song name`, `artist`, `album`, `where is the timer right now`, `media duration`.
Works totally fine for all music players, including Youtube Music Chrome App.

## How to use?
```
brew install nowplaying-cli
go run main.go
```
