package db

import (
	"dumbky/internal/log"
)

func GetAllCollections() ([]Collection, error) {
	rows, err := DB.Query(`SELECT id, name, created_at FROM collections`)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var out []Collection
	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			log.Error(err)
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func SaveRequest(request Request) error {
	_, err := DB.Exec(`INSERT INTO requests (collection_name, name, payload) VALUES (?, ?, ?)`,
		request.CollectionName, request.Name, request.Payload)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func LoadRequest(collectionName, requestName string) (Request, error) {
	const query = `
	SELECT r.id, r.collection_name, r.name, r.payload, r.created_at
	FROM requests r
	JOIN collections c ON r.collection_name = c.name
	WHERE c.name = ? AND r.name = ?
	LIMIT 1;
	`

	row := DB.QueryRow(query, collectionName, requestName)

	var request Request
	err := row.Scan(
		&request.ID,
		&request.CollectionName,
		&request.Name,
		&request.Payload,
		&request.CreatedAt,
	)
	if err != nil {
		log.Error(err)
		return Request{}, err
	}

	return request, nil
}
