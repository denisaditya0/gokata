package support

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
)

type S struct {
	Name    string
	Env     []string
	Project []string
	Serial  bool
	Steps   func()
}

func RunScenarios(scenarios []S) {
	for _, sc := range scenarios {
		sc := sc
		labels := buildLabels(sc.Env, sc.Project)

		runFunc := func() {
			Ctx().Clear()
			ResetScenarioLogs()
			Ctx().LoadTestData(sc.Name)

			EmitEvent(EventScenarioStart, map[string]interface{}{
				"name": sc.Name, "labels": labels,
			})

			if ListMode {
				fmt.Printf("\n%s\n", sc.Name)
			}
			sc.Steps()

			EmitEvent(EventScenarioEnd, map[string]interface{}{
				"name": sc.Name, "status": "passed",
			})

			if logs := GetScenarioLogs(); len(logs) > 0 {
				AddReportEntry("http_logs", ReportEntryVisibilityNever, logs)
			}
		}

		if sc.Serial {
			if len(labels) > 0 {
				It(sc.Name, Serial, Label(labels...), runFunc)
			} else {
				It(sc.Name, Serial, runFunc)
			}
		} else {
			if len(labels) > 0 {
				It(sc.Name, Label(labels...), runFunc)
			} else {
				It(sc.Name, runFunc)
			}
		}
	}
}

func buildLabels(groups ...[]string) []string {
	var labels []string
	for _, group := range groups {
		labels = append(labels, group...)
	}
	return labels
}
