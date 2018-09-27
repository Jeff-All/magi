package user

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

// func Put(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) error {
// 	body, err := util.IOUtil.ReadAll(r.Body)
// 	if err != nil {
// 		return errors.CodedError{
// 			Message:  "Unable to read the body",
// 			HTTPCode: http.StatusInternalServerError,
// 			Err:      err,
// 		}
// 	}
// 	var requests []*models.Request
// 	if err = util.Json.Unmarshal(body, &requests); err != nil {
// 		return errors.CodedError{
// 			Message:  "Unable to unmarshal request",
// 			HTTPCode: http.StatusInternalServerError,
// 			Err:      err,
// 		}
// 	}

// 	type ErrorResponse struct {
// 		Error string
// 		Sheet string
// 		Row   int
// 	}
// 	response := make([]interface{}, len(requests))
// 	for i, val := range requests {
// 		if err = validate.Struct(validation.Request(*val)); err != nil {
// 			log.WithFields(log.Fields{
// 				"error": err.Error(),
// 				"sheet": val.Sheet,
// 				"row":   val.Row,
// 			}).Error("Error validating request")
// 			response[i] = ErrorResponse{
// 				Error: err.Error(),
// 				Sheet: val.Sheet,
// 				Row:   val.Row,
// 			}
// 		} else if err = request.Create(val); err != nil {
// 			log.WithFields(log.Fields{
// 				"error": err.Error(),
// 				"sheet": val.Sheet,
// 				"row":   val.Row,
// 			}).Error("Error create request")
// 			response[i] = ErrorResponse{
// 				Error: err.Error(),
// 				Sheet: val.Sheet,
// 				Row:   val.Row,
// 			}
// 		} else {
// 			response[i] = val
// 		}
// 	}
// 	requestJSONString, err := util.Json.Marshal(response)
// 	if err != nil {
// 		return errors.CodedError{
// 			Message:  "Unable to marshal response",
// 			HTTPCode: http.StatusInternalServerError,
// 			Err:      err,
// 		}
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	if _, err = w.Write(requestJSONString); err != nil {
// 		return errors.CodedError{
// 			Message:  "Error while writing the response",
// 			HTTPCode: http.StatusInternalServerError,
// 			Err:      err,
// 		}
// 	}
// 	w.WriteHeader(http.StatusCreated)
// 	return nil
// }

// func GetPage(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) error {
// 	params := r.URL.Query()
// 	var limit int
// 	var err error
// 	if limit, err = strconv.Atoi(params.Get("limit")); err != nil || limit == 0 {
// 		limit = 20
// 	}
// 	var offset int
// 	if offset, err = strconv.Atoi(params.Get("offset")); err != nil {
// 		offset = 0
// 	}
// 	users := make([]models.User, 0, limit)
// 	if err := resources.DB.Limit(limit).Offset(offset).Find(&users).GetError(); err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return errors.CodedError{
// 				Message:  "endpoints.user.GetPage(): can't find users",
// 				HTTPCode: 404,
// 				Fields: logrus.Fields{
// 					"limit":  limit,
// 					"offset": offset,
// 				},
// 				Err: err,
// 			}
// 		}
// 		return errors.CodedError{
// 			Message:  "endpoints.user.GetPage(): Error querying users",
// 			HTTPCode: 500,
// 			Fields: logrus.Fields{
// 				"limit":  limit,
// 				"offset": offset,
// 			},
// 			Err: err,
// 		}
// 	}
// 	jsonapiRuntime := jsonapi.NewRuntime().Instrument("users.put")
// 	if err = jsonapiRuntime.MarshalPayload(w, []api.User(users)); err != nil {
// 		return errors.CodedError{
// 			Message:  "endpoints.user.GetPage(): Error while marshaling response",
// 			HTTPCode: 500,
// 			Err:      err,
// 		}
// 	}
// 	w.WriteHeader(http.StatusFound)
// 	return nil
// }
