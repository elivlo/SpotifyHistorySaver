package spotifySaver

import (
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/models"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"github.com/zmb3/spotify"
	"strings"
)

type FetchedSongs struct {
	db      *pop.Connection
	fetched []spotify.RecentlyPlayedItem

	history models.HistoryEntries
	tracks models.Tracks
	artists models.Artists
	connections models.ArtistsTracks
}


// CreateFetchedSongs will create FetchedSongs struct.
func CreateFetchedSongs(d *pop.Connection, songs []spotify.RecentlyPlayedItem) FetchedSongs {
	fetchedSongs := FetchedSongs{
		db: d,
		fetched: songs,
	}
	return fetchedSongs
}

// TransformAndInsertIntoDatabase will convert and insert recently played songs into database.
func (s *FetchedSongs) TransformAndInsertIntoDatabase() error {
	s.convertRecentlyToDBTables()
	err := s.db.Create(&s.tracks)
	if err != nil {
		return errors.Errorf("Could not insert tracks: %v", err)
	}
	err = s.db.Create(&s.artists)
	if err != nil {
		return errors.Errorf("Could not insert artists: %v", err)
	}
	err = s.db.Create(&s.history)
	if err != nil {
		return errors.Errorf("Could not insert history: %v", err)
	}
	err = s.db.Create(&s.connections)
	if err != nil {
		return errors.Errorf("Could not insert artist track connections: %v", err)
	}
	LOG.Debugf("Added %d new tracks, %d new artists and %d history tracks", len(s.tracks), len(s.artists), len(s.history))
	return nil
}

// convertRecentlyToDBTables will convert API json to database models.
// It will also exclude Tracks and Artists that already exists in database.
func (s *FetchedSongs) convertRecentlyToDBTables() {
	for _, song := range s.fetched {
		s.history = append(s.history, convertToHistoryEntry(song))

		track := convertToTrackEntry(song)
		trackInserted, err := s.trackAlreadyInserted(track.ID)
		if err != nil {
			LOG.Errorf("Song %v could not be added: %v\n", song, err)
			continue
		}
		if !trackInserted {
			s.tracks = append(s.tracks, track)
			arts, conn := convertToArtistEntries(song)
			s.connections = append(s.connections, conn...)

			for _, art := range arts {
				artInserted, err := s.artistAlreadyInserted(art.ID)
				if err != nil {
					LOG.Errorf("Artist %v could not be added: %v\n", art, err)
					continue
				}
				if !artInserted {
					s.artists = append(s.artists, art)
				}
			}
		}
	}
}

// trackAlreadyInserted check if database contains track.
func (s *FetchedSongs) trackAlreadyInserted(id string) (bool, error) {
	track := models.Track{}
	for _, t := range s.tracks {
		if t.ID == id {
			return true, nil
		}
	}

	err := s.db.Find(&track, id)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		return false, nil
	}
	if track.ID == id {
		return true, nil
	}
	return false, err
}

// artistAlreadyInserted check if database contains artist.
func (s *FetchedSongs) artistAlreadyInserted(id string) (bool, error) {
	artist := models.Artist{}
	for _, a := range s.artists {
		if a.ID == id {
			return true, nil
		}
	}

	err := s.db.Find(&artist, id)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		return false, nil
	}
	if artist.ID == id {
		return true, nil
	}
	return false, err
}

func convertToHistoryEntry (song spotify.RecentlyPlayedItem) models.HistoryEntry {
	return models.HistoryEntry{
		TrackId:  song.Track.ID.String(),
		PlayedAt: song.PlayedAt,
	}
}

func convertToTrackEntry (song spotify.RecentlyPlayedItem) models.Track {
	track := song.Track
	return models.Track{
		ID:          track.ID.String(),
		Name:        track.Name,
		TrackNumber: track.TrackNumber,
		DiscNumber:  track.DiscNumber,
		Explicit:    track.Explicit,
	}
}

// convertToArtistEntries created artists and the connection to a track.
func convertToArtistEntries (song spotify.RecentlyPlayedItem) (models.Artists, models.ArtistsTracks) {
	songID := song.Track.ID.String()
	var artists models.Artists
	var connection models.ArtistsTracks
	as := song.Track.Artists
	for _, a := range as {
		artists = append(artists, models.Artist{
			ID:       a.ID.String(),
			Name:     a.Name,
		})
		connection = append(connection, models.ArtistsTrack{
			ArtistID: a.ID.String(),
			TrackID:  songID,
		})
	}
	return artists, connection
}

func getLastHistoryEntry(db *pop.Connection) (models.HistoryEntry, error) {
	var last models.HistoryEntry
	err := db.Order("played_at DESC").First(&last)
	return last, err
}