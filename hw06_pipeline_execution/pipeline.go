package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func StageWrapper(in In, done In, stage Stage) Out {
	if done != nil {
		ch := make(Bi)
		go func() {
			defer close(ch)
			for {
				select {
				case <-done:
					return
				default:
				}

				select {
				case <-done:
					return
				case task, ok := <-in:
					if !ok {
						return
					}
					ch <- task
				}
			}
		}()
		return stage(ch)
	}
	return stage(in)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = StageWrapper(out, done, stage)
	}
	return out
}
