package app

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("./testdata")
	if err != nil {
		t.Error(err)
	}
}

func TestEnvOverridesCorrectly(t *testing.T) {
	want := "345"
	t.Setenv("TELEGRAM.TOKEN", want)
	conf, _ := LoadConfig("./testdata")
	if conf.Telegram.Token != want {
		t.Errorf("got %v want %v", conf.Telegram.Token, want)
	}
}
