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
## https://pypi.org/project/mecab-python3/
## NOTE: 辞書はmecab-ipadic-neologdの方がよさそう。
RUN pip install mecab-python3 unidic-lite && \
    apt-get install mecab-ipadic-utf8 -y

ENV FONT_PATH=/usr/share/fonts/truetype/fonts-japanese-gothic.ttf

# Run the web service on container startup. Here we use the gunicorn
# webserver, with one worker process and 8 threads.
# For environments with multiple CPU cores, increase the number of workers
# to be equal to the cores available.
CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 --timeout 0 main:app