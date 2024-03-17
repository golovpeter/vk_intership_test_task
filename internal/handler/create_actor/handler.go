package create_actor

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golovpeter/vk_intership_test_task/internal/common"
	"github.com/golovpeter/vk_intership_test_task/internal/service/actors"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger *logrus.Logger

	service actors.ActorService
}

func NewHandler(
	logger *logrus.Logger,
	service actors.ActorService,
) *handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) CreateActor(w http.ResponseWriter, r *http.Request) {
	var in *CreateActorIn

	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		h.logger.WithError(err).Error(common.BindJSONError)
		common.MakeErrorResponse(w, http.StatusBadRequest, common.BindJSONError)
		return
	}

	valid, errMsg := validateInParams(in)
	if !valid {
		h.logger.Error(errMsg)
		common.MakeErrorResponse(w, http.StatusBadRequest, errors.New(errMsg))
		return
	}

	err = h.service.CreateActor(r.Context(), &actors.ActorDataIn{
		Name:      in.Name,
		Gender:    in.Gender,
		BirthDate: in.BirthDate,
	})
	if err != nil {
		h.logger.WithError(err).Error(common.CreateActorError)
		common.MakeErrorResponse(w, http.StatusBadRequest, common.CreateActorError)
		return
	}
}
