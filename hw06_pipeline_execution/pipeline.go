package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func createProxyChannel(in In, done In) Out {
	proxyChannel := make(Bi)
	go func() {
		defer close(proxyChannel)
		for {
			select {
			case <-done:
				return
			case task, ok := <-in:
				if !ok {
					return
				}
				if task != nil {
					select {
					case <-done:
						return
					case proxyChannel <- task:
					}
				}
			}
		}
	}()
	return proxyChannel
}

func stageWrapper(in In, done In, stage Stage) Out {
	if done != nil {
		proxyChannel := createProxyChannel(in, done)
		return stage(proxyChannel)
	}
	return stage(in)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stageWrapper(out, done, stage)
	}
	return createProxyChannel(out, done)
}
