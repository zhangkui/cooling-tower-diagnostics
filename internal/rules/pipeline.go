package rules

import (
	"context"
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"time"
)

type Stage interface {
	Name() string
	Apply(context.Context, []model.Reading) ([]model.Reading, error)
}
type Pipeline struct{ stages []Stage }

func NewPipeline(stages ...Stage) *Pipeline {
	return &Pipeline{stages: append([]Stage(nil), stages...)}
}
func (p *Pipeline) Add(stage Stage) {
	if stage != nil {
		p.stages = append(p.stages, stage)
	}
}
func (p *Pipeline) Names() []string {
	out := []string{}
	for _, s := range p.stages {
		out = append(out, s.Name())
	}
	return out
}
func (p *Pipeline) Run(ctx context.Context, rows []model.Reading) ([]model.Reading, error) {
	current := append([]model.Reading(nil), rows...)
	for _, stage := range p.stages {
		if err := ctx.Err(); err != nil {
			return current, err
		}
		next, err := stage.Apply(ctx, current)
		if err != nil {
			return current, fmt.Errorf("stage %s: %w", stage.Name(), err)
		}
		current = next
	}
	return current, nil
}

type DropSuspect struct{}

func (DropSuspect) Name() string { return "drop-suspect" }
func (DropSuspect) Apply(ctx context.Context, rows []model.Reading) ([]model.Reading, error) {
	out := []model.Reading{}
	for _, r := range rows {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
		}
		if r.Quality != model.ReadingSuspect {
			out = append(out, r)
		}
	}
	return out, nil
}

type ClampValues struct {
	Min float64
	Max float64
}

func (ClampValues) Name() string { return "clamp-values" }
func (c ClampValues) Apply(ctx context.Context, rows []model.Reading) ([]model.Reading, error) {
	out := append([]model.Reading(nil), rows...)
	for i := range out {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		out[i].Value = Clamp(out[i].Value, c.Min, c.Max)
	}
	return out, nil
}

type FillMissing struct{ Interval time.Duration }

func (FillMissing) Name() string { return "fill-missing" }
func (f FillMissing) Apply(ctx context.Context, rows []model.Reading) ([]model.Reading, error) {
	if len(rows) < 2 {
		return rows, nil
	}
	out := append([]model.Reading(nil), rows...)
	for i := 1; i < len(rows); i++ {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		gap := out[i].RecordedAt.Sub(out[i-1].RecordedAt)
		if f.Interval > 0 && gap > f.Interval*2 {
			synthetic := out[i-1]
			synthetic.ID = ""
			synthetic.RecordedAt = out[i-1].RecordedAt.Add(f.Interval)
			out = append(out, synthetic)
		}
	}
	return out, nil
}

type SortByTime struct{}

func (SortByTime) Name() string { return "sort-time" }
func (SortByTime) Apply(ctx context.Context, rows []model.Reading) ([]model.Reading, error) {
	out := append([]model.Reading(nil), rows...)
	for i := 1; i < len(out); i++ {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		for j := i; j > 0 && out[j].RecordedAt.Before(out[j-1].RecordedAt); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}
