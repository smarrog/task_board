package grpc

import (
	"context"

	"github.com/rs/zerolog"
	taskdo "github.com/smarrog/task-board/core-service/internal/domain/task"
	taskuc "github.com/smarrog/task-board/core-service/internal/usecase/task"
	v1 "github.com/smarrog/task-board/shared/proto/base/v1"
)

type TasksHandler struct {
	v1.UnimplementedBoardsServiceServer

	log *zerolog.Logger

	createTask *taskuc.CreateTaskUseCase
	getTask    *taskuc.GetTaskUseCase
	updateTask *taskuc.UpdateTaskUseCase
	deleteTask *taskuc.DeleteTaskUseCase
}

func NewTasksHandler(
	log *zerolog.Logger,
	createTask *taskuc.CreateTaskUseCase,
	getTask *taskuc.GetTaskUseCase,
	updateTask *taskuc.UpdateTaskUseCase,
	deleteTask *taskuc.DeleteTaskUseCase,
) *TasksHandler {
	return &TasksHandler{
		log:        log,
		createTask: createTask,
		getTask:    getTask,
		updateTask: updateTask,
		deleteTask: deleteTask,
	}
}

func (h *TasksHandler) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.CreateTaskResponse, error) {
	input := taskuc.CreateTaskInput{
		ColumnId:    req.ColumnId,
		Position:    int(req.Position),
		Title:       req.Title,
		Description: req.Description,
		AssigneeId:  req.AssigneeId,
	}

	output, err := h.createTask.Execute(ctx, input)
	if err != nil {
		return nil, mapTasksErr(err)
	}

	return &v1.CreateTaskResponse{
		Task: toProtoTask(output.Task),
	}, nil
}

func (h *TasksHandler) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.GetTaskResponse, error) {
	input := taskuc.GetTaskInput{
		TaskId: req.TaskId,
	}

	output, err := h.getTask.Execute(ctx, input)
	if err != nil {
		return nil, mapTasksErr(err)
	}

	return &v1.GetTaskResponse{
		Task: toProtoTask(output.Task),
	}, nil
}

func (h *TasksHandler) UpdateTask(ctx context.Context, req *v1.UpdateTaskRequest) (*v1.UpdateTaskResponse, error) {
	input := taskuc.UpdateTaskInput{
		TaskId:      req.TaskId,
		ColumnId:    req.ColumnId,
		Position:    int(req.Position),
		Title:       req.Title,
		Description: req.Description,
		AssigneeId:  req.AssigneeId,
	}

	output, err := h.updateTask.Execute(ctx, input)
	if err != nil {
		return nil, mapTasksErr(err)
	}

	return &v1.UpdateTaskResponse{
		Task: toProtoTask(output.Task),
	}, nil
}

func (h *TasksHandler) DeleteTask(ctx context.Context, req *v1.DeleteTaskRequest) (*v1.DeleteTaskResponse, error) {
	input := taskuc.DeleteTaskInput{
		TaskId: req.TaskId,
	}

	_, err := h.deleteTask.Execute(ctx, input)
	if err != nil {
		return nil, mapTasksErr(err)
	}

	return &v1.DeleteTaskResponse{}, nil
}

func toProtoTask(b *taskdo.Task) *v1.Task {
	return &v1.Task{
		Id:          b.Id().String(),
		ColumnId:    b.ColumnId().String(),
		Position:    int32(b.Position().Int()),
		Title:       b.Title().String(),
		Description: b.Description().String(),
		AssigneeId:  b.AssigneeId().String(),
	}
}

func mapTasksErr(err error) error {
	switch {
	default:
		return mapCommonErr(err)
	}
}
