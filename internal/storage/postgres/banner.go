package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
	"github.com/lib/pq"
	"time"
)

func (p *Postgres) GetBannersWithDetails(params models.GetAllBannersParams) ([]models.BannerWithDetails, error) {
	const op = "postgres.GetBannersWithDetails"

	query := `
		SELECT b.id AS banner_id, bf.feature_id, array_agg(bt.tag_id) AS tag_ids, b.content, b.status, b.created_at, b.updated_at
		FROM banners b
		LEFT JOIN banner_features bf ON b.id = bf.banner_id
		LEFT JOIN banner_tags bt ON b.id = bt.banner_id
	`

	args := make([]interface{}, 0)
	argIndex := 1
	if params.FeatureID != nil {
		query += fmt.Sprintf(" WHERE bf.feature_id = $%d", argIndex)
		args = append(args, *params.FeatureID)
		argIndex++
	}
	if params.TagID != nil {
		if params.FeatureID != nil {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += fmt.Sprintf(" bt.tag_id = $%d", argIndex)
		args = append(args, *params.TagID)
		argIndex++
	}

	if params.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, *params.Limit)
		argIndex++
	}
	if params.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, *params.Offset)
	}

	query += " GROUP BY b.id, bf.feature_id;"

	rows, err := p.db.Queryx(query, args...)
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

	return bannersWithDetails, nil
}

func (p *Postgres) CreateBanner(newBanner models.CreateBannerRequest) error {
	const op = "postgres.CreateBanner"

	newBannerContent, err := json.Marshal(newBanner.NewBanner)

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	defer tx.Rollback()

	// TODO: insert is_active
	query := "INSERT INTO banners (content, status) VALUES ($1, $2) RETURNING id;"

	var bannerId int
	err = tx.QueryRow(query, newBannerContent, newBanner.IsActive).Scan(&bannerId)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if newBanner.FeatureID != "" {
		query = "INSERT INTO banner_features (banner_id, feature_id) VALUES ($1, $2);"
		_, err := tx.Exec(query, bannerId, newBanner.FeatureID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	}

	if len(newBanner.TagIds) != 0 {
		query := "INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2);"
		for _, tagId := range newBanner.TagIds {
			_, err = tx.Exec(query, bannerId, tagId)
			if err != nil {
				return fmt.Errorf("%s: %v", op, err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (p *Postgres) ChangeBanner(bannerID int64, req models.ChangeBannerRequest) error {
	const op = "postgres.ChangeBanner"

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	defer tx.Rollback()

	if req.NewContent != nil && (req.NewIsActive || !req.NewIsActive) {

		query := "UPDATE banners SET content = $1, status = $2 WHERE id = $3"
		newContentJSON, err := json.Marshal(req.NewContent)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}

		_, err = tx.Exec(query, newContentJSON, req.NewIsActive, bannerID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	} else if req.NewContent == nil && (req.NewIsActive || !req.NewIsActive) {
		query := "UPDATE banners SET status = $1 WHERE id = $2"
		_, err = tx.Exec(query, req.NewIsActive, bannerID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	}

	if len(req.NewTagIDs) > 0 {
		query := `
			DELETE FROM banner_tags
			WHERE banner_id = $1
		`
		_, err = tx.Exec(query, bannerID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}

		query = `
			INSERT INTO banner_tags (banner_id, tag_id)
			VALUES ($1, $2)
		`

		for _, tagID := range req.NewTagIDs {
			_, err := tx.Exec(query, bannerID, tagID)
			if err != nil {
				return fmt.Errorf("%s: %v", op, err)
			}
		}
	}

	if req.NewFeatureID > 0 {
		query := `
			DELETE FROM banner_features
			WHERE banner_id = $1
		`

		_, err = tx.Exec(query, bannerID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}

		query = `
			INSERT INTO banner_features (banner_id, feature_id)
			VALUES ($1, $2)
		`

		_, err = tx.Exec(query, bannerID, req.NewFeatureID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (p *Postgres) DeleteBanner(bannerID int64) error {
	const op = "postgres.DeleteBanner"

	query := `
		DELETE FROM banner_features
		WHERE banner_id = $1
	`

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, bannerID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	query = `
		DELETE FROM banner_tags
		WHERE banner_id = $1
	`

	_, err = tx.Exec(query, bannerID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	query = `
		DELETE FROM banners
		WHERE id = $1
	`

	_, err = tx.Exec(query, bannerID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (p *Postgres) DeleteBannerByParams(tagID, featureID int64) error {
	const op = "postgres.DeleteBannerByParams"

	query := `
		SELECT banners.id
		FROM banners
		JOIN banner_tags ON banners.id = banner_tags.banner_id
		JOIN banner_features ON banners.id = banner_features.banner_id
		WHERE banner_tags.tag_id = $1 AND banner_features.feature_id = $2;
	`

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	defer tx.Rollback()
	var bannerId int
	err = tx.QueryRow(query, tagID, featureID).Scan(&bannerId)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	query = `
		DELETE FROM banner_tags
		WHERE banner_id = $1
	`

	_, err = tx.Exec(query, bannerId)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	query = `
		DELETE FROM banner_features 
		WHERE banner_id = $1
	`

	_, err = tx.Exec(query, bannerId)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	query = `
		DELETE FROM banners
		WHERE id = $1;
	`

	_, err = tx.Exec(query, bannerId)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (p *Postgres) CheckBanner(bannerID, tagID, featureID int64) error {
	const op = "postgres.CheckBanner"

	if bannerID == 0 {
		query := `
		SELECT banners.id
		FROM banners
		JOIN banner_tags ON banners.id = banner_tags.banner_id
		JOIN banner_features ON banners.id = banner_features.banner_id
		WHERE banner_tags.tag_id = $1 AND banner_features.feature_id = $2;
		`
		var resultID int64
		err := p.db.Get(&resultID, query, tagID, featureID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
		return nil
	} else {
		query := `
		SELECT id FROM banners
		WHERE id = $1
	`
		var resultID int64
		err := p.db.Get(&resultID, query, bannerID)
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}

		return nil
	}
	return nil
}
