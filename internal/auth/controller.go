package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"sharePie-api/internal/types"
	"strings"
	"time"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Controller struct {
	userService types.IUserService
}

func NewController(service types.IUserService) *Controller {
	return &Controller{userService: service}
}

// Signup handles the registration of new users.
// @Summary Register a new user
// @Description Registers a new user in the system
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param input body services.CreateUserInput true "User information"
// @Success 200 {object} map[string]interface{} "Returns the created user"
// @Failure 400 {object} map[string]interface{} "Returns an error if the input is invalid"
// @Failure 500 {object} map[string]interface{} "Returns an error if the user creation fails"
// @Router /auth/signup [post]
func (controller *Controller) Signup(c *gin.Context) {
	var input types.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = string(hash)
	user, err := controller.userService.Create(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Login handles user logins.
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param input body LoginInput true "Login credentials"
// @Success 200 {object} map[string]interface{} "Returns a JWT token"
// @Failure 400 {object} map[string]interface{} "Returns an error if the login credentials are invalid"
// @Router /auth/login [post]
func (controller *Controller) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Email = strings.ToLower(input.Email)

	user, err := controller.userService.FindOneByEmail(input.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare sent in password with saved users password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Validate checks the validity of the current user.
// @Summary Validate a user
// @Description Checks if the current user is valid in the system
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "Returns the current user"
// @Router /auth/validate [get]
func (controller *Controller) Validate(c *gin.Context) {
	user, _ := c.Get("user")

	//if !exists {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
