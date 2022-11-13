// Filname: internal/data/entries.go

package data

import (
	"context"
	"time"
	"database/sql"
	"errors"
	"fmt"
	//"time"

	"kriol.camerontillett.net/internal/validator"
	//"kriol.camerontillett.net/internal/data"
	"github.com/lib/pq"
)

type Entry struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name string `json:"name"`
	Level string `json:"level"`
	Contact string `json:"contact"`
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`
	Address string `json:"address"`
	Mode []string `json:"mode"`
	Version int32 `json:"version"`
}

func ValidateEntries (v *validator.Validator, entries *Entry) {
	// Check() method to execute
	v.Check(entries.Name != "", "name", "must be provided")
	v.Check(len(entries.Name) <= 200, "name", "must not be more than 200 bytes long")
	
	v.Check(entries.Level != "", "level", "must be provided")
	v.Check(len(entries.Level) <= 200, "level", "must not be more than 200 bytes long")
	
	v.Check(entries.Contact != "", "contact", "must be provided")
	v.Check(len(entries.Contact) <= 200, "contact", "must not be more than 200 bytes long")
	
	v.Check(entries.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(entries.Phone, validator.PhoneRX), "phone", "must be a valid phone number")

	v.Check(entries.Email != "", "email", "must be provided")
	v.Check(validator.Matches(entries.Email, validator.EmailRX), "email", "must be a valid email")

	v.Check(entries.Website != "", "website", "must be provided")
	v.Check(validator.ValidWebsite(entries.Website), "website", "must be a valid url")

	v.Check(entries.Address != "", "address", "must be provided")
	v.Check(len(entries.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(entries.Mode != nil, "mode", "must be provided")
	v.Check(len(entries.Mode) >= 1, "mode", "must contain at least one entries")
	v.Check(len(entries.Mode) <= 5, "mode", "must contain at most 5 entries")
	v.Check(validator.Unique(entries.Mode), "mode", "must not contain duplicate entries")
}

// Define a Entries Model to wrap the sql.db connection pool
type EntryModel struct {
	DB *sql.DB
}

// Allows us to create a new Entry
func (m EntryModel) Insert(entries *Entry) error {
	query := `
		INSERT INTO entries (name, level, contact, phone, email, website, address, mode)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, version
	`
	// Create a context
	// Time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	//Collect the data fields into a slice
	args := []interface{}{
		entries.Name, entries.Level,
		entries.Contact, entries.Phone,
		entries.Email, entries.Website,
		entries.Address, pq.Array(entries.Mode),
	}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&entries.ID, &entries.CreatedAt, &entries.Version)
}

// Get () Allows us to retrieve a specific entry
func (m EntryModel) Get(id int64) (*Entry, error) {
	// Ensure that there is a valid ID
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query
	query := `
		SELECT id, created_at, name, level, contact, phone, email, website, address, mode, version
		FROM entries
		WHERE id = $1
	`
	// Declare a Entry variable to hold the returned data
	var entries Entry
	// Create a context
	// Time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&entries.ID,
		&entries.CreatedAt,
		&entries.Name,
		&entries.Level,
		&entries.Contact,
		&entries.Phone,
		&entries.Email,
		&entries.Website,
		&entries.Address,
		pq.Array(entries.Mode),
		&entries.Version,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &entries, nil
}

// Update() allows us to edit/alter a specific Entry
// KEY: Go's http.Server handles each request in its own goroutine
//Avoid data races
// OPtimistic locking (version number)
func (m EntryModel) Update(entries *Entry) error {
	// Create a query
	query := `
		UPDATE entries
		SET name = $1, 	  level = $2, contact = $3,
		    phone = $4,   email = $5, website = $6,
			address = $7, mode = $8,  version = version + 1
		WHERE id = $9
		AND version = $10
		RETURNING version
	`

	// Create a context
	// Time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	args := []interface{}{
		entries.Name,
		entries.Level,
		entries.Contact,
		entries.Phone,
		entries.Email,
		entries.Website,
		entries.Address,
		pq.Array(entries.Mode),
		entries.ID,
		entries.Version,
	}
	// Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&entries.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete() removes a specific Entry
func (m EntryModel) Delete(id int64) error {
	// Ensure that there is a valid ID
	if id < 1 {
		return ErrRecordNotFound
	}

	// Create the delete query
	query := `
		DELETE FROM entries
		WHERE id = $1
	`

	// Create a context
	// Time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by deleton
	// Call the RowsAffected()
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// The GetAll() method returns a list of all the schools sorted by id
func (m EntryModel) GetAll (name string, level string, mode []string, filters Filters) ([]*Entry, Metadata, error) {
	// Construst the query
	query := fmt.Sprintf (`
		SELECT COUNT(*) OVER(), id, created_at, name, level, 
				contact, phone, email, website, 
				address, mode, version
		FROM entries
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', level) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (mode @> $3 OR $3 = '{}' )
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())


	// Created a 3-second-timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	// Execute the query
	args := []interface{}{name, level, pq.Array(mode), filters.limit(), filters.offset()}
	// Execute the query
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Close the result set
	defer rows.Close()

	totalRecords := 0

	// Initialize an empty slice to 
	entry := []*Entry{}
	// iterate over the rpws in the resultset
	for rows.Next() {
		var entries Entry
		// Scan the values from the row into entry
		err := rows.Scan (
			&totalRecords,
			&entries.ID,
			&entries.CreatedAt,
			&entries.Name,
			&entries.Level,
			&entries.Contact,
			&entries.Phone,
			&entries.Email,
			&entries.Website,
			&entries.Address,
			pq.Array(&entries.Mode),
			&entries.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Entry to our slice
		entry = append(entry, &entries)
	}
	// Check for errors after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Return the slice of schools
	return entry, metadata, nil
}
