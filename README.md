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

This branch `profiling_1` focuses on improving CPU bound and IO operations for a single core.

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
time ./main -f ../fixture/1M_weather_stations.csv

real	0m0.242s
user	0m0.240s
sys	0m0.008s
```

### Changes: Round3

1. TreeMap to SwissMap: df968e00890879a6b8ad92aee42d8828d4e8ea2f: 

```shell
time ./main -f ../fixture/1M_weather_stations.csv

real	0m0.085s
user	0m0.068s
sys	0m0.001s
```

### Changes: Round4

1. Change convert.go to handle a single fractional digit: c13a5e880381cc7baf3a3ceb0ee3f460b3d090a1
2. Avoid the cost of uint16 conversion in convert.go: 955b97663012557163a25514adcec971d5fff2df

```shell
time ./main -f ../fixture/1M_weather_stations.csv

real	0m0.129s
user	0m0.102s
sys	0m0.008s
```

### Changes: Round 5

1. Optimizes stringify: 84617a2dcb2cc6be189e0176d6e8dd804bca51af
2. Optimizes PrintableResult: 95b67977fde89049efd97d0a40218e9e3ce87d2c

```shell
time ./main -f ../fixture/1M_weather_stations.csv

real	0m0.064s
user	0m0.064s
sys	0m0.000s
```

### Changes: Round 6

1. Optimizes stringify by using AppendFloat: 4b36cb3d9c844a6378e1c6857ef28f2baecf5d1b
2. Optimizes PrintableResult by using a common byte slice: a498ac67a57c2c1760e4ce6ae08116e3ce3e2269

```shell
time ./main -f ../fixture/1M_weather_stations.csv

real	0m0.089s
user	0m0.052s
sys	0m0.015s
```

```shell
time ./main -f ../fixture/10M_weather_stations.csv

real	0m0.532s
user	0m0.487s
sys	0m0.032s
```

Quick analysis:
Consider that we want to parse the file with one billion rows **(1000 million) in 5 seconds**.

This means the target is:
- **1000 million rows** in **5 seconds**
- **1000 million rows** in **5000000000 nanoseconds**
- **10^9 rows** in **5000000000 nanoseconds**
- **Each row** should be processed in **5 nanoseconds**

Based on the above numbers:
- **10 million rows** are being processed in **0.532 seconds**
- **10^7 rows** in **532000000 nanoseconds**
- **Each row** is being processed in **53.2 nanoseconds**

### Changes: Round 7

1. Optimizes convert: f078d76e953495ae618e6b3d6c61cf9a75e40cd0


### Changes: Round 8

1. Removes one bound check in convert: 68a9b2eed52f47f118d79e350aa16941a8e91b12
2. Modifies Parser (ParserV2) to use BufferedReader: 7875185495f428703e9f58421547685dcb19a543
3. Avoids calculation of average for each row: 3f76070fa309bc6b6206080e33acaa12ba658f13
4. Avoids the cost of byte slice to string conversion in stringify: e7088f65d46c23cef127c550cf0ff56275707cbe

```shell
time ./main -f ../fixture/10M_weather_stations.csv

real	0m0.442s
user	0m0.392s
sys	0m0.035s
```

- **10 million rows** are being processed in **0m0.442s seconds**
- **10^7 rows** in **442000000 nanoseconds**
- **Each row** is being processed in **44.2 nanoseconds**

### Changes: Round 9

1. Changes SplitIntoStationNameAndTemperature to iterate from end to begin: a2aa8cce955647786e41d52ebe2bdc00fd0fc63f
2. Represents temperature as int16: e3964256fafc738bbad832c36bc56983752cf594

```shell
time ./main -f ../fixture/10M_weather_stations.csv

real	0m0.422s
user	0m0.388s
sys	0m0.020s
```
- **10 million rows** are being processed in **0m0.422s seconds**
- **10^7 rows** in **422000000 nanoseconds**
- **Each row** is being processed in **42.2 nanoseconds**

### Changes: Round 9

1. Changes SplitIntoStationNameAndTemperature handle conversion of temperature: 27c53c8dd9e09d55517919cd09426f943fa95cca

```shell
time ./main -f ../fixture/10M_weather_stations.csv

real	0m0.405s
user	0m0.390s
sys	0m0.016s
```
- **10 million rows** are being processed in **0m0.405s seconds**
- **10^7 rows** in **405000000 nanoseconds**
- **Each row** is being processed in **40.5 nanoseconds**
