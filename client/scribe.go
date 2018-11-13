package scribe

import (
	"errors"
	"net"

	"../facebook/scribe"
	"git.apache.org/thrift.git/lib/go/thrift"
)

type ScribeLogger struct {
	transport *thrift.TFramedTransport
	client    *scribe.ScribeClient
}

func NewScribeLogger(host, port string) (*ScribeLogger, error) {
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
	return &ScribeLogger{
		transport: transport,
		client:    client,
	}, nil
}

func (s *ScribeLogger) SendOne(category, message string) (bool, error) {
	logEntry := &scribe.LogEntry{
		Category: category,
		Message: message,
	}
	result, err := s.client.Log([]*scribe.LogEntry{logEntry})
	if err != nil {
		return false, err
	}
	return s.dealResult(result)
}

func (s *ScribeLogger) SendArray(category string, messages []string) (bool, error) {
	var logEntrys []*scribe.LogEntry

	for _, message := range messages {
		logEntry := &scribe.LogEntry{
			Category: category,
			Message: message,
		}
		logEntrys = append(logEntrys, logEntry)
	}
	result, err := s.client.Log(logEntrys)
	if err != nil {
		return false, err
	}
	return s.dealResult(result)
}

func (s *ScribeLogger) dealResult(result scribe.ResultCode) (bool, error) {
	ok := false
	var err error
	switch result {
	case scribe.ResultCode_OK:
		ok = true
	case scribe.ResultCode_TRY_LATER:
		ok = false
	default:
		err = errors.New(result.String())
	}
	return ok, err
}

func (s *ScribeLogger) Close() error {
	return s.transport.Close()
}