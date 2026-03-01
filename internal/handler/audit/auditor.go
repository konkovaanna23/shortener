package audit

import "github.com/sirupsen/logrus"

func NewAuditor(auditFile, auditURL string) *Publisher {
	auditor := NewPublisher()
	if auditFile != "" {
		if fileObs, err := NewFileObserver(auditFile); err == nil {
			auditor.Subscribe(fileObs)
		} else {
			logrus.Error(err)
		}

	}
	if auditURL != "" {
		httpObs := NewHTTPObserver(auditURL)
		auditor.Subscribe(httpObs)
	}
	return auditor
}
