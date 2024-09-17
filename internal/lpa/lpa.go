package lpa

import (
	"context"
	"log/slog"

	"github.com/damonto/estkme-cloud/internal/config"
	"github.com/damonto/libeuicc-go"
)

type LPA struct {
	*libeuicc.Libeuicc
}

func New(driver libeuicc.APDU) (*LPA, error) {
	logLevel := libeuicc.LogInfoLevel
	if config.C.Verbose {
		logLevel = libeuicc.LogDebugLevel
	}
	libeuicc, err := libeuicc.New(driver, libeuicc.NewDefaultLogger(logLevel))
	if err != nil {
		return nil, err
	}
	return &LPA{
		Libeuicc: libeuicc,
	}, nil
}

func (l *LPA) Download(ctx context.Context, activationCode *libeuicc.ActivationCode, downloadOption *libeuicc.DownloadOption) error {
	return l.sendNotificationsAfterExecution(func() error {
		return l.DownloadProfile(ctx, activationCode, downloadOption)
	})
}

func (l *LPA) sendNotificationsAfterExecution(action func() error) error {
	currentNotifications, err := l.GetNotifications()
	if err != nil {
		return err
	}
	var lastSeqNumber int
	if len(currentNotifications) > 0 {
		lastSeqNumber = currentNotifications[len(currentNotifications)-1].SeqNumber
	}

	if err := action(); err != nil {
		slog.Error("failed to execute action", "error", err)
		return err
	}

	notifications, err := l.GetNotifications()
	if err != nil {
		slog.Error("failed to get notifications", "error", err)
		return err
	}
	for _, notification := range notifications {
		if notification.SeqNumber > lastSeqNumber {
			if err := l.ProcessNotification(
				notification.SeqNumber,
				notification.ProfileManagementOperation != libeuicc.NotificationProfileManagementOperationDelete,
			); err != nil {
				slog.Error("failed to process notification", "error", err, "notification", notification)
				return err
			}
			slog.Info("notification processed", "notification", notification)
		}
	}
	return nil
}
