package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

type LoopProcessor struct {
	// count - фиксированное число итераций (если >0).
	count int
	// predicate - условие для продолжения цикла (while predicate(m) == true).
	// Если nil, игнорируется.
	predicate func(message *camel.Message) bool
	// processors - процессоры, выполняемые в каждой итерации.
	processors []camel.Processor
	// copy - если true, каждая итерация работает на shallow копии Message.
	copy bool
}

func Loop(count int, predicate func(message *camel.Message) bool, processors ...camel.Processor) *LoopProcessor {

	return &LoopProcessor{
		count:      count,
		predicate:  predicate,
		processors: processors,
		copy:       true,
	}
}

func (p *LoopProcessor) Process(message *camel.Message) {
	if len(p.processors) == 0 {
		return // Нет процессоров - ничего не делаем.
	}

	iterations := 0
	for {
		// Проверяем условия выхода.
		if p.count > 0 && iterations >= p.count {
			break
		}
		if p.predicate != nil && !p.predicate(message) {
			break
		}

		// Если copy, создаём shallow копию.
		var current *camel.Message
		if p.copy {
			copy := *message // Shallow copy.
			current = &copy
		} else {
			current = message
		}

		// Выполняем процессоры в итерации.
		breakIteration := false
		for _, processor := range p.processors {
			if Invoke(processor, message) || current.Error != nil {
				breakIteration = true
				break
			}
		}

		// Если copy, копируем изменения обратно (включая Err).
		if p.copy {
			*message = *current
		}

		// Если ошибка/panic в итерации, прерываем весь цикл.
		if breakIteration || message.IsError() {
			break
		}

		iterations++
	}
}
