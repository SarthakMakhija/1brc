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

real	0m0.655s
user	0m0.786s
sys	0m0.037s
```

### Changes

1. 08ec5886e9a632fbc22dab0c32808fe2de05e618: Changed the map value to be a pointer to StationTemperatureStatistics to save the cost of Put operation.
2. 1277b5b155bd1d1e7b68dd5fea599718e8c349cb: Added custom Split of bytes (line is treated as bytes), but this commit has byte slice to string conversion.
3. 41d903a0959e2151d015a09d1c850e195175d272: Added a zero copy byte slice to string + does a slice copy in SplitIntoStationNameAndTemperature.

```shell
time ./main -f ../fixture/1M_weather_stations.csv

real	0m0.269s
user	0m0.266s
sys	0m0.013s
```