package application

import "time"

func TraceTimeExecutionForResultError[T any, E error](log *Logger, kind string, action func() (T, E)) (T, E) {
	start := time.Now()
	result, err := action()
	log.With(ExecutionTime, time.Since(start), ExecutionKing, kind).I("Operation execution finished")
	return result, err
}

func TraceTimeExecutionForResult[T any](log *Logger, kind string, action func() T) T {
	start, result := time.Now(), action()
	log.With(ExecutionTime, time.Since(start), ExecutionKing, kind).I("Operation execution finished")
	return result
}

func TraceTimeExecution(log *Logger, kind string, action func()) {
	start := time.Now()
	action()
	log.With(ExecutionTime, time.Since(start), ExecutionKing, kind).I("Operation execution finished")
}

const (
	TraceKindGatherCollection  = "gather.collection"
	TraceKindCreatePipeline    = "create.pipeline"
	TraceKindAggregatePipeline = "aggregate.pipeline"
	TraceKindInflatePipeline   = "inflate.pipeline"
)
