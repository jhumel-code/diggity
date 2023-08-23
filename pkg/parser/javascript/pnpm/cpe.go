package pnpm

import (
	"github.com/carbonetes/diggity/internal/cpe"
	"github.com/carbonetes/diggity/pkg/model"
)

func generateCPEs(pkg *model.Package) {
	cpe.NewCPE23(pkg, pkg.Name, pkg.Name, pkg.Version)
}