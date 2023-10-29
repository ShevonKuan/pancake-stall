import re
import json
import sqlite3
import tqdm
output = []
pattern = r'[。？！]'
def remove_symbols(text):
    # 使用正则表达式去除非中文字符
    pattern = r"[^\u4e00-\u9fa5 ]"  # 匹配非中文汉字字符和非空格字符

    cleaned_text = re.sub(pattern, "", text)

    return cleaned_text



with open('origin.txt','r',encoding='utf-8') as f :
    count = -1
    for line in f.readlines():
        if "<poetry>" in line:
            item = {}
            output.append(item)
            count += 1
            output[count]["title"] = line.replace("<poetry>","").replace("：","").strip()
            output[count]["paragraphs"] = []
           
        if line.strip() == "" or "<poetry>" in line:
            continue
        else:
            output[count]["paragraphs"]+=(re.split(remove_symbols(pattern), line.strip()))

# 创建数据库连接
conn = sqlite3.connect('高中必背.db')

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

# 字表
hanzi_index = {}
# 插入poetry_info
poetry_index = 0



for item in tqdm.tqdm(output):
    para_index = 0
    cursor.execute('''
        INSERT INTO poetry_info (title, author, paragraphs)
        VALUES (?, ?, ?)
    ''', (item['title'], "", json.dumps(item["paragraphs"],ensure_ascii=False)))
    # for sentence in item['paragraphs']:
    #     for w in remove_symbols(sentence):
    #         hanzi_index[w] = hanzi_index.get(w, []) + [(poetry_index,para_index)]
    #     para_index += 1
    poetry_index += 1

# 提交事务
conn.commit()

# 关闭连接
conn.close()
