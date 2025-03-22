package quark

type Variable struct {
	name  string
	value Object
}

func NewVariable(name string, value interface{}) (*Variable, error) {
	object, err := FromInterface(value)
	if err != nil {
		return nil, err
	}
	return &Variable{
		name:  name,
		value: object,
	}, nil
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Value() Object {
	return v.value
}
