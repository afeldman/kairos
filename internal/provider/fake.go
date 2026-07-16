package provider

import "context"

// FakeProvider implements Provider with canned responses for testing.
type FakeProvider struct {
	Response string
	Err      error
}

func (f *FakeProvider) Name() string { return "fake" }

func (f *FakeProvider) Generate(_ context.Context, _ Request) (string, error) {
	return f.Response, f.Err
}
