package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler { return &Handler{service: service} }

// Rutas REST consistentes bajo /api/v1
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/users", h.ListUsers)
	rg.GET("/users/:id", h.getUserByID)
	rg.POST("/users", h.createUser)
	rg.PATCH("/users/:id/usm_pesos", h.updateUserUSMPesos) // más explícito
}

// ===== DTOs de request/response =====

// request para crear usuario
type createRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"  binding:"required"`
	Email     string `json:"email"      binding:"required,email"`
	Password  string `json:"password"   binding:"required"`
	USMPesos  int64  `json:"usm_pesos"` // opcional; por defecto 0
}

// response sin password
type userResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	USMPesos  int64  `json:"usm_pesos"`
}

func toUserResponse(u *Usuario) userResponse {
	return userResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		USMPesos:  u.USMPesos,
	}
}

// ===== Handlers =====

func (h *Handler) ListUsers(c *gin.Context) {
	us, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]userResponse, 0, len(us))
	for _, u := range us { out = append(out, toUserResponse(u)) }
	c.JSON(http.StatusOK, gin.H{"users": out})
}

func (h *Handler) getUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}); return }

	u, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	if u == nil   { c.JSON(http.StatusNotFound,      gin.H{"error": "User not found"}); return }
	c.JSON(http.StatusOK, toUserResponse(u))
}

func (h *Handler) createUser(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	u := &Usuario{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password, // TODO: hashear si quieres ir más pro
		USMPesos:  req.USMPesos,
	}
	id, err := h.service.RegisterUser(c.Request.Context(), u)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }

	// devuelve el creado (sin password)
	u.ID = id
	c.JSON(http.StatusCreated, toUserResponse(u))
}

func (h *Handler) updateUserUSMPesos(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}); return }

	var body struct {
		Amount int64 `json:"amount"` // usa int64 para calzar con service/repo
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.service.UpdateUserUSMPesos(c.Request.Context(), id, body.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Entrega el valor actualizado de los pesos
	u, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(u))

}
