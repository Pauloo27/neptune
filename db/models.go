package db

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	MBID string `gorm:"unique"`
	Name string
}

type Album struct {
	gorm.Model
	MBID     string `gorm:"unique"`
	Title    string
	ArtistID uint
	Artist   Artist
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type TrackTag struct {
	gorm.Model
	TagID   uint
	Tag     Tag
	TrackID uint
	Track   Track
}

type Track struct {
	gorm.Model
	MBID         string `gorm:"unique"`
	YoutubeID    string
	AlbumID      uint
	Album        Album
	Title        string
	Length       int
	YoutubeTitle string
	Tags         []TrackTag
}
