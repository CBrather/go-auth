package album

import (
	"database/sql"
	"fmt"
	"log"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func Add(db *sql.DB, album Album) (int64, error) {
	row := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", album.Title, album.Artist, album.Price)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("Album :: Add :: Insertion: %v", err)
	}

	return id, nil
}

func GetByID(db *sql.DB, id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)

	if err := row.Scan(&alb.ID, &alb.Artist, &alb.Title, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("Album :: GetByID :: no album with id foung: %d", id)
		}

		return alb, fmt.Errorf("Album :: GetByID :: error for id %d: %v", id, err)
	}

	return alb, nil
}

func List(db *sql.DB) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("Album :: List :: Scan Rows: %v", err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Album :: List :: Rows: %v", err)
	}

	log.Printf("Albums: %v", albums)

	return albums, nil
}

func ListByArtist(db *sql.DB, name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByPrice %q: %v", name, err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	log.Printf("Albums: %v", albums)

	return albums, nil
}
