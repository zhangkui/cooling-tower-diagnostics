package diagnostics

import (
	"cooling-tower-diagnostics/internal/model"
	"fmt"
	"strings"
)

type Recommendation struct {
	Code     string
	Title    string
	Priority string
	Actions  []string
}

func Recommend(d model.Diagnosis) []Recommendation {
	out := []Recommendation{}
	for _, f := range d.Findings {
		low := strings.ToLower(f)
		if strings.Contains(low, "conductivity") {
			out = append(out, Recommendation{Code: "CHEM-01", Title: "检查排污与补水比例", Priority: d.Risk, Actions: []string{"确认电导率探头清洁", "核对排污阀动作", "复核补水水源"}})
		}
		if strings.Contains(low, "temp") {
			out = append(out, Recommendation{Code: "THERM-01", Title: "检查换热器与风机", Priority: d.Risk, Actions: []string{"检查风机频率", "查看进出水温差", "确认环境湿球温度"}})
		}
	}
	if len(out) == 0 {
		out = append(out, Recommendation{Code: "OBS-01", Title: "继续观察采样窗口", Priority: "info", Actions: []string{"保持传感器在线", "等待下一窗口"}})
	}
	return out
}
func RecommendationText(rs []Recommendation) string {
	parts := []string{}
	for _, r := range rs {
		parts = append(parts, fmt.Sprintf("[%s] %s", r.Code, r.Title))
	}
	return strings.Join(parts, "; ")
}
