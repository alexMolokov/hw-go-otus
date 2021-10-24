package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	if len(stages) == 0 {
		close(out)
		return out
	}

	go func() {
		defer close(out)
		pipeline := stages[0](in)

		for i := 1; i < len(stages); i++ {
			fn := stages[i]
			pipeline = fn(pipeline)
		}

		for {
			select {
			case value, ok := <-pipeline:
				if ok {
					out <- value
				}
				return

			case <-done:
				return
			}
		}
	}()

	return out
}
