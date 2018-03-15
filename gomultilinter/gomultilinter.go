package main

import (
	"context"
	"fmt"

	"github.com/liut0/gomultilinter-errcheck/internal/errcheck"
	"github.com/liut0/gomultilinter/api"
)

type errcheckLinterFactory struct {
}

type errcheckLinterConfig struct {
	Blank   bool     `json:"blank"`
	Asserts bool     `json:"asserts"`
	Exclude []string `json:"exclude"`
}

type errcheckLinter struct {
	chkr *errcheck.Checker
}

var LinterFactory api.LinterFactory = &errcheckLinterFactory{}

func (l *errcheckLinterFactory) NewLinterConfig() api.LinterConfig {
	return &errcheckLinterConfig{}
}

func (cfg *errcheckLinterConfig) NewLinter() (api.Linter, error) {
	chkr := errcheck.NewChecker()
	chkr.Asserts = cfg.Asserts
	chkr.Blank = cfg.Blank

	excludes := map[string]bool{}
	for _, e := range cfg.Exclude {
		excludes[e] = true
	}
	chkr.SetExclude(excludes)

	return &errcheckLinter{
		chkr: chkr,
	}, nil
}

func (*errcheckLinter) Name() string {
	return "errcheck"
}

func (l *errcheckLinter) LintPackage(ctx context.Context, pkg *api.Package, reporter api.IssueReporter) error {
	for _, err := range l.chkr.CheckParsedPackage(pkg.PkgInfo, pkg.FSet) {
		reporter.Report(&api.Issue{
			Position: err.Pos,
			Category: "unchecked-errors",
			Message:  fmt.Sprintf("unchecked error (%v)", err.FuncName),
			Severity: api.SeverityWarning,
		})
	}

	return nil
}
