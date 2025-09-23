package timings

import (
	"dickobrazz/application/logging"
	"time"
)

func ReportExecutionForResultError[T any, E error](log *logging.Logger, action func() (T, E), report func(l *logging.Logger)) (T, E) {
	start := time.Now()
	result, err := action()
	report(log.With(logging.ExecutionTime, time.Since(start).String()))
	return result, err
}

func ReportExecutionForResult[T any](log *logging.Logger, action func() T, report func(l *logging.Logger)) T {
	start, result := time.Now(), action()
	report(log.With(logging.ExecutionTime, time.Since(start).String()))
	return result
}

func ReportExecution(log *logging.Logger, action func(), report func(l *logging.Logger)) {
	start := time.Now()
	action()
	report(log.With(logging.ExecutionTime, time.Since(start).String()))
}
