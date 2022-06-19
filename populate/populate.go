package populate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arifirdaus3/stream_tv/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Category(db *gorm.DB) error {
	url := "https://iptv-org.github.io/api/categories.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var res []model.Category
	json.NewDecoder(resp.Body).Decode(&res)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&res)
	if result.Error != nil {
		return err
	}
	fmt.Println("Category populated with", result.RowsAffected, "rows")
	return nil
}

func Language(db *gorm.DB) error {
	url := "https://iptv-org.github.io/api/languages.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var res []model.Language
	json.NewDecoder(resp.Body).Decode(&res)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&res)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("Language populated with", result.RowsAffected, "rows")
	return nil
}
func Country(db *gorm.DB) error {
	url := "https://iptv-org.github.io/api/countries.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var iptvres []model.CountryIPTV
	var res []model.Country
	json.NewDecoder(resp.Body).Decode(&iptvres)
	for _, v := range iptvres {
		lang := sql.NullString{}
		if v.LanguageID != "" {
			lang.String = v.LanguageID
			lang.Valid = true
		}
		res = append(res, model.Country{
			Name:       v.Name,
			Code:       v.Code,
			LanguageID: lang,
			Flag:       v.Flag,
		})
	}

	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&res)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("Country populated with", result.RowsAffected, "rows")
	return nil
}
func Subdivision(db *gorm.DB) error {
	url := "https://iptv-org.github.io/api/subdivisions.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var res []model.Subdivision
	var iptv []model.SubdivisionIPTV
	json.NewDecoder(resp.Body).Decode(&iptv)
	for _, v := range iptv {
		country := sql.NullString{}
		if v.Country != "" {
			country.String = v.Country
			country.Valid = true
		}
		s := model.Subdivision{
			Code:      v.Code,
			Name:      v.Name,
			CountryID: country,
		}
		res = append(res, s)
	}
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&res)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("Subdivision populated with", result.RowsAffected, "rows")
	return nil
}

func Region(db *gorm.DB) error {
	url := "https://iptv-org.github.io/api/regions.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var res []model.RegionIPTV
	var in []model.Region
	json.NewDecoder(resp.Body).Decode(&res)
	for _, v := range res {
		for _, c := range v.CountriesArray {
			v.Countries = append(v.Countries, model.Country{Code: c})
		}
		in = append(in, model.Region{
			Code:      v.Code,
			Name:      v.Name,
			Countries: v.Countries,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			DeletedAt: v.DeletedAt,
		})
	}
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&in)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("Region populated with", result.RowsAffected, "rows")
	return nil
}
func Channel(db *gorm.DB) error {
	url := "https://iptv-org.github.io/api/channels.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var iptv []model.ChannelIPTV
	json.NewDecoder(resp.Body).Decode(&iptv)
	index := 0
	for index < len(iptv) {
		temp := []model.Channel{}
		for i := 0; i < 500 && index < len(iptv); i++ {
			lang := []model.Language{}
			for _, x := range iptv[index].Languages {
				lang = append(lang, model.Language{Code: x})
			}
			cat := []model.Category{}
			for _, x := range iptv[index].Categories {
				cat = append(cat, model.Category{ID: x})
			}
			v := model.Channel{
				Name:          iptv[index].Name,
				ID:            iptv[index].ID,
				NativeName:    iptv[index].NativeName,
				Country:       iptv[index].Country,
				Subdivision:   iptv[index].Subdivision,
				City:          iptv[index].City,
				BroadcastArea: iptv[index].BroadcastArea,
				Languages:     lang,
				Categories:    cat,
				IsNsfw:        iptv[index].IsNsfw,
				Launched:      iptv[index].Launched,
				Closed:        iptv[index].Closed,
				ReplacedBy:    iptv[index].ReplacedBy,
				Website:       iptv[index].Website,
				Logo:          iptv[index].Logo,
			}
			temp = append(temp, v)
			index++
		}
		result := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&temp)
		if result.Error != nil {
			return result.Error
		}
		fmt.Println("Channel populated with", result.RowsAffected, "rows")
	}
	fmt.Println("Channel populated with", index, "rows")
	return nil
}
func Guide(db *gorm.DB) error {
	url := "https://iptv-org.github.io/api/guides.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var res []model.GuideIPTV
	json.NewDecoder(resp.Body).Decode(&res)
	index := 0
	for index < len(res) {
		temp := []model.Guide{}
		for i := 0; i < 500 && index < len(res); i++ {
			cid := sql.NullString{}
			if res[index].ChannelID != "" {
				cid.String = res[index].ChannelID
				cid.Valid = true
			}
			v := model.Guide{
				Model: gorm.Model{
					ID: uint(index),
				},
				ChannelID: cid,
				Site:      res[index].Site,
				Lang:      res[index].Lang,
				URL:       res[index].URL,
			}
			temp = append(temp, v)
			index++
		}
		result := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&temp)
		if result.Error != nil {
			return result.Error
		}
		fmt.Println("Guide populated with", result.RowsAffected, "rows")
	}
	fmt.Println("Guide populated with", index, "rows")
	return nil
}

func Stream(db *gorm.DB) error {
	type s struct {
		Channel string `json:"channel"`
		URL     string `json:"url"`
		Status  string `json:"status"`
		urls    []string
		idx     int
	}
	streams := []s{}
	url := "https://iptv-org.github.io/api/streams.json"
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	json.NewDecoder(resp.Body).Decode(&streams)
	index := 0
	for index < len(streams) {
		temp := []model.Channel{}
		unique := map[string]s{}
		counter := 0
		for i := 0; i < 500 && index < len(streams); i++ {
			k, o := unique[streams[index].Channel]
			if o {
				urls := k.urls
				urls = append(urls, streams[index].URL)
				counter++
				temp[k.idx].URL = urls
			} else {
				streams[index].idx = i - counter
				streams[index].urls = []string{streams[index].URL}
				channel := model.Channel{
					ID:     streams[index].Channel,
					URL:    streams[index].urls,
					Status: streams[index].Status,
				}

				unique[streams[index].Channel] = streams[index]
				temp = append(temp, channel)
			}

			index++
		}
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "status"}),
		}).Create(&temp)
		if result.Error != nil {
			return result.Error
		}
		fmt.Println("Channel updated using Streams with", result.RowsAffected, "rows")
	}

	return nil
}
func All(db *gorm.DB) error {
	if err := Category(db); err != nil {
		return err
	}
	if err := Language(db); err != nil {
		return err
	}
	if err := Country(db); err != nil {
		return err
	}
	if err := Subdivision(db); err != nil {
		return err
	}
	if err := Region(db); err != nil {
		return err
	}
	if err := Channel(db); err != nil {
		return err
	}
	if err := Guide(db); err != nil {
		return err
	}
	if err := Stream(db); err != nil {
		return err
	}
	return nil
}
