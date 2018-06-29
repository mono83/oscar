package events

// EventRouter contains callbacks, that will be invoked when corresponding event
// is received by OnEvent func
type EventRouter struct {
	Assert          func(AssertDone, *Emitted)
	Log             func(LogEvent, *Emitted)
	Remote          func(RemoteRequest, *Emitted)
	Var             func(SetVar, *Emitted)
	Sleep           func(Sleep, *Emitted)
	Start           func(Start, *Emitted)
	Finish          func(Finish, *Emitted)
	RegistrationIn  func(RegistrationBegin, *Emitted)
	RegistrationOut func(RegistrationEnd, *Emitted)
}

// OnEvent runs corresponding callback
func (i EventRouter) OnEvent(o *Emitted) {
	if o == nil || o.Data == nil {
		return
	}

	e := o.Data

	switch e.(type) {
	case AssertDone:
		if i.Assert != nil {
			i.Assert(e.(AssertDone), o)
		}
	case LogEvent:
		if i.Log != nil {
			i.Log(e.(LogEvent), o)
		}
	case RemoteRequest:
		if i.Remote != nil {
			i.Remote(e.(RemoteRequest), o)
		}
	case SetVar:
		if i.Var != nil {
			i.Var(e.(SetVar), o)
		}
	case Sleep:
		if i.Sleep != nil {
			i.Sleep(e.(Sleep), o)
		}
	case Start:
		if i.Start != nil {
			i.Start(e.(Start), o)
		}
	case Finish:
		if i.Finish != nil {
			i.Finish(e.(Finish), o)
		}
	case RegistrationBegin:
		if i.RegistrationIn != nil {
			i.RegistrationIn(e.(RegistrationBegin), o)
		}
	case RegistrationEnd:
		if i.RegistrationOut != nil {
			i.RegistrationOut(e.(RegistrationEnd), o)
		}
	}
}
