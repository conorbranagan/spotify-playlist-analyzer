# Spotify Playlist Analyzer

## Overview

Analyzes all or a subset of playlists in your Spotify account and outputs some interesting statistics.

We will use the [Spotify Web API](https://developer.spotify.com/web-api/) to capture information about a user's playlists via the [Golang API wrapper](https://github.com/zmb3/spotify).

## Motivation

I often create a new playlist for every road trip that represents a snapshot of what I'm listning to at that moment. These are great because they provide a history of my music tastes over time.

I was interested in getting an idea of how my music taste change over time and what, if any, artists remained consistent.

## Running

* Create an `env` file with the following contents:

```
export SPOTIFY_ID=$YOUR_SPOTIFY_ID
export SPOTIFY_SECRET=$YOUR_SPOTIFY_SECRET
```

* Install the package.

```bash
$ go get
$ go install
```

* Run the app and authenticate.

```
$ spotify-playlist-analyzer
# Navigate to `localhost:8080` and copy the URL.
# Go to the relevant authentication URL.
```

Once authenticated, the stats should print out. Any further runs of the app should use a cached token on the disk.
