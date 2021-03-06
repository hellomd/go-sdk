package newrelic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	newrelic "github.com/newrelic/go-agent"
	"github.com/urfave/negroni"
)

func TestBasicNewRelic(t *testing.T) {
	fakeNewRelic := newFakeNewRelicApp()
	nrMid := NewMiddleware(fakeNewRelic)

	a := negroni.New(nrMid)
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/md", nil)

	a.ServeHTTP(response, req)

	if tx, ok := fakeNewRelic.transactions["/md"]; !ok {
		t.Errorf("Expected transaction to /md path, got %v", fakeNewRelic.transactions)
	} else {
		if tx.running {
			t.Errorf("Transaction still running")
		} else {
			t.Log(*fakeNewRelic)
		}
	}

}

type fakeNewRelicApp struct {
	newrelic.Application
	transactions map[string]*fakeNewRelicTx
}

func newFakeNewRelicApp() *fakeNewRelicApp {
	return &fakeNewRelicApp{transactions: map[string]*fakeNewRelicTx{}}
}

func (fnr *fakeNewRelicApp) StartTransaction(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	fnr.transactions[name] = &fakeNewRelicTx{running: true}
	return fnr.transactions[name]
}

type fakeNewRelicTx struct {
	newrelic.Transaction
	running bool
}

func (ftx *fakeNewRelicTx) End() error {
	ftx.running = false
	return nil
}
