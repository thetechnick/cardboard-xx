package run

import "context"

type MyTestType struct{}

func (m *MyTestType) Test1(a string)                {}
func (m *MyTestType) Test2(a, b string)             {}
func (m *MyTestType) Test3(a string, b int, c bool) {}
func (m *MyTestType) Test4(a, b, c, d string)       {}
func (m *MyTestType) Test5(a, b, c, d, e string)    {}
func (m *MyTestType) Test6(a, b, c, d, e, f string) {}

type MyTestTypeErr struct{}

func (m *MyTestTypeErr) Test1(a string) error {
	return nil
}
func (m *MyTestTypeErr) Test2(a, b string) error {
	return nil
}
func (m *MyTestTypeErr) Test3(a string, b int, c bool) error {
	return nil
}
func (m *MyTestTypeErr) Test4(a, b, c, d string) error {
	return nil
}
func (m *MyTestTypeErr) Test5(a, b, c, d, e string) error {
	return nil
}
func (m *MyTestTypeErr) Test6(a, b, c, d, e, f string) error {
	return nil
}

type MyTestTypeCtxErr struct{}

func (m *MyTestTypeCtxErr) Test1(ctx context.Context, a string) error {
	return nil
}
func (m *MyTestTypeCtxErr) Test2(ctx context.Context, a, b string) error {
	return nil
}
func (m *MyTestTypeCtxErr) Test3(ctx context.Context, a string, b int, c bool) error {
	return nil
}
func (m *MyTestTypeCtxErr) Test4(ctx context.Context, a, b, c, d string) error {
	return nil
}
func (m *MyTestTypeCtxErr) Test5(ctx context.Context, a, b, c, d, e string) error {
	return nil
}
func (m *MyTestTypeCtxErr) Test6(ctx context.Context, a, b, c, d, e, f string) error {
	return nil
}

type MyTestTypeCtx struct{}

func (m *MyTestTypeCtx) Test1(ctx context.Context, a string)                {}
func (m *MyTestTypeCtx) Test2(ctx context.Context, a, b string)             {}
func (m *MyTestTypeCtx) Test3(ctx context.Context, a string, b int, c bool) {}
func (m *MyTestTypeCtx) Test4(ctx context.Context, a, b, c, d string)       {}
func (m *MyTestTypeCtx) Test5(ctx context.Context, a, b, c, d, e string)    {}
func (m *MyTestTypeCtx) Test6(ctx context.Context, a, b, c, d, e, f string) {}

func myFunc() {}
func myFuncErr() error {
	return nil
}
func myFuncCtx(ctx context.Context) {}
func myFuncCtxErr(ctx context.Context) error {
	return nil
}
