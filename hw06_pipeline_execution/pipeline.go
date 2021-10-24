package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

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

	out := make(Bi)
	go func() {
		defer close(out)
		pipeline := stagesFn[0](in)

		for i := 1; i < len(stagesFn); i++ {
			pipeline = stagesFn[i](pipeline)
		}

		if done == nil {
			for value := range pipeline {
				out <- value
			}
			return
		}

		for {
			select {
			case value, ok := <-pipeline:
				if !ok {
					return
				}
				out <- value

			case <-done:
				return
			}
		}
	}()

	return out
}
