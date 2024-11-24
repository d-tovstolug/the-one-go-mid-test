package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/d.tovstoluh/the-one-go-mid-test/rest-api/model"
	mockstorage "github.com/d.tovstoluh/the-one-go-mid-test/rest-api/storage/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testNewTask = model.Task{
		Name:   "test-name-1",
		Status: 2,
	}
	testExistingTask = model.Task{
		ID:     "test-id-1",
		Name:   "test-name-1",
		Status: 2,
	}
	testExistingTaskJSON      = `{"id":"test-id-1","name":"test-name-1","status":2}`
	testErrInvalidRequestJSON = `{"error":"invalid request"}`
)

func prepare(ctrl *gomock.Controller) (storageMock *mockstorage.MockTaskStorage, router *gin.Engine) {
	storageMock = mockstorage.NewMockTaskStorage(ctrl)

	r := gin.Default()
	NewTaskController(storageMock).Register(r)
	return storageMock, r
}

func TestTaskController_SaveTask(t *testing.T) {
	type args struct {
		request *model.Task
	}
	tests := []struct {
		name           string
		args           args
		expectedCode   int
		expectedResult string
		mockFunc       func(mock *mockstorage.MockTaskStorage, args args)
	}{
		{
			name: "positive create",
			args: args{
				request: &testNewTask,
			},
			expectedCode:   http.StatusOK,
			expectedResult: testExistingTaskJSON,
			mockFunc: func(mock *mockstorage.MockTaskStorage, args args) {
				mock.EXPECT().Save(gomock.Any(), args.request).Return(&testExistingTask, nil)
			},
		},
		{
			name: "positive update",
			args: args{
				request: &testExistingTask,
			},
			expectedCode:   http.StatusOK,
			expectedResult: testExistingTaskJSON,
			mockFunc: func(mock *mockstorage.MockTaskStorage, args args) {
				mock.EXPECT().Save(gomock.Any(), args.request).Return(&testExistingTask, nil)
			},
		},
		{
			name: "negative invalid request",
			args: args{
				request: &model.Task{},
			},
			expectedCode:   http.StatusBadRequest,
			expectedResult: testErrInvalidRequestJSON,
			mockFunc:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			// prepare services
			ctrl := gomock.NewController(t1)
			mock, router := prepare(ctrl)
			w := httptest.NewRecorder()

			// setup mock
			if tt.mockFunc != nil {
				tt.mockFunc(mock, tt.args)
			}

			// prepare and execute http request
			reqData, err := json.Marshal(&tt.args.request)
			if err != nil {
				t1.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/v1/tasks/", bytes.NewBuffer(reqData))
			if err != nil {
				t1.Fatal(err)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t1, tt.expectedCode, w.Code, "not expected response status code")
			assert.Equal(t1, tt.expectedResult, w.Body.String(), "not expected response")
		})
	}
}

func TestTaskController_GetTask(t *testing.T) {
	type args struct {
		request string
	}
	tests := []struct {
		name           string
		args           args
		expectedCode   int
		expectedResult string
		mockFunc       func(mock *mockstorage.MockTaskStorage, args args)
	}{
		{
			name: "positive get",
			args: args{
				request: "test-id-1",
			},
			expectedCode:   http.StatusOK,
			expectedResult: testExistingTaskJSON,
			mockFunc: func(mock *mockstorage.MockTaskStorage, args args) {
				mock.EXPECT().Get(gomock.Any(), args.request).Return(&testExistingTask, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			// prepare services
			ctrl := gomock.NewController(t1)
			mock, router := prepare(ctrl)
			w := httptest.NewRecorder()

			// setup mock
			if tt.mockFunc != nil {
				tt.mockFunc(mock, tt.args)
			}

			// prepare and execute http request
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/tasks/%s", tt.args.request), nil)
			if err != nil {
				t1.Fatal(err)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t1, tt.expectedCode, w.Code, "not expected response status code")
			assert.Equal(t1, tt.expectedResult, w.Body.String(), "not expected response")
		})
	}
}

func TestTaskController_GetAllTasks(t *testing.T) {
	tests := []struct {
		name           string
		expectedCode   int
		expectedResult string
		mockFunc       func(mock *mockstorage.MockTaskStorage)
	}{
		{
			name:           "positive get all",
			expectedCode:   http.StatusOK,
			expectedResult: fmt.Sprintf("[%s]", testExistingTaskJSON),
			mockFunc: func(mock *mockstorage.MockTaskStorage) {
				mock.EXPECT().GetAll(gomock.Any()).Return([]*model.Task{&testExistingTask}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			// prepare services
			ctrl := gomock.NewController(t1)
			mock, router := prepare(ctrl)
			w := httptest.NewRecorder()

			// setup mock
			if tt.mockFunc != nil {
				tt.mockFunc(mock)
			}

			// prepare and execute http request
			req, err := http.NewRequest(http.MethodGet, "/v1/tasks/", nil)
			if err != nil {
				t1.Fatal(err)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t1, tt.expectedCode, w.Code, "not expected response status code")
			assert.Equal(t1, tt.expectedResult, w.Body.String(), "not expected response")
		})
	}
}

func TestTaskController_DeleteTask(t *testing.T) {
	type args struct {
		request string
	}
	tests := []struct {
		name           string
		args           args
		expectedCode   int
		expectedResult string
		mockFunc       func(mock *mockstorage.MockTaskStorage, args args)
	}{
		{
			name: "positive delete",
			args: args{
				request: testExistingTask.ID,
			},
			expectedCode:   http.StatusNoContent,
			expectedResult: "",
			mockFunc: func(mock *mockstorage.MockTaskStorage, args args) {
				mock.EXPECT().Delete(gomock.Any(), args.request).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			// prepare services
			ctrl := gomock.NewController(t1)
			mock, router := prepare(ctrl)
			w := httptest.NewRecorder()

			// setup mock
			if tt.mockFunc != nil {
				tt.mockFunc(mock, tt.args)
			}

			// prepare and execute http request
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/tasks/%s", tt.args.request), nil)
			if err != nil {
				t1.Fatal(err)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t1, tt.expectedCode, w.Code, "not expected response status code")
			assert.Equal(t1, tt.expectedResult, w.Body.String(), "not expected response")
		})
	}
}
