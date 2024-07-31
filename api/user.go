package api

import (
	"errors"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
	"simplebank/worker"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

// createUserRequest defines the request body for creating a new user
// @Description Request body for creating a new user
// @Param username body string true "Username of the new user" example("johndoe")
// @Param password body string true "Password for the new user" example("password123")
// @Param role body string true "Role of the new user" example("user")
// @Param full_name body string true "Full name of the new user" example("John Doe")
// @Param email body string true "Email of the new user" example("johndoe@example.com")
// @Accept json
// @Produce json
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,role"` // binding:role is the custom Validator from validator.go
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// userResponse defines the response body for user details
// @Description Response body containing user details
// @Property username string "Username of the user" example("johndoe")
// @Property full_name string "Full name of the user" example("John Doe")
// @Property email string "Email of the user" example("johndoe@example.com")
// @Property password_changed_at string "Timestamp when the password was last changed" example("2024-07-31T12:00:00Z")
// @Property created_at string "Timestamp when the user was created" example("2024-07-31T12:00:00Z")
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// newUserResponse converts a database user to a userResponse
func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// createUser creates a new user in the system
// @Summary Create a new user
// @Description Create a new user with the provided details. A verification email will be sent after user creation.
// @Tags users
// @Accept json
// @Produce json
// @Param request body createUserRequest true "Create User Request"
// @Success 200 {object} userResponse "User created successfully"
// @Failure 400 {object} gin.H "Bad Request - Invalid input data"
// @Failure 403 {object} gin.H "Forbidden - User with the same username already exists"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	// Validating the request from the body JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.Username,
			HashedPassword: hashedPassword,
			FullName:       req.FullName,
			Role:           req.Role,
			Email:          req.Email,
		},
		// Create a Redis task
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second), // Necessary delay to wait for transaction to finish
				asynq.Queue(worker.QueueCritical), // Set the priority queue
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	// Use DB transaction to create a user and send email in a single transaction
	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create a response body without the hash password
	rsp := newUserResponse(txResult.User)
	ctx.JSON(http.StatusOK, rsp)
}

// loginUserRequest defines the request body for user login
// @Description Request body for user login
// @Param username body string true "Username of the user" example("johndoe")
// @Param password body string true "Password of the user" example("password123")
// @Accept json
// @Produce json
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

// loginUserResponse defines the response body for user login
// @Description Response body for user login including session and token details
// @Property session_id string "Session ID" example("d3f1e8e0-30ae-4d34-8b35-69115b43d8a9")
// @Property access_token string "Access Token for the user" example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqb2huZG9lIiwiaWF0IjoxNjA0Nzk1MzEyLCJleHBpcmN5X2NsaWVudF9pZCI6ImQzZjFlOGUwLTMwYWUtNGQzNC04YjM1LTY5MTU1YjQzZDhhOSIsImlhdCI6MTYwNDc5NTMxMn0.7L7zxAkD2I6kC60zkb_KQfc1Hw4JS4iJgb65hv5kAGk")
// @Property access_token_expires_at string "Expiration time of the access token" example("2024-07-31T12:00:00Z")
// @Property refresh_token string "Refresh Token for the user" example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqb2huZG9lIiwiaWF0IjoxNjA0Nzk1MzEyLCJleHBpcmN5X2NsaWVudF9pZCI6ImQzZjFlOGUwLTMwYWUtNGQzNC04YjM1LTY5MTU1YjQzZDhhOSIsImlhdCI6MTYwNDc5NTMxMn0.7L7zxAkD2I6kC60zkb_KQfc1Hw4JS4iJgb65hv5kAGk")
// @Property refresh_token_expires_at string "Expiration time of the refresh token" example("2024-08-30T12:00:00Z")
// @Property user userResponse "Details of the logged-in user"
// @Accept json
// @Produce json
type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

// loginUser handles user login and returns tokens and session details
// @Summary User Login
// @Description Authenticate the user and return access and refresh tokens, along with user details.
// @Tags users
// @Accept json
// @Produce json
// @Param request body loginUserRequest true "Login User Request"
// @Success 200 {object} loginUserResponse "Login successful, returns user details and tokens"
// @Failure 400 {object} gin.H "Bad Request - Invalid input data"
// @Failure 401 {object} gin.H "Unauthorized - Invalid credentials"
// @Failure 404 {object} gin.H "Not Found - User not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users/login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, user.Role, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create a session to store the refresh token
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     accessPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

// updateUserRequest defines the request body for updating user information
// @Description Request body for updating user information
// @Param username body string true "Username of the user to update" example("johndoe")
// @Param password body string true "New password for the user" example("newpassword123")
// @Param full_name body string true "New full name for the user" example("John Doe")
// @Param email body string true "New email address for the user" example("johndoe@example.com")
// @Accept json
// @Produce json
type updateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// updateUserResponse defines the response body for updating user information
// @Description Response body for updating user information
// @Property username string "Username of the user" example("johndoe")
// @Property full_name string "Full name of the user" example("John Doe")
// @Property email string "Email address of the user" example("johndoe@example.com")
// @Property role string "Role of the user" example("user")
// @Property password_changed_at string "Timestamp when the password was last changed" example("2024-07-31T12:00:00Z")
type updateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

// updateUser handles updating user details
// @Summary Update User Information
// @Description Update user details such as password, full name, and email. Only users with sufficient roles can update user information.
// @Tags users
// @Accept json
// @Produce json
// @Param request body updateUserRequest true "Update User Request"
// @Success 200 {object} updateUserResponse "User details updated successfully"
// @Failure 400 {object} gin.H "Bad Request - Invalid input data"
// @Failure 401 {object} gin.H "Unauthorized - Insufficient permissions or invalid credentials"
// @Failure 404 {object} gin.H "Not Found - User not found"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /users [put]
func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get auth payload and validate RBAC
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err := RBAC(ctx, authPayload.Role, []string{util.DepositorRole, util.BankerRole})
	if err != nil {
		return
	}

	if authPayload.Role != util.BankerRole && authPayload.Username != req.Username {
		err := errors.New("invalid user name")
		// Abort the API call and return 401 to the user
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		Username: req.Username,
		FullName: pgtype.Text{
			String: req.FullName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: req.Email,
			Valid:  true,
		},
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err == nil {
		// Passwords are the same, no need to update HashedPassword or PasswordChangedAt
		arg.HashedPassword = pgtype.Text{
			String: "",
			Valid:  false,
		}
		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Time{},
			Valid: false,
		}
	} else {
		// Passwords are different, update HashedPassword and PasswordChangedAt
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	result, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := updateUserResponse{
		Username:          result.Username,
		FullName:          result.FullName,
		Email:             result.Email,
		Role:              result.Role,
		PasswordChangedAt: result.PasswordChangedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
