package git

import (
	"reflect"
	"testing"
)

func TestFakeRepository_ImplementsRepository(t *testing.T) {
	var _ Repository = &FakeRepository{}
}

func TestFakeRepository_LogTruncatesToN(t *testing.T) {
	f := &FakeRepository{Commits: []string{"c1", "c2", "c3"}}

	got, err := f.Log(2)
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}
	want := []string{"c1", "c2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Log(2) = %v, want %v", got, want)
	}
}

func TestFakeRepository_LogReturnsAllWhenNExceedsLen(t *testing.T) {
	f := &FakeRepository{Commits: []string{"c1"}}

	got, err := f.Log(5)
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}
	if len(got) != 1 || got[0] != "c1" {
		t.Fatalf("Log(5) = %v, want [\"c1\"]", got)
	}
}
