### 1 billion row challenge

[![build](https://github.com/SarthakMakhija/1brc/actions/workflows/build.yml/badge.svg)](https://github.com/SarthakMakhija/1brc/actions/workflows/build.yml)

```shell
time ./main -f ../fixture/44K_weather_stations.csv

real	0m0.116s
user	0m0.105s
sys	0m0.025s
```

```shell
time ./main -f ../fixture/1M_weather_stations.csv

real	0m1.010s
user	0m1.242s
sys	0m0.018s
```