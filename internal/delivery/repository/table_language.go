package repository

import (
	"context"
	"database/sql"
	"fmt"
)

func (d *Database) GetLanguages(ctx context.Context) ([]*Language, error) {
	query := `SELECT code, lang FROM language ORDER BY lang ASC`
	data, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("database.GetLanguages.QueryContext: %v", err)
	}
	defer data.Close()

	languages, err := scanRowsLanguage(data)
	if err != nil {
		return nil, fmt.Errorf("database.GetLanguages.scanRowsLanguage: %v", err)
	}
	return languages, nil
}

func scanRowsLanguage(rows *sql.Rows) ([]*Language, error) {
	var languages []*Language
	for rows.Next() {
		var language Language
		err := rows.Scan(
			&language.Code,
			&language.Lang,
		)
		if err != nil {
			return nil, err
		}

		languages = append(languages, &language)
	}

	return languages, nil
}
