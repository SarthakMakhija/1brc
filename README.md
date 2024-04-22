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

### Changes: Round1

1. 08ec5886e9a632fbc22dab0c32808fe2de05e618: Changed the map value to be a pointer to StationTemperatureStatistics to save the cost of Put operation.
2. 1277b5b155bd1d1e7b68dd5fea599718e8c349cb: Added custom Split of bytes (line is treated as bytes), but this commit has byte slice to string conversion.
3. 41d903a0959e2151d015a09d1c850e195175d272: Added a zero copy byte slice to string + does a slice copy in SplitIntoStationNameAndTemperature.

```shell
time ./main -f ../fixture/1M_weather_stations.csv

real	0m0.269s
user	0m0.266s
sys	0m0.013s
```

### Changes: Round2

1. Custom toFloat64: 8fe65963c49d0a4000be185fe44b7cb2a2f081af, 4b1fcdb5805ee88a2a7b3c68ce1a3c168868b728, 85a4e07a0169041200ccffd9a59362e9c25947a4
2. Avoid temperature copy: 6d55773195e8a61f730b3294658e3edf83949591

```shell
real	0m0.242s
user	0m0.240s
sys	0m0.008s
```

### Changes: Round3

1. TreeMap to SwissMap: df968e00890879a6b8ad92aee42d8828d4e8ea2f: 

```shell
real	0m0.085s
user	0m0.068s
sys	0m0.001s
```

### Changes: Round4

1. Change convert.go to handle a single fractional digit: c13a5e880381cc7baf3a3ceb0ee3f460b3d090a1
2. Avoid the cost of uint16 conversion in convert.go: 955b97663012557163a25514adcec971d5fff2df

```shell
real	0m0.129s
user	0m0.102s
sys	0m0.008s
```

### Changes: Round 5

1. Optimizes stringify: 84617a2dcb2cc6be189e0176d6e8dd804bca51af
2. Optimizes PrintableResult: 95b67977fde89049efd97d0a40218e9e3ce87d2c

```shell
real	0m0.064s
user	0m0.064s
sys	0m0.000s
```