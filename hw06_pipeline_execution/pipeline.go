package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stagesLen := len(stages)

	if stagesLen == 0 {
		return nil
	}

	resOut := make(Bi)

	pipeline := in
	for _, stage := range stages {
		pipeline = stage(pipeline)
	}

	go func() {
		defer close(resOut)

		for {
			select {
			case <-done:
				return
			case val, ok := <-pipeline:
				if !ok {
					return
				}
				resOut <- val
			}
		}
	}()

	return resOut
}
