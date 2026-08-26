基于 Go 实现的工业冷却塔水质与能效联动诊断服务交付镜像

The image runs the durable SQLite-backed HTTP service. Build it with:

```sh
./build_benzhi_docker.sh cooling-tower-diagnostics linux/amd64
```

The runtime database is created below `DATA_DIR`; no database files are part of the source snapshot.
