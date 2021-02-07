import MeCab

from wordcloud import WordCloud

def parse_noun(text:str) -> str:
    tagger = MeCab.Tagger('-d /usr/lib/x86_64-linux-gnu/mecab/dic/mecab-ipadic-neologd')

    parsed_nouns=""
    node = tagger.parseToNode(text)
    while node:
        word = node.surface
        pos = node.feature.split(",")[0]
        print('{0} , {1}'.format(word, pos))
        if pos == '名詞':
            parsed_nouns+=word+","

        node = node.next
    return parsed_nouns

def gen_wordcloud_png(text:str):
    font_path = "/mnt/c/Windows/Fonts/YuGothR.ttc"

    wordcloud = WordCloud(
        font_path=font_path,
        background_color="black",
        width=900,
        height=500,
        contour_width=1,
        contour_color="black",
    )

    wordcloud.generate(text)
    wordcloud.to_file("genereted_by_mecab.png")

text="ある日の暮方の事である。一人の下人げにんが、羅生門らしょうもんの下で雨やみを待>っていた。"
gen_wordcloud_png(parse_noun(text))
