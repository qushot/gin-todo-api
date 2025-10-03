package usecase_test

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/google/go-cmp/cmp"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	mock_repository "github.com/qushot/gin-todo-api/internal/mocks/repository"
	"github.com/qushot/gin-todo-api/internal/usecase"
)

func Test_createTodo_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock_repository.NewMockTodo(ctrl)
	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, todo model.Todo) (*model.Todo, error) {
			return &todo, nil
		}).AnyTimes()

	tests := []struct {
		name    string
		todo    model.Todo
		want    *model.Todo
		wantErr bool
	}{
		{
			name:    "success",
			todo:    model.Todo{ID: "1", Title: "Test Todo", Done: false},
			want:    &model.Todo{ID: "1", Title: "Test Todo", Done: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := usecase.NewCreateTodo(mockTodoRepo)
			got, gotErr := uc.Execute(context.Background(), tt.todo)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Execute() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Execute() succeeded unexpectedly")
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
