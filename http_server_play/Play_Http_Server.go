package httpserverplay

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

func PlayHttpServer() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]string{"pong": "i am pong", "_t": time.Now().String()}

		var bytes, err = json.Marshal(data)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(200)
		w.Write(bytes)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("http_server_play/template/index.html")
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}

		reqHeaderJson, _ := json.MarshalIndent(r.Header, "", "\t")
		data := map[string]string{
			"title":   "This is title",
			"content": fmt.Sprintf("Ip: %s \nHeaders: %s", r.RemoteAddr, reqHeaderJson),
			"_t":      time.Now().String(),
		}
		t.Execute(w, data)
	})

	// https://gist.github.com/paulmach/7271283?permalink_comment_id=2245985#gistcomment-2245985
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("http_server_play/public"))))

	var addr = "0.0.0.0:3000"
	log.Println("http listern at ", addr)

	httpServer := http.Server{
		Addr: addr,
	}
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
