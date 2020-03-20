package app

import (
	"github.com/FRahimov84/Mux/pkg/mux"
	"github.com/FRahimov84/PurchaseService/pkg/core/purchase"
	"github.com/FRahimov84/PurchaseService/pkg/core/token"
	jwt2 "github.com/FRahimov84/PurchaseService/pkg/mux/middleware/jwt"
	"github.com/FRahimov84/myJwt/pkg/jwt"
	"github.com/FRahimov84/rest/pkg/rest"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

type Server struct {
	router      *mux.ExactMux
	pool        *pgxpool.Pool
	purchaseSvc *purchase.Service
	secret      jwt.Secret
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, productSvc *purchase.Service, secret jwt.Secret) *Server {
	return &Server{router: router, pool: pool, purchaseSvc: productSvc, secret: secret}
}

func (s Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s Server) Start() {
	s.InitRoutes()
}

func (s Server) handlePurchaseList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.purchaseSvc.PurchaseList(s.pool)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &list)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handlePurchaseByUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload := jwt2.FromContext(request.Context()).(*token.Payload)
		prod, err := s.purchaseSvc.PurchaseByUserID(payload.Id, s.pool)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &prod)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handleNewPurchase() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		pur := purchase.Purchase{}
		err := rest.ReadJSONBody(request, &pur)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		payload := jwt2.FromContext(request.Context()).(*token.Payload)
		pur.User_id = payload.Id
		err = s.purchaseSvc.AddNewPurchase(pur, s.pool)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("New Purchase Added!"))
		if err != nil {
			log.Print(err)
		}
	}
}
