package shared

// SpinnerStep represents a step of a multistep spinner.
type SpinnerStep struct {
	Callback func() error
	Name     string
}

// TestRunner represents a simple test that has a name and can be run via a callback.
type TestRunner struct {
	Name     string
	Callback func() error
}
