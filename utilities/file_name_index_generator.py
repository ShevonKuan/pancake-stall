import os
import json

# 请在项目根目录执行该程序

# 设置诗词仓库路径
CHINESE_POETRY_REPO = './chinese-poetry'

CATALOGUE = {}

def traverse_directory(path):
    '''遍历并输出所有包含诗词的json路径'''
    for root, _, files in os.walk(path):
        for file in files:
            if file.endswith('.json'):
                file_path = os.path.join(root, file)
                with open(file_path, 'r', encoding='utf-8') as file_stream:
                    content = ''.join(file_stream.readlines())
                    if ('para' in content) or ('content' in content):
                        category = file_path.split('\\')[1]
                        if category != 'loader':
                            CATALOGUE[category] = CATALOGUE.get(category, []) + [file_path]

traverse_directory(CHINESE_POETRY_REPO)

with open('PoetryCategory.json', "w", encoding='utf-8') as file:
    json.dump(CATALOGUE, file, ensure_ascii=False)