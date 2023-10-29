import json
import sqlite3
import tqdm
import zhconv
import re
with open('PoetryCategory.json', 'r', encoding='utf-8') as json_file:
    index = json.load(json_file)


def jsons_reader(path):
    r = []
    for json_file in path:
        with open(json_file, 'r', encoding='utf-8') as j:
            r += json.load(j)
    return r

def general_parser(path):
    r = jsons_reader(path)
    return [
        {
            "title": ii["title"],
            "author": ii["author"],
            "paragraphs": ii.get("paragraphs",""),
        }
        for ii in r
    ]

def sc(path):
    r = jsons_reader(path)
    return [
        {
            "title": ii["rhythmic"],
            "author": ii["author"],
            "paragraphs": ii["paragraphs"],
        }
        for ii in r
    ]

def ccsj(path):
    r = jsons_reader(path)
    return [
        {
            "title": ii["title"],
            "author": "曹操",
            "paragraphs": ii["paragraphs"],
        }
        for ii in r
    ]

def nlxd(path):
    r = jsons_reader(path)
    return [
        {
            "title": ii["title"],
            "author": ii["author"],
            "paragraphs": ii["para"],
        }
        for ii in r
    ]


def sj(path):
    r = jsons_reader(path)
    return [
        {
            "title": ii["title"],
            "author": "",
            "paragraphs": ii["content"],
        }
        for ii in r
    ]

def remove_symbols(text):
    # 使用正则表达式去除非中文字符
    cleaned_text = re.sub(r'[^\u4e00-\u9fff\s]', '', text)

    # 去除多余的空白字符
    cleaned_text = re.sub(r'\s+', ' ', cleaned_text)

    return cleaned_text


# handler
data = []
for category in index:
    match category:
        # 不喜欢哪个可以注释掉
        case '五代诗词':
            data += general_parser(index[category])
        case '全唐诗':
            data += general_parser(index[category])
        case "宋词":
          data += sc(index[category])
        case "御定全唐詩":
            data += general_parser(index[category])
        case "曹操诗集":
            data += ccsj(index[category])
        case "水墨唐诗":
            data += general_parser(index[category])
        case "纳兰性德":
            data += nlxd(index[category])
        case "诗经":
            data += sj(index[category])
    

    
# 创建数据库连接
conn = sqlite3.connect('poetry.db')

# 创建游标对象
cursor = conn.cursor()

# 创建表
cursor.execute('''
    CREATE TABLE IF NOT EXISTS poetry_info (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        author TEXT,
        paragraphs TEXT
    )
''')
# cursor.execute('''
#     CREATE TABLE IF NOT EXISTS hanzi_index (
#         hanzi TEXT,
#         hanzi_index TEXT
#     )
# ''')

# 字表
hanzi_index = {}
# 插入poetry_info
poetry_index = 0
for item in tqdm.tqdm(data):
    para_index = 0
    cursor.execute('''
        INSERT INTO poetry_info (title, author, paragraphs)
        VALUES (?, ?, ?)
    ''', (item['title'], item['author'], zhconv.convert(json.dumps(item['paragraphs'], ensure_ascii=False), 'zh-tw')))
    # for sentence in item['paragraphs']:
    #     for w in remove_symbols(sentence):
    #         hanzi_index[w] = hanzi_index.get(w, []) + [(poetry_index,para_index)]
    #     para_index += 1
    poetry_index += 1

# for item in tqdm.tqdm(hanzi_index):
#     cursor.execute('''
#         INSERT INTO hanzi_index (hanzi, hanzi_index)
#         VALUES (?, ?)
#     ''', (item, json.dumps(hanzi_index[item])))

# 提交事务
conn.commit()

# 关闭连接
conn.close()
