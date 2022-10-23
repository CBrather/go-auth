package album_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/CBrather/go-auth/internal/repositories/album"
)

func TestShouldAddAlbum(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error trying to create database mock: %v", err)
	}

	defer db.Close()

	newAlbum := album.Album{
		Title:  "Appetite for Destruction",
		Artist: "Guns 'n Roses",
		Price:  19.99,
	}

	expectedRows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	mock.ExpectQuery("INSERT INTO album").WithArgs(newAlbum.Title, newAlbum.Artist, newAlbum.Price).WillReturnRows(expectedRows)

	id, err := album.Add(db, newAlbum)
	if err != nil {
		t.Errorf("Unexpected error testing album.Add: %v", err)
	}

	if id != 1 {
		t.Errorf("Not all expectations were met:")
	}
}
