package requestbodyview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/ui/validators"
	"dumbky/internal/ui/views/keyvalueeditorview"
	"dumbky/internal/utils"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type RequestBodyView struct {
	UI                 *fyne.Container
	BodyKeyValueEditor keyvalueeditorview.KeyValueEditorView
	BodyTypeBinding    binding.String
	BodyRawBinding     binding.String

	bodyRawEntry   *widget.Entry
	bodyTypeSelect *widget.Select
}

type RequestBodyState struct {
	BodyType string
	BodyForm keyvalueeditorview.KeyValueEditorState
	BodyRaw  string
}

func (rbv RequestBodyView) ToState() (RequestBodyState, error) {
	bodyType, bodyTypeErr := rbv.BodyTypeBinding.Get()
	if bodyTypeErr != nil {
		log.Error(bodyTypeErr)
		return RequestBodyState{}, bodyTypeErr
	}
	bodyForm, bodyFormErr := rbv.BodyKeyValueEditor.ToState()
	if bodyFormErr != nil {
		log.Error(bodyFormErr)
		return RequestBodyState{}, bodyFormErr
	}
	bodyRaw, bodyRawErr := rbv.BodyRawBinding.Get()
	if bodyRawErr != nil {
		log.Error(bodyRawErr)
		return RequestBodyState{}, bodyRawErr
	}
	return RequestBodyState{
		BodyType: bodyType,
		BodyForm: bodyForm,
		BodyRaw:  bodyRaw,
	}, nil
}

func (rbv RequestBodyView) LoadState(requestBodyState RequestBodyState) error {
	bodyType := requestBodyState.BodyType
	if !utils.ElementInSlice(constants.UIBodyTypes(), bodyType) {
		bodyType = constants.UI_BODY_TYPE_DEFAULT
	}
	bodyTypeErr := rbv.BodyTypeBinding.Set(bodyType)
	if bodyTypeErr != nil {
		log.Error(bodyTypeErr)
		return bodyTypeErr
	}
	bodyFormErr := rbv.BodyKeyValueEditor.LoadState(requestBodyState.BodyForm)
	if bodyFormErr != nil {
		log.Error(bodyFormErr)
		return bodyFormErr
	}
	bodyRawErr := rbv.BodyRawBinding.Set(requestBodyState.BodyRaw)
	if bodyRawErr != nil {
		log.Error(bodyRawErr)
		return bodyRawErr
	}
	return nil
}

func (rbv RequestBodyView) ValidateBodyRaw() error {
	return rbv.bodyRawEntry.Validate()
}

func (rbv RequestBodyView) EnableBodyTypeSelect() {
	rbv.bodyTypeSelect.Enable()
}

func (rbv RequestBodyView) DisableBodyTypeSelect() {
	rbv.bodyTypeSelect.Disable()
}

func ComposeRequestBodyView() RequestBodyView {
	bodyTypeBind := binding.NewString()
	bodyRawBind := binding.NewString()

	bodyKeyValueEditorView := keyvalueeditorview.ComposeKeyValueEditorView(validators.ValidateFormBodyKey, validators.ValidateFormBodyValue)
	bodyTypeSelect := widget.NewSelect(constants.UIBodyTypes(), nil)
	bodyRawEntry := widget.NewMultiLineEntry()
	bodyRawEntry.TextStyle.Monospace = true
	bodyRawEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_BODY_TYPE_RAW)

	bodyContentStack := container.NewStack(bodyKeyValueEditorView.UI, bodyRawEntry)

	bodyTypeSelect.SetSelected(constants.UI_BODY_TYPE_DEFAULT)
	bodyTypeErr := bodyTypeBind.Set(constants.UI_BODY_TYPE_DEFAULT)
	if bodyTypeErr != nil {
		log.Error(bodyTypeErr)
	}

	bodyTypeSelect.Bind(bodyTypeBind)
	bodyRawEntry.Bind(bodyRawBind)

	bodyRawEntry.Validator = validators.ValidateRawBodyContent

	bodyTypeBind.AddListener(binding.NewDataListener(func() {
		bodyType, bodyTypeErr := bodyTypeBind.Get()
		if bodyTypeErr != nil {
			log.Error(bodyTypeErr)
			return
		}
		if bodyType == constants.UI_BODY_TYPE_FORM {
			bodyRawEntry.Hide()
			bodyKeyValueEditorView.UI.Show()
			bodyContentStack.Refresh()
		} else if bodyType == constants.UI_BODY_TYPE_RAW {
			bodyKeyValueEditorView.UI.Hide()
			bodyRawEntry.Show()
			bodyContentStack.Refresh()
		} else if bodyType == constants.UI_BODY_TYPE_NONE {
			bodyKeyValueEditorView.UI.Hide()
			bodyRawEntry.Hide()
			bodyContentStack.Refresh()
		} else {
			log.Error(errors.New("invalid body type"))
		}
	}))

	ui := container.NewBorder(bodyTypeSelect, nil, nil, nil, bodyContentStack)

	return RequestBodyView{
		ui,
		bodyKeyValueEditorView,
		bodyTypeBind,
		bodyRawBind,
		bodyRawEntry,
		bodyTypeSelect,
	}
}
