package section

import "time"

type Processor struct {
	WebServer ProcessorWebServer `split_words:"true"`
}

type ProcessorWebServer struct {
	ListenPort   uint32        `split_words:"true" default:"8080"`
	ReadTimeout  time.Duration `split_words:"true" default:"30s"`
	WriteTimeout time.Duration `split_words:"true" default:"30s"`
}
