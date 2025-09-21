package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/adapter/controller"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/adapter/presenter"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port"
	ucase "github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/usecase"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/auth"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/datasource"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/logger"
)

type appDeps struct {
	ctrl port.UserController
	pres port.Presenter
	jwt  port.JWTSigner
}

var app appDeps

func build(ctx context.Context) (appDeps, error) {
	cfg := config.Load()
	log := logger.NewLogger(cfg.Environment)
	log.Info("api: building dependencies")
	repo, err := datasource.NewDynamoUserRepository(ctx, cfg)
	if err != nil {
		return appDeps{}, err
	}
	jwtSigner := auth.NewJWTSigner(cfg)
	uc := ucase.NewUserUseCase(repo, jwtSigner)
	ctrl := controller.NewUserController(uc)
	pres := presenter.NewJSONPresenter()
	return appDeps{ctrl: ctrl, pres: pres, jwt: jwtSigner}, nil
}

func respond(status int, payload any) (events.APIGatewayV2HTTPResponse, error) {
	b, _ := json.Marshal(payload)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
		Body: string(b),
	}, nil
}

func parseBody[T any](body string, v *T) error {
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func extractBearerToken(hdr string) string {
	parts := strings.SplitN(hdr, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func normalizePath(p string) string {
	if p == "" {
		return p
	}
	p = strings.ToLower(p)
	if strings.HasSuffix(p, "/") && len(p) > 1 {
		p = strings.TrimRight(p, "/")
	}
	return p
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if app.ctrl == nil {
		deps, err := build(ctx)
		if err != nil {
			return respond(500, map[string]string{"error": "internal error"})
		}
		app = deps
	}

	method := strings.ToUpper(req.RequestContext.HTTP.Method)
	path := req.RawPath
	if path == "" {
		path = req.RequestContext.HTTP.Path
	}
	path = normalizePath(path)

	switch {
	case method == "POST" && path == "/prod/users/register":
		var in dto.RegisterInput
		if err := parseBody(req.Body, &in); err != nil {
			return respond(400, map[string]string{"error": "invalid body"})
		}
		b, err := app.ctrl.Register(ctx, app.pres, in)
		if err != nil {
			status := 400
			if errors.Is(err, ucase.ErrEmailAlreadyExists) {
				status = 409
			}
			return respond(status, map[string]string{"error": err.Error()})
		}
		var out any
		_ = json.Unmarshal(b, &out)
		return respond(201, out)

	case method == "POST" && path == "/prod/users/login":
		var in dto.LoginInput
		if err := parseBody(req.Body, &in); err != nil {
			return respond(400, map[string]string{"error": "invalid body"})
		}
		b, err := app.ctrl.Login(ctx, app.pres, in)
		if err != nil {
			status := 400
			if errors.Is(err, ucase.ErrInvalidCredentials) || errors.Is(err, ucase.ErrInvalidInput) {
				status = 401
			}
			return respond(status, map[string]string{"error": err.Error()})
		}
		var out any
		_ = json.Unmarshal(b, &out)
		return respond(200, out)

	case method == "GET" && path == "/prod/users/me":
		auth := req.Headers["authorization"]
		tok := extractBearerToken(auth)
		if tok == "" {
			return respond(401, map[string]string{"error": "missing bearer token"})
		}
		userID, err := app.jwt.Verify(tok)
		if err != nil {
			return respond(401, map[string]string{"error": "invalid token"})
		}
		b, err := app.ctrl.GetMe(ctx, app.pres, userID)
		if err != nil {
			status := 400
			if errors.Is(err, ucase.ErrUserNotFound) {
				status = 404
			}
			return respond(status, map[string]string{"error": err.Error()})
		}
		var out any
		_ = json.Unmarshal(b, &out)
		return respond(200, out)
	}

	log.Printf("DEBUG: No route matched - method=%s, path=%s", method, path)
	return respond(404, map[string]string{"error": "not found"})
}

func main() {
	lambda.Start(handler)
}
