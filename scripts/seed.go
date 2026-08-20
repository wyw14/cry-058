package main

import (
	"context"
	"fmt"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/repository/memory"
)

func main() {
	p := memory.NewPorts()
	e := p.Projects.Create(context.Background(), &domain.GrantProject{ID: "demo", Code: "DEMO-2026", Name: "社区关怀", Year: 2026, AnnualCapCents: 1000000})
	fmt.Println("seed", e)
}
