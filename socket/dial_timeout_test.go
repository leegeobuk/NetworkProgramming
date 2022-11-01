package socket

import (
    "net"
    "syscall"
    "testing"
    "time"
)

// DialTimeout is equivalent to net.DialTimeout.
// Only difference is that it uses custom net.Dialer with Control function.
func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
    d := net.Dialer{
        Control: func(_, addr string, _ syscall.RawConn) error {
            return &net.DNSError{
                Err:         "connection timed out",
                Name:        addr,
                Server:      "127.0.0.1",
                IsTimeout:   true,
                IsTemporary: true,
            }
        },
        Timeout: timeout,
    }
    return d.Dial(network, address)
}

// TestDialTimeout dials with custom timeout which overrides the default timeout set by OS.
func TestDialTimeout(t *testing.T) {
    const (
        address = "10.0.0.1:http"
        timeout = 5 * time.Second
    )

    c, err := DialTimeout("tcp", address, timeout)
    if err == nil {
        c.Close()
        t.Fatal("connection did not time out")
    }
    nErr, ok := err.(net.Error)
    if !ok {
        t.Fatal(err)
    }
    if !nErr.Timeout() {
        t.Fatal("error is not a timeout")
    }
}
