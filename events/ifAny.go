package events

// IfAny contains callbacks, that will be invoked when corresponding event
// is received by OnEvent fucn
type IfAny struct {
	Assert func(AssertDone)
	Log    func(LogEvent)
	Remote func(RemoteRequest)
	Var    func(SetVar)
	Sleep  func(Sleep)
	Start  func(Start)
	Finish func(Finish)
}

// OnEvent runs corresponding callback
func (i IfAny) OnEvent(e interface{}) {
	if e == nil {
		return
	}

	switch e.(type) {
	case AssertDone:
		if i.Assert != nil {
			i.Assert(e.(AssertDone))
		}
	case LogEvent:
		if i.Log != nil {
			i.Log(e.(LogEvent))
		}
	case RemoteRequest:
		if i.Remote != nil {
			i.Remote(e.(RemoteRequest))
		}
	case SetVar:
		if i.Var != nil {
			i.Var(e.(SetVar))
		}
	case Sleep:
		if i.Sleep != nil {
			i.Sleep(e.(Sleep))
		}
	case Start:
		if i.Start != nil {
			i.Start(e.(Start))
		}
	case Finish:
		if i.Finish != nil {
			i.Finish(e.(Finish))
		}
	}
}
