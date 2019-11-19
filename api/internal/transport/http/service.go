package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/koind/banner-rotation/api/internal/domain/repository"
	"github.com/koind/banner-rotation/api/internal/domain/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// HTTP rotation service
type RotationService struct {
	service.RotationService
	logger *zap.Logger
}

// Will return new http rotation service
func NewHTTPRotationService(rotation service.RotationService, logger *zap.Logger) *RotationService {
	return &RotationService{
		RotationService: rotation,
		logger:          logger,
	}
}

// Adds a banner in the rotation
func (s *RotationService) AddBannerHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	rotation := repository.Rotation{}

	err := decoder.Decode(&rotation)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		return
	}

	rotation.SetDatetimeOfCreate()
	newRotation, err := s.Add(r.Context(), rotation)
	if err != nil {
		s.logger.Error(
			"An error occurred while adding a banner to the rotation",
			zap.Error(err),
		)

		w.Write([]byte(err.Error()))
	} else {
		s.logger.Info(
			"Banner added to rotation",
			zap.Any("rotation", newRotation),
		)

		json.NewEncoder(w).Encode(newRotation)
	}
}

// Sets the transition on the banner
func (s *RotationService) SetTransitionHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var rotationForm struct {
		BannerID int `json:"bannerId"`
		GroupID  int `json:"groupId"`
	}

	err := decoder.Decode(&rotationForm)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		return
	}

	rotation, err := s.RotationRepository.FindOneByBannerID(r.Context(), rotationForm.BannerID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		return
	}

	err = s.SetTransition(r.Context(), *rotation, rotationForm.GroupID)
	if err != nil {
		s.logger.Error(
			"Error when set the transition on the banner",
			zap.Error(err),
		)

		w.Write([]byte(err.Error()))
	} else {
		s.logger.Info(
			"Was set the transition on the banner",
			zap.Any("bannerID", rotationForm.BannerID),
			zap.Any("groupID", rotationForm.GroupID),
		)

		w.Write([]byte("ok"))
	}
}

// Selects a banner to display
func (s *RotationService) SelectBannerHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var rotationForm struct {
		SlotID  int `json:"slotId"`
		GroupID int `json:"groupId"`
	}

	err := decoder.Decode(&rotationForm)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))

		return
	}

	bannerID, err := s.SelectBanner(r.Context(), rotationForm.SlotID, rotationForm.GroupID)
	if err != nil {
		s.logger.Error(
			"Error when select banner",
			zap.Error(err),
		)

		w.Write([]byte(err.Error()))
	} else {
		s.logger.Info(
			"Was selected the banner to view",
			zap.Any("slotID", rotationForm.SlotID),
			zap.Any("groupID", rotationForm.GroupID),
			zap.Any("bannerID", bannerID),
		)

		json.NewEncoder(w).Encode(bannerID)
	}
}

// Removes the banner from the rotation
func (s *RotationService) RemoveBannerHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if id, ok := vars["id"]; ok {
		bannerID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))

			return
		}

		err = s.Remove(r.Context(), bannerID)
		if err != nil {
			s.logger.Error(
				"Error removing banner from rotation",
				zap.Error(err),
			)

			w.Write([]byte(err.Error()))
		} else {
			s.logger.Info(
				"The banner has been removed from rotation",
				zap.Int("bannerID", bannerID),
			)

			w.Write([]byte("ok"))
		}
	} else {
		w.WriteHeader(400)
		w.Write([]byte("Banner id not found"))
	}
}
