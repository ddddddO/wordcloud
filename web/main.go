package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ddddddO/extxt"
	tmpl "github.com/ddddddO/wordcloud/web/templates"
)

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}

func runServer() error {
	log.Println("start")

	http.HandleFunc("/", indexHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}

const src = "src_file"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if !basicAuthenticated(r) {
		w.Header().Add("WWW-Authenticate", `Basic realm="secret xxxx"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		t, err := template.New("index").Parse(tmpl.IndexHTML)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, src); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		f, header, err := r.FormFile(src)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_ = header

		buf := &bytes.Buffer{}
		if err := extxt.RunByServer(buf, f); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// TODO:
		// 1. ブラウザから入力したテキスト・画像から抽出したテキストをパブリッシュする先のTopicと受信しpushするSubscriptionをterraformから作成する。こちらまず読んでから実装-> https://cloud.google.com/run/docs/triggering/pubsub-push?hl=ja#run_pubsub_handler-python
		// 2. サブスクライブしたテキストデータから生成されるWordCloudイメージデータ(CloudRun(Python)の処理)、のパブリッシュ先のTopicを作成する。
		// 3. WordCloudイメージデータをPullする先のSubscriptionを作成する。
		// 4. 1で作成したTopicに対してテキストデータをパブリッシュ
		// 5. 3で作成したSubscriptionからWordCloudイメージデータを同期Pullする処理
		// 6. 5の終了後、2で作成したTopic/3で作成したSubscriptionを削除する
		// 7. レスポンスにWordCloudイメージを返却

		tmp := &struct {
			Text  string   `json:"Text"`
			Words []string `json:"Words"`
		}{}

		if err := json.Unmarshal(buf.Bytes(), tmp); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		t, err := template.New("extxt").Parse(tmpl.ExtxtHTML)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, tmp); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

var (
	validNames     = strings.Split(os.Getenv("BASIC_AUTH_NAMES"), ",")
	validPasswords = strings.Split(os.Getenv("BASIC_AUTH_PASSWORDS"), ",")
)

func basicAuthenticated(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}

	for i := range validNames {
		if username != validNames[i] {
			continue
		}
		if password == validPasswords[i] {
			return true
		}
	}
	return false
}
