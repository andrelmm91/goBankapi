package mail

import (
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "a test email"
	content := `
	<h1> test </h1>
	`
	to := []string{"andre.lmm91@gmail.com"}
	attachFiles := []string{"../godependencies.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
