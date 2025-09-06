package domain

import "time"

type FakeHasher struct{}

func (FakeHasher) Hash(raw string) (string, error) { return "hash:" + raw, nil }
func (FakeHasher) Compare(hash, raw string) bool   { return hash == "hash:"+raw }

type FakeClock struct{ t time.Time }

func (c FakeClock) Now() time.Time {
	if !c.t.IsZero() {
		return c.t
	}
	return time.Now()
}

// Expose constructors used by integration tests
func NewMockPasswordHasher() FakeHasher { return FakeHasher{} }
func NewMockClock() FakeClock           { return FakeClock{} }
