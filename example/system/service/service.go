package service

import (
	"github.com/social-network/subscan-plugin/example/system/dao"
	"github.com/social-network/subscan-plugin/example/system/model"
	"github.com/social-network/subscan-plugin/storage"
	"github.com/social-network/subscan-plugin/tools"
	"github.com/social-network/substrate-api-rpc"
)

type Service struct {
	d storage.Dao
}

func New(d storage.Dao) *Service {
	return &Service{
		d: d,
	}
}

func (s *Service) GetExtrinsicError(hash string) *model.ExtrinsicError {
	return dao.ExtrinsicError(s.d.DB(), hash)
}

func (s *Service) ExtrinsicFailed(spec, blockTimestamp int, blockHash string, event *storage.Event, paramEvent []storage.EventParam) {

	type DispatchErrorModule struct {
		Index int `json:"index"`
		Error int `json:"error"`
	}

	for _, param := range paramEvent {

		if param.Type == "DispatchError" {

			var dr map[string]interface{}
			tools.UnmarshalToAnything(&dr, param.Value)

			if _, ok := dr["Error"]; ok {
				_ = dao.CreateExtrinsicError(s.d.DB(),
					event.ExtrinsicHash,
					dao.CheckExtrinsicError(spec, s.d.SpecialMetadata(spec),
						tools.IntFromInterface(dr["Module"]),
						tools.IntFromInterface(dr["Error"])))

			} else if _, ok := dr["Module"]; ok {
				var module DispatchErrorModule
				tools.UnmarshalToAnything(&module, dr["Module"])

				_ = dao.CreateExtrinsicError(s.d.DB(),
					event.ExtrinsicHash,
					dao.CheckExtrinsicError(spec, s.d.SpecialMetadata(spec), module.Index, module.Error))

			} else if _, ok := dr["BadOrigin"]; ok {
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "BadOrigin"})

			} else if _, ok := dr["CannotLookup"]; ok {
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "CannotLookup"})

			} else if _, ok := dr["Other"]; ok {
				_ = dao.CreateExtrinsicError(s.d.DB(), event.ExtrinsicHash, &substrate.MetadataModuleError{Name: "Other"})

			}
			break
		}
	}
}
