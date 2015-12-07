package stats

import (
	resSt "github.com/thoas/stats"
)

var Metrics *resSt.Stats

func init() {
	Metrics = resSt.New()
}
