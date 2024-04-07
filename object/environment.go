package object

func NewEnv() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.store[name]
	if !ok && e.outer != nil {
		o, ok = e.outer.Get(name)
	}
	return o, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func NewScope(outer *Environment) *Environment {
	env := NewEnv()
	env.outer = outer
	return env
}
