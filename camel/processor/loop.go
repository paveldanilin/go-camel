package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

type LoopProcessor struct {
	// count - фиксированное число итераций (если >0).
	count int
	// predicate - условие для продолжения цикла (while predicate(m) == true).
	// Если nil, игнорируется.
	predicate func(exchange *camel.Exchange) bool
	// processors - процессоры, выполняемые в каждой итерации.
	processors []camel.Processor
	// copy - если true, каждая итерация работает на shallow копии Message.
	copy bool
}

func Loop(count int, predicate func(exchange *camel.Exchange) bool, processors ...camel.Processor) *LoopProcessor {

	return &LoopProcessor{
		count:      count,
		predicate:  predicate,
		processors: processors,
		copy:       true,
	}
}

func (p *LoopProcessor) Process(exchange *camel.Exchange) {
	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	if len(p.processors) == 0 {
		return // Нет процессоров - ничего не делаем.
	}

	iterations := 0
	for {
		// Проверяем условия выхода.
		if p.count > 0 && iterations >= p.count {
			break
		}
		if p.predicate != nil && !p.predicate(exchange) {
			break
		}

		// Если copy, создаём shallow копию.
		var current *camel.Exchange
		if p.copy {
			copy := *exchange // Shallow copy.
			current = &copy
		} else {
			current = exchange
		}

		// Выполняем процессоры в итерации.
		breakIteration := false
		for _, processor := range p.processors {
			if InvokeWithRecovery(processor, exchange) || current.Error != nil {
				breakIteration = true
				break
			}
		}

		// Если copy, копируем изменения обратно (включая Err).
		if p.copy {
			*exchange = *current
		}

		// Если ошибка/panic в итерации, прерываем весь цикл.
		if breakIteration || exchange.IsError() {
			break
		}

		iterations++
	}
}
