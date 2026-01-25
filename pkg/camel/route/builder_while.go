package route

import (
	"github.com/paveldanilin/go-camel/pkg/camel/route/step"
)

type WhileBuilder struct {
	builder  *Builder
	loopStep *step.WhileLoop
}

func (wb *WhileBuilder) WithName(stepName string) *WhileBuilder {
	step.ApplyOptions(wb.loopStep, step.WithName(stepName))
	return wb
}

func (wb *WhileBuilder) WithCopyExchange(copyExchange bool) *WhileBuilder {
	wb.loopStep.CopyExchange = copyExchange
	return wb
}

func (wb *WhileBuilder) LoopBody(block StepBuilder) *Builder {
	if wb.builder.err != nil {
		return wb.builder
	}

	loopBody := step.Pipeline{Steps: []step.RouteStep{}, StoOnError: true}

	wb.builder.pushStack(&loopBody.Steps)
	block(wb.builder)
	wb.builder.popStack()

	wb.loopStep.Steps = append(wb.loopStep.Steps, &loopBody)

	return wb.builder
}
