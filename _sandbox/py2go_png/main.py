from wordcloud import WordCloud
import MeCab
from PIL import Image
import io
import sys

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
    img_byte_arr = io.BytesIO()
    image.save(img_byte_arr, format="PNG")
    return img_byte_arr.getvalue()

data = image_to_byte_array(gen_wordcloud())
sys.stdout.buffer.write(data)
