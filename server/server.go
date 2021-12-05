package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/observer/driver"
	"github.com/whoisnian/observer/serial"
)

func viewHandler(store *httpd.Store) {
	store.Respond200(generatedHtml)
}

func createEventHandler(port *serial.Port, encodeFunc driver.EncodeFunc) httpd.HandlerFunc {
	return func(store *httpd.Store) {
		var event driver.JSKeyEvent
		json.NewDecoder(store.R.Body).Decode(&event)
		store.R.Body.Close()
		isCombo := store.R.URL.Query().Get("combo") == "true"

		code, isCombo, isExit := driver.DecodeFromJS(event, isCombo)
		if !isCombo && !isExit {
			if res := encodeFunc(code); len(res) > 0 {
				port.Push(res)
				port.Push(encodeFunc(driver.EmptyKeycodes))
			}
		}

		store.RespondJson(struct {
			Code  string
			Combo bool
			Exit  bool
		}{code.String(), isCombo, isExit})
	}
}

func createStreamHandler(upStream string) httpd.HandlerFunc {
	return func(store *httpd.Store) {
		store.R.URL.Scheme = "http"
		store.R.URL.Host = upStream
		res, err := http.DefaultTransport.RoundTrip(store.R)
		if err != nil {
			logger.Error(err)
			return
		}
		defer res.Body.Close()

		for k, vv := range res.Header {
			for _, v := range vv {
				store.W.Header().Add(k, v)
			}
		}
		store.W.WriteHeader(res.StatusCode)
		io.Copy(store.W, res.Body)
	}
}

func Start(listenAt string, upStream string, port *serial.Port, encodeFunc driver.EncodeFunc) error {
	stop := port.GoWaitAndSend()
	defer stop()

	mux := httpd.NewMux()
	mux.Handle("/stream", "GET", createStreamHandler(upStream))
	mux.Handle("/api/event", "POST", createEventHandler(port, encodeFunc))
	mux.Handle("/*", "GET", viewHandler)

	return http.ListenAndServe(listenAt, logger.Req(mux))
}
