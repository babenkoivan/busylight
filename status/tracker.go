package status

type Transition struct {
	From Status
	To   Status
}

type Tracker struct {
	current Status
	C       chan Transition
}

func (t *Tracker) Record(status Status) {
	if t.current == status {
		return
	}

	t.C <- Transition{t.current, status}
	t.current = status
}

func NewTracker() *Tracker {
	return &Tracker{
		C: make(chan Transition),
	}
}
