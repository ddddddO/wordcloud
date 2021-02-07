from google.cloud import pubsub
from wordcloud import WordCloud
import MeCab

import time
from PIL import Image
import io

def gen_wordcloud():
    font_path = "/mnt/c/Windows/Fonts/YuGothR.ttc"

    wordcloud = WordCloud(
        font_path=font_path,
        background_color="black",
        width=900,
        height=500,
        contour_width=1,
        contour_color="black",
    )

    text = "こんにちは おはよう,ssああ、なんで,神奈川、おはよう、おはよう"
    wordcloud.generate(text)

    return wordcloud.to_image()

def image_to_byte_array(image:Image):
    imgByteArr = io.BytesIO()
    image.save(imgByteArr, format="PNG")
    imgByteArr = imgByteArr.getvalue()
    return imgByteArr

project_id = "test-emu-555"
topic_name = "my_topic"

topic_path = "projects/{project}/topics/{topic}".format(project=project_id, topic=topic_name)
publish_client = pubsub.PublisherClient()

publish_client.publish(topic_path, b'This is my message.', foo='bar', buzz='fizz')

# img = gen_wordcloud()
# print("send image...")
# publish_client.publish(topic_path, image_to_byte_array(img))
