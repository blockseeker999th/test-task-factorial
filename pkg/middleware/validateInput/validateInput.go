package validateinput

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/blockseeker999th/test-task-factorial/pkg/models"
	"github.com/blockseeker999th/test-task-factorial/pkg/utils"
	"github.com/blockseeker999th/test-task-factorial/pkg/utils/validation"
	"github.com/julienschmidt/httprouter"
)

func ValidateInput(log *slog.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var input models.Factorial
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Error(utils.ErrFailedToDecode)
			utils.WriteJSON(w, http.StatusBadRequest, utils.NewErrorResponse(utils.ErrFailedToDecode))
			return
		}

		if err := validation.ValidationStruct(input); err != nil {
			log.Error(utils.ErrValidationFailed)

			utils.WriteJSON(w, http.StatusBadRequest, utils.NewErrorResponse(utils.ErrValidationFailed))
			return
		}

		ctx := context.WithValue(r.Context(), "validatedInput", input)
		next(w, r.WithContext(ctx), p)
	}
}
