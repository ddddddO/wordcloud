import io
import os
import sys
import base64

from flask import Flask, request
from google.cloud import pubsub_v1
#import MeCab
from PIL import Image
from wordcloud import WordCloud


def gen_wordcloud(text:str):
    # TODO: エラーハンドリング
    font_path = os.getenv("FONT_PATH")
    print(font_path)

    wordcloud = WordCloud(
        font_path=font_path,
        background_color="black",
        width=900,
        height=500,
        contour_width=1,
        contour_color="black",
    )

    wordcloud.generate(text)

    return wordcloud.to_image()

def image_to_byte_array(image:Image):
    imgByteArr = io.BytesIO()
    image.save(imgByteArr, format="PNG")
    imgByteArr = imgByteArr.getvalue()
    return imgByteArr

# data = image_to_byte_array(gen_wordcloud())
# sys.stdout.buffer.write(data)

# f = open("tmp.png", "wb")
# f.write(image_to_byte_array(gen_wordcloud()))
# f.close()

app = Flask(__name__)
publisher = pubsub_v1.PublisherClient()

@app.route("/", methods=["POST"])
def index():
    envelope = request.get_json()
    if not envelope:
        msg = "no Pub/Sub message received"
        print(f"error: {msg}")
        return f"Bad Request: {msg}", 400

    if not isinstance(envelope, dict) or "message" not in envelope:
        msg = "invalid Pub/Sub message format"
        print(f"error: {msg}")
        return f"Bad Request: {msg}", 400

    pubsub_message = envelope["message"]

    text = "World"
    if isinstance(pubsub_message, dict) and "data" in pubsub_message:
        text = base64.b64decode(pubsub_message["data"]).decode("utf-8").strip()

    print(f"Hello {text}!")

    PROJECT_ID = "wordcloud-304009"
    TOPIC_WORD_CLOUD_NAME = "receive-word-cloud-topic"

    # word cloud topicにpush
    topic_path = publisher.topic_path(
        PROJECT_ID,
        TOPIC_WORD_CLOUD_NAME)

    data = image_to_byte_array(gen_wordcloud(text))
    publisher.publish(topic_path, data=data)

    return ("", 204)

if __name__ == "__main__":
    PORT = int(os.getenv("PORT")) if os.getenv("PORT") else 8080

    # This is used when running locally. Gunicorn is used to run the
    # application on Cloud Run. See entrypoint in Dockerfile.
    app.run(host="127.0.0.1", port=PORT, debug=True)
