package logging

import (
	"encoding/json"
	"github.com/alxzoomer/clickhouse-explorer/pkg/appinfo"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestInit(t *testing.T) {
	Init(true)
	fw := &fakeWriter{}
	l := log.Output(fw)
	msg := "test message"

	l.Info().Msg(msg)

	r := struct {
		Message   string `json:"message"`
		Service   string `json:"service"`
		Version   string `json:"version"`
		Branch    string `json:"branch"`
		BuildTime string `json:"buildtime"`
		Host      string `json:"host"`
	}{}
	json.Unmarshal(fw.buff, &r)

	if r.Message != msg {
		t.Errorf("TestInit want = %s, got = %s", msg, r.Message)
	}
	if r.Service != appinfo.Service {
		t.Errorf("TestInit want = %s, got = %s", appinfo.Service, r.Service)
	}
	if r.Version != appinfo.Version {
		t.Errorf("TestInit want = %s, got = %s", appinfo.Version, r.Version)
	}
	if r.Branch != appinfo.Branch {
		t.Errorf("TestInit want = %s, got = %s", appinfo.Branch, r.Branch)
	}
	if r.BuildTime != appinfo.BuildTime {
		t.Errorf("TestInit want = %s, got = %s", appinfo.BuildTime, r.BuildTime)
	}
	if r.Host != appinfo.Hostname {
		t.Errorf("TestInit want = %s, got = %s", appinfo.Hostname, r.Host)
	}
}
