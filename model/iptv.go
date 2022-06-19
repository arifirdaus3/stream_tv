package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CountryIPTV struct {
	Code       string `msgpack:"code" json:"code" `
	Name       string `msgpack:"name" json:"name"`
	LanguageID string `msgpack:"lang" json:"lang"`
	Flag       string `msgpack:"flag" json:"flag"`
}

type RegionIPTV struct {
	Code           string         `json:"code"`
	Name           string         `json:"name"`
	Countries      []Country      ` gorm:"many2many:region_country;"`
	CountriesArray []string       `json:"countries"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type SubdivisionIPTV struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type ChannelIPTV struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name"`
	NativeName    string         `json:"native_name"`
	Network       string         `json:"network"`
	Country       string         `json:"country"`
	Subdivision   string         `json:"subdivision"`
	City          string         `json:"city"`
	BroadcastArea pq.StringArray `json:"broadcast_area" gorm:"type:text[]"`
	Languages     []string       `json:"languages" gorm:"many2many:channel_language;"`
	Categories    []string       `json:"categories" gorm:"many2many:channel_category;"`
	IsNsfw        bool           `json:"is_nsfw"`
	Launched      string         `json:"launched"`
	Closed        string         `json:"closed"`
	ReplacedBy    string         `json:"replaced_by"`
	Website       string         `json:"website"`
	Logo          string         `json:"logo"`
}

type GuideIPTV struct {
	gorm.Model
	ChannelID string `json:"channel"`
	Site      string `json:"site"`
	Lang      string `json:"lang"`
	URL       string `json:"url"`
}

var (
	CategoryKey    = "category"
	LanguageKey    = "language"
	CountryKey     = "country"
	RegionKey      = "region"
	SubDivisionKey = "subdivision"
	ChannelKey     = "channel"
	GuideKey       = "guide"
)
