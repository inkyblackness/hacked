package movies

import (
	"context"

	"github.com/inkyblackness/hacked/ss1/content/movie"
)

type compressionTask struct {
	width      int
	height     int
	input      movie.Scene
	ctx        context.Context
	ctxCancel  context.CancelFunc
	resultChan chan compressionResult
}

type compressionResult interface{}

type compressionAborted struct{}

type compressionFailed struct{ err error }

type compressionFinished struct{ scene movie.HighResScene }

func newCompressionTask(scene movie.Scene, width, height int) *compressionTask {
	task := &compressionTask{
		width:      width,
		height:     height,
		input:      scene,
		resultChan: make(chan compressionResult),
	}
	task.ctx, task.ctxCancel = context.WithCancel(context.Background())

	return task
}

func (task *compressionTask) run() {
	defer close(task.resultChan)
	highResScene, err := movie.HighResSceneFrom(task.ctx, task.input, task.width, task.height)
	switch {
	case task.ctx.Err() != nil:
		task.resultChan <- compressionAborted{}
	case err != nil:
		task.resultChan <- compressionFailed{err: err}
	default:
		task.resultChan <- compressionFinished{scene: highResScene}
	}
}

func (task *compressionTask) update() compressionResult {
	select {
	case result, ok := <-task.resultChan:
		if ok {
			return result
		} else {
			return compressionAborted{}
		}
	default:
		return nil
	}
}

func (task *compressionTask) cancel() {
	task.ctxCancel()
}
