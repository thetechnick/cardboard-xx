package run

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var dr = NewDependencyRun()

// Executes dependencies one after the other.
func SerialDeps(ctx context.Context, deps ...Dep) error {
	return dr.Serial(ctx, deps...)
}

// Executes dependencies in parallel.
func ParallelDeps(ctx context.Context, deps ...Dep) error {
	return dr.Parallel(ctx, deps...)
}

// Represents a dependency.
type Dep interface {
	// Unique Identifier to ensure this dependency only executes once.
	ID() string
	// Executes the dependency.
	Run(ctx context.Context) error
}

// Returns a new DependencyRun context.
func NewDependencyRun() *depRun {
	return &depRun{
		ran: map[string]*depOnce{},
	}
}

type depRun struct {
	// remembers functions that have already been executed.
	ran map[string]*depOnce
	mux sync.RWMutex
}

func (r *depRun) Reset() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.ran = map[string]*depOnce{}
}

// Executes dependencies in parallel.
func (r *depRun) Parallel(ctx context.Context, deps ...Dep) error {
	var (
		wg      sync.WaitGroup
		errs    []error
		errsMux sync.Mutex
	)
	wg.Add(len(deps))
	for _, dep := range deps {
		dep := r.get(dep)
		go func() {
			defer wg.Done()
			if err := dep.Run(ctx); err != nil {
				errsMux.Lock()
				errs = append(errs, fmt.Errorf("running %s: %w", dep.ID(), err))
				errsMux.Unlock()
			}
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}

// Executes dependencies one after the other.
func (r *depRun) Serial(ctx context.Context, deps ...Dep) error {
	for _, dep := range deps {
		dep := r.get(dep)
		if err := dep.Run(ctx); err != nil {
			return fmt.Errorf("running %s: %w", dep.ID(), err)
		}
	}
	return nil
}

func (r *depRun) get(dep Dep) Dep {
	r.mux.RLock()
	defer r.mux.RUnlock()
	out, ok := r.ran[dep.ID()]
	if !ok {
		out = newOnce(dep)
		r.ran[dep.ID()] = out
	}
	return out
}

type dep struct {
	id  string
	run func(ctx context.Context) error
}

func (d *dep) ID() string {
	return d.id
}

func (d *dep) Run(ctx context.Context) error {
	return d.run(ctx)
}

type fn interface {
	func() | func() error | func(ctx context.Context) | func(ctx context.Context) error
}

// Wraps a function with no parameters for the dependency handler.
func Fn[T fn](fn T) Dep {
	return &dep{
		id: funcID(fn),
		run: func(ctx context.Context) error {
			switch v := any(fn).(type) {
			case func():
				v()
			case func() error:
				return v()
			case func(context.Context):
				v(ctx)
			case func(context.Context) error:
				return v(ctx)
			}
			return nil
		},
	}
}

type fn1[A any] interface {
	func(A) | func(A) error | func(context.Context, A) | func(context.Context, A) error
}

// Wraps a function with one parameter for the dependency handler.
func Fn1[T fn1[A], A any](fn T, a1 A) Dep {
	return &dep{
		id: funcID(fn, a1),
		run: func(ctx context.Context) error {
			switch v := any(fn).(type) {
			case func(A):
				v(a1)
			case func(A) error:
				return v(a1)
			case func(context.Context, A):
				v(ctx, a1)
			case func(context.Context, A) error:
				return v(ctx, a1)
			}
			return nil
		},
	}
}

type fn2[A, B any] interface {
	func(A, B) | func(A, B) error | func(context.Context, A, B) | func(context.Context, A, B) error
}

// Wraps a function with two parameters for the dependency handler.
func Fn2[T fn2[A, B], A, B any](fn T, a1 A, a2 B) Dep {
	return &dep{
		id: funcID(fn, a1, a2),
		run: func(ctx context.Context) error {
			switch v := any(fn).(type) {
			case func(A, B):
				v(a1, a2)
			case func(A, B) error:
				return v(a1, a2)
			case func(context.Context, A, B):
				v(ctx, a1, a2)
			case func(context.Context, A, B) error:
				return v(ctx, a1, a2)
			}
			return nil
		},
	}
}

type fn3[A, B, C any] interface {
	func(A, B, C) | func(A, B, C) error | func(context.Context, A, B, C) | func(context.Context, A, B, C) error
}

// Wraps a function with three parameters for the dependency handler.
func Fn3[T fn3[A, B, C], A, B, C any](fn T, a1 A, a2 B, a3 C) Dep {
	return &dep{
		id: funcID(fn, a1, a2, a3),
		run: func(ctx context.Context) error {
			switch v := any(fn).(type) {
			case func(A, B, C):
				v(a1, a2, a3)
			case func(A, B, C) error:
				return v(a1, a2, a3)
			case func(context.Context, A, B, C):
				v(ctx, a1, a2, a3)
			case func(context.Context, A, B, C) error:
				return v(ctx, a1, a2, a3)
			}
			return nil
		},
	}
}

type fn4[A, B, C, D any] interface {
	func(A, B, C, D) | func(A, B, C, D) error | func(context.Context, A, B, C, D) | func(context.Context, A, B, C, D) error
}

// Wraps a function with four parameters for the dependency handler.
func Fn4[T fn4[A, B, C, D], A, B, C, D any](fn T, a1 A, a2 B, a3 C, a4 D) Dep {
	return &dep{
		id: funcID(fn, a1, a2, a3, a4),
		run: func(ctx context.Context) error {
			switch v := any(fn).(type) {
			case func(A, B, C, D):
				v(a1, a2, a3, a4)
			case func(A, B, C, D) error:
				return v(a1, a2, a3, a4)
			case func(context.Context, A, B, C, D):
				v(ctx, a1, a2, a3, a4)
			case func(context.Context, A, B, C, D) error:
				return v(ctx, a1, a2, a3, a4)
			}
			return nil
		},
	}
}

type fn5[A, B, C, D, E any] interface {
	func(A, B, C, D, E) | func(A, B, C, D, E) error | func(context.Context, A, B, C, D, E) | func(context.Context, A, B, C, D, E) error
}

// Wraps a function with five parameters for the dependency handler.
func Fn5[T fn5[A, B, C, D, E], A, B, C, D, E any](fn T, a1 A, a2 B, a3 C, a4 D, a5 E) Dep {
	return &dep{
		id: funcID(fn, a1, a2, a3, a4, a5),
		run: func(ctx context.Context) error {
			switch v := any(fn).(type) {
			case func(A, B, C, D, E):
				v(a1, a2, a3, a4, a5)
			case func(A, B, C, D, E) error:
				return v(a1, a2, a3, a4, a5)
			case func(context.Context, A, B, C, D, E):
				v(ctx, a1, a2, a3, a4, a5)
			case func(context.Context, A, B, C, D, E) error:
				return v(ctx, a1, a2, a3, a4, a5)
			}
			return nil
		},
	}
}

type fn6[A, B, C, D, E, F any] interface {
	func(A, B, C, D, E, F) | func(A, B, C, D, E, F) error | func(context.Context, A, B, C, D, E, F) | func(context.Context, A, B, C, D, E, F) error
}

// Wraps a function with six parameters for the dependency handler.
func Fn6[T fn6[A, B, C, D, E, F], A, B, C, D, E, F any](fn T, a1 A, a2 B, a3 C, a4 D, a5 E, a6 F) Dep {
	return &dep{
		id: funcID(fn, a1, a2, a3, a4, a5, a6),
		run: func(ctx context.Context) error {
			switch v := any(fn).(type) {
			case func(A, B, C, D, E, F):
				v(a1, a2, a3, a4, a5, a6)
			case func(A, B, C, D, E, F) error:
				return v(a1, a2, a3, a4, a5, a6)
			case func(context.Context, A, B, C, D, E, F):
				v(ctx, a1, a2, a3, a4, a5, a6)
			case func(context.Context, A, B, C, D, E, F) error:
				return v(ctx, a1, a2, a3, a4, a5, a6)
			}
			return nil
		},
	}
}

// returns a string that can be used to identify the given function and arguments.
func funcID(fn any, args ...any) string {
	fnV := reflect.ValueOf(fn)
	fnR := runtime.FuncForPC(fnV.Pointer())
	name := strings.TrimSuffix(fnR.Name(), "-fm")
	if len(args) > 0 {
		argStrings := make([]string, len(args))
		for i, arg := range args {
			argStrings[i] = fmt.Sprintf("%#v", arg)
		}
		return fmt.Sprintf("%s(%s)", name, strings.Join(argStrings, ", "))
	}
	return name + "()"
}

// container type to ensure a dependency only runs once.
type depOnce struct {
	once *sync.Once
	dep  Dep
	err  error
}

func newOnce(dep Dep) *depOnce {
	return &depOnce{
		once: &sync.Once{},
		dep:  dep,
	}
}

func (o *depOnce) ID() string {
	return o.dep.ID()
}

func (o *depOnce) Run(ctx context.Context) error {
	o.once.Do(func() {
		o.err = o.dep.Run(ctx)
	})
	return o.err
}
