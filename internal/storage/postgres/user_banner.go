package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
	"github.com/lib/pq"
	"time"
)

func (p *Postgres) GetUserBanner(params models.UserBanner) (*models.BannerWithDetails, error) {
	const op = "postgres.GetUserBanner"

	query := `
		SELECT b.id AS banner_id, bf.feature_id, array_agg(bt.tag_id) AS tag_ids, b.content, b.status, b.created_at, b.updated_at
		FROM banners b
		LEFT JOIN banner_features bf ON b.id = bf.banner_id
		LEFT JOIN banner_tags bt ON b.id = bt.banner_id
		WHERE bf.feature_id = $1 AND bt.tag_id = $2
		GROUP BY b.id, bf.feature_id;
	`

	rows, err := p.db.Queryx(query, params.FeatureID, params.TagID)
	if err != nil {
		return nil, fmt.Errorf("%s:%d", op, err)
	}
	defer rows.Close()

	var bannersWithDetails []models.BannerWithDetails
	for rows.Next() {
		var banner models.BannerWithDetails
		var contentBytes []byte
		var tagIDs []int64
		var status bool
		var createdAt, updatedAt time.Time
		err := rows.Scan(&banner.BannerID, &banner.FeatureID, pq.Array(&tagIDs), &contentBytes, &status, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s:%d", op, err)
		}
		err = json.Unmarshal(contentBytes, &banner.Content)
		if err != nil {
			return nil, fmt.Errorf("%s:%d", op, err)
		}

		banner.TagIDs = tagIDs
		banner.Status = &status
		banner.CreatedAt = &createdAt
		banner.UpdatedAt = &updatedAt
		bannersWithDetails = append(bannersWithDetails, banner)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%d", op, err)
	}

	return &bannersWithDetails[0], nil
}
