package state

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/buyoio/runr/pkg/scm"
)

type GoTemplate struct {
	*Runner
	scm.Provider
	Args []string
}
type GoTemplateOptions struct {
	Runner      *Runner
	SCMPlatform *SCMPlatform
}

func (d *Data) Content() (string, error) {
	content, err := d.File.Content()
	if err != nil {
		return "", err
	}
	if d.tmpl == nil {
		return content, nil
	}

	args, err := d.Args.Parse()
	if err != nil {
		return "", err
	}
	d.tmpl.Args = args

	t, err := template.New("").Funcs(sprig.FuncMap()).Parse(content)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := t.Execute(&b, d.tmpl); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (a Args) Parse() ([]string, error) {
	var args []string
	for _, arg := range a {
		if arg.Value != nil {
			args = append(args, *arg.Value)
			continue
		}
		if arg.Exec != nil {
			out, err := arg.Exec.Output(arg.Envs)
			if err != nil {
				return nil, err
			}
			args = append(args, strings.TrimSpace(string(out)))
			continue
		}
		if arg.File != nil {
			out, err := arg.File.Content()
			if err != nil {
				return nil, err
			}
			args = append(args, strings.TrimSpace(out))
			continue
		}
	}
	return args, nil
}

func (d *Data) SetLabels(labels []string) {
	if d.tmpl != nil {
		d.tmpl.Provider.SetLabels(labels)
	}
}

func (d *Data) setGoTemplate(ctx context.Context, logger *slog.Logger, options *GoTemplateOptions) error {
	d.tmpl = &GoTemplate{}
	if options.SCMPlatform != nil {
		p, err := scm.NewProvider(ctx, logger, &scm.ProviderOptions{
			SCMPlatform: scm.SCMPlatform(*options.SCMPlatform),
		})
		if err != nil {
			return err
		}
		d.tmpl.Provider = p
	}
	if options.Runner != nil {
		d.tmpl.Runner = options.Runner
	}
	return nil
}
