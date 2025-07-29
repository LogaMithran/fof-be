package connectors

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		limit, _ := strconv.Atoi(q.Get("limit"))

		if page <= 0 {
			page = 1
		}

		if limit <= 0 {
			limit = 10
		}

		return db.Offset((page - 1) * limit).Limit(limit)
	}
}
