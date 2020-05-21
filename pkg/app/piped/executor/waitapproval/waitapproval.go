// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package waitapproval

import (
	"context"
	"time"

	"github.com/kapetaniosci/pipe/pkg/app/piped/executor"
	"github.com/kapetaniosci/pipe/pkg/model"
)

type Executor struct {
	executor.Input
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageWaitApproval, f)
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	var (
		originalStatus = e.Stage.Status
		ctx            = sig.Context()
		ticker         = time.NewTicker(5 * time.Second)
	)
	defer ticker.Stop()

	e.LogPersister.AppendInfo("Waiting for an approval...")
	for {
		select {
		case <-ticker.C:
			if ok := e.checkApproval(ctx); !ok {
				continue
			}
			e.LogPersister.AppendInfo("Got an approval from abc")
			return model.StageStatus_STAGE_SUCCESS

		case s := <-sig.Ch():
			switch s {
			case executor.StopSignalCancel:
				return model.StageStatus_STAGE_CANCELLED
			case executor.StopSignalTerminate:
				return originalStatus
			default:
				return model.StageStatus_STAGE_FAILURE
			}
		}
	}
}

func (e *Executor) checkApproval(ctx context.Context) bool {
	commands := e.CommandLister.ListCommands()

	for _, cmd := range commands {
		c := cmd.GetApproveStage()
		if c == nil {
			continue
		}

		if err := cmd.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil); err == nil {
			return true
		}
		return false
	}

	return false
}
