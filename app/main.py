from wordcloud import WordCloud
import MeCab

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

    wordcloud.to_file("genereted.png")

    # img = wordcloud.to_image()
    # print(img)

gen_wordcloud()