package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CountryResponse struct {
	Code      string         `msgpack:"code" json:"code"`
	Name      string         `msgpack:"name" json:"name"`
	Language  Language       `msgpack:"language" json:"language" gorm:"foreignKey:LanguageID;references:Code"`
	Flag      string         `msgpack:"flag" json:"flag"`
	CreatedAt time.Time      `msgpack:"created_at" json:"created_at"`
	UpdatedAt time.Time      `msgpack:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}

type RegionResponse struct {
	Code           string         `msgpack:"code" json:"code"`
	Name           string         `msgpack:"name" json:"name"`
	CountriesArray []string       `msgpack:"countries" json:"countries"`
	CreatedAt      time.Time      `msgpack:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `msgpack:"updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}

type SubdivisionResponse struct {
	Code      string         `msgpack:"code" json:"code"`
	Name      string         `msgpack:"name" json:"name"`
	CountryID string         `msgpack:"country_id" json:"country_id"`
	CreatedAt time.Time      `msgpack:"created_at" json:"created_at"`
	UpdatedAt time.Time      `msgpack:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}

type ChannelResponse struct {
	ID            string                    `msgpack:"id" json:"id"`
	Name          string                    `msgpack:"name" json:"name"`
	NativeName    string                    `msgpack:"native_name" json:"native_name"`
	Network       string                    `msgpack:"network" json:"network"`
	Country       string                    `msgpack:"country" json:"country"`
	Subdivision   string                    `msgpack:"subdivision" json:"subdivision"`
	City          string                    `msgpack:"city" json:"city"`
	BroadcastArea pq.StringArray            `msgpack:"broadcast_area" json:"broadcast_area"`
	Languages     []ChannelLanguageResponse `msgpack:"languages" json:"languages"`
	Categories    []ChannelCategoryResponse `msgpack:"categories" json:"categories"`
	IsNsfw        bool                      `msgpack:"is_nsfw" json:"is_nsfw"`
	Launched      string                    `msgpack:"launched" json:"launched"`
	Closed        string                    `msgpack:"closed" json:"closed"`
	ReplacedBy    string                    `msgpack:"replaced_by" json:"replaced_by"`
	Website       string                    `msgpack:"website" json:"website"`
	Logo          string                    `msgpack:"logo" json:"logo"`
	URL           pq.StringArray            `msgpack:"url" json:"url"`
	Status        string                    `msgpack:"status" json:"status"`
	CreatedAt     time.Time                 `msgpack:"created_at" json:"created_at"`
	UpdatedAt     time.Time                 `msgpack:"updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt            `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}
type ChannelLanguageResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type ChannelCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type GuideResponse struct {
	ID        string         `json:"id"`
	ChannelID string         `json:"channel_id"`
	Site      string         `json:"site"`
	Lang      string         `json:"lang"`
	URL       string         `json:"url"`
	CreatedAt time.Time      `msgpack:"created_at" json:"created_at"`
	UpdatedAt time.Time      `msgpack:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}
