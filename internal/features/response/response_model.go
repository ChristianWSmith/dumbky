package response

import "fyne.io/fyne/v2/data/binding"

type model struct {
	status binding.String
	time   binding.String
	body   binding.String
}

func newModel() *model {
	return &model{
		status: binding.NewString(),
		time:   binding.NewString(),
		body:   binding.NewString(),
	}
}

func (m *model) setResponse(status, time, body string) {
	m.status.Set(status)
	m.time.Set(time)
	m.body.Set(body)
}
