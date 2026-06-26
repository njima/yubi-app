package controller

import "testing"

func TestRequiredBodyRejectsNil(t *testing.T) {
	body, err := requiredBody((*struct{})(nil))

	if err == nil {
		t.Fatalf("expected error for nil body")
	}
	if body != nil {
		t.Fatalf("expected nil body, got %#v", body)
	}
}

func TestRequiredBodyReturnsBody(t *testing.T) {
	want := &struct{ Name string }{Name: "demo"}

	body, err := requiredBody(want)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if body != want {
		t.Fatalf("expected original body pointer")
	}
}
