package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arifirdaus3/stream_tv/model"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type handler struct {
	db    *gorm.DB
	cache *cache.Cache
}

func (h *handler) handleCategory(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := []model.Category{}
	var count int64
	page, size := getPageAndSize(q)
	h.db.Scopes(paginate(page, size)).
		Where("name ILIKE ? AND id ILIKE ?", searcher(q.Get("name")), searcher(q.Get("id"))).
		Find(&result)

	c, f := h.cache.Get(model.CategoryKey)
	if !f {
		h.db.Model(&model.Category{}).Count(&count)
		c = count
		h.cache.SetDefault(model.CategoryKey, count)
	}
	count = c.(int64)

	useMsgPack, _ := strconv.ParseBool(q.Get("msgpack"))
	if useMsgPack {
		paginationSuccessMsgPackResponse(w, result, hasNext(page, size, count), count)
		return
	}
	paginationSuccessJSONResponse(w, result, hasNext(page, size, count), count)
}

func (h *handler) handleLanguage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := []model.Language{}
	var count int64
	page, size := getPageAndSize(q)
	h.db.Scopes(paginate(page, size)).
		Where("name ILIKE ? AND code ILIKE ?", searcher(q.Get("name")), searcher(q.Get("code"))).
		Find(&result)

	c, f := h.cache.Get(model.LanguageKey)
	if !f {
		h.db.Model(&model.Language{}).Count(&count)
		c = count
		h.cache.SetDefault(model.LanguageKey, count)
	}
	count = c.(int64)

	useMsgPack, _ := strconv.ParseBool(q.Get("msgpack"))
	if useMsgPack {
		paginationSuccessMsgPackResponse(w, result, hasNext(page, size, count), count)
		return
	}
	paginationSuccessJSONResponse(w, result, hasNext(page, size, count), count)
}

func (h *handler) handleCountry(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := []model.CountryResponse{}
	var count int64
	page, size := getPageAndSize(q)
	rows, _ := h.db.Model(&model.Country{}).
		Scopes(paginate(page, size)).
		Select("countries.code, countries.name, countries.flag, countries.created_at, countries.updated_at, countries.deleted_at, languages.name, languages.code, languages.created_at, languages.updated_at, languages.deleted_at").
		Joins("JOIN languages ON countries.lang = languages.code").
		Where("countries.name ILIKE ? AND countries.code ILIKE ?", searcher(q.Get("name")), searcher(q.Get("code"))).
		Rows()

	defer rows.Close()
	for rows.Next() {
		var c model.CountryResponse
		rows.Scan(&c.Code, &c.Name, &c.Flag, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &c.Language.Name, &c.Language.Code, &c.Language.CreatedAt, &c.Language.UpdatedAt, &c.Language.DeletedAt)
		result = append(result, c)
	}
	c, f := h.cache.Get(model.CountryKey)
	if !f {
		h.db.Model(&model.Country{}).Count(&count)
		c = count
		h.cache.SetDefault(model.CountryKey, count)
	}
	count = c.(int64)

	useMsgPack, _ := strconv.ParseBool(q.Get("msgpack"))
	if useMsgPack {
		paginationSuccessMsgPackResponse(w, result, hasNext(page, size, count), count)
		return
	}
	paginationSuccessJSONResponse(w, result, hasNext(page, size, count), count)
}

func (h *handler) handleRegion(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := []model.RegionResponse{}
	var count int64
	page, size := getPageAndSize(q)
	subQuery := h.db.Model(&model.Region{}).Scopes(paginate(page, size))
	rows, _ := h.db.Table("(?) as r", subQuery).Select("r.code, r.name, r.created_at, r.updated_at, r.deleted_at, rc.country_code").
		Joins("JOIN region_country rc on rc.region_code = r.code").
		Where("r.name ILIKE ? AND r.code ILIKE ?", searcher(q.Get("name")), searcher(q.Get("code"))).
		Rows()

	index := 0
	unique := make(map[string]int)
	defer rows.Close()
	for rows.Next() {
		var c model.RegionResponse
		var countryCode string
		rows.Scan(&c.Code, &c.Name, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &countryCode)
		if _, ok := unique[c.Code]; !ok {
			unique[c.Code] = index
			index++
			result = append(result, c)
			continue
		}
		result[unique[c.Code]].CountriesArray = append(result[unique[c.Code]].CountriesArray, countryCode)
	}
	c, f := h.cache.Get(model.RegionKey)
	if !f {
		h.db.Model(&model.Region{}).Count(&count)
		c = count
		h.cache.SetDefault(model.RegionKey, count)
	}
	count = c.(int64)

	useMsgPack, _ := strconv.ParseBool(q.Get("msgpack"))
	if useMsgPack {
		paginationSuccessMsgPackResponse(w, result, hasNext(page, size, count), count)
		return
	}
	paginationSuccessJSONResponse(w, result, hasNext(page, size, count), count)
}

func (h *handler) handleSubDivision(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := []model.SubdivisionResponse{}
	var count int64
	page, size := getPageAndSize(q)
	h.db.Model(&model.Subdivision{}).
		Scopes(paginate(page, size)).
		Where("name ILIKE ? AND code ILIKE ?", searcher(q.Get("name")), searcher(q.Get("code"))).
		Find(&result)

	c, f := h.cache.Get(model.SubDivisionKey)
	if !f {
		h.db.Model(&model.Subdivision{}).Count(&count)
		c = count
		h.cache.SetDefault(model.SubDivisionKey, count)
	}
	count = c.(int64)

	useMsgPack, _ := strconv.ParseBool(q.Get("msgpack"))
	if useMsgPack {
		paginationSuccessMsgPackResponse(w, result, hasNext(page, size, count), count)
		return
	}
	paginationSuccessJSONResponse(w, result, hasNext(page, size, count), count)
}

func (h *handler) handleChannel(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := []model.ChannelResponse{}
	var count int64
	page, size := getPageAndSize(q)
	nsfw, _ := strconv.ParseBool(q.Get("nsfw"))
	hasURL, _ := strconv.ParseBool(q.Get("has_url"))
	where := "name ILIKE ? AND id ILIKE ?"
	if nsfw {
		where += " AND is_nsfw = true"
	}
	if hasURL {
		where += " AND url IS NOT NULL"
	}
	subQuery := h.db.Scopes(paginate(page, size)).Model(&model.Channel{}).Select("*").Where(where, searcher(q.Get("name")), searcher(q.Get("id")))
	selectColumn := "c.id, c.name, c.native_name, c.network, c.country, c.subdivision, c.city, c.broadcast_area, c.is_nsfw, c.launched, c.closed, c.replaced_by, c.website, c.logo, c.created_at, c.updated_at, c.deleted_at, c.url, c.status, cat.id, cat.name, l.code, l.name"
	rows, _ := h.db.Table("(?) as c", subQuery).Select(selectColumn).
		Joins("LEFT JOIN channel_category cc on cc.channel_id = c.id").
		Joins("LEFT JOIN categories cat on cat.id = cc.category_id").
		Joins("LEFT JOIN channel_language cl on cl.channel_id = c.id").
		Joins("LEFT JOIN languages l on l.code = cl.language_code").
		Rows()

	index := 0
	unique := make(map[string]int)
	uniqueCat := make(map[string]bool)
	uniqueLan := make(map[string]bool)
	defer rows.Close()
	for rows.Next() {
		var c model.ChannelResponse
		var cc model.ChannelCategoryResponse
		var cl model.ChannelLanguageResponse

		rows.Scan(&c.ID, &c.Name, &c.NativeName, &c.Network, &c.Country, &c.Subdivision, &c.City, &c.BroadcastArea, &c.IsNsfw, &c.Launched, &c.Closed, &c.ReplacedBy, &c.Website, &c.Logo, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &c.URL, &c.Status, &cc.ID, &cc.Name, &cl.Code, &cl.Name)

		uc := fmt.Sprintf("%s-%s", c.ID, cc.ID)
		ul := fmt.Sprintf("%s-%s", c.ID, cl.Code)

		if _, ok := unique[c.ID]; !ok {
			unique[c.ID] = index
			index++
			if cc.ID != "" {
				c.Categories = append(c.Categories, cc)
				uniqueCat[uc] = true
			}
			if cl.Code != "" {
				c.Languages = append(c.Languages, cl)
				uniqueLan[ul] = true
			}
			result = append(result, c)
			continue
		}
		if cc.ID != "" && !uniqueCat[uc] {
			uniqueCat[uc] = true
			result[unique[c.ID]].Categories = append(result[unique[c.ID]].Categories, cc)
		}
		if cl.Code != "" && !uniqueLan[ul] {
			uniqueLan[ul] = true
			result[unique[c.ID]].Languages = append(result[unique[c.ID]].Languages, cl)
		}
	}
	c, f := h.cache.Get(model.ChannelKey)
	if !f {
		h.db.Model(&model.Channel{}).Count(&count)
		c = count
		h.cache.SetDefault(model.ChannelKey, count)
	}
	count = c.(int64)

	useMsgPack, _ := strconv.ParseBool(q.Get("msgpack"))
	if useMsgPack {
		paginationSuccessMsgPackResponse(w, result, hasNext(page, size, count), count)
		return
	}
	paginationSuccessJSONResponse(w, result, hasNext(page, size, count), count)
}

func (h *handler) handleGuide(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	result := []model.GuideResponse{}
	var count int64
	page, size := getPageAndSize(q)
	rows, _ := h.db.Model(&model.Guide{}).
		Scopes(paginate(page, size)).
		Select("id,channel_id,site,lang,url,created_at, updated_at, deleted_at").
		Where("channel_id ILIKE ?", searcher(q.Get("channel"))).
		Rows()

	defer rows.Close()
	for rows.Next() {
		var g model.GuideResponse
		rows.Scan(&g.ID, &g.ChannelID, &g.Site, &g.Lang, &g.URL, &g.CreatedAt, &g.UpdatedAt, &g.DeletedAt)
		result = append(result, g)
	}

	c, f := h.cache.Get(model.GuideKey)
	if !f {
		h.db.Model(&model.Guide{}).Count(&count)
		c = count
		h.cache.SetDefault(model.GuideKey, count)
	}
	count = c.(int64)

	useMsgPack, _ := strconv.ParseBool(q.Get("msgpack"))
	if useMsgPack {
		paginationSuccessMsgPackResponse(w, result, hasNext(page, size, count), count)
		return
	}
	paginationSuccessJSONResponse(w, result, hasNext(page, size, count), count)
}
