package app

import (
	"github.com/FRahimov84/Mux/pkg/mux"
	"github.com/FRahimov84/PurchaseService/pkg/core/purchase"
	"github.com/FRahimov84/myJwt/pkg/jwt"
	"github.com/FRahimov84/rest/pkg/rest"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"strconv"
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

func (s Server) handleProductList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.purchaseSvc.ProductList(s.pool)
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

func (s Server) handleProductByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		prod, err := s.purchaseSvc.ProductByID(int64(id), s.pool)
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

func (s Server) handleNewProduct() http.HandlerFunc {
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
		err = s.purchaseSvc.AddNewProduct(pur, s.pool)
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

func (s Server) handleDeleteProduct() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = s.purchaseSvc.RemoveByID(int64(id), s.pool)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		_, err = writer.Write([]byte("Purchase removed!"))
		if err != nil {
			log.Print(err)
		}
	}
}
