// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package scribe

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/apuckey/scribe-logger-go/facebook/fb303"
	"math"
)

// (needed to ensure safety because of naive import list construction.)
var _ = math.MinInt32
var _ = thrift.ZERO
var _ = fmt.Printf

var _ = fb303.GoUnusedProtection__

type Scribe interface {
	fb303.FacebookService

	// Parameters:
	//  - Messages
	Log(messages []*LogEntry) (r ResultCode, err error)
}

type ScribeClient struct {
	*fb303.FacebookServiceClient
}

func NewScribeClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ScribeClient {
	return &ScribeClient{FacebookServiceClient: fb303.NewFacebookServiceClientFactory(t, f)}
}

func NewScribeClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ScribeClient {
	return &ScribeClient{FacebookServiceClient: fb303.NewFacebookServiceClientProtocol(t, iprot, oprot)}
}

// Parameters:
//  - Messages
func (p *ScribeClient) Log(messages []*LogEntry) (r ResultCode, err error) {
	if err = p.sendLog(messages); err != nil {
		return
	}
	return p.recvLog()
}

func (p *ScribeClient) sendLog(messages []*LogEntry) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("Log", thrift.CALL, p.SeqId)
	args0 := NewLogArgs()
	args0.Messages = messages
	err = args0.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ScribeClient) recvLog() (value ResultCode, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error2 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error3 error
		error3, err = error2.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error3
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result1 := NewLogResult()
	err = result1.Read(iprot)
	iprot.ReadMessageEnd()
	value = result1.Success
	return
}

type ScribeProcessor struct {
	*fb303.FacebookServiceProcessor
}

func NewScribeProcessor(handler Scribe) *ScribeProcessor {
	self4 := &ScribeProcessor{fb303.NewFacebookServiceProcessor(handler)}
	self4.AddToProcessorMap("Log", &scribeProcessorLog{handler: handler})
	return self4
}

type scribeProcessorLog struct {
	handler Scribe
}

func (p *scribeProcessorLog) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewLogArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("Log", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewLogResult()
	if result.Success, err = p.handler.Log(args.Messages); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing Log: "+err.Error())
		oprot.WriteMessageBegin("Log", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("Log", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

type LogArgs struct {
	Messages []*LogEntry `thrift:"messages,1"`
}

func NewLogArgs() *LogArgs {
	return &LogArgs{}
}

func (p *LogArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *LogArgs) readField1(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return fmt.Errorf("error reading list being: %s")
	}
	p.Messages = make([]*LogEntry, 0, size)
	for i := 0; i < size; i++ {
		_elem5 := NewLogEntry()
		if err := _elem5.Read(iprot); err != nil {
			return fmt.Errorf("%T error reading struct: %s", _elem5)
		}
		p.Messages = append(p.Messages, _elem5)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return fmt.Errorf("error reading list end: %s")
	}
	return nil
}

func (p *LogArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Log_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *LogArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if p.Messages != nil {
		if err := oprot.WriteFieldBegin("messages", thrift.LIST, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:messages: %s", p, err)
		}
		if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Messages)); err != nil {
			return fmt.Errorf("error writing list begin: %s")
		}
		for _, v := range p.Messages {
			if err := v.Write(oprot); err != nil {
				return fmt.Errorf("%T error writing struct: %s", v)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return fmt.Errorf("error writing list end: %s")
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:messages: %s", p, err)
		}
	}
	return err
}

func (p *LogArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LogArgs(%+v)", *p)
}

type LogResult struct {
	Success ResultCode `thrift:"success,0"`
}

func NewLogResult() *LogResult {
	return &LogResult{
		Success: math.MinInt32 - 1, // unset sentinal value
	}
}

func (p *LogResult) IsSetSuccess() bool {
	return int64(p.Success) != math.MinInt32-1
}

func (p *LogResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *LogResult) readField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 0: %s")
	} else {
		p.Success = ResultCode(v)
	}
	return nil
}

func (p *LogResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Log_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	default:
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *LogResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.I32, 0); err != nil {
			return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
		}
		if err := oprot.WriteI32(int32(p.Success)); err != nil {
			return fmt.Errorf("%T.success (0) field write error: %s", p)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 0:success: %s", p, err)
		}
	}
	return err
}

func (p *LogResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LogResult(%+v)", *p)
}
