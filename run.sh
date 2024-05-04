#!/bin/bash

echo "Running"
for (( i = 0; i < 5; i++ )); do
  { time ./cmd/main -f ./fixture/10M_weather_stations.csv; } 2>&1 | tee output.txt | tail -3
  sleep 2
  echo ---
done