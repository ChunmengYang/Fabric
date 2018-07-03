package http

import (
	"net/http"
	"encoding/json"
	"time"
	//_ "net/http/pprof" //ip:port//debug/pprof/
	//_ "github.com/mkevac/debugcharts"  //ip:port/debug/charts/

	"github.com/Sirupsen/logrus"
)

func init() {
	configRoutes()
}

func configRoutes() {
	http.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var body string
		err := decoder.Decode(&body)
		if err != nil {
			http.Error(w, "connot decode body", http.StatusBadRequest)
			return
		}

		// TODO

		w.Write([]byte("success"))
	})
}

func Start() {
	addr := "1997"

	server := &http.Server{
		Addr:           addr,
		ReadTimeout:	time.Second * 10,
		WriteTimeout: 	time.Second * 30,
	}

	logrus.Infoln("listening", addr)
	logrus.Infoln(server.ListenAndServe())
}