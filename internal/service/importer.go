package service

import (
	"bufio"
	"context"
	"cooling-tower-diagnostics/internal/model"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

type Importer struct{ Service *Service }

func (i Importer) CSV(ctx context.Context, r io.Reader) (int, error) {
	if i.Service == nil {
		return 0, model.ErrInvalid
	}
	reader := csv.NewReader(bufio.NewReader(r))
	reader.ReuseRecord = false
	count := 0
	for {
		record, e := reader.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			return count, e
		}
		if len(record) < 4 || strings.EqualFold(record[0], "tower_id") {
			continue
		}
		value, e := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		if e != nil {
			return count, e
		}
		req := model.ReadingRequest{TowerID: strings.TrimSpace(record[0]), Sensor: strings.TrimSpace(record[1]), Value: value, Unit: strings.TrimSpace(record[3])}
		if _, e := i.Service.ValidateAndIngest(ctx, req); e != nil {
			return count, e
		}
		count++
	}
	return count, nil
}
func (i Importer) JSONLines(ctx context.Context, r io.Reader) (int, error) {
	if i.Service == nil {
		return 0, model.ErrInvalid
	}
	count := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		req, e := i.Service.ParsePayload(scanner.Bytes())
		if e != nil {
			return count, e
		}
		if _, e = i.Service.ValidateAndIngest(ctx, req); e != nil {
			return count, e
		}
		count++
	}
	return count, scanner.Err()
}
