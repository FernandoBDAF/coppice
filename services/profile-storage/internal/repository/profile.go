package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"microservices/services/profile-storage/internal/logger"
	"microservices/services/profile-storage/internal/models"
)

// ProfileRepository handles database operations for profiles
type ProfileRepository struct {
	db  *sqlx.DB
	log *zap.Logger
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{
		db:  db,
		log: logger.Get(),
	}
}

// Create creates a new profile with its addresses and contacts
func (r *ProfileRepository) Create(ctx context.Context, profile *models.Profile) error {
	r.log.Info("Creating new profile",
		logger.String("profile_id", profile.ID.String()),
		logger.String("email", profile.Email),
	)

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.log.Error("Failed to begin transaction",
			logger.ErrorField(err),
			logger.String("profile_id", profile.ID.String()),
		)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			r.log.Error("Transaction panic",
				zap.Any("panic", p),
				logger.String("profile_id", profile.ID.String()),
			)
			panic(p)
		}
	}()

	// Create profile
	query := `
		INSERT INTO profiles (id, first_name, last_name, email, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	err = tx.QueryRowxContext(ctx, query,
		profile.ID, profile.FirstName, profile.LastName, profile.Email, profile.Phone,
	).Scan(&profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.log.Error("Transaction rollback failed",
				logger.ErrorField(err),
				logger.ErrorField(rbErr),
				logger.String("profile_id", profile.ID.String()),
			)
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		r.log.Error("Failed to create profile",
			logger.ErrorField(err),
			logger.String("profile_id", profile.ID.String()),
		)
		return fmt.Errorf("error creating profile: %w", err)
	}

	// Create addresses
	for i := range profile.Addresses {
		addr := &profile.Addresses[i]
		addr.ProfileID = profile.ID
		if err := r.createAddressTx(ctx, tx, addr); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				r.log.Error("Transaction rollback failed",
					logger.ErrorField(err),
					logger.ErrorField(rbErr),
					logger.String("profile_id", profile.ID.String()),
				)
				return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
			return err
		}
	}

	// Create contacts
	for i := range profile.Contacts {
		contact := &profile.Contacts[i]
		contact.ProfileID = profile.ID
		if err := r.createContactTx(ctx, tx, contact); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				r.log.Error("Transaction rollback failed",
					logger.ErrorField(err),
					logger.ErrorField(rbErr),
					logger.String("profile_id", profile.ID.String()),
				)
				return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.log.Error("Failed to commit transaction",
			logger.ErrorField(err),
			logger.String("profile_id", profile.ID.String()),
		)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Info("Successfully created profile",
		logger.String("profile_id", profile.ID.String()),
		logger.String("email", profile.Email),
	)
	return nil
}

// GetByID retrieves a profile by its ID
func (r *ProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Profile, error) {
	r.log.Debug("Getting profile by ID",
		logger.String("profile_id", id.String()),
	)

	query := `
		SELECT id, first_name, last_name, email, phone, created_at, updated_at
		FROM profiles
		WHERE id = $1`

	var profile models.Profile
	err := r.db.GetContext(ctx, &profile, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("Profile not found",
				logger.String("profile_id", id.String()),
			)
			return nil, nil
		}
		r.log.Error("Failed to get profile",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return nil, fmt.Errorf("error getting profile: %w", err)
	}

	// Get addresses
	addresses, err := r.getAddresses(ctx, id)
	if err != nil {
		return nil, err
	}
	profile.Addresses = addresses

	// Get contacts
	contacts, err := r.getContacts(ctx, id)
	if err != nil {
		return nil, err
	}
	profile.Contacts = contacts

	r.log.Debug("Successfully retrieved profile",
		logger.String("profile_id", id.String()),
		logger.String("email", profile.Email),
	)
	return &profile, nil
}

// Update updates an existing profile
func (r *ProfileRepository) Update(ctx context.Context, profile *models.Profile) error {
	r.log.Info("Updating profile",
		logger.String("profile_id", profile.ID.String()),
		logger.String("email", profile.Email),
	)

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.log.Error("Failed to begin transaction",
			logger.ErrorField(err),
			logger.String("profile_id", profile.ID.String()),
		)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			r.log.Error("Transaction panic",
				zap.Any("panic", p),
				logger.String("profile_id", profile.ID.String()),
			)
			panic(p)
		}
	}()

	// Update profile
	query := `
		UPDATE profiles
		SET first_name = $1, last_name = $2, email = $3, phone = $4
		WHERE id = $5
		RETURNING updated_at`

	err = tx.QueryRowxContext(ctx, query,
		profile.FirstName, profile.LastName, profile.Email, profile.Phone, profile.ID,
	).Scan(&profile.UpdatedAt)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.log.Error("Transaction rollback failed",
				logger.ErrorField(err),
				logger.ErrorField(rbErr),
				logger.String("profile_id", profile.ID.String()),
			)
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		r.log.Error("Failed to update profile",
			logger.ErrorField(err),
			logger.String("profile_id", profile.ID.String()),
		)
		return fmt.Errorf("error updating profile: %w", err)
	}

	// Update addresses
	if err := r.updateAddressesTx(ctx, tx, profile.ID, profile.Addresses); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.log.Error("Transaction rollback failed",
				logger.ErrorField(err),
				logger.ErrorField(rbErr),
				logger.String("profile_id", profile.ID.String()),
			)
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Update contacts
	if err := r.updateContactsTx(ctx, tx, profile.ID, profile.Contacts); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.log.Error("Transaction rollback failed",
				logger.ErrorField(err),
				logger.ErrorField(rbErr),
				logger.String("profile_id", profile.ID.String()),
			)
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.log.Error("Failed to commit transaction",
			logger.ErrorField(err),
			logger.String("profile_id", profile.ID.String()),
		)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Info("Successfully updated profile",
		logger.String("profile_id", profile.ID.String()),
		logger.String("email", profile.Email),
	)
	return nil
}

// Delete deletes a profile and its related data
func (r *ProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.log.Info("Deleting profile",
		logger.String("profile_id", id.String()),
	)

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.log.Error("Failed to begin transaction",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			r.log.Error("Transaction panic",
				zap.Any("panic", p),
				logger.String("profile_id", id.String()),
			)
			panic(p)
		}
	}()

	// Delete profile (cascade will handle addresses and contacts)
	query := `DELETE FROM profiles WHERE id = $1`
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.log.Error("Transaction rollback failed",
				logger.ErrorField(err),
				logger.ErrorField(rbErr),
				logger.String("profile_id", id.String()),
			)
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		r.log.Error("Failed to delete profile",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return fmt.Errorf("error deleting profile: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.log.Error("Failed to get rows affected",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		r.log.Warn("Profile not found for deletion",
			logger.String("profile_id", id.String()),
		)
		return nil
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.log.Error("Failed to commit transaction",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Info("Successfully deleted profile",
		logger.String("profile_id", id.String()),
	)
	return nil
}

// Helper methods for addresses and contacts with transaction support

func (r *ProfileRepository) createAddressTx(ctx context.Context, tx *sqlx.Tx, addr *models.Address) error {
	query := `
		INSERT INTO addresses (id, profile_id, street, city, state, country, postal_code, is_primary)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`

	return tx.QueryRowxContext(ctx, query,
		addr.ID, addr.ProfileID, addr.Street, addr.City, addr.State,
		addr.Country, addr.PostalCode, addr.IsPrimary,
	).Scan(&addr.CreatedAt, &addr.UpdatedAt)
}

func (r *ProfileRepository) createContactTx(ctx context.Context, tx *sqlx.Tx, contact *models.Contact) error {
	query := `
		INSERT INTO contacts (id, profile_id, type, value, is_primary)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	return tx.QueryRowxContext(ctx, query,
		contact.ID, contact.ProfileID, contact.Type, contact.Value, contact.IsPrimary,
	).Scan(&contact.CreatedAt, &contact.UpdatedAt)
}

func (r *ProfileRepository) getAddresses(ctx context.Context, profileID uuid.UUID) ([]models.Address, error) {
	query := `
		SELECT id, profile_id, street, city, state, country, postal_code, is_primary, created_at, updated_at
		FROM addresses
		WHERE profile_id = $1`

	var addresses []models.Address
	err := r.db.SelectContext(ctx, &addresses, query, profileID)
	if err != nil {
		return nil, fmt.Errorf("error getting addresses: %w", err)
	}

	return addresses, nil
}

func (r *ProfileRepository) getContacts(ctx context.Context, profileID uuid.UUID) ([]models.Contact, error) {
	query := `
		SELECT id, profile_id, type, value, is_primary, created_at, updated_at
		FROM contacts
		WHERE profile_id = $1`

	var contacts []models.Contact
	err := r.db.SelectContext(ctx, &contacts, query, profileID)
	if err != nil {
		return nil, fmt.Errorf("error getting contacts: %w", err)
	}

	return contacts, nil
}

func (r *ProfileRepository) updateAddressesTx(ctx context.Context, tx *sqlx.Tx, profileID uuid.UUID, addresses []models.Address) error {
	// Delete existing addresses
	query := `DELETE FROM addresses WHERE profile_id = $1`
	if _, err := tx.ExecContext(ctx, query, profileID); err != nil {
		return fmt.Errorf("error deleting addresses: %w", err)
	}

	// Create new addresses
	for i := range addresses {
		addr := &addresses[i]
		addr.ProfileID = profileID
		if err := r.createAddressTx(ctx, tx, addr); err != nil {
			return err
		}
	}

	return nil
}

func (r *ProfileRepository) updateContactsTx(ctx context.Context, tx *sqlx.Tx, profileID uuid.UUID, contacts []models.Contact) error {
	// Delete existing contacts
	query := `DELETE FROM contacts WHERE profile_id = $1`
	if _, err := tx.ExecContext(ctx, query, profileID); err != nil {
		return fmt.Errorf("error deleting contacts: %w", err)
	}

	// Create new contacts
	for i := range contacts {
		contact := &contacts[i]
		contact.ProfileID = profileID
		if err := r.createContactTx(ctx, tx, contact); err != nil {
			return err
		}
	}

	return nil
}

// List retrieves a list of profiles with pagination
func (r *ProfileRepository) List(ctx context.Context, page, pageSize int) ([]*models.Profile, error) {
	r.log.Debug("Listing profiles",
		logger.Int("page", page),
		logger.Int("page_size", pageSize),
	)

	offset := (page - 1) * pageSize

	query := `
		SELECT id, first_name, last_name, email, phone, created_at, updated_at
		FROM profiles
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	var profiles []*models.Profile
	err := r.db.SelectContext(ctx, &profiles, query, pageSize, offset)
	if err != nil {
		r.log.Error("Failed to list profiles",
			logger.ErrorField(err),
		)
		return nil, fmt.Errorf("error listing profiles: %w", err)
	}

	// Get addresses and contacts for each profile
	for _, profile := range profiles {
		addresses, err := r.getAddresses(ctx, profile.ID)
		if err != nil {
			return nil, err
		}
		profile.Addresses = addresses

		contacts, err := r.getContacts(ctx, profile.ID)
		if err != nil {
			return nil, err
		}
		profile.Contacts = contacts
	}

	r.log.Debug("Successfully listed profiles",
		logger.Int("count", len(profiles)),
	)
	return profiles, nil
}
