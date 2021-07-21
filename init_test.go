package gokisscore

import (
	"testing"
)

func TestParseScript(t *testing.T) {

	filler := NewFiller()
	filler.ParseScript()

	if filler.Error != nil {
		t.Fatal(filler.Error.Error())
	}
}

func TestGetBody(t *testing.T) {

	filler := NewFiller()
	filler.GetVersion()
	filler.GetBody()

	if filler.Error != nil {
		t.Fatal()
	}

}

func TestGetVerstion(t *testing.T) {

	data := []struct {
		url              string
		is_empty_version bool
	}{
		{
			"https://inspin.me/version.json",
			false,
		},
		{
			"https://inspin.me/version",
			true,
		},
		{
			"",
			true,
		},
	}

	for _, v := range data {

		filler := NewFiller()
		filler.URL_VERTION = v.url
		filler.GetVersion()

		is_empty_version := len(filler.game_version) == 0

		if is_empty_version != v.is_empty_version {
			t.Logf("error TODO: %s\n", filler.Error.Error())
			t.Fatalf("error get game version. url: %s,want is_empty_version: %v\n", v.url, v.is_empty_version)
		}

	}
}
