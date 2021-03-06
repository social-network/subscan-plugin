package main

import (
	"github.com/social-network/subscan-plugin/example/system/http"
	"github.com/social-network/subscan-plugin/example/system/model"
	"github.com/social-network/subscan-plugin/example/system/service"
	"github.com/social-network/subscan-plugin/router"
	"github.com/social-network/subscan-plugin/storage"
	"github.com/social-network/subscan-plugin/tools"
	"github.com/shopspring/decimal"
)

var srv *service.Service

type System struct {
	d storage.Dao
}

func New() *System {
	return &System{}
}

func (a *System) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *System) InitHttp() (routers []router.Http) {
	return http.Router(srv)
}

func (a *System) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, events []storage.Event) error {
	return nil
}

func (a *System) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	var paramEvent []storage.EventParam
	tools.UnmarshalToAnything(&paramEvent, event.Params)
	switch event.EventId {
	case "ExtrinsicFailed":
		srv.ExtrinsicFailed(block.SpecVersion, block.BlockTimestamp, block.Hash, event, paramEvent)
	}
	return nil
}

func (a *System) Migrate() {
	db := a.d.DB()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.ExtrinsicError{},
	)
	db.Model(model.ExtrinsicError{}).AddUniqueIndex("extrinsic_hash", "extrinsic_hash")
}

func (a *System) Version() string {
	return "0.1"
}
