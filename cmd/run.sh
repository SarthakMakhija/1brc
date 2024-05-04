#!/bin/bash

cd ../
echo "Performing build"
go build -v ./...

cd cmd/
echo "Creating executable"
rm -rf main
go build -o main main.go

echo "Running"
for (( i = 0; i < 5; i++ )); do
  { time ./main -f ../fixture/10M_weather_stations.csv; } 2>&1 | tee output.txt | tail -3
  sleep 2
  echo ---
done