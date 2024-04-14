package postgres

import (
	"fmt"
	"github.com/fishmanDK/avito_test_task/models"
)

func (p *Postgres) CreateTag(tag models.Tag) error {
	const op = "postgres.CreateTag"

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	defer tx.Rollback()

	query := "INSERT INTO tags (name) VALUES ($1);"

	_, err = tx.Exec(query, tag.Name)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (p *Postgres) CreateFeature(feature models.Feature) error {
	const op = "postgres.CreateFeature"

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	defer tx.Rollback()

	query := "INSERT INTO features (name) VALUES ($1);"

	_, err = tx.Exec(query, feature.Name)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
