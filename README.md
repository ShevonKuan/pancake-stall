# 摊煎饼小游戏🥞

小红书上很火的摊煎饼小游戏golang生成实现🥞
## demo
![](demo/Demo_20x50_2.jpg)
## 最简单的使用方法
### 系统需求：
1. LaTeX编译环境：XeLaTeX编译器，tikz,ctex宏包
2. windows系统，最近事情有点多，就先不做交叉编译了

### 开始
1. 从release下载已经嵌入全唐诗数据库的`pancake-stall.exe`文件
2. 以`cmd`或`powershell`执行，其中命令行参数
   
    ```bash
    Flags:
        --enableLog          [高级选项] 在stdIO输出画布 (default true)
    -h, --help               help for 摊煎饼
        --output string      tex文件输出路径 (default "output.tex")
        --poetryCounts int   注入诗句条数，<=0表示使用注入位置（选字）次数为终止条件
        --resultCounts int   [高级选项] 数据库查询结果返回计数 (default 5000)
        --split              是否允许切断诗句 (default true)
        --sql string         sqlite数据库文件路径，不指定时使用自带的全唐诗数据库
        --start string       用于开始的第一个字
        --trialCounts int    [高级选项] 注入位置（选字）次数，设定停止条件，当尝试多少个注入位置依然无法注入时停止随机选择注入位置 (default 10)
        --xlength int        画布行数 (default 20)
        --ylength int        画布列数 (default 20)
    ```

    > 例如生成一个24行高，11列宽，以"花"字开始匹配，总共填入20条诗句的煎饼。

    ```bash
    .\pancake-stall.exe --xlength 24 --ylength 11 --start 花 --poetryCounts 20 --output output.tex
    ```

3. 使用`xelatex`编译生成出来的`tex`文件

   ```bash
    xelatex output.tex
   ```

4. 当生成的tex文件节点非常多时，可能会遇到以下报错，
   ```bash
    ! TeX capacity exceeded, sorry [main memory size=5000000].
   ```
    这时我们需要修改一下latex的编译内存设置`C:\texlive\2023\texmf-dist\web2c\texmf.cnf`(请更改为自己的路径)
    ```latex
    main_memory = 8000000 % words of inimemory available; also applies to inimf&mp
    extra_mem_top = 8000000     % extra high memory for chars, tokens, etc.
    extra_mem_bot = 8000000     % extra low memory for boxes, glue, breakpoints, etc.

    % ConTeXt needs lots of memory.
    extra_mem_top.context = 8000000
    extra_mem_bot.context = 8000000

    % Words of font info for TeX (total size of all TFM files, approximately).
    % Must be >= 20000 and <= 147483647 (without tex.ch changes).
    font_mem_size = 8000000
    ```

## 从头开始
### 同步上游`chinese-poetry`仓库<https://github.com/chinese-poetry/chinese-poetry>
```bash
git submodule init
git submodule update
```
### 生成数据库文件
请根据自己需求修改以下文件
```bash
python utilities/file_name_index_generator.py
python utilities/poetry_database_genrator.py
```
### 编译可执行文件
```bash
go build .
```
## TODO
由于时间有限，最近需要备考，因此没有太多时间，请见谅
- 使用web前端来可视化，代替latex
- 整理接地气的数据库
- 横向和纵向可能存在重复，特别是像高中数据库这种，条目非常少的如
  ![](demo/Demo_高中数据库_1.jpg)
