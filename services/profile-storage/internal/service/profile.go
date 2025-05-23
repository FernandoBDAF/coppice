package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"microservices/services/profile-storage/internal/logger"
	"microservices/services/profile-storage/internal/models"
	"microservices/services/profile-storage/internal/repository"
)

// Service errors
var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrProfileNotFound   = errors.New("profile not found")
	ErrDuplicateEmail    = errors.New("email already in use")
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrTimeout           = errors.New("operation timed out")
)

// ProfileService handles business logic for profile operations
type ProfileService struct {
	repo               *repository.ProfileRepository
	transactionTimeout time.Duration
	log                *zap.Logger
}

// NewProfileService creates a new profile service
func NewProfileService(repo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		repo:               repo,
		transactionTimeout: 30 * time.Second,
		log:                logger.Get(),
	}
}

// CreateProfile creates a new profile with validation and business rules
func (s *ProfileService) CreateProfile(ctx context.Context, req *models.ProfileRequest) (*models.Profile, error) {
	startTime := time.Now()
	s.log.Info("Creating new profile",
		logger.String("email", req.Email),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, s.transactionTimeout)
	defer cancel()

	// Validate request
	if err := req.Validate(); err != nil {
		s.log.Error("Invalid profile request",
			logger.ErrorField(err),
			logger.String("email", req.Email),
		)
		return nil, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	// Check if email is already in use
	// TODO: Implement email uniqueness check

	// Create profile model
	profile := &models.Profile{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Addresses: req.Addresses,
		Contacts:  req.Contacts,
	}

	// Set timestamps
	now := time.Now()
	profile.CreatedAt = now
	profile.UpdatedAt = now

	// Generate UUIDs for addresses and contacts
	for i := range profile.Addresses {
		profile.Addresses[i].ID = uuid.New()
	}
	for i := range profile.Contacts {
		profile.Contacts[i].ID = uuid.New()
	}

	// Create profile in database
	if err := s.repo.Create(ctx, profile); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Error("Profile creation timed out",
				logger.ErrorField(err),
				logger.String("profile_id", profile.ID.String()),
				logger.String("email", profile.Email),
				logger.Duration("timeout", s.transactionTimeout),
			)
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		s.log.Error("Failed to create profile",
			logger.ErrorField(err),
			logger.String("profile_id", profile.ID.String()),
			logger.String("email", profile.Email),
		)
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	s.log.Info("Successfully created profile",
		logger.String("profile_id", profile.ID.String()),
		logger.String("email", profile.Email),
		logger.Duration("duration", time.Since(startTime)),
	)
	return profile, nil
}

// GetProfile retrieves a profile by ID with business rules
func (s *ProfileService) GetProfile(ctx context.Context, id uuid.UUID) (*models.Profile, error) {
	startTime := time.Now()
	s.log.Debug("Getting profile",
		logger.String("profile_id", id.String()),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, s.transactionTimeout)
	defer cancel()

	profile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Error("Profile retrieval timed out",
				logger.ErrorField(err),
				logger.String("profile_id", id.String()),
				logger.Duration("timeout", s.transactionTimeout),
			)
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		if errors.Is(err, repository.ErrNotFound) {
			s.log.Debug("Profile not found",
				logger.String("profile_id", id.String()),
			)
			return nil, fmt.Errorf("%w: %v", ErrProfileNotFound, err)
		}
		s.log.Error("Failed to get profile",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if profile == nil {
		s.log.Debug("Profile not found",
			logger.String("profile_id", id.String()),
		)
		return nil, ErrProfileNotFound
	}

	s.log.Debug("Successfully retrieved profile",
		logger.String("profile_id", id.String()),
		logger.String("email", profile.Email),
		logger.Duration("duration", time.Since(startTime)),
	)
	return profile, nil
}

// UpdateProfile updates an existing profile with validation and business rules
func (s *ProfileService) UpdateProfile(ctx context.Context, id uuid.UUID, req *models.ProfileRequest) (*models.Profile, error) {
	startTime := time.Now()
	s.log.Info("Updating profile",
		logger.String("profile_id", id.String()),
		logger.String("email", req.Email),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, s.transactionTimeout)
	defer cancel()

	// Validate request
	if err := req.Validate(); err != nil {
		s.log.Error("Invalid profile update request",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
			logger.String("email", req.Email),
		)
		return nil, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	// Get existing profile
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Error("Profile retrieval timed out during update",
				logger.ErrorField(err),
				logger.String("profile_id", id.String()),
				logger.Duration("timeout", s.transactionTimeout),
			)
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		if errors.Is(err, repository.ErrNotFound) {
			s.log.Debug("Profile not found for update",
				logger.String("profile_id", id.String()),
			)
			return nil, fmt.Errorf("%w: %v", ErrProfileNotFound, err)
		}
		s.log.Error("Failed to get profile for update",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	// Update profile fields
	existing.FirstName = req.FirstName
	existing.LastName = req.LastName
	existing.Email = req.Email
	existing.Phone = req.Phone
	existing.Addresses = req.Addresses
	existing.Contacts = req.Contacts

	// Generate UUIDs for new addresses and contacts
	for i := range existing.Addresses {
		if existing.Addresses[i].ID == uuid.Nil {
			existing.Addresses[i].ID = uuid.New()
		}
	}
	for i := range existing.Contacts {
		if existing.Contacts[i].ID == uuid.Nil {
			existing.Contacts[i].ID = uuid.New()
		}
	}

	// Update profile in database
	if err := s.repo.Update(ctx, existing); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Error("Profile update timed out",
				logger.ErrorField(err),
				logger.String("profile_id", id.String()),
				logger.String("email", existing.Email),
				logger.Duration("timeout", s.transactionTimeout),
			)
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		s.log.Error("Failed to update profile",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
			logger.String("email", existing.Email),
		)
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	s.log.Info("Successfully updated profile",
		logger.String("profile_id", id.String()),
		logger.String("email", existing.Email),
		logger.Duration("duration", time.Since(startTime)),
	)
	return existing, nil
}

// DeleteProfile deletes a profile with business rules
func (s *ProfileService) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	startTime := time.Now()
	s.log.Info("Deleting profile",
		logger.String("profile_id", id.String()),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, s.transactionTimeout)
	defer cancel()

	// Check if profile exists
	exists, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Error("Profile retrieval timed out during delete",
				logger.ErrorField(err),
				logger.String("profile_id", id.String()),
				logger.Duration("timeout", s.transactionTimeout),
			)
			return fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		if !errors.Is(err, repository.ErrNotFound) {
			s.log.Error("Failed to check profile existence",
				logger.ErrorField(err),
				logger.String("profile_id", id.String()),
			)
			return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
		}
	}

	if exists == nil {
		s.log.Debug("Profile not found for deletion",
			logger.String("profile_id", id.String()),
		)
		return ErrProfileNotFound
	}

	// Delete profile
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Error("Profile deletion timed out",
				logger.ErrorField(err),
				logger.String("profile_id", id.String()),
				logger.Duration("timeout", s.transactionTimeout),
			)
			return fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		s.log.Error("Failed to delete profile",
			logger.ErrorField(err),
			logger.String("profile_id", id.String()),
		)
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	s.log.Info("Successfully deleted profile",
		logger.String("profile_id", id.String()),
		logger.Duration("duration", time.Since(startTime)),
	)
	return nil
}

// ListProfiles retrieves a list of profiles with pagination
func (s *ProfileService) ListProfiles(ctx context.Context, page, pageSize int) ([]*models.Profile, error) {
	startTime := time.Now()
	s.log.Info("Listing profiles",
		logger.Int("page", page),
		logger.Int("page_size", pageSize),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, s.transactionTimeout)
	defer cancel()

	profiles, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Error("Profile listing timed out",
				logger.ErrorField(err),
				logger.Duration("timeout", s.transactionTimeout),
			)
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		s.log.Error("Failed to list profiles",
			logger.ErrorField(err),
		)
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	s.log.Info("Successfully listed profiles",
		logger.Int("count", len(profiles)),
		logger.Duration("duration", time.Since(startTime)),
	)
	return profiles, nil
}

// SearchProfiles searches profiles by various criteria
func (s *ProfileService) SearchProfiles(ctx context.Context, query string) ([]*models.Profile, error) {
	// TODO: Implement search functionality
	return nil, fmt.Errorf("not implemented")
}
