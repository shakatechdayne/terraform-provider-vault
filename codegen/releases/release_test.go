package releases

import (
	"testing"
)

var testTags = []string{
	"v1.0.0-beta1",
	"v0.11.4",
	"v0.11.3",
	"v0.11.2",
	"v0.11.1",
	"v0.11.0",
	"v0.11.0-beta1",
	"v0.10.4",
	"v0.10.3",
	"v0.10.2",
	"v0.10.1",
	"v0.10.0",
	"v0.10.0-rc1",
	"v0.9.6",
	"v0.9.5",
	"v0.9.4",
	"v0.9.3",
	"v0.9.2",
	"v0.9.1",
	"v0.9.0",
	"v0.8.3",
	"v0.8.2",
	"v0.8.1",
	"v0.8.0",
	"v0.8.0-rc1",
	"v0.8.0-beta1",
	"v0.7.3",
	"v0.7.2",
	"v0.7.1",
	"v0.7.0",
}

func TestParse(t *testing.T) {
	release, err := Parse("v0.8.2-beta1")
	if err != nil {
		t.Fatal(err)
	}
	if release.Major != 0 {
		t.Fatalf("expected 0 but received %d", release.Major)
	}
	if release.Minor != 8 {
		t.Fatalf("expected 8 but received %d", release.Minor)
	}
	if release.Patch != 2 {
		t.Fatalf("expected 2 but received %d", release.Patch)
	}
	if release.ProductName != "beta" {
		t.Fatalf("expected beta but received %s", release.ProductName)
	}
	if release.ProductNum != 1 {
		t.Fatalf("expected 1 but received %d", release.ProductNum)
	}
	for _, testTag := range testTags {
		release, err = Parse(testTag)
		if err != nil {
			t.Fatalf("failed on %s with %s", testTag, err)
		}
	}
}

func TestIsAfter(t *testing.T) {
	// Test major versions.
	this := &Release{
		Major:       1,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that := &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if !this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       1,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if this.IsAfter(that) {
		t.Fatalf("%+v is not after %+v", this, that)
	}

	// Test minor versions.
	this = &Release{
		Major:       0,
		Minor:       1,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if !this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       1,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if this.IsAfter(that) {
		t.Fatalf("%+v is not after %+v", this, that)
	}

	// Test patch versions.
	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       1,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if !this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       1,
		ProductName: "",
		ProductNum:  0,
	}
	if this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	// Test product numbers.
	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  1,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if !this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  1,
	}
	if this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	// Test same version.
	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if this.IsAfter(that) {
		t.Fatalf("%+v is not after %+v", this, that)
	}

	// Test "" is after rc.
	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "rc",
		ProductNum:  0,
	}
	if !this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "rc",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "",
		ProductNum:  0,
	}
	if this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	// Test rc is after beta.
	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "rc",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "beta",
		ProductNum:  0,
	}
	if !this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}

	this = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "beta",
		ProductNum:  0,
	}
	that = &Release{
		Major:       0,
		Minor:       0,
		Patch:       0,
		ProductName: "rc",
		ProductNum:  0,
	}
	if this.IsAfter(that) {
		t.Fatalf("%+v is after %+v", this, that)
	}
}
