package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/zmb3/spotify"
)

type pair struct {
	key string
	val int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].val > p[j].val }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type playlistStats struct {
	c           *spotify.Client
	artistCount map[string]int
}

func newPlaylistStats(c *spotify.Client) *playlistStats {
	return &playlistStats{
		c:           c,
		artistCount: make(map[string]int),
	}
}

func (s *playlistStats) analyzeUser(userID string) error {
	playlists := make([]spotify.SimplePlaylist, 0)
	ppage, err := s.c.GetPlaylistsForUser(userID)
	if err != nil {
		return err
	}
	playlists = append(playlists, ppage.Playlists...)

	offset := ppage.Limit
	for offset < ppage.Total {
		opts := &spotify.Options{
			Offset: &offset,
			Limit:  &ppage.Limit,
		}
		ppage, err = s.c.GetPlaylistsForUserOpt(userID, opts)
		if err != nil {
			return err
		}
		for _, p := range ppage.Playlists {
			// Only consider collaborative or playlists created by this user
			if p.Owner.ID == userID || p.Collaborative {
				playlists = append(playlists, p)
			}
		}
		offset += ppage.Limit
	}

	analyzed := 0
	log.Printf("Analyzing %d playlists...", len(playlists))
	for _, p := range playlists {
		s.updateFromPlaylist(p.Owner.ID, p.ID)
		analyzed++

		pct := (float32(analyzed) / float32(len(playlists))) * 100
		fmt.Printf("\r%1.1f%% of playlists analyzed", pct)
	}
	return nil
}

func (s *playlistStats) updateFromPlaylist(owner string, playlistID spotify.ID) error {
	tracks := make([]spotify.PlaylistTrack, 0)
	tpage, err := s.c.GetPlaylistTracksOpt(owner, playlistID, nil, "")
	if err != nil {
		return err
	}

	tracks = append(tracks, tpage.Tracks...)
	for _, t := range tracks {
		var artists []string
		for _, a := range t.Track.Artists {
			artists = append(artists, a.Name)
			s.artistCount[a.Name]++
		}
	}
	return nil
}

func (s *playlistStats) printStats() {
	fmt.Printf("\n\n--------------------------------\n")
	fmt.Printf("Playlist Statistics")
	fmt.Printf("\n--------------------------------\n\n")

	aCounts := make([]pair, 0, len(s.artistCount))
	for a, c := range s.artistCount {
		aCounts = append(aCounts, pair{a, c})
	}
	sort.Sort(pairList(aCounts))
	for _, p := range aCounts {
		if p.val < 5 {
			break
		}
		fmt.Printf("* %s - %d\n", p.key, p.val)
	}
}
