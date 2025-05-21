package keyvalueview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type KeyValueView struct {
	UI             *fyne.Container
	DestroyButton  *widget.Button
	EnabledBinding binding.Bool
	KeyBinding     binding.String
	ValueBinding   binding.String

	keyEntry   *widget.Entry
	valueEntry *widget.Entry
}

type KeyValueState struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

func (kvv KeyValueView) ValidateKey() error {
	return kvv.keyEntry.Validate()
}

func (kvv KeyValueView) ValidateValue() error {
	return kvv.valueEntry.Validate()
}

func (kvv KeyValueView) ToState() (KeyValueState, error) {
	enabled, enabledErr := kvv.EnabledBinding.Get()
	if enabledErr != nil {
		log.Error(enabledErr)
		return KeyValueState{}, enabledErr
	}
	key, keyErr := kvv.KeyBinding.Get()
	if keyErr != nil {
		log.Error(keyErr)
		return KeyValueState{}, keyErr
	}
	value, valueErr := kvv.ValueBinding.Get()
	if valueErr != nil {
		log.Error(valueErr)
		return KeyValueState{}, valueErr
	}
	return KeyValueState{
		Enabled: enabled,
		Key:     key,
		Value:   value,
	}, nil
}

func (kvv KeyValueView) LoadState(kvs KeyValueState) error {
	enabledErr := kvv.EnabledBinding.Set(kvs.Enabled)
	if enabledErr != nil {
		log.Error(enabledErr)
		return enabledErr
	}
	keyErr := kvv.KeyBinding.Set(kvs.Key)
	if keyErr != nil {
		log.Error(keyErr)
		return keyErr
	}
	valueErr := kvv.ValueBinding.Set(kvs.Value)
	if valueErr != nil {
		log.Error(valueErr)
		return valueErr
	}
	return nil
}

func ComposeKeyValueView(keyValidator, valueValidator func(val string) error) KeyValueView {
	keyBinding := binding.NewString()
	valueBinding := binding.NewString()
	enabledBinding := binding.NewBool()

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_KEY)
	keyEntry.TextStyle.Monospace = true
	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_VALUE)
	valueEntry.TextStyle.Monospace = true
	enabledCheck := widget.NewCheck("", nil)
	destroyButton := widget.NewButtonWithIcon("", nil, nil)
	destroyButton.Icon = destroyButton.Theme().Icon(theme.IconNameContentClear)

	keyEntry.Bind(keyBinding)
	valueEntry.Bind(valueBinding)
	enabledCheck.Bind(enabledBinding)

	keyEntry.Validator = keyValidator
	valueEntry.Validator = valueValidator

	err := enabledBinding.Set(true)
	if err != nil {
		log.Error(err)
	}
	enabledBinding.AddListener(binding.NewDataListener(func() {
		enabled, enabledErr := enabledBinding.Get()
		if enabledErr != nil {
			log.Error(enabledErr)
			return
		}
		if enabled {
			keyEntry.Enable()
			valueEntry.Enable()
		} else {
			keyEntry.Disable()
			valueEntry.Disable()
		}
	}))

	keyValue := container.NewGridWithColumns(2, keyEntry, valueEntry)

	ui := container.NewBorder(nil, nil, enabledCheck, destroyButton, keyValue)

	return KeyValueView{
		ui,
		destroyButton,
		enabledBinding,
		keyBinding,
		valueBinding,
		keyEntry,
		valueEntry,
	}
}
