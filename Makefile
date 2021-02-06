# ビルド構成ファイルによるビルド
# https://cloud.google.com/cloud-build/docs/building/build-containers?hl=ja#use-buildconfig
build_app:
	gcloud builds submit --config app/cloudbuild.yaml

build_web:
	gcloud builds submit --config web/cloudbuild.yaml

# NOTE: 18 選択
deploy_app: build_app
	gcloud run deploy --image gcr.io/wordcloud-304009/gen-wordcloud --platform managed

# NOTE: 18 選択
deploy_web: build_web
	gcloud run deploy --image gcr.io/wordcloud-304009/web --platform managed

localpubsub:
	export CLOUDSDK_PYTHON=/usr/bin/python2 && \
	gcloud beta emulators pubsub start --project=test-emu-555

	# エミュレータ起動後、以下を実行してアプリケーションが Pub/Sub エミュレータに接続できるようにする。
	# export CLOUDSDK_PYTHON=/usr/bin/python2 && gcloud beta emulators pubsub env-init
	# -> export PUBSUB_EMULATOR_HOST=::1:<利用ポート> が実行されるが、実行環境が悪いからか、このままだとtopicが生成できない。
	#    なので、各ターミナルで以下を実行する。
	# export PUBSUB_EMULATOR_HOST=localhost:<利用ポート>
