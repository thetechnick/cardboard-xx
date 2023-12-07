package sh

import (
	"io"
	"log/slog"
)

type Option interface {
	ApplyToRunner(r *Runner)
}

type Logger struct{ *slog.Logger }

func (l Logger) ApplyToRunner(r *Runner) {
	r.logger = l.Logger
}

type Environment map[string]string

func (e Environment) ApplyToRunner(r *Runner) {
	r.env = e
}

type WorkDir string

func (wd WorkDir) ApplyToRunner(r *Runner) {
	r.workDir = string(wd)
}

type Stdout struct{ io.Writer }

func (stdout Stdout) ApplyToRunner(r *Runner) {
	r.stdout = stdout.Writer
}

type Stderr struct{ io.Writer }

func (stderr Stderr) ApplyToRunner(r *Runner) {
	r.stderr = stderr.Writer
}
