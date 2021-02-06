import io
import os
import sys
import base64

from flask import Flask, request
from google.cloud import pubsub_v1
#import MeCab
# from PIL import Image
# from wordcloud import WordCloud


# def gen_wordcloud():
#     font_path = "/mnt/c/Windows/Fonts/YuGothR.ttc"

#     wordcloud = WordCloud(
#         font_path=font_path,
#         background_color="black",
#         width=900,
#         height=500,
#         contour_width=1,
#         contour_color="black",
#     )

#     # text = "こんにちは おはよう,ssああ、なんで,神奈川、おはよう、おはよう"
#     text = "なお、コンピュータの発達などで、これまで安全だった暗号が解読されるようになることを暗号の危殆化といいます。 NIST SP800-57 では、* 1024 ビットの RSAを新規用途には使わない* 2048 ビットの RSA を新規用途に使うのは2030年まで. 4096 ビットの RSAは 2031 年以降も新規用途に使えるき たいかレいう方針になっています。 第13章で紹介するGnuPG 2.1.4 では、 RSAの鍵長は2048 ビットがデフォルト値です。この章のまとめこの章では、公開鍵暗号と、その代表的実現方法である RSAについて学びました。公開鍵暗号を使えば、鍵配送問題が解決します。公開鍵暗号は、暗号における革命的な発明であり、現在のコンピュータやインターネットで使われている暗号技術は、公開建暗ませんてどット号から大きな恩恵を得ています。対称暗号は、平文を複雑な形に変換して機密性を保ちます。 一方、 公開鍵暗号は数学的に困難な問題を元にして機密性を保ちます。 たとえばRSAでは、大きな数の素因数分解を利用していました。 対称暗号と公開鍵暗号は根本的に異なる発想から生まれているのです。公開鍵暗号によって鍵配送問題は解決しましたが、公開鍵暗号に対してはman-in-the-middle攻撃が可能です。 これを防ぐためには「この公開鍵は正当な受信者の鍵なのか」 という問いに答えられなければなりません。これについては、第9章および第10章で解説します。公開鍵暗号が登場したからといって、対称暗号がなくなるわけではありません。公開暗号の実行スピードは対称暗号のそれよりもはるかに遅いからです。 ですから通常は、 対が暗号と公開鍵略号を組み合わせた通信が行われます。対称暗号によって処理スピードをクプし、公開鍵略号を使って鍵配送問題を解決するのです。 これがハイブリッド暗号システムです。 ハイプリッド暗号システムについては次の章で詳しく紹介します。ににビットRSH151"
#     wordcloud.generate(text)

#     # wordcloud.to_file("genereted.png")

#     # img = wordcloud.to_image()
#     # print(img)

#     return wordcloud.to_image()

# def image_to_byte_array(image:Image):
#     imgByteArr = io.BytesIO()
#     image.save(imgByteArr, format="PNG")
#     imgByteArr = imgByteArr.getvalue()
#     return imgByteArr

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

    name = "World"
    if isinstance(pubsub_message, dict) and "data" in pubsub_message:
        name = base64.b64decode(pubsub_message["data"]).decode("utf-8").strip()

    print(f"Hello {name}!")

    PROJECT_ID = "wordcloud-304009"
    TOPIC_WORD_CLOUD_NAME = "receive-word-cloud-topic"

    # word cloud topicにpush
    topic_path = publisher.topic_path(
        PROJECT_ID,
        TOPIC_WORD_CLOUD_NAME)

    data = b'to go!!'
    publisher.publish(topic_path, data=data)

    return ("", 204)

if __name__ == "__main__":
    PORT = int(os.getenv("PORT")) if os.getenv("PORT") else 8080

    # This is used when running locally. Gunicorn is used to run the
    # application on Cloud Run. See entrypoint in Dockerfile.
    app.run(host="127.0.0.1", port=PORT, debug=True)
