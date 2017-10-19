package tagc

func Port(port int) func(*Tagc) {
	return func(t *Tagc) {
		t.port = port
	}
}
