[OPENLOGI Advent Calendar 2022](https://qiita.com/advent-calendar/2022/openlogi) 16日目の記事です。
https://qiita.com/advent-calendar/2022/openlogi

この記事では、気圧・温度・湿度・CO2濃度ツイートシステムについて紹介します。完全に趣味の話で、以下のようにツイートしてくれるシステムについての記事です。
![](https://storage.googleapis.com/zenn-user-upload/4a955b1efa61-20221211.png)

皆様は環境と体調が連動して変化していると感じることはないでしょうか？「なんだか今日は身体が重く感じる」「頭痛がする」けど、原因に心当たりがない。そんな時、自分を取り巻く環境にヒントがあるかもしれません（ないかもしれません）。

# きっかけ
このシステムを作ろう&改造していこうと思ったきっかけがいくつかあります。

- 低気圧だと頭痛がする気がする。気圧を計測できるなら確認してみたい。
- 冬場だとCO2濃度が上昇して具合が悪くなるそう？CO2濃度も出してみたい。
- Raspberry Piが自宅にあり、センサー類があれば繋げて環境に関わる数値を計測できる。楽しそう。
- はんだ付けするのね。中学1年生以来で懐かしい。
- パッと把握したいな。
- AWSのサービス使ってみたい。

いろいろありますね（まだあってissueを切っています）。

それでは、ここからシステムの紹介に入っていきたいと思います。

# 機能
本システムの機能は以下です。
1. BME280から気圧・温度・湿度 / MH_Z19からCO2濃度を取得
2. 取得したデータを使って線グラフを生成
3. データ取得と線グラフ生成は以下タイミングで実行し、ツイート
	a. 毎日09:00/18:00
	b. Bluetoothリモコンシャッターのボタン押下

MH_Z19はCO2濃度を計測するセンサーです(金色)。BME280は気圧・温度・湿度が計測できるセンサーです(紫色)。
![](https://storage.googleapis.com/zenn-user-upload/1150314ca550-20221211.jpg)

これらセンサー類を直接見た時、「ちっさ！これで計測できるのか凄！」と思ったのを覚えています。

# 構成
本システムの構成は以下です。矢印はデータの流れを表しています。
![](https://storage.googleapis.com/zenn-user-upload/209d6c1cbfa6-20221211.png)

...モノリスでも機能は実現できます。ですが、ばらけさせたり、AWSのサービスを使いたかったのでこうしました。

# 処理の流れ
## トリガー
まず入口から見ていきます。
このシステムを起動するトリガーになるのが、構成の左下の部分です。2つのトリガーがあります。
1. 「Execute Dagu at 9 and 18」
	- 毎日09:00/18:00にシステムを起動します
2. 「Bluetooth remocon」
	- Bluetoothリモコンシャッターのボタン押下時にシステムを起動します

1は、[Dagu](https://github.com/yohamta/dagu)というツールを使わせていただいています🙇。
https://zenn.dev/yohamta/articles/959f91fb1b505b

2は、リモコンのボタンを押下することでシステムを起動します。（近くの100均で購入）
![](https://storage.googleapis.com/zenn-user-upload/00106212ad1c-20221211.jpg)
[こちらのツイート](https://twitter.com/sozoraemon/status/1574969255208325125?s=20&t=i2r26kbyVx00GJu13hcrzw)に触発されてシステムに組み込みました🙇（これでいつでも確認できる！）
https://twitter.com/sozoraemon/status/1574969255208325125?s=20&t=i2r26kbyVx00GJu13hcrzw

ここまでで2つのトリガーを紹介しました。これらトリガー後の処理は共通です。

## センサー値取得
トリガーによってシステムの起動後、まず各センサーから計測値を取得します。構成の左上の部分です。
MH_Z19（CO2濃度）、BME280（気圧・温度・湿度）用に、それぞれで専用のプログラムがあります。
取得したセンサー値はSQLite3に保存します。

## 線グラフ画像の生成
この節は構成の右側の説明です。
SQLite3に蓄積された各種センサー値（CO2濃度・気圧・温度・湿度）は、直近の日時から10レコード取得され、Amazon SNSにpublishされます。
そして、subscribeしているAWS Lambdaが、センサー値をもとに線グラフ画像を生成します。
また、このLambdaでの処理は、生成した画像データをbase64エンコードし、Amazon SQSへキューイングします。

## そして、ツイート
構成の左側中心あたりに、「Pull and Decode image data and Post tweet」というプログラムがいるのですが、その通りの処理をします。

キューイングされた画像データ（base64エンコード済み）をAmazon SQSからpullして、base64でデコードし画像データとして復元後、ツイートしています。

# 発見
私は、「低気圧だと頭痛がする気がする」と思っていましたが、低気圧に限らず「気圧差があると頭痛がする」傾向がありそうと気付きました。（そうでない日もありますが...）

あと、台風が接近したときの気圧差を見れたのが面白かったです。（以下は2022/08/13 18:00時点の画像。台風8号（メアリー））
![](https://storage.googleapis.com/zenn-user-upload/fa3dde4956cd-20221211.png)


# 課題

いろいろあるかと思うのですが、特にと言うと、今の実装は全く疎結合に出来ていない事です。なので、新しくセンサーを追加したい場合、既存プログラムへの修正が必要になってくるため、極力新規にプログラムを追加するだけでシステムが動くようにはしたいと思っています。

あと、もっと短いスパンでセンサー値を取得し、蓄積するデータ量を増やして何かしたいなと思っています。例えば、
[こういうの](https://zenn.dev/thorie/articles/548iot_room_condition_aws_timestream_grafana_cloud)や、
https://zenn.dev/thorie/articles/548iot_room_condition_aws_timestream_grafana_cloud

[こういうの](https://twitter.com/siroitori0413/status/1582936266307678208?s=20&t=cVtHml8x26W5weixqnys6g)
https://twitter.com/siroitori0413/status/1582936266307678208?s=20&t=cVtHml8x26W5weixqnys6g
をやりたいですね。

# 最後に
よければぜひ！SRE / CRE も絶賛募集中です！（2024/11/18 現在）
https://herp.careers/v1/openlogi/requisition-groups/486b8b01-6cf9-4434-8601-381c9c092e0d

# 参考
- [本システムを含むリポジトリ](https://github.com/ddddddO/sensor-pi)
	- READMEやコードにメモしているのでご興味があれば是非！

- 「構成」章のイメージは、https://github.com/mingrammer/diagrams を使って作成しています。
- 線グラフ画像は、https://github.com/gonum/plot で生成しています。
