% Unzip4Win 高度な使い方  
%  
% 2019/03/04  

Unzip4Winの高度な使い方を紹介します。
以下の説明では、unzip4win.exe（Macの場合はunzip4win）にパスが通っている前提で説明します。


# Debugログの出力

`-debug`オプションを有効にすることで、Debugログの出力を行います。

```bash
> unzip4win.exe -debug C:¥path¥to¥file.zip
```

# 設定の上書き

Unzip4Winはパスワードや展開後ファイルの出力先をアプリケーションに内包していますが、設定ファイルを読み込ませることで上書きすることができます。

`-config`オプションに続けて設定ファイル名を指定することで有効になります。

```bash
> unzip4win.exe -config .¥config.toml C:¥path¥to¥file.zip
```

設定ファイルは[TOML](https://ja.wikipedia.org/wiki/TOML)で記述します。
以下に完全なUnzip4Winの設定ファイルの例を記載します。

##### config.toml {#configtoml}

```toml
[output]
saveCurrent = true
outputPath = ""

[password]
tryDays = 10

[[spec]]
format = "password20060102"
startDate = 2019-02-01

[[spec]]
format = "20060102password"
startDate = 2019-01-01
```

### 各項目について

* `[output]` : 展開後ファイルの出力先を設定します
  * `saveCurrent` : `true`を指定した場合、zipファイルと同じディレクトリにファイルを展開します。
    `false`を指定した場合、`outputPath`で指定したディレクトリにファイルを展開します。
  * `outputPath` : `saveCurrent`に`false`を指定した場合にファイルを展開するディレクトリを指定します。
    環境変数の展開等は行わないため、フルパスで指定するようにしてください。
    作業フォルダからの相対パス指定も可能ですが、作業フォルダの設定で事故る可能性があるのでフルパスを推奨します。

* `[password]` : zipファイルの展開用パスワードについて設定します
  * `tryDays` : パスワード探索の際に、当日から数えて遡る日数を指定します。

* `[[spec]]` : パスワード探索の際に使用する設定を行います
    本セクションは複数指定することが可能です。
  * `format` : パスワードルールを指定します。
    日付のフォーマット指定文字は[Golang標準](https://golang.org/pkg/time/#Time.Format)のものを使用します。
        * 年 : `2006`(yyyy) or `06`(yy)
        * 月 : `01`(MM) or `1`(M)
        * 日 : `02`(dd) or `2`(d)
        * 時 : `15`(24h) or `3` or `03`(12h)
        * 分 : `04` or `4`
        * 秒 : `05` or `5`
        * TZ : `MST` or `-0700`
        * 曜日 : `Mon` or `Monday`
  * `startDate` : `format`で指定したパスワードルールの開始日付を指定します。


# 設定例

### Zipファイルの展開先を指定する

```toml
[output]
saveCurrent = false
outputPath = "C:¥Users¥<USER-NAME>¥Unzip"
```

### パスワードの探索を1年前まで実施する

```toml
[password]
tryDays = 365
```

### 独自のパスワードルールを指定する

※ パスワードルールを上書きしたプロセスでは、既存のパスワードルールはクリアされます。
例えば、[config.tomlの例](#configtoml)がアプリケーションに埋め込まれている場合に以下の指定を行った場合、
2019年2月1日以降の日付でも以下のパスワードルールで解析を試行します。

```toml
[[spec]]
format = "my-password-060102"
startDate = 2019-01-01
```
