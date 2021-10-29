package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func stageDone(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()
		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stagesFn := make([]Stage, 0, len(stages))
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		stagesFn = append(stagesFn, stage)
	}

	if len(stagesFn) == 0 {
		return in
	}

	pipeline := in
	if done == nil {
		for _, stageFn := range stagesFn {
			pipeline = stageFn(pipeline)
		}
		return pipeline
	}

	pipeline = stageDone(pipeline, done)
	for _, stageFn := range stagesFn {
		pipeline = stageFn(stageDone(pipeline, done))
	}

	return pipeline
}
