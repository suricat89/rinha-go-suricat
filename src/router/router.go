package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/suricat89/rinha-2024-q1/src/controller"
)

type Router struct {
	controller *controller.CustomerController
}

func NewRouter(controller *controller.CustomerController) *Router {
	return &Router{controller}
}

func (r *Router) Load(f *fiber.App) {
	f.Post("/clientes/:id/transacoes", r.controller.NewTransaction)
	f.Get("/clientes/:id/extrato", r.controller.GetTransactions)
}
