package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Category struct {
	ID        string         `msgpack:"id" json:"id" gorm:"primaryKey"`
	Name      string         `msgpack:"name" json:"name" gorm:"index"`
	CreatedAt time.Time      `msgpack:"created_at" json:"created_at"`
	UpdatedAt time.Time      `msgpack:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}

type Language struct {
	Code      string         `msgpack:"code" json:"code" gorm:"primaryKey"`
	Name      string         `msgpack:"name" json:"name" gorm:"index"`
	CreatedAt time.Time      `msgpack:"created_at" json:"created_at"`
	UpdatedAt time.Time      `msgpack:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}

type Country struct {
	Code       string         `msgpack:"code" json:"code" gorm:"primaryKey"`
	Name       string         `msgpack:"name" json:"name"`
	LanguageID sql.NullString `msgpack:"lang" json:"lang"`
	Language   Language       `msgpack:"language" json:"language" gorm:"foreignKey:LanguageID;references:Code"`
	Flag       string         `msgpack:"flag" json:"flag"`
	CreatedAt  time.Time      `msgpack:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `msgpack:"updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `msgpack:"deleted_at" json:"deleted_at" gorm:"index"`
}

type Region struct {
	Code      string         `json:"code" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Countries []Country      `gorm:"many2many:region_country;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Subdivision struct {
	Code      string         `json:"code" gorm:"primaryKey"`
	Name      string         `json:"name"`
	CountryID sql.NullString `json:"country_id"`
	Country   Country        `json:"country" gorm:"foreignKey:CountryID;references:Code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Channel struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name"`
	NativeName    string         `json:"native_name"`
	Network       string         `json:"network"`
	Country       string         `json:"country"`
	Subdivision   string         `json:"subdivision"`
	City          string         `json:"city"`
	BroadcastArea pq.StringArray `json:"broadcast_area" gorm:"type:text[]"`
	Languages     []Language     `gorm:"many2many:channel_language;"`
	Categories    []Category     `gorm:"many2many:channel_category;"`
	IsNsfw        bool           `json:"is_nsfw"`
	Launched      string         `json:"launched"`
	Closed        string         `json:"closed"`
	ReplacedBy    string         `json:"replaced_by"`
	Website       string         `json:"website"`
	Logo          string         `json:"logo"`
	URL           pq.StringArray `json:"url" gorm:"type:text[]"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Guide struct {
	gorm.Model
	Channel   Channel        `json:"channel" gorm:"foreignKey:ChannelID;references:ID"`
	ChannelID sql.NullString `json:"channel_id"`
	Site      string         `json:"site"`
	Lang      string         `json:"lang"`
	URL       string         `json:"url"`
}
