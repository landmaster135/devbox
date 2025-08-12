# Pythonのタブサイズがどうしても「4」にならない_20241215_test_sample

priority: 5
status: 3. Done (https://www.notion.so/3-Done-eeeeeeeeeeeeeeeeeeeeeee?pvs=21)
linking_tasks: あゆあゆのDailyReport_2024-12-08(Sun) (%E3%81%82%E3%82%86%E3%81%82%E3%82%86%E3%81%AEDailyReport_2024-12-08(Sun)%ppppppppppppppppppppppppppppp.md)
tags: blog (https://www.notion.so/blog-vvvvvvvvvvvvvvvvvvvvvvvvvv?pvs=21), WordPress (https://www.notion.so/WordPress-xxxxxxxxxxxxxxxxxxxxxxxxxx?pvs=21), GoogleAnalytics (https://www.notion.so/GoogleAnalytics-llllllllllllllllllllllllllllll?pvs=21), GoogleSearchConsole (https://www.notion.so/GoogleSearchConsole-kkkkkkkkkkkkkkkkkkkkkkkkkkk?pvs=21), test_sample (https://www.notion.so/test_sample-yyyyyyyyyyyyyyyyyyyyyy?pvs=21)
estimated_dev_minutes: 240
actual_dev_minutes: 150
done_at: December 8, 2024
artifacts: ブログ活動（投稿部）：エンドルに浸かる (https://www.notion.so/bbbbbbbbbbbbbbbbbbbbbbbbb?pvs=21), プログラミングハッカー活動：test (https://www.notion.so/test-ggggggggggggggggggggggggg?pvs=21)
created_time: December 8, 2024 1:36 PM
id: TK-1999

# Draft

- [Pythonのタブサイズがどうしても「4」にならない｜エンドルに浸かる。](https://www.example.com/python-tabsize-not-4-but-2)
- [【Python】Cloud RunによるFastAPIアプリからCloud SQLのPostgreSQLに接続する（第2回）_20241115_test_sample](%E3%80%90Python%E3%80%91Cloud%20Run%E3%81%AB%E3%82%88%E3%82%8BFastAPI%E3%82%A2%E3%83%95%E3%82%9A%E3%83%AA%E3%81%8B%E3%82%89Cloud%20SQL%E3%81%AEPostgre%333333333333333333333333333333333333.md)
- [How can I customize the tab-to-space conversion factor in VS Code? - Stack Overflow](https://stackoverflow.com/questions/29972396/how-can-i-customize-the-tab-to-space-conversion-factor-in-vs-code)
- [x]  ブログ記事を書く。・・・合計105分掛かった。
- [x]  「test.code-profile」をAWSに保存する。
- [x]  バックアップを取る。・・・合計15分掛かった。
- [x]  Closed.

# Materials

## Memo 1: Pythonのタブサイズがどうしても「2」になる。

変更前

```python
  "[python]": {
    "diffEditor.ignoreTrimWhitespace": false,
    "editor.tabSize": 4,
    "editor.formatOnPaste": true,
    "editor.formatOnSave": true
  },
```

変更後

```python
  "[python]": {
    "diffEditor.ignoreTrimWhitespace": false,
    "editor.tabSize": 4,
    "editor.indentSize": 4, // <- New
    "editor.formatOnPaste": true,
    "editor.formatOnSave": true
  },
```

# Content

## はじまり

サンプルちゃん

コイツ、いつまで経っても「2」だな・・・

135ml

いつまでも「4」だな。

<br>

## 解決法その1：設定のTab Size

最も基本的な項目です。

ここの設定では、基本的なタブのサイズが設定されます。しかし、僕は基本的にタブのサイズは「2」が良いと思っているので、ここは「2」のままで、Pythonの時だけ「4」にして欲しい。

![](https://www.example.com/wp-content/uploads/2024/12/sample17.webp)

<br>

ちなみに、settings.jsonだと`"editor.tabSize"`の項目に該当します。

<br>

## 解決法その2：「Detect Indentation」の設定

「Detect Indentation」という設定項目も見ておきます。

この設定がオンになっていると、ファイルのインデントを見て、そのファイルにおけるインデントのサイズを設定します。オフになっていれば、エディタで設定してあるインデントのサイズしか反映されないわけです。

個人的には、settings.jsonで設定した値だけを反映してほしいので、オフにしておきます。

![](https://www.example.com/wp-content/uploads/2024/12/sample18.webp)

<br>

ちなみに、settings.jsonだと`"editor.detectIndentation"`の項目に該当します。Falseにすればオフになります。

<br>

## 解決法その4：.editorconfigのindent_size

「Editorconfig」は、コードエディタ間でファイル(コード)のフォーマットを一貫させるための機能です。そして、「.editorconfig」は、その設定ファイルの形式です。JetBrainsの各IDEや、AtomとかXcodeでも使えるみたいですね。

その.editorconfigをこんな風に設定してみます。

```yaml
# Editor configuration, see https://editorconfig.org
root = true

[*.py]
indent_size = 4

```

<br>

設定してみると・・・やりました！

Pythonファイルの時にだけ、インデントのサイズが4になりました！

![](https://www.example.com/wp-content/uploads/2024/12/sample20.webp)

<br>

うーんでもなあ・・・、インデントのサイズを設定するためにわざわざ1つファイルを増やさなければならないのが気に入りません・・・。他に方法は無いのか？

## 解決法その5：settings.jsonでどちらか片方しか設定していない。

ここまで記事を書いてきて、何か気になります・・・。タブサイズ？　インデントサイズ？

そう、もしかしたら、settings.jsonの中で、タブサイズもしくはインデントサイズの片方しか設定していないのではないでしょうか・・・。

僕のVSCodeがPythonのインデントをずっと「2」で表示し続けていた時のsettings.jsonのPythonにおける設定がこれです。

```json
...

  "[python]": {
    "diffEditor.ignoreTrimWhitespace": false,
    "editor.tabSize": 4,
    "editor.formatOnPaste": true,
    "editor.formatOnSave": true
  },

...
```

<br>

そこで、このようにsettings.jsonを変更します。

```json
...

  "[python]": {
    "diffEditor.ignoreTrimWhitespace": false,
    "editor.tabSize": 4,
    "editor.indentSize": 4, // <- New
    "editor.formatOnPaste": true,
    "editor.formatOnSave": true
  },

...
```

<br>

設定してみると・・・やりました！

Pythonファイルの時にだけ、タブおよびインデントのサイズが4になりました！

![](https://www.example.com/wp-content/uploads/2024/12/sample21.webp)

<br>

ちなみに、settings.jsonだと`"editor.tabSize"`もしくは`"editor.indentSize"`の項目に該当します。

## まとめ

この記事を見ることで、古今東西の人のPython利用者のタブおよびインデントのサイズが「4」になることを願っています。

<br>

タブおよびインデントのサイズをどうにかしたいと思っている人は、他にもいらっしゃるみたいですね。（8スペースのタブを4スペースのインデントを設定する？　そんなことあるの？）

https://stackoverflow.com/questions/29972396/how-can-i-customize-the-tab-to-space-conversion-factor-in-vs-code

<br>

## Python関連の書籍

Python2年生 デスクトップアプリ開発のしくみ 体験してわかる！会話でまなべる！

[https://amzn.to/5555555](https://amzn.to/5555555)

<br>

動かして学ぶ！Python FastAPI開発入門

[https://amzn.to/666666666666](https://amzn.to/666666666666)

<br>

動かして学ぶ! Pythonサーバレスアプリ開発入門

[https://amzn.to/777777777777777](https://amzn.to/77777777777777)

<br>

## おしまい

サンプルちゃん

やっと「4」になった・・・！！

135ml

一生「2」にならないで欲しいな。

<br>

以上になります！
