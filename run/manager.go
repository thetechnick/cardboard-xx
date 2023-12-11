package run

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Manager struct {
	targets map[string]func(ctx context.Context, args []string) error
}

func New() *Manager {
	return &Manager{
		targets: map[string]func(ctx context.Context, args []string) error{},
	}
}

func (m *Manager) Register(things ...any) error {
	for _, thing := range things {
		if err := m.register(thing); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Run(things ...any) {
	if err := m.Register(things...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	m.run()
}

func (m *Manager) run() {
	ctx := context.Background()
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "unexpected number of arguments")
		os.Exit(1)
	}
	if err := m.Call(ctx, args[1], args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (m *Manager) Call(ctx context.Context, id string, args []string) error {
	fn, ok := m.targets[id]
	if ok {
		return fn(ctx, args)
	}
	return nil
}

func (m *Manager) register(thing any) error {
	thingType := reflect.TypeOf(thing)
	thingValue := reflect.ValueOf(thing)
	typeID := strings.ToLower(thingType.Elem().Name())
	for i := 0; i < thingType.NumMethod(); i++ {
		method := thingType.Method(i)
		if !method.IsExported() {
			continue
		}

		// check params
		if method.Type.NumIn() != 3 && method.Type.NumOut() != 1 ||
			!(method.Type.In(1).String() == "context.Context") ||
			!(method.Type.In(2).String() == "[]string") ||
			!(method.Type.Out(0).String() == "error") {
			return fmt.Errorf(
				"%s.%s() must have signature like func(context.Context, []string) error",
				thingType.Elem().Name(), method.Name)
		}

		methodID := strings.ToLower(method.Name)
		m.targets[typeID+":"+methodID] = func(ctx context.Context, args []string) error {
			out := thingValue.MethodByName(method.Name).Call([]reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(args),
			})
			errI := out[0].Interface()
			if errI == nil {
				return nil
			}
			return errI.(error)
		}
	}
	return nil
}
