# bus-timetable-generator
[盛岡市バスロケーションシステム](http://gps.iwatebus.or.jp/pc/index.htm)に掲載されている運行情報を元にJSON形式で時刻表を出力します。  
このツールで出力したJSONは [c-ardinal/saketoba.countdown-timer](https://github.com/c-ardinal/saketoba.countdown-timer) で利用可能です。

# Usage
```
git clone https://github.com/c-ardinal/bus-timetable-generator
cd bus-timetable-generator
go get -d -v
go build
./bus-timetable-generator {FROM} {TO}
```

# Example
 - 盛岡駅から岩手県立大学までの時刻表生成
```
./bus-timetable-generator 盛岡駅前 岩手県立大学入口
```

- 岩手大学から盛岡バスセンターまでの時刻表生成
```
./bus-timetable-generator 岩手大学前 盛岡バスセンター
```

# Caution
 - 要はスクレイピングツールなのでご利用は計画的に。悪用厳禁。
 - バス停の名前は間違えないように気をつけて下さい。補完機能は有りません。(例: ×盛岡駅 -> ○盛岡駅前)

# License
MIT
