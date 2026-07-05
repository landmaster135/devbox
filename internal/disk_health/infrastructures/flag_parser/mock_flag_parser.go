package flag_parser

type MockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	ParseFunc    func() error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: map[string]string{},
		boolValues:   map[string]bool{},
	}
}

func (p *MockFlagParser) SetString(name string, value string) {
	p.stringValues[name] = value
}

func (p *MockFlagParser) SetBool(name string, value bool) {
	p.boolValues[name] = value
}

func (p *MockFlagParser) StringVar(value *string, name string, defaultValue string, usage string) {
	if presetValue, ok := p.stringValues[name]; ok {
		*value = presetValue
		return
	}
	if *value != "" {
		return
	}
	*value = defaultValue
}

func (p *MockFlagParser) BoolVar(value *bool, name string, defaultValue bool, usage string) {
	if presetValue, ok := p.boolValues[name]; ok {
		*value = presetValue
		return
	}
	if *value {
		return
	}
	*value = defaultValue
}

func (p *MockFlagParser) Parse() error {
	if p.ParseFunc == nil {
		return nil
	}
	return p.ParseFunc()
}
