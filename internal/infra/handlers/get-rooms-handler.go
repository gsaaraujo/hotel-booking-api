package handlers

import (
	"context"

	"github.com/google/uuid"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type GetRoomsHandlerOutput struct {
	Id       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	Number   string    `json:"number"`
	Capacity uint8     `json:"capacity"`
	Price    uint64    `json:"price"`
}

type GetRoomsHandler struct {
	Conn              *pgx.Conn
	HttpLogger        webhttp.HttpLogger
	HttpAuthorization webhttp.HttpAuthorization
}

func (g *GetRoomsHandler) Handle(c echo.Context) error {
	authorizationToken := c.Request().Header.Get("Authorization")

	if authorizationToken == "" {
		return webhttp.NewUnauthorized(c, "missing or invalid authorization token")
	}

	if !g.HttpAuthorization.IsCustomer(authorizationToken) {
		return webhttp.NewForbidden(c, "you do not have permission to access this resource")
	}

	type RoomSchema struct {
		Id       uuid.UUID
		Type     string
		Number   string
		Capacity uint8
		Price    uint64
	}

	rows, err := g.Conn.Query(context.Background(), "SELECT id, number, type, capacity, price FROM rooms")

	if err != nil {
		g.HttpLogger.Log(c, err)
		return webhttp.NewInternalServerError(c)
	}

	var roomsSchema []RoomSchema
	for rows.Next() {
		var roomSchema RoomSchema
		err := rows.Scan(&roomSchema.Id, &roomSchema.Number, &roomSchema.Type, &roomSchema.Capacity, &roomSchema.Price)

		if err != nil {
			g.HttpLogger.Log(c, err)
			return webhttp.NewInternalServerError(c)
		}

		roomsSchema = append(roomsSchema, roomSchema)
	}

	if len(roomsSchema) == 0 {
		return webhttp.NewOk(c, []GetRoomsHandlerOutput{})
	}

	var getRoomsHandlerOutput []GetRoomsHandlerOutput
	for _, roomSchema := range roomsSchema {
		getRoomsHandlerOutput = append(getRoomsHandlerOutput, GetRoomsHandlerOutput(roomSchema))
	}

	return webhttp.NewOk(c, getRoomsHandlerOutput)
}
