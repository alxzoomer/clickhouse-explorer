package logging

import (
	"github.com/rs/zerolog/log"
	"strings"
	"testing"
)

type fakeWriter struct {
	buff []byte
}

func (w *fakeWriter) Write(p []byte) (n int, err error) {
	w.buff = p
	return len(p), nil
}

func (w *fakeWriter) text() string {
	return string(w.buff)
}

func TestZeroLogErrorLogAdapter_Write(t *testing.T) {
	var fw = &fakeWriter{}
	log.Logger = log.Logger.Output(fw)
	l := NewErrorLog()
	msg := "test message"

	l.Printf("%s", msg)

	if !strings.Contains(fw.text(), msg) {
		t.Errorf("Write() gotText = %s, wantText %s", fw.text(), msg)
	}
}
