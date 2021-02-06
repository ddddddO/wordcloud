package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// "github.com/ddddddO/extxt"
	tmpl "github.com/ddddddO/wordcloud/web/templates"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

const (
	src                       = "src_file"
	projectID                 = "wordcloud-304009"
	topicTxtName              = "receive-text-topic"
	topicWordCloudName        = "receive-word-cloud-topic"
	subscriptionWordCloudName = "pull-word-cloud-subscription"
)

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
		_ = f

		ctx := context.Background()
		pubsubClient, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("pubsub.NewClient: %v", err)
			return
		}

		wordCloudTopic, err := pubsubClient.CreateTopic(ctx, topicWordCloudName)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("pubsub.CreateTopic: %v", err)
			return
		}
		defer wordCloudTopic.Delete(ctx)

		wordCloudSub, err := pubsubClient.CreateSubscription(ctx, subscriptionWordCloudName, pubsub.SubscriptionConfig{Topic: wordCloudTopic})
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("pubsub.CreateSubscription: %v", err)
			return
		}
		defer wordCloudSub.Delete(ctx)

		// Turn on synchronous mode. This makes the subscriber use the Pull RPC rather
		// than the StreamingPull RPC, which is useful for guaranteeing MaxOutstandingMessages,
		// the max number of messages the client will hold in memory at a time.
		wordCloudSub.ReceiveSettings.Synchronous = true
		wordCloudSub.ReceiveSettings.MaxOutstandingMessages = 10

		// Receive messages for 20 seconds.
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		// Create a channel to handle messages to as they come in.
		cm := make(chan *pubsub.Message)
		defer close(cm)
		// Handle individual messages in a goroutine.
		go func(w http.ResponseWriter) {
			for msg := range cm {
				fmt.Printf("Got message\n")

				reader := bytes.NewReader(msg.Data)
				img, _, err := image.Decode(reader)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-type", "image/png")
				err = png.Encode(w, img)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				msg.Ack()
			}
		}(w)

		txtTopic := pubsubClient.Topic(topicTxtName)
		res := txtTopic.Publish(r.Context(), &pubsub.Message{Data: []byte("to py!!!")})
		if _, err := res.Get(r.Context()); err != nil {
			log.Printf("Publish.Get: %v", err)
			http.Error(w, "Error requesting translation", http.StatusInternalServerError)
			return
		}

		// Receive blocks until the passed in context is done.
		err = wordCloudSub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			cm <- msg
		})
		if err != nil && status.Code(err) != codes.Canceled {
			fmt.Errorf("Receive: %v", err)
			http.Error(w, "Error requesting translation", http.StatusInternalServerError)
			return
		}

		// buf := &bytes.Buffer{}
		// if err := extxt.RunByServer(buf, f); err != nil {
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }

		// TODO:
		// 1. ブラウザから入力したテキスト・画像から抽出したテキストをパブリッシュする先のTopicと受信しpushするSubscriptionをterraformから作成する。こちらまず読んでから実装-> https://cloud.google.com/run/docs/triggering/pubsub-push?hl=ja#run_pubsub_handler-python
		// 2. サブスクライブしたテキストデータから生成されるWordCloudイメージデータ(CloudRun(Python)の処理)、のパブリッシュ先のTopicを作成する。
		// 3. WordCloudイメージデータをPullする先のSubscriptionを作成する。
		// 4. 1で作成したTopicに対してテキストデータをパブリッシュ
		// 5. 3で作成したSubscriptionからWordCloudイメージデータを同期Pullする処理
		// 6. 5の終了後、2で作成したTopic/3で作成したSubscriptionを削除する
		// 7. レスポンスにWordCloudイメージを返却

		// tmp := &struct {
		// 	Text  string   `json:"Text"`
		// 	Words []string `json:"Words"`
		// }{}

		// if err := json.Unmarshal(buf.Bytes(), tmp); err != nil {
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }

		// t, err := template.New("extxt").Parse(tmpl.ExtxtHTML)
		// if err != nil {
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }

		// debug
		// tmp.Text = "succeeded!"
		// tmp.Words = []string{"aa", "ssd"}

		// if err := t.Execute(w, tmp); err != nil {
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }

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
