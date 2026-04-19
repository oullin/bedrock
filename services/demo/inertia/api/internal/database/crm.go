package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrNotFound is returned when an UPDATE or DELETE affects zero rows.

// CursorPage holds cursor-paginated results.

// OffsetPage holds offset-paginated results.

type CursorPage[T any] struct {
	Data       []T     `json:"data"`
	NextCursor *string `json:"next_cursor"`
	PrevCursor *string `json:"prev_cursor"`
}

type OffsetPage[T any] struct {
	Data        []T `json:"data"`
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string `json:"-"`
	VerifiedAt   *time.Time
}

type Organization struct {
	ID            int64
	Name          string
	ContactsCount int
}

type Contact struct {
	ID               int64
	OrganizationID   *int64
	OrganizationName string
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	IsFavorite       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Note struct {
	ID          int64
	ContactID   int64
	ContactName string
	UserID      int64
	UserName    string
	Body        string
	CreatedAt   time.Time
}

var ErrNotFound = errors.New("database: record not found")

func checkRowsAffected(result sql.Result) error {
	n, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func CreateUser(db *sql.DB, name, email, password string, verifiedAt *time.Time) (int64, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return 0, err
	}

	result, err := db.Exec(
		"INSERT INTO users (name, email, password, verified_at) VALUES (?, ?, ?, ?)",
		name,
		email,
		string(hashedPassword),
		verifiedAt,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func FindUserByEmail(db *sql.DB, email string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	row := db.QueryRow(`
		SELECT id, name, email, password, verified_at
		FROM users
		WHERE lower(email) = lower(?)
		LIMIT 1
	`, email)

	return scanUser(row)
}

func FindUserByID(db *sql.DB, id int64) (*User, error) {
	row := db.QueryRow(`
		SELECT id, name, email, password, verified_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id)

	return scanUser(row)
}

func ListOrganizations(db *sql.DB, search string) ([]Organization, error) {
	query := `
		SELECT o.id, o.name, COUNT(c.id) AS contacts_count
		FROM organizations o
		LEFT JOIN contacts c ON c.organization_id = o.id
	`

	var args []any

	if search = strings.TrimSpace(search); strings.TrimSpace(search) != "" {
		query += " WHERE lower(o.name) LIKE lower(?)"
		args = append(args, "%"+search+"%")
	}

	query += " GROUP BY o.id, o.name ORDER BY o.name ASC"

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Organization, error) {
		var org Organization

		if err := scan(&org.ID, &org.Name, &org.ContactsCount); err != nil {
			return Organization{}, err
		}

		return org, nil
	}), nil
}

func GetOrganization(db *sql.DB, id int64) (*Organization, error) {
	row := db.QueryRow(`
		SELECT o.id, o.name, COUNT(c.id) AS contacts_count
		FROM organizations o
		LEFT JOIN contacts c ON c.organization_id = o.id
		WHERE o.id = ?
		GROUP BY o.id, o.name
	`, id)

	var org Organization

	if err := row.Scan(&org.ID, &org.Name, &org.ContactsCount); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &org, nil
}

func CreateOrganization(db *sql.DB, name string) (int64, error) {
	result, err := db.Exec("INSERT INTO organizations (name) VALUES (?)", strings.TrimSpace(name))

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateOrganization(db *sql.DB, id int64, name string) error {
	result, err := db.Exec(
		"UPDATE organizations SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		strings.TrimSpace(name),
		id,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(result)
}

func CountOrganizations(db *sql.DB) (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM organizations").Scan(&count)

	return count, err
}

func ListContacts(db *sql.DB, search string, favoritesOnly bool) ([]Contact, error) {
	query := `
		SELECT
			c.id,
			c.organization_id,
			COALESCE(o.name, ''),
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.is_favorite,
			c.created_at,
			c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE 1 = 1
	`

	var args []any

	if search = strings.TrimSpace(search); strings.TrimSpace(search) != "" {
		query += `
			AND (
				lower(c.first_name) LIKE lower(?)
				OR lower(c.last_name) LIKE lower(?)
				OR lower(c.email) LIKE lower(?)
				OR lower(COALESCE(o.name, '')) LIKE lower(?)
			)
		`
		like := "%" + search + "%"
		args = append(args, like, like, like, like)
	}

	if favoritesOnly {
		query += " AND c.is_favorite = 1"
	}

	query += " ORDER BY c.first_name ASC, c.last_name ASC"

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Contact, error) {
		return scanContact(scan)
	}), nil
}

func ListRecentContacts(db *sql.DB, limit int) ([]Contact, error) {
	rows, err := db.Query(`
		SELECT
			c.id,
			c.organization_id,
			COALESCE(o.name, ''),
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.is_favorite,
			c.created_at,
			c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		ORDER BY c.first_name ASC, c.last_name ASC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Contact, error) {
		return scanContact(scan)
	}), nil
}

func ListContactsByOrganization(db *sql.DB, organizationID int64) ([]Contact, error) {
	rows, err := db.Query(`
		SELECT
			c.id,
			c.organization_id,
			COALESCE(o.name, ''),
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.is_favorite,
			c.created_at,
			c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE c.organization_id = ?
		ORDER BY c.first_name ASC, c.last_name ASC
	`, organizationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Contact, error) {
		return scanContact(scan)
	}), nil
}

func GetContact(db *sql.DB, id int64) (*Contact, error) {
	row := db.QueryRow(`
		SELECT
			c.id,
			c.organization_id,
			COALESCE(o.name, ''),
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.is_favorite,
			c.created_at,
			c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE c.id = ?
		LIMIT 1
	`, id)

	contact, err := scanContact(row.Scan)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &contact, nil
}

func CreateContact(db *sql.DB, contact Contact) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO contacts (
			organization_id,
			first_name,
			last_name,
			email,
			phone,
			is_favorite
		) VALUES (?, ?, ?, ?, ?, ?)
	`,
		nullableInt64(contact.OrganizationID),
		strings.TrimSpace(contact.FirstName),
		strings.TrimSpace(contact.LastName),
		strings.TrimSpace(contact.Email),
		strings.TrimSpace(contact.Phone),
		boolToInt(contact.IsFavorite),
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateContact(db *sql.DB, id int64, contact Contact) error {
	result, err := db.Exec(`
		UPDATE contacts
		SET organization_id = ?,
			first_name = ?,
			last_name = ?,
			email = ?,
			phone = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		nullableInt64(contact.OrganizationID),
		strings.TrimSpace(contact.FirstName),
		strings.TrimSpace(contact.LastName),
		strings.TrimSpace(contact.Email),
		strings.TrimSpace(contact.Phone),
		id,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(result)
}

func ToggleContactFavorite(db *sql.DB, id int64) error {
	result, err := db.Exec(`
		UPDATE contacts
		SET is_favorite = CASE WHEN is_favorite = 1 THEN 0 ELSE 1 END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)

	if err != nil {
		return err
	}

	return checkRowsAffected(result)
}

func DeleteContact(db *sql.DB, id int64) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM notes WHERE contact_id = ?", id); err != nil {
		return err
	}

	result, err := tx.Exec("DELETE FROM contacts WHERE id = ?", id)

	if err != nil {
		return err
	}

	if err := checkRowsAffected(result); err != nil {
		return err
	}

	return tx.Commit()
}

func ListContactsPaginated(db *sql.DB, search string, favoritesOnly bool, cursor *string, direction string, perPage int) (CursorPage[Contact], error) {
	if perPage <= 0 {
		return CursorPage[Contact]{}, fmt.Errorf("perPage must be greater than 0, got %d", perPage)
	}

	query := `
		SELECT
			c.id, c.organization_id, COALESCE(o.name, ''),
			c.first_name, c.last_name, c.email, c.phone,
			c.is_favorite, c.created_at, c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE 1 = 1
	`

	var args []any

	if search = strings.TrimSpace(search); strings.TrimSpace(search) != "" {
		query += `
			AND (
				lower(c.first_name) LIKE lower(?)
				OR lower(c.last_name) LIKE lower(?)
				OR lower(c.email) LIKE lower(?)
				OR lower(COALESCE(o.name, '')) LIKE lower(?)
			)
		`
		like := "%" + search + "%"
		args = append(args, like, like, like, like)
	}

	if favoritesOnly {
		query += " AND c.is_favorite = 1"
	}

	if cursor != nil && strings.TrimSpace(*cursor) != "" {
		cursorID, err := strconv.ParseInt(*cursor, 10, 64)

		if err != nil || cursorID <= 0 {
			return CursorPage[Contact]{}, fmt.Errorf("invalid cursor: %q", *cursor)
		}

		if direction == "prev" {
			query += " AND c.id < ?"
		} else {
			query += " AND c.id > ?"
		}

		args = append(args, cursorID)
	}

	if direction == "prev" {
		query += " ORDER BY c.id DESC LIMIT ?"
	} else {
		query += " ORDER BY c.id ASC LIMIT ?"
	}

	args = append(args, perPage+1)

	rows, err := db.Query(query, args...)

	if err != nil {
		return CursorPage[Contact]{}, err
	}

	defer rows.Close()

	all := scanRows(rows, func(scan func(...any) error) (Contact, error) {
		return scanContact(scan)
	})

	hasExtra := len(all) > perPage

	if hasExtra {
		all = all[:perPage]
	}

	if direction == "prev" {
		for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
			all[i], all[j] = all[j], all[i]
		}
	}

	page := CursorPage[Contact]{Data: all}

	if len(all) > 0 {
		if (direction != "prev" && hasExtra) || (direction == "prev" && cursor != nil && strings.TrimSpace(*cursor) != "") {
			next := fmt.Sprintf("%d", all[len(all)-1].ID)
			page.NextCursor = &next
		}

		if (direction == "prev" && hasExtra) || (direction != "prev" && cursor != nil && strings.TrimSpace(*cursor) != "") {
			prev := fmt.Sprintf("%d", all[0].ID)
			page.PrevCursor = &prev
		}
	}

	return page, nil
}

func ListOrganizationsPaginated(db *sql.DB, search string, page int, perPage int) (OffsetPage[Organization], error) {
	if perPage <= 0 {
		return OffsetPage[Organization]{}, fmt.Errorf("perPage must be greater than 0, got %d", perPage)
	}

	countQuery := `
		SELECT COUNT(*)
		FROM organizations o
	`

	query := `
		SELECT o.id, o.name, COUNT(c.id) AS contacts_count
		FROM organizations o
		LEFT JOIN contacts c ON c.organization_id = o.id
	`

	var args []any

	var countArgs []any

	if search = strings.TrimSpace(search); strings.TrimSpace(search) != "" {
		countQuery += " WHERE lower(o.name) LIKE lower(?)"
		query += " WHERE lower(o.name) LIKE lower(?)"
		countArgs = append(countArgs, "%"+search+"%")
		args = append(args, "%"+search+"%")
	}

	var total int

	if err := db.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		return OffsetPage[Organization]{}, err
	}

	if page < 1 {
		page = 1
	}

	lastPage := (total + perPage - 1) / perPage

	if lastPage < 1 {
		lastPage = 1
	}

	if page > lastPage {
		page = lastPage
	}

	offset := (page - 1) * perPage
	query += " GROUP BY o.id, o.name ORDER BY o.name ASC LIMIT ? OFFSET ?"
	args = append(args, perPage, offset)

	rows, err := db.Query(query, args...)

	if err != nil {
		return OffsetPage[Organization]{}, err
	}

	defer rows.Close()

	orgs := scanRows(rows, func(scan func(...any) error) (Organization, error) {
		var org Organization

		if err := scan(&org.ID, &org.Name, &org.ContactsCount); err != nil {
			return Organization{}, err
		}

		return org, nil
	})

	return OffsetPage[Organization]{
		Data:        orgs,
		Total:       total,
		PerPage:     perPage,
		CurrentPage: page,
		LastPage:    lastPage,
	}, nil
}

func ListContactsByOrgPaginated(db *sql.DB, organizationID int64, cursor *string, direction string, perPage int) (CursorPage[Contact], error) {
	if perPage <= 0 {
		return CursorPage[Contact]{}, fmt.Errorf("perPage must be greater than 0, got %d", perPage)
	}

	query := `
		SELECT
			c.id, c.organization_id, COALESCE(o.name, ''),
			c.first_name, c.last_name, c.email, c.phone,
			c.is_favorite, c.created_at, c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE c.organization_id = ?
	`

	args := []any{organizationID}

	if cursor != nil && strings.TrimSpace(*cursor) != "" {
		cursorID, err := strconv.ParseInt(*cursor, 10, 64)

		if err != nil || cursorID <= 0 {
			return CursorPage[Contact]{}, fmt.Errorf("invalid cursor: %q", *cursor)
		}

		if direction == "prev" {
			query += " AND c.id < ?"
		} else {
			query += " AND c.id > ?"
		}

		args = append(args, cursorID)
	}

	if direction == "prev" {
		query += " ORDER BY c.id DESC LIMIT ?"
	} else {
		query += " ORDER BY c.id ASC LIMIT ?"
	}

	args = append(args, perPage+1)

	rows, err := db.Query(query, args...)

	if err != nil {
		return CursorPage[Contact]{}, err
	}

	defer rows.Close()

	all := scanRows(rows, func(scan func(...any) error) (Contact, error) {
		return scanContact(scan)
	})

	hasExtra := len(all) > perPage

	if hasExtra {
		all = all[:perPage]
	}

	if direction == "prev" {
		for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
			all[i], all[j] = all[j], all[i]
		}
	}

	result := CursorPage[Contact]{Data: all}

	if len(all) > 0 {
		if (direction != "prev" && hasExtra) || (direction == "prev" && cursor != nil && strings.TrimSpace(*cursor) != "") {
			next := fmt.Sprintf("%d", all[len(all)-1].ID)
			result.NextCursor = &next
		}

		if (direction == "prev" && hasExtra) || (direction != "prev" && cursor != nil && strings.TrimSpace(*cursor) != "") {
			prev := fmt.Sprintf("%d", all[0].ID)
			result.PrevCursor = &prev
		}
	}

	return result, nil
}

func CountContacts(db *sql.DB) (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&count)

	return count, err
}

func CreateNote(db *sql.DB, contactID, userID int64, body string) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO notes (contact_id, user_id, body) VALUES (?, ?, ?)",
		contactID,
		userID,
		strings.TrimSpace(body),
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func CreateNoteAt(db *sql.DB, contactID, userID int64, body string, createdAt time.Time) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO notes (contact_id, user_id, body, created_at) VALUES (?, ?, ?, ?)",
		contactID,
		userID,
		strings.TrimSpace(body),
		createdAt,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func ListContactNotes(db *sql.DB, contactID int64) ([]Note, error) {
	rows, err := db.Query(`
		SELECT
			n.id,
			n.contact_id,
			(c.first_name || ' ' || c.last_name) AS contact_name,
			n.user_id,
			u.name,
			n.body,
			n.created_at
		FROM notes n
		INNER JOIN contacts c ON c.id = n.contact_id
		INNER JOIN users u ON u.id = n.user_id
		WHERE n.contact_id = ?
		ORDER BY n.created_at DESC
	`, contactID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Note, error) {
		var note Note

		if err := scan(
			&note.ID,
			&note.ContactID,
			&note.ContactName,
			&note.UserID,
			&note.UserName,
			&note.Body,
			&note.CreatedAt,
		); err != nil {
			return Note{}, err
		}

		return note, nil
	}), nil
}

func ListRecentNotes(db *sql.DB, limit int) ([]Note, error) {
	rows, err := db.Query(`
		SELECT
			n.id,
			n.contact_id,
			(c.first_name || ' ' || c.last_name) AS contact_name,
			n.user_id,
			u.name,
			n.body,
			n.created_at
		FROM notes n
		INNER JOIN contacts c ON c.id = n.contact_id
		INNER JOIN users u ON u.id = n.user_id
		ORDER BY n.created_at DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Note, error) {
		var note Note

		if err := scan(
			&note.ID,
			&note.ContactID,
			&note.ContactName,
			&note.UserID,
			&note.UserName,
			&note.Body,
			&note.CreatedAt,
		); err != nil {
			return Note{}, err
		}

		return note, nil
	}), nil
}

func CountNotes(db *sql.DB) (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&count)

	return count, err
}

func scanUser(row interface{ Scan(...any) error }) (*User, error) {
	var user User

	var verifiedAt sql.NullTime

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &verifiedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if verifiedAt.Valid {
		user.VerifiedAt = &verifiedAt.Time
	}

	return &user, nil
}

func scanContact(scan func(...any) error) (Contact, error) {
	var contact Contact

	var organizationID sql.NullInt64

	var favorite int

	if err := scan(
		&contact.ID,
		&organizationID,
		&contact.OrganizationName,
		&contact.FirstName,
		&contact.LastName,
		&contact.Email,
		&contact.Phone,
		&favorite,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	); err != nil {
		return Contact{}, err
	}

	if organizationID.Valid {
		id := organizationID.Int64
		contact.OrganizationID = &id
	}

	contact.IsFavorite = favorite == 1

	return contact, nil
}

func nullableInt64(v *int64) any {
	if v == nil {
		return nil
	}

	return *v
}

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}
