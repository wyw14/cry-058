package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wyw14/cry058/internal/application"
	"github.com/wyw14/cry058/internal/domain"
	"github.com/wyw14/cry058/internal/middleware"
	"net/http"
)

type Router struct {
	p           *domain.Ports
	projects    *application.ProjectService
	claims      *application.ClaimService
	rules       *application.RuleService
	settlements *application.SettlementService
}

func NewRouter(p *domain.Ports) *Router {
	return &Router{p: p, projects: application.NewProjectService(p), claims: application.NewClaimService(p), rules: application.NewRuleService(p), settlements: application.NewSettlementService(p)}
}
func (r *Router) Handler() *gin.Engine {
	g := gin.New()
	g.Use(gin.Recovery(), middleware.RequestID(), middleware.Secure())
	g.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	g.GET("/readyz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ready"}) })
	v := g.Group("/api/v1")
	v.GET("/projects/:id", r.getProject)
	v.POST("/claims", r.submitClaim)
	v.POST("/rules", r.publishRule)
	v.POST("/settlements/preview/:claim", r.preview)
	return g
}
func (r *Router) getProject(c *gin.Context) {
	v, e := r.projects.Get(c, c.Param("id"))
	if e != nil {
		c.JSON(404, gin.H{"error": e.Error()})
		return
	}
	c.JSON(200, v)
}
func (r *Router) submitClaim(c *gin.Context) {
	var v domain.ExpenseClaim
	if e := c.ShouldBindJSON(&v); e != nil {
		c.JSON(400, gin.H{"error": e.Error()})
		return
	}
	x, e := r.claims.Submit(c, v)
	if e != nil {
		c.JSON(422, gin.H{"error": e.Error()})
		return
	}
	c.JSON(201, x)
}
func (r *Router) publishRule(c *gin.Context) {
	var v domain.RuleVersion
	if e := c.ShouldBindJSON(&v); e != nil {
		c.JSON(400, gin.H{"error": e.Error()})
		return
	}
	x, e := r.rules.Publish(c, v)
	if e != nil {
		c.JSON(422, gin.H{"error": e.Error()})
		return
	}
	c.JSON(201, x)
}
func (r *Router) preview(c *gin.Context) {
	x, w, e := r.settlements.Preview(c, c.Param("claim"))
	if e != nil {
		c.JSON(422, gin.H{"error": e.Error()})
		return
	}
	c.JSON(200, gin.H{"settlement": x, "warnings": w})
}
