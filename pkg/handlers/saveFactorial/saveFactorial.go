package savefactorial

//go:generate mockgen -source=saveFactorial.go -destination=../../mocks/saveFactorialMock.go

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/blockseeker999th/test-task-factorial/pkg/models"
	"github.com/blockseeker999th/test-task-factorial/pkg/utils"
	slogerr "github.com/blockseeker999th/test-task-factorial/pkg/utils/logger/slogErr"
	"github.com/julienschmidt/httprouter"
)

type Result struct {
	FactorialA int `json:"a"`
	FactorialB int `json:"b"`
}

type CalculationResult struct {
	ValueType string
	Result    int
}

type CalculationSaver interface {
	SaveCalculations(valueA int, valueB int) (int64, error)
}

func New(log *slog.Logger, calcSaver CalculationSaver) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		input := r.Context().Value("validatedInput").(models.Factorial)

		calcChan := make(chan CalculationResult, 2)

		var wg sync.WaitGroup
		wg.Add(2)

		calculate := func(n int, valueType string) {
			defer wg.Done()
			result := 1
			if n == 0 {
				calcChan <- CalculationResult{ValueType: valueType, Result: 1}
				return
			}

			for i := 2; i <= n; i++ {
				result *= i
			}
			calcChan <- CalculationResult{ValueType: valueType, Result: result}
		}

		go calculate(input.ValueA, "a")
		go calculate(input.ValueB, "b")

		go func() {
			wg.Wait()
			close(calcChan)
		}()

		result := Result{}
		for res := range calcChan {
			if res.ValueType == "a" {
				result.FactorialA = res.Result
			} else if res.ValueType == "b" {
				result.FactorialB = res.Result
			}
		}

		id, err := calcSaver.SaveCalculations(result.FactorialA, result.FactorialB)
		if err != nil {
			log.Error(utils.ErrSavingCalc, slogerr.Err(err))

			utils.WriteJSON(w, http.StatusInternalServerError, utils.NewErrorResponse(err.Error()))
			return
		}

		log.Info("calculations successfully saved", slog.Int64("id", id))

		utils.WriteJSON(w, http.StatusCreated, result)
	}
}
