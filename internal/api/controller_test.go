package api

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockLoaderService struct {
	mock.Mock
}

func (m *MockLoaderService) LoadComponent(cType string, version string, hash []byte) ([]byte, []byte, error) {
	args := m.Called(cType, version)
	var content []byte
	var err error

	if args.Get(0) != nil {
		content = args.Get(0).([]byte)
		hash = args.Get(1).([]byte)
		err = nil
	} else {
		content = nil
		hash = nil
		err = args.Error(2)
	}
	return content, hash, err
}

func (m *MockLoaderService) AssertExpectations(t *testing.T) {
	m.Mock.AssertExpectations(t)
}

func TestLoaderController_LoadComponent(t *testing.T) {
	mockService := &MockLoaderService{}
	mockService.On("LoadComponent", "not_default", "1.0.1", mock.Anything).Return([]byte{1}, []byte{1}, nil).Once()
	mockService.On("LoadComponent", "core", "1.0.0", mock.Anything).Return([]byte{1}, []byte{1}, nil).Once()
	mockService.On("LoadComponent", "non_exist", "1.0.0", mock.Anything).Return(nil, nil, errors.New("service error")).Once()

	tests := []struct {
		name    string
		payload *Request
		reply   *Response
		wantErr bool
	}{
		{
			name:    "Success - Valid Payload",
			payload: &Request{Type: "not_default", Version: "1.0.1"},
			reply:   &Response{},
			wantErr: false,
		},
		{
			name:    "Success - No Payload",
			payload: &Request{Type: "core", Version: "1.0.0"},
			reply:   &Response{},
			wantErr: false,
		},
		{
			name:    "Error - Service Returns Error",
			payload: &Request{Type: "non_exist", Version: "1.0.0"},
			reply:   &Response{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LoaderController{
				loaderService: mockService,
			}

			err := lc.LoadComponent(tt.payload, tt.reply)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadComponent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	mockService.AssertExpectations(t)
}
