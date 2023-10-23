package service

import (
	"fmt"
	"net/http"

	"github.com/SAMBA-Research/microservice-shared/tracing"
	"github.com/kjasuquo/srv-SAMBA/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

	"github.com/kjasuquo/srv-SAMBA/internal/config"
	srvdb "github.com/kjasuquo/srv-SAMBA/internal/db"
	"github.com/kjasuquo/srv-SAMBA/version"
)

type Microservice struct {
	cfg *config.Config
	dbm *srvdb.DbConnectionManager
}

func NewMicroservice(cfg *config.Config, dbm *srvdb.DbConnectionManager) (srv *Microservice, err error) {
	srv = &Microservice{
		cfg: cfg,
		dbm: dbm,
	}
	return
}

func (srv *Microservice) Run() {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.GET("/", srv.indexHandler)
	router.GET("/sum", srv.sumHandler)

	log.Info().Msgf("Microservice %s listening on %s:%d", version.ServiceName, srv.cfg.ServiceBind, srv.cfg.ServicePort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", srv.cfg.ServiceBind, srv.cfg.ServicePort), router)
	if err != nil {
		panic(err)
	}
}

func (srv *Microservice) indexHandler(w http.ResponseWriter, r bunrouter.Request) (err error) {
	_, span := tracing.Tracer().Start(r.Context(), "service.indexHandler")
	defer span.End()

	w.Write([]byte("This is an API server"))
	return
}

func (srv *Microservice) sumHandler(w http.ResponseWriter, r bunrouter.Request) (err error) {
	_, span := tracing.Tracer().Start(r.Context(), "service.sumHandler")
	defer span.End()

	amount, err := utils.Sum(srv.cfg.SumURL)
	if err != nil {
		return utils.ReturnJSON(w, err.Error(), http.StatusInternalServerError)
	}

	return utils.ReturnJSON(w, map[string]any{
		"ok":  true,
		"sum": amount,
	}, http.StatusOK)

}
