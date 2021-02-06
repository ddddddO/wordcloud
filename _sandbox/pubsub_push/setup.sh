#!/bin/bash

export CLOUDSDK_PYTHON=/usr/bin/python

# pubsubでメッセージを受信するCloudRunをデプロイする。
# https://cloud.google.com/run/docs/setup?hl=ja

#gcloud auth configure-docker
#gcloud components install docker-credential-gcr


# # # project切り替え
PROJECT_ID="wordcloud-304009"
# gcloud config set project ${PROJECT_ID}

# https://cloud.google.com/run/docs/quickstarts/build-and-deploy?hl=ja#python
SERVICE="gen-wordcloud"

## Cloud Build を使用してコンテナ イメージをビルドします。
gcloud builds submit --tag gcr.io/${PROJECT_ID}/${SERVICE}

## Cloud Runへデプロイ。
gcloud run deploy --image gcr.io/${PROJECT_ID}/${SERVICE} --platform managed

# https://cloud.google.com/run/docs/triggering/pubsub-push?hl=ja#command-line

## サブスクリプション用サービスアカウント作成。
SERVICE_ACCOUNT_NAME="go-to-py-cloudrun-invoker"
DISPLAYED_SERVICE_ACCOUNT_NAME="Go to Py Cloud Run Invoker"

gcloud iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
   --display-name "${DISPLAYED_SERVICE_ACCOUNT_NAME}"

## 上記で作成したサービスアカウントにCloudRun起動権限を付与する。
gcloud run services add-iam-policy-binding ${SERVICE} \
   --member=serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com \
   --role=roles/run.invoker

## Topicを作成
TOPIC_NAME="receive-text-topic"
gcloud pubsub topics create ${TOPIC_NAME}

## Pub/Sub がプロジェクトで認証トークンを作成できるようにします。?
PROJECT_NUMBER="65304505662"

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
     --member=serviceAccount:service-${PROJECT_NUMBER}@gcp-sa-pubsub.iam.gserviceaccount.com \
     --role=roles/iam.serviceAccountTokenCreator

## 必要な権限付きで(上記で)作成されたサービスアカウントで、Pub Sub サブスクリプションを作成します。
SUBSCRIPTION_ID="push-text-subscription"
### NOTE: Cloud Runデプロイ後でないとわからない
SERVICE_URL="https://gen-wordcloud-jixpif5lqq-uc.a.run.app"

gcloud beta pubsub subscriptions create ${SUBSCRIPTION_ID} --topic ${TOPIC_NAME} \
   --push-endpoint=${SERVICE_URL}/ \
   --push-auth-service-account=${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com

## Topicにメッセージを送信
#gcloud pubsub topics publish ${TOPIC_NAME} --message "hello"