# toll-calculator


```
Run Apache Kafka
docker compose up
```


```
Run Prometheus
!!!!!----You should in the same folder as your prometheus.yml-----!!!!!
docker run --name prometheus -d --network host -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```


```
Run Gafana
docker run --name grafana -d --network host -p 3000:3000 grafana/grafana
```