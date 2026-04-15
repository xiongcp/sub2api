//go:build unit

package service

import (
	"bufio"
	"context"
	"encoding/base64"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type emailSettingRepoStub struct {
	values map[string]string
}

func (s *emailSettingRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
}

func (s *emailSettingRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	panic("unexpected GetValue call")
}

func (s *emailSettingRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
}

func (s *emailSettingRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		out[key] = s.values[key]
	}
	return out, nil
}

func (s *emailSettingRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
}

func (s *emailSettingRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
}

func (s *emailSettingRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
}

func TestNormalizeSMTPSecurityModeWithLegacy(t *testing.T) {
	require.Equal(t, SMTPSecurityModeStartTLS, NormalizeSMTPSecurityModeWithLegacy("", false))
	require.Equal(t, SMTPSecurityModeImplicitTLS, NormalizeSMTPSecurityModeWithLegacy("", true))
	require.Equal(t, SMTPSecurityModePlain, NormalizeSMTPSecurityModeWithLegacy("plain", true))
	require.Equal(t, SMTPSecurityModeStartTLS, NormalizeSMTPSecurityModeWithLegacy("starttls", true))
}

func TestEmailService_GetSMTPConfig_PrefersSecurityModeAndFallsBackToLegacy(t *testing.T) {
	t.Run("prefers explicit mode", func(t *testing.T) {
		repo := &emailSettingRepoStub{
			values: map[string]string{
				SettingKeySMTPHost:         "smtp.example.com",
				SettingKeySMTPPort:         "465",
				SettingKeySMTPUsername:     "user",
				SettingKeySMTPPassword:     "secret",
				SettingKeySMTPSecurityMode: "plain",
				SettingKeySMTPUseTLS:       "true",
			},
		}
		svc := NewEmailService(repo, nil)

		cfg, err := svc.GetSMTPConfig(context.Background())
		require.NoError(t, err)
		require.Equal(t, SMTPSecurityModePlain, cfg.SecurityMode)
		require.False(t, cfg.UseTLS)
	})

	t.Run("falls back to legacy bool", func(t *testing.T) {
		repo := &emailSettingRepoStub{
			values: map[string]string{
				SettingKeySMTPHost:   "smtp.example.com",
				SettingKeySMTPUseTLS: "true",
			},
		}
		svc := NewEmailService(repo, nil)

		cfg, err := svc.GetSMTPConfig(context.Background())
		require.NoError(t, err)
		require.Equal(t, SMTPSecurityModeImplicitTLS, cfg.SecurityMode)
		require.True(t, cfg.UseTLS)
		require.Equal(t, 587, cfg.Port)
	})
}

func TestEmailService_TestSMTPConnectionWithConfig_StartTLSRequiresExtension(t *testing.T) {
	addr := startPlainSMTPServer(t, false)
	host, port := splitHostPort(t, addr)
	svc := NewEmailService(nil, nil)

	err := svc.TestSMTPConnectionWithConfig(&SMTPConfig{
		Host:         host,
		Port:         port,
		SecurityMode: SMTPSecurityModeStartTLS,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "server does not support STARTTLS")
}

func TestEmailService_TestSMTPConnectionWithConfig_ImplicitTLSFailsAgainstPlainServer(t *testing.T) {
	addr := startPlainSMTPServer(t, true)
	host, port := splitHostPort(t, addr)
	svc := NewEmailService(nil, nil)

	err := svc.TestSMTPConnectionWithConfig(&SMTPConfig{
		Host:         host,
		Port:         port,
		SecurityMode: SMTPSecurityModeImplicitTLS,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "tls connection failed")
}

func TestEmailService_TestSMTPConnectionWithConfig_FallsBackToLoginWhenPlainIsRejected(t *testing.T) {
	addr, done := startAuthSMTPServer(t, authServerConfig{
		expectedConnections: 2,
		authLine:            "PLAIN LOGIN",
		handleAuth: func(conn net.Conn, reader *bufio.Reader, session int, line string) error {
			if session == 0 {
				if line != "AUTH PLAIN AHVzZXIAc2VjcmV0" {
					return unexpectedSMTPLineError("AUTH PLAIN AHVzZXIAc2VjcmV0", line)
				}
				_, _ = conn.Write([]byte("504 5.7.4 Unrecognized authentication type\r\n"))
				return nil
			}

			if strings.ToUpper(line) != "AUTH LOGIN" {
				return unexpectedSMTPLineError("AUTH LOGIN", line)
			}
			_, _ = conn.Write([]byte("334 VXNlcm5hbWU6\r\n"))

			line, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if strings.TrimSpace(line) != base64.StdEncoding.EncodeToString([]byte("user")) {
				return unexpectedSMTPLineError(base64.StdEncoding.EncodeToString([]byte("user")), strings.TrimSpace(line))
			}
			_, _ = conn.Write([]byte("334 UGFzc3dvcmQ6\r\n"))

			line, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			if strings.TrimSpace(line) != base64.StdEncoding.EncodeToString([]byte("secret")) {
				return unexpectedSMTPLineError(base64.StdEncoding.EncodeToString([]byte("secret")), strings.TrimSpace(line))
			}
			_, _ = conn.Write([]byte("235 2.7.0 Authentication successful\r\n"))

			line, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
			if strings.ToUpper(strings.TrimSpace(line)) != "QUIT" {
				return unexpectedSMTPLineError("QUIT", strings.TrimSpace(line))
			}
			_, _ = conn.Write([]byte("221 Bye\r\n"))
			return nil
		},
	})
	host, port := splitHostPort(t, addr)
	svc := NewEmailService(nil, nil)

	err := svc.TestSMTPConnectionWithConfig(&SMTPConfig{
		Host:         host,
		Port:         port,
		Username:     "user",
		Password:     "secret",
		SecurityMode: SMTPSecurityModePlain,
	})
	require.NoError(t, err)
	require.NoError(t, <-done)
}

func TestEmailService_TestSMTPConnectionWithConfig_UnsupportedAuthMechanism(t *testing.T) {
	addr, done := startAuthSMTPServer(t, authServerConfig{
		expectedConnections: 1,
		authLine:            "XOAUTH2",
	})
	host, port := splitHostPort(t, addr)
	svc := NewEmailService(nil, nil)

	err := svc.TestSMTPConnectionWithConfig(&SMTPConfig{
		Host:         host,
		Port:         port,
		Username:     "user",
		Password:     "secret",
		SecurityMode: SMTPSecurityModePlain,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "supported smtp auth mechanism")
	require.Contains(t, err.Error(), "XOAUTH2")
	require.NoError(t, <-done)
}

func startPlainSMTPServer(t *testing.T, allowQuit bool) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write([]byte("220 localhost ESMTP ready\r\n"))
		reader := bufio.NewReader(conn)
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "EHLO") {
			_, _ = conn.Write([]byte("250 localhost\r\n"))
			if allowQuit {
				line, _ = reader.ReadString('\n')
				if strings.ToUpper(strings.TrimSpace(line)) == "QUIT" {
					_, _ = conn.Write([]byte("221 Bye\r\n"))
				}
			}
		}
	}()

	return ln.Addr().String()
}

type authServerConfig struct {
	expectedConnections int
	authLine            string
	handleAuth          func(conn net.Conn, reader *bufio.Reader, session int, firstAuthLine string) error
}

func startAuthSMTPServer(t *testing.T, cfg authServerConfig) (string, <-chan error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	done := make(chan error, 1)
	expectedConnections := cfg.expectedConnections
	if expectedConnections <= 0 {
		expectedConnections = 1
	}

	go func() {
		defer close(done)
		var runErr error
		defer func() {
			done <- runErr
		}()

		for session := 0; session < expectedConnections; session++ {
			conn, err := ln.Accept()
			if err != nil {
				runErr = err
				return
			}
			if runErr = runAuthSMTPSession(conn, cfg, session); runErr != nil {
				return
			}
		}
	}()

	return ln.Addr().String(), done
}

func runAuthSMTPSession(conn net.Conn, cfg authServerConfig, session int) error {
	defer conn.Close()

	_, _ = conn.Write([]byte("220 localhost ESMTP ready\r\n"))
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(line)), "EHLO") {
		return unexpectedSMTPLineError("EHLO *", strings.TrimSpace(line))
	}

	if cfg.authLine != "" {
		_, _ = conn.Write([]byte("250-localhost\r\n250-AUTH " + cfg.authLine + "\r\n250 OK\r\n"))
	} else {
		_, _ = conn.Write([]byte("250 localhost\r\n"))
	}

	line, err = reader.ReadString('\n')
	if err != nil {
		return nil
	}
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(strings.ToUpper(trimmed), "AUTH ") && cfg.handleAuth != nil {
		return cfg.handleAuth(conn, reader, session, trimmed)
	}

	if strings.ToUpper(trimmed) == "QUIT" {
		_, _ = conn.Write([]byte("221 Bye\r\n"))
	}
	return nil
}

func unexpectedSMTPLineError(expected, actual string) error {
	return &smtpScriptError{expected: expected, actual: actual}
}

type smtpScriptError struct {
	expected string
	actual   string
}

func (e *smtpScriptError) Error() string {
	return "unexpected smtp line: expected " + e.expected + ", got " + e.actual
}

func splitHostPort(t *testing.T, addr string) (string, int) {
	t.Helper()
	host, portStr, err := net.SplitHostPort(addr)
	require.NoError(t, err)
	port, err := net.LookupPort("tcp", portStr)
	require.NoError(t, err)
	return host, port
}
