package db

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
)

func CreateCollection(collectionName string) error {
	_, err := DB.Exec(`REPLACE INTO collections (name) VALUES (?)`,
		collectionName)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil

}

func CreateDefaultCollection() error {
	err := CreateCollection(constants.DB_DEFAULT_COLLECTION_NAME)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

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

func FetchCollectionNames() []string {
	rows, err := DB.Query("SELECT name FROM collections ORDER BY created_at ASC")
	if err != nil {
		log.Error(err)
		return []string{}
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Error(err)
			continue
		}
		names = append(names, name)
	}
	return names
}

func FetchRequestNames(collectionName string) []string {
	rows, err := DB.Query("SELECT name FROM requests WHERE collection_name = ? ORDER BY created_at ASC", collectionName)
	if err != nil {
		log.Error(err)
		return []string{}
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Error(err)
			continue
		}
		names = append(names, name)
	}
	return names
}

func SaveRequest(request Request) error {
	_, err := DB.Exec(`REPLACE INTO requests (collection_name, name, payload) VALUES (?, ?, ?)`,
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
