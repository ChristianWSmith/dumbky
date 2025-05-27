#!/usr/bin/env bash

set -e

FEATURE_PATH="internal/features"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

FEATURE_NAME=$1

if [ -z "$FEATURE_NAME" ]; then
  echo "Usage: $0 <feature-name>"
  exit 1
fi


if [ -d "${SCRIPT_DIR}/${FEATURE_PATH}/${FEATURE_NAME}" ]; then
  echo "Feature '${FEATURE_NAME}' already exists."
  exit 1
fi


mkdir -p "${SCRIPT_DIR}/${FEATURE_PATH}/${FEATURE_NAME}"

cat > "${SCRIPT_DIR}/${FEATURE_PATH}/${FEATURE_NAME}/${FEATURE_NAME}_controller.go" << EOF
package ${FEATURE_NAME}

import (
	"dumbky/internal/features"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model *model
	view  *view
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	model := newModel()
	view := newView()
	return &Controller{
		model: model,
		view:  view,
	}
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

EOF

cat > "${SCRIPT_DIR}/${FEATURE_PATH}/${FEATURE_NAME}/${FEATURE_NAME}_model.go" << EOF
package ${FEATURE_NAME}

type model struct {
}

func newModel() *model {
	return &model{}
}

EOF

cat > "${SCRIPT_DIR}/${FEATURE_PATH}/${FEATURE_NAME}/${FEATURE_NAME}_view.go" << EOF
package ${FEATURE_NAME}

import (
	"fyne.io/fyne/v2"
)

type view struct {
	ui fyne.CanvasObject
}

func newView() *view {
	return &view{
		ui: nil,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}

EOF

cat > "${SCRIPT_DIR}/${FEATURE_PATH}/${FEATURE_NAME}/${FEATURE_NAME}_test.go" << EOF
package ${FEATURE_NAME}

EOF

