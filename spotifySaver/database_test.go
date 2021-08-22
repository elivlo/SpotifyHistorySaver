package spotifySaver

import (
	"github.com/elivlo/SpotifyHistorySaver/models"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
	"testing"
	"time"
)

func TestNewFetchedSongs(t *testing.T) {
	songs := NewFetchedSongs(DB, []spotify.RecentlyPlayedItem{})
	assert.Equal(t, DB, songs.db)
	assert.Equal(t, 0, len(songs.fetched))

	assert.Equal(t, 0, len(songs.history))
	assert.Equal(t, 0, len(songs.tracks))
	assert.Equal(t, 0, len(songs.artists))
	assert.Equal(t, 0, len(songs.connections))
}

func TestFetchedSongs_TransformAndInsertIntoDatabase(t *testing.T) {
	_, log := getTestLogger()
	songs := NewFetchedSongs(DB, []spotify.RecentlyPlayedItem{})

	err := songs.TransformAndInsertIntoDatabase(log)
	assert.Nil(t, err)
}

func TestFetchedSongs_convertRecentlyToDBTables(t *testing.T) {
	hook, log := getTestLogger()

	songs := NewFetchedSongs(DB, []spotify.RecentlyPlayedItem{{
		Track: spotify.SimpleTrack{
			Artists: []spotify.SimpleArtist{{
				Name: "a_name",
				ID:   "a_id",
			}},
			ID: "t_id",
		},
		PlayedAt: time.Now(),
		PlaybackContext: spotify.PlaybackContext{
			ExternalURLs: nil,
			Endpoint:     "endpoint",
			Type:         "type",
			URI:          "uri",
		},
	}})

	songs.convertRecentlyToDBTables(log)
	assert.Equal(t, 0, len(hook.AllEntries()))
}

func TestFetchedSongs_trackAlreadyInserted(t *testing.T) {
	songs := NewFetchedSongs(DB, []spotify.RecentlyPlayedItem{})

	b, err := songs.trackAlreadyInserted("test_track")
	assert.False(t, b)
	assert.NoError(t, err)

	songs = NewFetchedSongs(DB, []spotify.RecentlyPlayedItem{})
	songs.tracks = append(songs.tracks, models.Track{
		ID: "t_id",
	})

	b, err = songs.trackAlreadyInserted("t_id")
	assert.True(t, b)
	assert.NoError(t, err)

	err = DB.Create(&songs.tracks)
	assert.NoError(t, err)
	songs.tracks = models.Tracks{}

	b, err = songs.trackAlreadyInserted("t_id")
	assert.True(t, b)
	assert.NoError(t, err)
}

func TestFetchedSongs_artistAlreadyInserted(t *testing.T) {
	songs := NewFetchedSongs(DB, []spotify.RecentlyPlayedItem{})

	b, err := songs.artistAlreadyInserted("a_id")
	assert.False(t, b)
	assert.NoError(t, err)

	songs = NewFetchedSongs(DB, []spotify.RecentlyPlayedItem{})
	songs.artists = append(songs.artists, models.Artist{
		ID: "a_id",
	})

	b, err = songs.artistAlreadyInserted("a_id")
	assert.True(t, b)
	assert.NoError(t, err)

	err = DB.Create(&songs.artists)
	assert.NoError(t, err)
	songs.artists = models.Artists{}

	b, err = songs.artistAlreadyInserted("a_id")
	assert.True(t, b)
	assert.NoError(t, err)
}

func TestConvertToHistoryEntry(t *testing.T) {
	now := time.Now()
	song := spotify.RecentlyPlayedItem{
		Track: spotify.SimpleTrack{
			ID: "t_id",
		},
		PlayedAt: now,
	}

	entry := convertToHistoryEntry(song)
	assert.Equal(t, now, entry.PlayedAt)
	assert.Equal(t, "t_id", entry.TrackID)
}

func TestConvertToTrackEntry(t *testing.T) {
	song := spotify.RecentlyPlayedItem{
		Track: spotify.SimpleTrack{
			DiscNumber:  1,
			Explicit:    true,
			ID:          "t_id",
			Name:        "t_name",
			TrackNumber: 1,
		},
	}

	entry := convertToTrackEntry(song)
	assert.Equal(t, "t_id", entry.ID)
	assert.Equal(t, "t_name", entry.Name)
	assert.Equal(t, 1, entry.TrackNumber)
	assert.Equal(t, 1, entry.DiscNumber)
	assert.True(t, entry.Explicit)
}

func TestConvertToArtistEntries(t *testing.T) {
	song := spotify.RecentlyPlayedItem{
		Track: spotify.SimpleTrack{
			Artists: []spotify.SimpleArtist{{
				Name: "a_name",
				ID:   "a_id",
			}},
			ID: "t_id",
		},
	}

	artists, tracks := convertToArtistEntries(song)
	assert.Equal(t, 1, len(artists))
	assert.Equal(t, 1, len(tracks))

	assert.Equal(t, "a_name", artists[0].Name)
	assert.Equal(t, "a_id", artists[0].ID)

	assert.Equal(t, "t_id", tracks[0].TrackID)
	assert.Equal(t, "a_id", tracks[0].ArtistID)
}

func TestGetLastHistoryEntry(t *testing.T) {
	now := time.Now()
	entry := models.HistoryEntry{
		TrackID:  "t_id",
		PlayedAt: now,
	}
	err := DB.Create(&entry)
	assert.NoError(t, err)

	e, err := getLastHistoryEntry(DB)
	assert.NoError(t, err)

	assert.Equal(t, "t_id", e.TrackID)
	assert.Equal(t, 1, e.ID)
}
