package services

import (
	"context"
	"fmt"
	"time"

	"github.com/fernandobarroso/profile-service/microservices/services/profile-api/internal/config"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-api/internal/logger"
	"github.com/fernandobarroso/profile-service/microservices/services/profile-api/internal/models"
	"go.uber.org/zap"
)

// ProfileError represents a profile service error
type ProfileError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *ProfileError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// ProfileServiceInterface defines the interface for profile-related operations
type ProfileServiceInterface interface {
	GetProfiles(ctx context.Context) ([]*models.Profile, error)
	GetProfile(ctx context.Context, id string) (*models.Profile, error)
	CreateProfile(ctx context.Context, req *models.ProfileRequest) (*models.Profile, error)
	UpdateProfile(ctx context.Context, id string, req *models.ProfileRequest) (*models.Profile, error)
	DeleteProfile(ctx context.Context, id string) error
}

// ProfileService handles profile-related business logic
type ProfileService struct {
	storageClient *StorageClient
}

// NewProfileService creates a new profile service
func NewProfileService(cfg *config.Config, storageClient *StorageClient) *ProfileService {
	return &ProfileService{
		storageClient: storageClient,
	}
}

// GetProfiles retrieves all profiles
func (s *ProfileService) GetProfiles(ctx context.Context) ([]*models.Profile, error) {
	logger.LogInfo(ctx, "Getting all profiles")
	profiles, err := s.storageClient.GetProfiles(ctx)
	if err != nil {
		logger.LogError(ctx, "Error getting profiles", err)
		return nil, &ProfileError{
			Code:    500,
			Message: "Failed to get profiles",
			Err:     err,
		}
	}
	logger.LogInfo(ctx, "Successfully retrieved profiles",
		zap.Int("count", len(profiles)))
	return profiles, nil
}

// GetProfile retrieves a profile by ID
func (s *ProfileService) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	if id == "" {
		return nil, &ProfileError{
			Code:    400,
			Message: "Profile ID is required",
		}
	}

	logger.LogInfo(ctx, "Getting profile",
		zap.String("id", id))
	profile, err := s.storageClient.GetProfile(ctx, id)
	if err != nil {
		logger.LogError(ctx, "Error getting profile", err,
			zap.String("id", id))
		return nil, &ProfileError{
			Code:    500,
			Message: fmt.Sprintf("Failed to get profile %s", id),
			Err:     err,
		}
	}
	logger.LogInfo(ctx, "Successfully retrieved profile",
		zap.String("id", id))
	return profile, nil
}

// CreateProfile creates a new profile
func (s *ProfileService) CreateProfile(ctx context.Context, req *models.ProfileRequest) (*models.Profile, error) {
	if req == nil {
		return nil, &ProfileError{
			Code:    400,
			Message: "Profile request is required",
		}
	}

	if err := req.Validate(); err != nil {
		logger.LogError(ctx, "Invalid profile request", err)
		return nil, &ProfileError{
			Code:    400,
			Message: "Invalid profile request",
			Err:     err,
		}
	}

	logger.LogInfo(ctx, "Creating new profile",
		zap.String("email", req.Email))
	profile := &models.Profile{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Bio:       req.Bio,
		ImageURLs: req.ImageURLs,
		Address:   req.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdProfile, err := s.storageClient.CreateProfile(ctx, profile)
	if err != nil {
		logger.LogError(ctx, "Error creating profile", err,
			zap.String("email", req.Email))
		return nil, &ProfileError{
			Code:    500,
			Message: "Failed to create profile",
			Err:     err,
		}
	}
	logger.LogInfo(ctx, "Successfully created profile",
		zap.String("id", createdProfile.ID.String()),
		zap.String("email", req.Email))
	return createdProfile, nil
}

// UpdateProfile updates an existing profile
func (s *ProfileService) UpdateProfile(ctx context.Context, id string, req *models.ProfileRequest) (*models.Profile, error) {
	if id == "" {
		return nil, &ProfileError{
			Code:    400,
			Message: "Profile ID is required",
		}
	}

	if req == nil {
		return nil, &ProfileError{
			Code:    400,
			Message: "Profile request is required",
		}
	}

	if err := req.Validate(); err != nil {
		logger.LogError(ctx, "Invalid profile request", err)
		return nil, &ProfileError{
			Code:    400,
			Message: "Invalid profile request",
			Err:     err,
		}
	}

	logger.LogInfo(ctx, "Updating profile",
		zap.String("id", id))
	// First get the existing profile
	existingProfile, err := s.storageClient.GetProfile(ctx, id)
	if err != nil {
		logger.LogError(ctx, "Error getting existing profile", err,
			zap.String("id", id))
		return nil, &ProfileError{
			Code:    500,
			Message: fmt.Sprintf("Failed to get existing profile %s", id),
			Err:     err,
		}
	}

	// Update the fields
	existingProfile.FirstName = req.FirstName
	existingProfile.LastName = req.LastName
	existingProfile.Email = req.Email
	existingProfile.Phone = req.Phone
	existingProfile.Bio = req.Bio
	existingProfile.ImageURLs = req.ImageURLs
	existingProfile.Address = req.Address
	existingProfile.UpdatedAt = time.Now()

	updatedProfile, err := s.storageClient.UpdateProfile(ctx, id, existingProfile)
	if err != nil {
		logger.LogError(ctx, "Error updating profile", err,
			zap.String("id", id))
		return nil, &ProfileError{
			Code:    500,
			Message: fmt.Sprintf("Failed to update profile %s", id),
			Err:     err,
		}
	}
	logger.LogInfo(ctx, "Successfully updated profile",
		zap.String("id", id))
	return updatedProfile, nil
}

// DeleteProfile deletes a profile
func (s *ProfileService) DeleteProfile(ctx context.Context, id string) error {
	if id == "" {
		return &ProfileError{
			Code:    400,
			Message: "Profile ID is required",
		}
	}

	logger.LogInfo(ctx, "Deleting profile",
		zap.String("id", id))
	err := s.storageClient.DeleteProfile(ctx, id)
	if err != nil {
		logger.LogError(ctx, "Error deleting profile", err,
			zap.String("id", id))
		return &ProfileError{
			Code:    500,
			Message: fmt.Sprintf("Failed to delete profile %s", id),
			Err:     err,
		}
	}
	logger.LogInfo(ctx, "Successfully deleted profile",
		zap.String("id", id))
	return nil
}
