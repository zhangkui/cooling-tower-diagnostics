# 工业冷却塔水质与能效联动诊断服务

一个可本地运行的 Go 单体服务，用于冷却塔传感器采集、质量校准、窗口诊断、告警升级与审计追踪。

## 启动

```powershell
go run ./cmd/server
```

默认监听 `http://localhost:8080`，数据写入 `./data/cooling.db`。可通过 `DATA_DIR` 指定目录。

## 示例

```powershell
curl http://localhost:8080/healthz
curl http://localhost:8080/api/towers
curl -X POST http://localhost:8080/api/towers -H 'Content-Type: application/json' -d '{"name":"北区一号塔","site":"North","design_kw":420}'
curl -X POST http://localhost:8080/api/readings -H 'Content-Type: application/json' -d '{"tower_id":"tower-1","sensor":"conductivity","value":1800,"unit":"uS/cm"}'
curl http://localhost:8080/api/alerts
curl 'http://localhost:8080/api/water-balance?tower_id=tower-1&from=2026-08-26T00:00:00Z&to=2026-08-26T01:00:00Z&tolerance_m3=0.5'
```

## 验证

```powershell
go test ./...
go build ./...
```

前端页面位于 `/`，仅作为后端服务的操作入口，核心逻辑全部在 Go 服务端。

`/api/water-balance` 使用 `makeup_flow`、`blowdown_flow` 和 `evaporation_flow`（单位 `m3/h`）的相邻读数进行梯形积分。超过采样间隔的读数不会被外推，而会作为采集缺口返回，避免把网关离线误判为稳定水量。
