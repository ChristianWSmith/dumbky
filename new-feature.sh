#!/bin/bash

# Check if feature name is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <FeatureName>"
    exit 1
fi

FEATURE_NAME_TITLE_CASE=$1
FEATURE_NAME_LOWER_CASE=$(echo "$FEATURE_NAME_TITLE_CASE" | tr '[:upper:]' '[:lower:]')
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
FEATURE_DIR="$SCRIPT_DIR/internal/features/$FEATURE_NAME_LOWER_CASE"

# Check if the feature directory already exists
if [ -d "$FEATURE_DIR" ]; then
    echo "Error: Directory $FEATURE_DIR already exists."
    exit 1
fi

# Create the feature directory
mkdir -p "$FEATURE_DIR/controller"
mkdir -p "$FEATURE_DIR/view"
mkdir -p "$FEATURE_DIR/model"

# Function to write content to file
write_file() {
    local file_path=$1
    local content=$2

    echo "$content" > "$file_path"
}

# Create controller file
write_file "$FEATURE_DIR/controller/${FEATURE_NAME_LOWER_CASE}_controller.go" "package controller

import (
        \"dumbky/internal/features/$FEATURE_NAME_LOWER_CASE/model\"
        \"dumbky/internal/features/$FEATURE_NAME_LOWER_CASE/view\"

        \"fyne.io/fyne/v2\"
)

type controllerImpl struct {
        model model.Model
        view  view.View
}

func NewController() *controllerImpl {
        c := &controllerImpl{
                model: model.NewModel(),
                view:  view.NewView(),
        }

        c.bindAll()

        return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
        return c.view.CanvasObject()
}

func (c *controllerImpl) bindAll() {
        // bindings := c.model.GetBindings()
        // bindables := c.view.GetBindables()
}
"

# Create view file
write_file "$FEATURE_DIR/view/${FEATURE_NAME_LOWER_CASE}_view.go" "package view

import (
        \"fyne.io/fyne/v2\"
)

type viewImpl struct {
        ui fyne.CanvasObject
}

type Bindables struct {
}

type View interface {
        GetBindables() Bindables
        CanvasObject() fyne.CanvasObject
}

func NewView() View {
        v := &viewImpl{
                ui: nil,
        }

        return v
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
        return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
        return Bindables{}
}
"

# Create model file
write_file "$FEATURE_DIR/model/${FEATURE_NAME_LOWER_CASE}_model.go" "package model

type modelImpl struct {
}

type Bindings struct {
}

type Model interface {
        GetBindings() Bindings
}

func NewModel() Model {
        m := &modelImpl{}
        return m
}

func (m *modelImpl) GetBindings() Bindings {
        return Bindings{}
}
"

# Create main feature file
write_file "$FEATURE_DIR/${FEATURE_NAME_LOWER_CASE}.go" "package response

import (
        \"dumbky/internal/features\"
        \"dumbky/internal/features/$FEATURE_NAME_LOWER_CASE/controller\"
)

type ${FEATURE_NAME_TITLE_CASE}Controller interface {
        features.Controller
}

func New() ${FEATURE_NAME_TITLE_CASE}Controller {
        return controller.NewController()
}
"

echo "Feature $FEATURE_NAME_TITLE_CASE has been created successfully."
