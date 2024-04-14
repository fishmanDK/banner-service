package models

import "time"

type GetAllBannersParams struct {
	FeatureID *int
	TagID     *int
	Limit     *int
	Offset    *int
}

type BannerWithDetails struct {
	BannerID  int                    `db:"id"`
	FeatureID *int                   `db:"tag_ids"`
	TagIDs    []int64                `db:"feature_id"`
	Content   map[string]interface{} `db:"content"`
	Status    *bool                  `db:"status"`
	CreatedAt *time.Time             `db:"created_at"`
	UpdatedAt *time.Time             `db:"updated_at"`
}

//type BannerWithDetails struct {
//	BannerID  int                    `db:"id"`         // Исправлено на "id" для соответствия столбцу banners.id
//	TagIDs    []int64                `db:"tag_ids"`    // Исправлено для соответствия результату array_agg
//	FeatureID *int                   `db:"feature_id"` // Исправлено для соответствия столбцу banner_features.feature_id
//	Content   map[string]interface{} `db:"content"`    // Исправлено для соответствия столбцу banners.content
//	Status    *bool                  `db:"status"`     // Исправлено для соответствия столбцу banners.status
//	CreatedAt *time.Time             `db:"created_at"` // Исправлено для соответствия столбцу banners.created_at
//	UpdatedAt *time.Time             `db:"updated_at"` // Исправлено для соответствия столбцу banners.updated_at
//}

type ChangeBannerRequest struct {
	NewTagIDs    []int64                 `json:"new_tag_ids,omitempty"`
	NewFeatureID int64                   `json:"new_feature_id,omitempty"`
	NewContent   *map[string]interface{} `json:"new_content,omitempty"`
	NewIsActive  bool                    `json:"new_is_active,omitempty"`
}

type CreateBannerRequest struct {
	TagIds    []string               `json:"tag_ids"`
	FeatureID string                 `json:"feature_id"`
	NewBanner map[string]interface{} `json:"new_banner"`
	IsActive  string                 `json:"is_active"`
}

type Banner struct {
	ID        int
	TagIDs    []int
	FeatureID int
	Content   map[string]interface{}
	IsActive  bool
	CreatedAt string
	UpdatedAt string
}
