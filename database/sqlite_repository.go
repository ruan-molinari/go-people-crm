package database

import (
	"database/sql"
	"errors"
	"go-people-crm/model"
)

var (
    ErrDuplicate    = errors.New("record already exists")
    ErrNotExists    = errors.New("row not exists")
    ErrUpdateFailed = errors.New("update failed")
    ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS person(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone TEXT,
		last_time_met TEXT,
		meeting_frequency INTEGER
	);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(person model.Person) (*model.Person, error) {
	res, err := r.db.Exec(`
		INSERT INTO person(
			name,
			phone,
			last_time_met,
			meeting_frequency
		) values(?, ?, ?, ?)
	`, person.Name, person.Phone, person.LastTiemMet, person.MeetingFrequecy)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	person.ID = id

	return &person, nil
}

func (r *SQLiteRepository) All() ([]model.Person, error) {
	rows, err := r.db.Query("SELECT * FROM person")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []model.Person
	for rows.Next() {
		var person model.Person
		if err := rows.Scan(
				&person.ID,
				&person.Name,
				&person.Phone,
				&person.LastTiemMet,
				&person.MeetingFrequecy,
			); err != nil {
			return nil, err
		}

		all = append(all, person)
	}

	return all, nil
}

func (r *SQLiteRepository) GetByName(name string) (*model.Person, error) {
	row := r.db.QueryRow("SELECT * FROM person WHERE name = ?", name)

	var person model.Person
	if err := row.Scan(
			&person.ID,
			&person.Name,
			&person.Phone,
			&person.LastTiemMet,
			&person.MeetingFrequecy,
		); err != nil {
		return nil, err
	}

	return &person, nil
}

func (r *SQLiteRepository) Update(updatedPerson model.Person) (*model.Person, error) {
	if updatedPerson.ID == 0 {
		return nil, errors.New("invalid updated ID")
	}
	res, err := r.db.Exec(`
		UPDATE person SET
			name = ?,
			phone = ?,
			last_time_met = ?,
			meeting_frequency = ?
		WHERE id = ?
		`,
		updatedPerson.Name,
		updatedPerson.Phone,
		updatedPerson.LastTiemMet,
		updatedPerson.MeetingFrequecy,
		updatedPerson.ID,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updatedPerson, nil
}

func (r *SQLiteRepository) Upsert(person model.Person) (*model.Person, error) {
	if person.ID == 0 {
		return r.Create(person)
	} else {
		return r.Update(person)
	}
}

func (r *SQLiteRepository) Delete(id int64) error {
	res, err := r.db.Exec("DELETE FROM person WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
