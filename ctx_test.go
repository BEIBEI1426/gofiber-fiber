package fiber

import (
	"testing"
)

func Test_Ctx_IP_Malformed_XForwardedFor(t *testing.T) {
	tests := []struct {
		name       string
		headerVal  string
		remoteIP   string
		expected   string
	}{
		{
			name:      "leading comma and space",
			headerVal: ", 203.0.113.195",
			remoteIP:  "127.0.0.1",
			expected:  "203.0.113.195",
		},
		{
			name:      "middle empty entry",
			headerVal: "203.0.113.195, , 70.41.3.18",
			remoteIP:  "127.0.0.1",
			expected:  "203.0.113.195",
		},
		{
			name:      "leading spaces and commas",
			headerVal: "  ,  203.0.113.195  , 70.41.3.18",
			remoteIP:  "127.0.0.1",
			expected:  "203.0.113.195",
		},
		{
			name:      "all empty entries fallback to remote",
			headerVal: ",,",
			remoteIP:  "10.0.0.1",
			expected:  "10.0.0.1",
		},
		{
			name:      "single valid IP",
			headerVal: "192.168.1.1",
			remoteIP:  "127.0.0.1",
			expected:  "192.168.1.1",
		},
		{
			name:      "trailing comma",
			headerVal: "203.0.113.195,",
			remoteIP:  "127.0.0.1",
			expected:  "203.0.113.195",
		},
		{
			name:      "empty header fallback",
			headerVal: "",
			remoteIP:  "10.0.0.2",
			expected:  "10.0.0.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockRequestHandler{
				headers:  map[string]string{"X-Forwarded-For": tt.headerVal},
				remoteIP: tt.remoteIP,
			}
			c := &Ctx{
				config: &Config{
					ProxyHeader:        "X-Forwarded-For",
					EnableIPValidation: true,
				},
				handler: mock,
			}

			got := c.IP()
			if got != tt.expected {
				t.Errorf("IP() = %q, want %q", got, tt.expected)
			}
		})
	}
}