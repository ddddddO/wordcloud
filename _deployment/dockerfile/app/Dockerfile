FROM python:3.9-slim

# Allow statements and log messages to immediately appear in the Knative logs
ENV PYTHONUNBUFFERED True

# Copy local code to the container image.
ENV APP_HOME /app
WORKDIR $APP_HOME
COPY ./app ./

RUN apt-get update && \
    apt-get install fonts-takao-gothic && \
    pip install Flask gunicorn google-cloud-pubsub wordcloud

# mecabインストール
## https://qiita.com/Sak1361/items/47e9ec464ccc770cd65c
RUN apt-get install git make curl xz-utils file sudo mecab libmecab-dev mecab-ipadic mecab-ipadic-utf8 -y && \
    git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git && \
    cd mecab-ipadic-neologd && ./bin/install-mecab-ipadic-neologd -n -a -y && \
    pip install mecab-python3 unidic-lite

ENV FONT_PATH=/usr/share/fonts/truetype/fonts-japanese-gothic.ttf
#ENV DICTIONARY_PATH=/usr/lib/x86_64-linux-gnu/mecab/dic/mecab-ipadic-neologd

# Run the web service on container startup. Here we use the gunicorn
# webserver, with one worker process and 8 threads.
# For environments with multiple CPU cores, increase the number of workers
# to be equal to the cores available.
CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 --timeout 0 main:app