localpubsub:
	export CLOUDSDK_PYTHON=/usr/bin/python2 && \
	gcloud beta emulators pubsub start --project=test-emu-555

	# エミュレータ起動後、以下を実行してアプリケーションが Pub/Sub エミュレータに接続できるようにする。
	# gcloud beta emulators pubsub env-init
	# -> export PUBSUB_EMULATOR_HOST=::1:<利用ポート> が実行されるが、実行環境が悪いからか、このままだとtopicが生成できない。
	#    なので、各ターミナルで以下を実行する。
	# export PUBSUB_EMULATOR_HOST=localhost:<利用ポート>