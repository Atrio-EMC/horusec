// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package eslint

import (
	"github.com/ZupIT/horusec/development-kit/pkg/entities/analyser/eslint"
	"github.com/ZupIT/horusec/development-kit/pkg/enums/languages"
	"github.com/ZupIT/horusec/development-kit/pkg/enums/tools"
	jsonUtils "github.com/ZupIT/horusec/development-kit/pkg/utils/json"
	"github.com/ZupIT/horusec/development-kit/pkg/utils/logger"
	dockerEntities "github.com/ZupIT/horusec/horusec-cli/internal/entities/docker"
	"github.com/ZupIT/horusec/horusec-cli/internal/helpers/messages"
	"github.com/ZupIT/horusec/horusec-cli/internal/services/formatters"
)

type Formatter struct {
	formatters.IService
}

func NewFormatter(service formatters.IService) formatters.IFormatter {
	return &Formatter{
		service,
	}
}

func (f *Formatter) StartAnalysis(projectSubPath string) {
	if f.ToolIsToIgnore(tools.Eslint) {
		logger.LogDebugWithLevel(messages.MsgDebugToolIgnored+tools.Eslint.ToString(), logger.DebugLevel)
		return
	}

	err := f.executeDockerContainer(projectSubPath)
	f.LogAnalysisError(err, tools.Eslint, projectSubPath)

	f.SetLanguageIsFinished()
}

func (f *Formatter) executeDockerContainer(projectSubPath string) error {
	f.LogDebugWithReplace(messages.MsgDebugToolStartAnalysis, tools.Eslint)

	output, err := f.ExecuteContainer(f.getDockerConfig(projectSubPath))
	if err != nil {
		f.SetAnalysisError(err)
		return err
	}

	f.processOutput(output)
	f.LogDebugWithReplace(messages.MsgDebugToolFinishAnalysis, tools.Eslint)

	return nil
}

func (f *Formatter) getDockerConfig(projectSubPath string) *dockerEntities.AnalysisData {
	return &dockerEntities.AnalysisData{
		Image:    ImageName,
		Tag:      ImageTag,
		CMD:      f.AddWorkDirInCmd(ImageCmd, projectSubPath, tools.Eslint),
		Language: languages.Javascript,
	}
}

func (f *Formatter) processOutput(output string) {
	if output == "" {
		logger.LogDebugWithLevel(
			messages.MsgDebugOutputEmpty, logger.DebugLevel, map[string]interface{}{"tool": tools.Eslint.ToString()})
		return
	}

	eslintOutput, err := f.parseOutput(output)
	if err != nil {
		return
	}

	logger.LogInfo("", eslintOutput, *eslintOutput, &eslintOutput)
}

func (f *Formatter) parseOutput(output string) (eslintOutput *[]eslint.Output, err error) {
	err = jsonUtils.ConvertStringToOutput(output, &eslintOutput)
	logger.LogErrorWithLevel(f.GetAnalysisIDErrorMessage(tools.GoSec, output), err, logger.ErrorLevel)
	return eslintOutput, err
}
