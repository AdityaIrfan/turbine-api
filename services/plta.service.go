package services

import (
	"errors"
	"net/http"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"
)

type pltaService struct {
	pltaRepo contract.IPltaRepository
	userRepo contract.IUserRepository
}

func NewPltaService(pltaRepo contract.IPltaRepository, userRepo contract.IUserRepository) contract.IPltaService {
	return &pltaService{
		pltaRepo: pltaRepo,
		userRepo: userRepo,
	}
}

func (t *pltaService) Create(c echo.Context, in *models.PltaCreateRequest) error {
	pltaByName, err := t.pltaRepo.GetByEqualNameWithSelectedFields(in.Name, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if !pltaByName.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "nama plta sudah digunakan")
	}

	plta := in.ToModelCreate()
	if err := t.pltaRepo.Create(plta); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	// get created by
	user, err := t.userRepo.GetByIdWithSelectedFields(plta.CreatedBy, "name")
	if err != nil {
		log.Error().Err(errors.New("FAILED TO GET USER AFTER CREATING PLTA : " + err.Error())).Msg("")
	} else {
		plta.CreatedByUser = user
	}

	return helpers.Response(c, http.StatusOK, "berhasil membuat plta baru", plta.ToResponse())
}

func (t *pltaService) Update(c echo.Context, in *models.PltaUpdateRequest) error {
	plta, err := t.pltaRepo.GetByIdWithSelectedFields(in.Id, "id")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if plta.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "plta tidak ditemukan")
	}

	if in.Name != plta.Name {
		pltaByName, err := t.pltaRepo.GetByEqualNameWithSelectedFields(in.Name, "id")
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if !pltaByName.IsEmpty() && pltaByName.Id != in.Id {
			return helpers.Response(c, http.StatusBadRequest, "nama plta sudah digunakan")
		}
	}

	plta = in.ToModelUpdate()
	if err := t.pltaRepo.Update(plta); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	// get created by
	userCreated, err := t.userRepo.GetByIdWithSelectedFields(plta.CreatedBy, "name")
	if err != nil {
		log.Error().Err(errors.New("FAILED TO GET USER AFTER UPDATING PLTA : " + err.Error())).Msg("")
	} else {
		plta.CreatedByUser = userCreated
	}

	// get updated by
	if plta.CreatedBy != plta.UpdatedBy {
		userUpdated, err := t.userRepo.GetByIdWithSelectedFields(plta.UpdatedBy, "name")
		if err != nil {
			log.Error().Err(errors.New("FAILED TO GET USER AFTER UPDATING PLTA : " + err.Error())).Msg("")
		} else {
			plta.UpdatedByUser = userUpdated
		}
	} else {
		plta.UpdatedByUser = userCreated
	}

	return helpers.Response(c, http.StatusOK, "berhasil mengubah plta", plta.ToResponse())
}

func (t *pltaService) Detail(c echo.Context, id string) error {
	plta, err := t.pltaRepo.GetByIdWithPreloads(id, "PltaUnits")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if plta.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "plta tidak ditemukan")
	}

	// get created by
	userCreated, err := t.userRepo.GetByIdWithSelectedFields(plta.CreatedBy, "name")
	if err != nil {
		log.Error().Err(errors.New("FAILED TO GET USER BY ON GETTING DETAIL PLTA : " + err.Error())).Msg("")
	} else {
		plta.CreatedByUser = userCreated
	}

	// get updated by
	if plta.UpdatedAt != nil {
		if plta.CreatedBy != plta.UpdatedBy {
			userUpdated, err := t.userRepo.GetByIdWithSelectedFields(plta.UpdatedBy, "name")
			if err != nil {
				log.Error().Err(errors.New("FAILED TO GET USER ON GETTING DETAIL PLTA : " + err.Error())).Msg("")
			} else {
				plta.UpdatedByUser = userUpdated
			}
		} else {
			plta.UpdatedByUser = userCreated
		}

	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan plta", plta.ToResponse())
}

func (t *pltaService) GetListMaster(c echo.Context, in *models.PltaGetListMasterRequest) error {
	user, err := t.userRepo.GetByIdWithSelectedFields(in.UserId, "id, location_radius_status")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() {
		return helpers.ResponseForbiddenAccess(c)
	}

	listPlta, err := t.pltaRepo.GetAll(in.Search)
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var pltaResponse []*models.PltaResponseMaster
	for _, plta := range listPlta {
		pltaResponse = append(pltaResponse, plta.ToResponseMaster(user.RadiusStatus)...)
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua plta master", pltaResponse)
}

func (t *pltaService) Delete(c echo.Context, in *models.PltaDeleteRequest) error {
	plta, err := t.pltaRepo.GetByIdWithSelectedFields(in.Id, "*")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	} else if plta.IsEmpty() {
		return helpers.Response(c, http.StatusBadRequest, "plta tidak ditemukan")
	}

	plta.DeletedBy = in.DeletedBy
	if err := t.pltaRepo.Delete(plta); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	return helpers.Response(c, http.StatusOK, "berhasil menghapus plta")
}

func (t *pltaService) GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error {
	listPlta, pagination, err := t.pltaRepo.GetListWithPaginate(cursor, "*")
	if err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	var pltaRes = []*models.PltaListResponse{}
	var usersMap = make(map[string]string)
	for _, plta := range listPlta {
		// get created by
		if value, ok := usersMap[plta.CreatedBy]; !ok {
			userCreated, err := t.userRepo.GetByIdWithSelectedFields(plta.CreatedBy, "name")
			if err != nil {
				log.Error().Err(errors.New("FAILED TO GET USER ON GETTING LIST PLTA WITH PAGINATE : " + err.Error())).Msg("")
				usersMap[plta.CreatedBy] = ""
			} else {
				plta.CreatedByUser = userCreated
				usersMap[plta.CreatedBy] = userCreated.Name
			}
		} else {
			plta.CreatedByUser = &models.User{
				Name: value,
			}
		}

		pltaRes = append(pltaRes, plta.ToResponseList())
	}

	return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua plta", pltaRes, pagination)
}
