package onec

type step struct {
	name    string
	handler func() error
}

type steps struct {
	steps []step
	log   logWriter
	state map[string]any
}

func (s *steps) add(name string, handler func() error) {
	s.steps = append(s.steps, step{
		name:    name,
		handler: handler,
	})
}

func (s *steps) run() error {
	for i, item := range s.steps {
		s.log("log", item.name)
		if err := item.handler(); err != nil {
			return err
		}

		s.log("progress", byte((float32(i)/float32(len(s.steps)))*100))
	}

	return nil
}
