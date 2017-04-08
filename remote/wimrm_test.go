package remote

// We use a forked version of net/http and crypto/tls
// because the standard libs do now support renegotiation

import "testing"

func Test_WinRM(t *testing.T) {
	s := NewWinRmSession()
	config := NewConfig()
	config.Host = "192.168.0.239"
	config.UseWinRm = true
	config.User = "administrator"
	config.Password = "xxxxxx"
	s.Init(config)
	err := s.Open()
	if err != nil {
		panic(err)
	}
	s.Run("ipconfig /all")

}
