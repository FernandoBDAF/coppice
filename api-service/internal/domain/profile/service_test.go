package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

var errCacheMiss = errors.New("cache miss")

// mockRepository is an in-memory Repository used for unit tests.
type mockRepository struct {
	profiles  map[uuid.UUID]*Profile
	createErr error
}

func newMockRepository() *mockRepository {
	return &mockRepository{profiles: map[uuid.UUID]*Profile{}}
}

func (m *mockRepository) Create(_ context.Context, p *Profile) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.profiles[p.ID] = p
	return nil
}

func (m *mockRepository) GetByID(_ context.Context, id uuid.UUID) (*Profile, error) {
	if p, ok := m.profiles[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (m *mockRepository) GetByEmail(_ context.Context, email string) (*Profile, error) {
	for _, p := range m.profiles {
		if p.Email == email {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) Update(_ context.Context, p *Profile) error {
	m.profiles[p.ID] = p
	return nil
}

func (m *mockRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.profiles, id)
	return nil
}

func (m *mockRepository) List(_ context.Context, _, _ int) ([]*Profile, error) {
	out := make([]*Profile, 0, len(m.profiles))
	for _, p := range m.profiles {
		out = append(out, p)
	}
	return out, nil
}

// mockCache is an in-memory Cache used to assert invalidation behavior.
type mockCache struct {
	profiles             map[uuid.UUID]*Profile
	lists                map[int][]*Profile
	invalidateCalls      int
	invalidateListsCalls int
}

func newMockCache() *mockCache {
	return &mockCache{
		profiles: map[uuid.UUID]*Profile{},
		lists:    map[int][]*Profile{},
	}
}

func (m *mockCache) GetProfile(_ context.Context, id uuid.UUID) (*Profile, error) {
	if p, ok := m.profiles[id]; ok {
		return p, nil
	}
	return nil, errCacheMiss
}

func (m *mockCache) SetProfile(_ context.Context, p *Profile) error {
	m.profiles[p.ID] = p
	return nil
}

func (m *mockCache) InvalidateProfile(_ context.Context, id uuid.UUID) error {
	m.invalidateCalls++
	delete(m.profiles, id)
	return nil
}

func (m *mockCache) GetProfileList(_ context.Context, page int) ([]*Profile, error) {
	if list, ok := m.lists[page]; ok {
		return list, nil
	}
	return nil, errCacheMiss
}

func (m *mockCache) SetProfileList(_ context.Context, page int, profiles []*Profile) error {
	m.lists[page] = profiles
	return nil
}

func (m *mockCache) InvalidateProfileLists(_ context.Context) error {
	m.invalidateListsCalls++
	m.lists = map[int][]*Profile{}
	return nil
}

func validReq() *ProfileRequest {
	return &ProfileRequest{
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
	}
}

func TestService_Create_PopulatesCacheAndInvalidatesLists(t *testing.T) {
	repo := newMockRepository()
	cache := newMockCache()
	// Pre-seed a stale list page to prove it gets dropped on write.
	cache.lists[1] = []*Profile{{ID: uuid.New(), Email: "stale@example.com"}}
	svc := NewService(repo, cache, errCacheMiss)

	created, err := svc.Create(context.Background(), validReq())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatalf("expected generated profile ID")
	}
	if _, ok := cache.profiles[created.ID]; !ok {
		t.Errorf("expected profile to be cached after create")
	}
	if cache.invalidateListsCalls != 1 {
		t.Errorf("expected list cache invalidation on create, got %d calls", cache.invalidateListsCalls)
	}
	if _, stillCached := cache.lists[1]; stillCached {
		t.Errorf("expected stale list page to be evicted after create")
	}
}

func TestService_Update_InvalidatesListsAndRefreshesEntry(t *testing.T) {
	repo := newMockRepository()
	cache := newMockCache()
	svc := NewService(repo, cache, errCacheMiss)

	created, err := svc.Create(context.Background(), validReq())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	cache.lists[1] = []*Profile{created}
	cache.invalidateListsCalls = 0 // reset counter from Create

	req := validReq()
	req.FirstName = "Grace"
	if _, err := svc.Update(context.Background(), created.ID, req); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if cache.invalidateListsCalls != 1 {
		t.Errorf("expected list cache invalidation on update, got %d calls", cache.invalidateListsCalls)
	}
	if cache.profiles[created.ID].FirstName != "Grace" {
		t.Errorf("expected cached profile to reflect update, got %q", cache.profiles[created.ID].FirstName)
	}
}

func TestService_Delete_InvalidatesEntryAndLists(t *testing.T) {
	repo := newMockRepository()
	cache := newMockCache()
	svc := NewService(repo, cache, errCacheMiss)

	created, err := svc.Create(context.Background(), validReq())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	cache.lists[1] = []*Profile{created}

	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if cache.invalidateCalls != 1 {
		t.Errorf("expected single-entry cache invalidation on delete, got %d calls", cache.invalidateCalls)
	}
	if cache.invalidateListsCalls != 2 { // one from Create, one from Delete
		t.Errorf("expected list cache invalidation on delete, got %d calls", cache.invalidateListsCalls)
	}
	if _, ok := repo.profiles[created.ID]; ok {
		t.Errorf("expected profile removed from repository")
	}
}

func TestService_GetByID_CacheHitAvoidsRepository(t *testing.T) {
	repo := newMockRepository()
	cache := newMockCache()
	svc := NewService(repo, cache, errCacheMiss)

	id := uuid.New()
	cache.profiles[id] = &Profile{ID: id, Email: "cached@example.com"}

	got, err := svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if got.Email != "cached@example.com" {
		t.Errorf("expected cached profile to be returned, got %+v", got)
	}
	if _, inRepo := repo.profiles[id]; inRepo {
		t.Errorf("repository should not have been populated by the test")
	}
}

func TestService_Create_ValidationError(t *testing.T) {
	repo := newMockRepository()
	cache := newMockCache()
	svc := NewService(repo, cache, errCacheMiss)

	req := validReq()
	req.Email = "not-an-email"

	if _, err := svc.Create(context.Background(), req); !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}
