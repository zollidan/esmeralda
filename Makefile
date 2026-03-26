.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go --output docs