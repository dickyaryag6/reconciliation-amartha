package main

import "amartha-test/entities/http_handlers"

type module struct {
	httpHandler httphandlers.Handlers
}

func loadModules(handler httphandlers.Handlers) module {
	modules := module{
		httpHandler: handler,
	}

	return modules
}