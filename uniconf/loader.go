package uniconf

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/nikmy/algo/reflex"
)

type Loader[T any] interface {
	Load(ctx context.Context) (*T, error)
}

type OverrideStrategy int

const (
	AcceptFirst = OverrideStrategy(iota)
	AcceptLast
)

func PipelineWithLogger[T any](strategy OverrideStrategy, logger infoLogger, sources ...Loader[T]) Loader[T] {
	return pipelineLoader[T]{
		accept: strategy,
		logger: logger,
		stages: sources,
	}
}

func Pipeline[T any](strategy OverrideStrategy, sources ...Loader[T]) Loader[T] {
	return pipelineLoader[T]{
		accept: strategy,
		stages: sources,
	}
}

type pipelineLoader[T any] struct {
	accept OverrideStrategy
	logger infoLogger
	stages []Loader[T]
}

func (l pipelineLoader[T]) Load(ctx context.Context) (*T, error) {
	var (
		loaded T
		errs   []error
	)

	for _, source := range l.stages {
		updated, err := source.Load(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		switch l.accept {
		case AcceptFirst:
			l.override(&loaded, updated)
		case AcceptLast:
			l.override(updated, &loaded)
			loaded = *updated
		}
	}

	if len(errs) == len(l.stages) {
		return nil, errors.Join(errs...)
	}

	if l.logger != nil {
		for _, err := range errs {
			l.logger.Info(fmt.Sprintf("skip stage due to error: %s", err.Error()))
		}
	}

	return &loaded, nil
}

func (pipelineLoader[T]) override(cfg *T, overrider *T) {
	reflex.Override(reflect.ValueOf(cfg), reflect.ValueOf(overrider))
}
