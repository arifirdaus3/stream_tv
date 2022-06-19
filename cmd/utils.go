package main

import (
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

func getPageAndSize(q url.Values) (int, int) {
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(q.Get("size"))
	if size < 1 {
		size = 1000
	}
	if size > 1000 {
		size = 1000
	}
	return page, size
}

func paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

func searcher(in string) string {
	return "%" + in + "%"
}

func hasNext(page, size int, count int64) bool {
	return int64((page * size)) < count
}
