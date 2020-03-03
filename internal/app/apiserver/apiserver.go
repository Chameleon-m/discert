package apiserver

import (
	"database/sql"
	"github.com/Chameleon-m/discert/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	//// todo path to config
	//sessionStore := sessions.NewFilesystemStore("", []byte(config.SessionKey))
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	//sessionStore.MaxAge(86400 * 30 * 12) // 360 days
	// todo server -> handler | srv := http.Server | srv.ListenAndServe()
	handler := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, handler)
}

func newDB(databaseDriver string, databaseURL string) (*sql.DB, error) {
	db, err := sql.Open(databaseDriver, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
