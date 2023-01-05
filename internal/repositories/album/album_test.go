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

	expectedInsertRows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	mock.ExpectQuery("INSERT INTO album").WithArgs(newAlbum.Title, newAlbum.Artist, newAlbum.Price).WillReturnRows(expectedInsertRows)

	expectedQueryRows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).AddRow("1", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	mock.ExpectQuery("SELECT \\* FROM album WHERE id \\= \\$1").WithArgs(1).WillReturnRows(expectedQueryRows)

	addedAlbum, err := album.Add(db, newAlbum)
	if err != nil {
		t.Errorf("Unexpected error testing album.Add: %v", err)
	}

	if addedAlbum.ID != "1" || addedAlbum.Artist != newAlbum.Artist || addedAlbum.Title != newAlbum.Title || addedAlbum.Price != newAlbum.Price {
		t.Errorf("The returned album doesn't match the expected data")
	}
}
