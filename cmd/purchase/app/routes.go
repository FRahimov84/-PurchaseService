package app

import (
	"context"
	"errors"
	"github.com/FRahimov84/PurchaseService/pkg/core/token"
	"github.com/FRahimov84/PurchaseService/pkg/mux/middleware/authenticated"
	"github.com/FRahimov84/PurchaseService/pkg/mux/middleware/authorized"
	"github.com/FRahimov84/PurchaseService/pkg/mux/middleware/jwt"
	"github.com/FRahimov84/PurchaseService/pkg/mux/middleware/logger"
	"reflect"
)

func (s Server) InitRoutes() {

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		panic(errors.New("can't create database"))
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `
create table if not exists purchase (
    id BIGSERIAL primary key,
    user_id integer not null,
    product_id integer not null,
    price integer not null check ( price>=0 ),
    kol integer not null check ( kol>0 ),
    purchase_date date default current_date,
    removed BOOLEAN DEFAULT FALSE
);
`)
	if err != nil {
		panic(errors.New("can't create database"))
	}

	s.router.GET(
		"/api/purchase",
		s.handlePurchaseList(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get list"),
	)

	s.router.GET(
		"/api/purchase/me",
		s.handlePurchaseByUser(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get purchase by id"),
	)

	s.router.POST(
		"/api/purchase/new",
		s.handleNewPurchase(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("post new purchase"),
	)


}