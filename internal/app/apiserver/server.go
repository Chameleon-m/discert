package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Chameleon-m/discert/internal/app/model"
	"github.com/Chameleon-m/discert/internal/app/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

const (
	sessionName          = "discert"
	ctxKeyAccount ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

func newServer(store store.Store, sessionStore sessions.Store) *server {
	srv := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	srv.configureRouter()

	return srv
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *server) configureRouter() {
	srv.router.Use(
		srv.setRequestID,
		srv.logRequest,
		handlers.CORS(handlers.AllowedOrigins([]string{"*"})),
	)
	srv.router.HandleFunc("/accounts", srv.handleAccountsCreate()).Methods("POST")
	srv.router.HandleFunc("/sessions", srv.handleSessionsCreate()).Methods("POST")

	// todo
	//srv.router.NotFoundHandler
	//srv.router.MethodNotAllowedHandler

	// /private/***
	private := srv.router.PathPrefix("/private").Subrouter()
	private.Use(srv.authenticateAccount)
	private.HandleFunc("/whoami", srv.handleWhoami()).Methods("GET")

}

// middleware
func (srv *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (srv *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := srv.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (srv *server) authenticateAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := srv.sessionStore.Get(r, sessionName)
		if err != nil {
			srv.errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["account_id"]
		if !ok {
			srv.errorRespond(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		account, err := srv.store.AccountRepository().Find(id.(int64))
		if err != nil {
			srv.errorRespond(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyAccount, account)))
	})
}

// routing handles

func (srv *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyAccount).(*model.Account))
	}
}

func (srv *server) handleAccountsCreate() http.HandlerFunc {
	type requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &requestData{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.errorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		account := &model.Account{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := srv.store.AccountRepository().Create(account); err != nil {
			srv.errorRespond(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		srv.respond(w, r, http.StatusCreated, account)
	}
}

func (srv *server) handleSessionsCreate() http.HandlerFunc {

	type requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &requestData{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.errorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		account, err := srv.store.AccountRepository().FindByEmail(req.Email)
		if err != nil || !account.ComparePassword(req.Password) {
			srv.errorRespond(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := srv.sessionStore.Get(r, sessionName)
		if err != nil {
			srv.errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["account_id"] = account.ID
		if err := srv.sessionStore.Save(r, w, session); err != nil {
			srv.errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		srv.respond(w, r, http.StatusOK, nil)
	}
}

func (srv *server) errorRespond(w http.ResponseWriter, r *http.Request, code int, err error) {
	srv.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (srv *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
