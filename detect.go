package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

type BuildPlanMetadata struct {
	Launch bool `toml:"launch"`
}

//go:generate faux --interface Parser --output fakes/parser.go
type Parser interface {
	Parse(path string) (hasMri bool, err error)
}

func Detect(gemfileParser Parser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		_, err := os.Stat(filepath.Join(context.WorkingDir, "config.ru"))
		if err != nil {
			if os.IsNotExist(err) {
				return packit.DetectResult{}, packit.Fail
			}
			return packit.DetectResult{}, fmt.Errorf("failed to stat config.ru: %w", err)
		}

		hasMri, err := gemfileParser.Parse(filepath.Join(context.WorkingDir, "Gemfile"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("failed to parse Gemfile: %w", err)
		}

		if !hasMri {
			return packit.DetectResult{}, packit.Fail
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "gems",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "bundler",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "mri",
						Metadata: BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			},
		}, nil
	}
}
