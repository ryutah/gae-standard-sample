# Cloud Scheduler, Cloud Tasksに関して
## Cloud Scheduler
まだBeta版ということもあると思うが、AppEngineでのみ利用することを考えるとAppEngine Cron Jobを使ったほうが制御が楽そう

### 問題点
* **X-Cloudscheduler** ヘッダが手動で上書き可能なため、Cloud Schedulerからのリクエストであるという識別には利用できない
  - Cron Jobの場合は **X-Appengine-Cron** を利用可能
* **X-Appengine-User-Ip** で判定できそうではあるが、Cloud SchedulerのIPアドレスに関するドキュメントが見当たらないため、ちょっと使用がためらわれる
  - Cron Jobの場合は[ドキュメント](https://cloud.google.com/appengine/docs/standard/python/config/cronref?hl=ja#originating_ip_address)がある

#### Cloud Schedulerから送信されたリクエストヘッダ
```console
X-Appengine-Country: [ZZ] 
X-Cloudscheduler: [true] 
User-Agent: [AppEngine-Google; (+http://code.google.com/appengine)] 
X-Appengine-Default-Version-Hostname: [sandbox-hara.appspot.com] 
X-Appengine-Request-Log-Id: [5c78c16c00ff0cb7ff003a3fd6db0001627e73616e64626f782d6861726100017461736b733a323031393033303174313431373434000100] 
X-Appengine-Api-Ticket: [9e4e51960a6f7790] 
X-Appengine-User-Ip: [0.1.0.2] 
X-Cloud-Trace-Context: [57c88a903baaa3765688ce5a710e6a29/2006728061747462277;o=1] 
Content-Length: [0] 
X-Forwarded-Proto: [http] 
X-Forwarded-For: [0.1.0.2, 169.254.1.1] 
X-Appengine-Https: [off] 
```

#### Curlで叩いた場合のリクエストヘッダ
```console
$ curl -X GET -H "X-Cloudscheduler: true" -H "X-Appengine-User-Ip: 0.1.0.2" "https://tasks-dot-sandbox-hara.appspot.com/cron"

X-Appengine-Request-Log-Id: [5c78c23900ff0e99e9667ec6270001627e73616e64626f782d6861726100017461736b733a323031393033303174313431373434000100] 
Accept: [*/*] 
X-Forwarded-For: [118.238.209.228, 169.254.1.1] 
X-Appengine-Https: [on] 
X-Appengine-Citylatlong: [35.708068,139.752167] 
X-Appengine-Region: [13] 
X-Appengine-User-Ip: [118.238.209.228] 
X-Appengine-Country: [JP] 
X-Appengine-City: [bunkyo] 
X-Cloudscheduler: [true] <- 上書き可能
User-Agent: [curl/7.54.0] 
X-Forwarded-Proto: [https] 
X-Appengine-Default-Version-Hostname: [sandbox-hara.appspot.com] 
X-Appengine-Api-Ticket: [e9af8feca2eb9c52] 
X-Cloud-Trace-Context: [086fa0c30f7a1aadc5bfce844e2654ae/12744968701236989972;o=1] 
```

#### Cron Jobから送信されたリクエストヘッダ
```console
-Appengine-Taskexecutioncount: [0] 
User-Agent: [AppEngine-Google; (+http://code.google.com/appengine)] 
X-Forwarded-For: [0.1.0.1, 169.254.1.1] 
X-Appengine-Api-Ticket: [4fc4916492d14ad5] 
X-Appengine-Request-Log-Id: [5c78c5e300ff0b81fbc7b2cb170001627e73616e64626f782d6861726100017461736b733a323031393033303174313431373434000100] 
X-Appengine-Country: [ZZ] 
X-Forwarded-Proto: [http] 
X-Cloud-Trace-Context: [d5f67d469ff55947a3a0b2eefb986973/17410141928712907195;o=1] 
X-Appengine-Tasketa: [1551423600.6787338] 
X-Appengine-Queuename: [__cron] 
X-Appengine-Cron: [true] 
X-Appengine-User-Ip: [0.1.0.1] 
X-Appengine-Default-Version-Hostname: [sandbox-hara.appspot.com] 
X-Appengine-Https: [off] 
X-Appengine-Taskretrycount: [0] 
X-Appengine-Taskname: [d5155a15703831adea4a21ec3007816c] 
```

#### Curlで叩いた場合のリクエストヘッダ
```console
$ curl -X GET -H "X-Appengine-Cron: true" -H "X-Appengine-User-Ip: 0.1.0.1" "https://tasks-dot-sandbox-hara.appspot.com/cron"

X-Appengine-Region: [13] 
X-Appengine-City: [bunkyo] 
User-Agent: [curl/7.54.0] 
X-Forwarded-Proto: [https] 
X-Appengine-Request-Log-Id: [5c78c66e00ff07057f0d582d270001627e73616e64626f782d6861726100017461736b733a323031393033303174313431373434000100] 
X-Appengine-Country: [JP] 
X-Cloud-Trace-Context: [dce74b0c732261930d063075254cd5f3/3736680919647464296;o=1] 
X-Appengine-Api-Ticket: [e87ce70350ec2e7f] 
X-Appengine-User-Ip: [118.238.209.228] 
Accept: [*/*] 
X-Forwarded-For: [118.238.209.228, 169.254.1.1] 
X-Appengine-Default-Version-Hostname: [sandbox-hara.appspot.com] 
X-Appengine-Https: [on] 
X-Appengine-Citylatlong: [35.708068,139.752167] 
```

## Cloud Tasks
Taskqueueと同じ様な感覚で使えそうだが、 **ローカルエミュレータがない** ので今までと比べると、開発時にちょっと一工夫必要そう
また、 `login: admin` が利用できなくなったため、アクセス制御にリクエストヘッダの情報( **X-Appengine-Queuename** など)を使うことになると思うのだが、検証の際に手動でタスクを実行したいというようなケースで多少手間が増えそう。
(gcloud コマンドで代用できそうではある)

### Cloud Tasksから送信されたリクエストヘッダ
```console
X-Appengine-Default-Version-Hostname: [sandbox-hara.appspot.com] 
X-Cloud-Trace-Context: [8b0881964d49ac61a17094b79f3d458e/7518162044529025862;o=1] 
X-Forwarded-Proto: [http] 
X-Forwarded-For: [0.1.0.2, 169.254.1.1] 
X-Appengine-User-Ip: [0.1.0.2] 
X-Appengine-Https: [off] 
X-Appengine-Taskexecutioncount: [0] 
X-Appengine-Taskretrycount: [0] 
User-Agent: [AppEngine-Google; (+http://code.google.com/appengine)] 
X-Appengine-Api-Ticket: [5bed68292d2226e1] 
X-Appengine-Country: [ZZ] 
X-Appengine-Tasketa: [1551419303.010438] 
X-Appengine-Taskname: [0724593178424063128] 
X-Appengine-Queuename: [example-queue] 
X-Appengine-Request-Log-Id: [5c78c7a700ff018b5aa8ab23a10001627e73616e64626f782d6861726100017461736b733a323031393033303174313431373434000100] 
```

### Curlで叩いた場合のリクエストヘッダ
```console
$ curl -X GET -H "X-Appengine-Queuename: example-queue" "https://tasks-dot-sandbox-hara.appspot.com/tasks"

X-Forwarded-Proto: [https] 
User-Agent: [curl/7.54.0] 
X-Appengine-Citylatlong: [35.708068,139.752167] 
X-Appengine-Region: [13] 
X-Cloud-Trace-Context: [5a7fe2538e54ad730f8d59bc27eb92b2/12726832335119258647] 
X-Appengine-Request-Log-Id: [5c78c83b00ff07626c412a26650001627e73616e64626f782d6861726100017461736b733a323031393033303174313431373434000100] 
X-Appengine-Api-Ticket: [d2e97614b9949d43] 
X-Appengine-Country: [JP] 
Accept: [*/*] 
X-Forwarded-For: [118.238.209.228, 169.254.1.1] 
X-Appengine-Default-Version-Hostname: [sandbox-hara.appspot.com] 
X-Appengine-User-Ip: [118.238.209.228] 
X-Appengine-Https: [on] 
X-Appengine-City: [bunkyo] 
```
