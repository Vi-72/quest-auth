package domain

import "time"

type fakeHasher struct{}

func (fakeHasher) Hash(raw string) (string, error) { return "hash:" + raw, nil }
func (fakeHasher) Compare(hash, raw string) bool   { return hash == "hash:"+raw }

type fakeClock struct{ t time.Time }

func (c fakeClock) Now() time.Time {
	if !c.t.IsZero() {
		return c.t
	}
	return time.Now()
}
