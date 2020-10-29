package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// chanWithDoneCheck перекладывает данные из одного канала в другой с проверкой на необходимость завершения.
func chanWithDoneCheck(done In, in In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inStream := in
	for _, stage := range stages {
		inStream = stage(chanWithDoneCheck(done, inStream)) // выходной поток одного, становится входным другого
	}
	return inStream // результат ExecutePipeline будет читаться из выходного стрима последнего стейджа
}
