package events

import "context"

type PipelineHandler interface {
	Process(context.Context, *Event, HandlerFunc)
}

type PipelineHandlerFunc func(context.Context, *Event, HandlerFunc)

type Pipeline struct {
	first *delegatingHandler
	last  *delegatingHandler
}

func (p *Pipeline) Use(handler PipelineHandler) {
	p.UseFunc(handler.Process)
}

func (p *Pipeline) UseFunc(handler PipelineHandlerFunc) {
	dh := &delegatingHandler{inner: handler}
	if p.first == nil {
		p.first = dh
	}

	if p.last != nil {
		p.last.next = dh.Process
	}
	p.last = dh
}

func (p *Pipeline) Process(ctx context.Context, event *Event) {
	p.first.Process(ctx, event)
}

type delegatingHandler struct {
	inner PipelineHandlerFunc
	next  HandlerFunc
}

func (h *delegatingHandler) Process(ctx context.Context, event *Event) {
	h.inner(ctx, event, h.next)
}
