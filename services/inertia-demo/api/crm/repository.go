package crm

import (
	"database/sql"
	"errors"

	"github.com/bedrock/services/inertia-demo/api/internal/database"
)

type databaseRepository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) (*databaseRepository, error) {
	if db == nil {
		return nil, errors.New("crm: database connection must not be nil")
	}

	return &databaseRepository{db: db}, nil
}

func (r databaseRepository) ListRecentNotes(limit int) ([]database.Note, error) {
	return database.ListRecentNotes(r.db, limit)
}

func (r databaseRepository) CountContacts() (int, error) {
	return database.CountContacts(r.db)
}

func (r databaseRepository) CountOrganizations() (int, error) {
	return database.CountOrganizations(r.db)
}

func (r databaseRepository) CountNotes() (int, error) {
	return database.CountNotes(r.db)
}

func (r databaseRepository) ListContactsPaginated(search string, favoritesOnly bool, cursor *string, direction string, perPage int) (database.CursorPage[database.Contact], error) {
	return database.ListContactsPaginated(r.db, search, favoritesOnly, cursor, direction, perPage)
}

func (r databaseRepository) GetContact(id int64) (*database.Contact, error) {
	return database.GetContact(r.db, id)
}

func (r databaseRepository) CreateContact(contact database.Contact) (int64, error) {
	return database.CreateContact(r.db, contact)
}

func (r databaseRepository) UpdateContact(id int64, contact database.Contact) error {
	return database.UpdateContact(r.db, id, contact)
}

func (r databaseRepository) DeleteContact(id int64) error {
	return database.DeleteContact(r.db, id)
}

func (r databaseRepository) ToggleContactFavorite(id int64) error {
	return database.ToggleContactFavorite(r.db, id)
}

func (r databaseRepository) ListContactNotes(contactID int64) ([]database.Note, error) {
	return database.ListContactNotes(r.db, contactID)
}

func (r databaseRepository) CreateNote(contactID, userID int64, body string) (int64, error) {
	return database.CreateNote(r.db, contactID, userID, body)
}

func (r databaseRepository) ListOrganizations(search string) ([]database.Organization, error) {
	return database.ListOrganizations(r.db, search)
}

func (r databaseRepository) ListOrganizationsPaginated(search string, page int, perPage int) (database.OffsetPage[database.Organization], error) {
	return database.ListOrganizationsPaginated(r.db, search, page, perPage)
}

func (r databaseRepository) GetOrganization(id int64) (*database.Organization, error) {
	return database.GetOrganization(r.db, id)
}

func (r databaseRepository) UpdateOrganization(id int64, name string) error {
	return database.UpdateOrganization(r.db, id, name)
}

func (r databaseRepository) ListContactsByOrganization(organizationID int64) ([]database.Contact, error) {
	return database.ListContactsByOrganization(r.db, organizationID)
}

func (r databaseRepository) ListContactsByOrgPaginated(organizationID int64, cursor *string, direction string, perPage int) (database.CursorPage[database.Contact], error) {
	return database.ListContactsByOrgPaginated(r.db, organizationID, cursor, direction, perPage)
}
