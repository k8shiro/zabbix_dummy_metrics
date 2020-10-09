# zabbix_dummy_metrics

zabbixにhost・item・triggerを登録。登録itemに疑似的にメトリクスを登録する。


- zabbix_init/main.go: zabbix初期登録スクリプト、オプションの指定に沿ってzabbixへのhost・item・trigger登録を実行する

オプション

```
-host 登録host数、203.0.113.1から順に203.0.113.2, 203.0.113.3...と指定数分hostを登録する
-item hostごとのitem数、各hostごとにZabbix trapper、Numeric(float)のitemをKeyをkey.0, key.1, ....として登録する。同時にtriggerも{203.0.113.1:key.0.last()}>90のようにitemごとに作られる
-zabbix zabbixのAPIをたたくためのURL、http://～～
-user zabbix user
-pass zabbix password
```

- zabbix_sender/main.go: zabbixへのメトリクス登録スクリプト、オプションで指定されたitemに対して0.0～100.0の値をランダムに毎秒登録

オプション

```
-host 対象ホスト数、zabbix_initで生成したhostが203.0.113.2, 203.0.113.3...のように対象ホスト数分選択される
-item hostごとの対象item数、対象ホストごとにkey.0, key.1, ...のように対象item数分選択される
-zabbix zabbixのtrapper登録先(ポートは50051で固定している) ※ httpはつけない
```



- 実行例

以下を実行すると

- host
    - 203.0.113.1、203.0.113.2、203.0.113.3
- item
    - 203.0.113.1:key.0、203.0.113.1:key.1
    - 203.0.113.2:key.0、203.0.113.1:key.1
    - 203.0.113.3:key.0、203.0.113.3:key.1
- trigger
    - すべてのitemに90より大きいとエラーになるtrigger

が登録され、毎秒0.0～100.0の値がランダムに登録される。

```
# go runコマンドをdocker containerで実行する
docker image build -t zabbix_dummy_metrics .
docker run --rm -it -v $(pwd):/zabbix_dummy_metrics zabbix_dummy_metrics bash

# 500host 20item/hostを登録
go run zabbix_init/main.go -host=3 -item=2 -zabbix=http://<zabbix ip> -user=Admin -pass=zabbix

# 500×20item/s のメトリクス(float)を登録
go run zabbix_sender/main.go -host=3 -item=2 -zabbix=<zabbix ip>
```
