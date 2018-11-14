package scribe

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/apuckey/scribe-logger-go"
	"github.com/apuckey/scribe-logger-go/facebook/scribe"
	"net"
	"os"
)

type ScribeLogger struct {
	transport *thrift.TFramedTransport
	client    *scribe.ScribeClient
	category  string
	channel   chan *scribe.LogEntry
	formatter logger.Formatter
}

func NewScribeLogger(host, port, category string, formatter logger.Formatter, bufferSize int) (*ScribeLogger, error) {
	Ttransport, err := thrift.NewTSocket(net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	transport := thrift.NewTFramedTransport(Ttransport)

	protocol := thrift.NewTBinaryProtocol(transport, false, false)

	client := scribe.NewScribeClientProtocol(transport, protocol, protocol)
	if err := transport.Open(); err != nil {
		return nil, err
	}

	l := &ScribeLogger{
		transport: transport,
		client:    client,
		category:  category,
		channel:   make(chan *scribe.LogEntry, bufferSize),
		formatter: formatter,
	}

	go l.sendLoop()

	return l, nil
}

func (s *ScribeLogger) sendLoop() {

	defer func() {
		e := recover()
		if e != nil {
			fmt.Fprintf(os.Stderr, "Restarting sender go routine.")
			go s.sendLoop()
		}

	}()

	for msg := range s.channel {
		if msg != nil {
			//send to the server
			result, err := s.client.Log([]*scribe.LogEntry{msg})
			if err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("[ScribeError]: %s", err.Error()))
				fmt.Fprintf(os.Stderr, fmt.Sprintf("            -> %s", msg.Message))
			}
			if result != scribe.ResultCode_OK {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("[ScribeDown]: %s", msg.Message))
			}
		}
	}
}

func (s *ScribeLogger) sendOne(message string) {
	logEntry := &scribe.LogEntry{
		Category: s.category,
		Message:  message,
	}
	s.channel <- logEntry
}

func (s *ScribeLogger) sendArray(category string, messages []string) {
	for _, message := range messages {
		logEntry := &scribe.LogEntry{
			Category: s.category,
			Message:  message,
		}
		s.channel <- logEntry
	}
}

func (s *ScribeLogger) Close() error {
	return s.transport.Close()
}

func (s *ScribeLogger) SetFormatter(f logger.Formatter) {
	s.formatter = f
}

func (s *ScribeLogger) Emit(ctx *logger.MessageContext, message string) error {
	s.sendOne(s.formatter.Format(ctx, message))
	return nil
}
