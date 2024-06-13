package savefactorial

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	validateinput "github.com/blockseeker999th/test-task-factorial/pkg/middleware/validateInput"
	mock_saveHandle "github.com/blockseeker999th/test-task-factorial/pkg/mocks"
	"github.com/blockseeker999th/test-task-factorial/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestHandlerSaveFactorial(t *testing.T) {
	type mockBehavior func(h *mock_saveHandle.MockCalculationSaver, factorial *models.Factorial)

	testTable := []struct {
		name                 string
		inputBody            string
		inputFactorial       models.Factorial
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "save factorial",
			inputBody: `{"a":5,"b":3}`,
			inputFactorial: models.Factorial{
				ValueA: 5,
				ValueB: 5,
			},
			mockBehavior: func(h *mock_saveHandle.MockCalculationSaver, factorial *models.Factorial) {
				h.EXPECT().SaveCalculations(gomock.AssignableToTypeOf(0), gomock.AssignableToTypeOf(0)).Return(int64(1), nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"a":120,"b":6}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			result := mock_saveHandle.NewMockCalculationSaver(ctrl)
			testCase.mockBehavior(result, &testCase.inputFactorial)

			log := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)

			handler := New(log, result)

			r := httprouter.New()
			r.POST("/calculate", validateinput.ValidateInput(log, handler))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)

			if w.Code == http.StatusCreated {
				var response models.Factorial
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResponseBody, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
