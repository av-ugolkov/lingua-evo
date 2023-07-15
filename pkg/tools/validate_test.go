package tools

import "testing"

func TestIsEmailValid(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Validate email 1", args: struct{ e string }{e: "qwer"}, want: false},
		{name: "Validate email 2", args: struct{ e string }{e: "qwer@"}, want: false},
		{name: "Validate email 3", args: struct{ e string }{e: "@asdf"}, want: false},
		{name: "Validate email 4", args: struct{ e string }{e: "asdf@asdf"}, want: false},
		{name: "Validate email 5", args: struct{ e string }{e: "asdf$@asdf.ru"}, want: false},
		{name: "Validate email 6", args: struct{ e string }{e: "asdf@#asdf.ru"}, want: false},
		{name: "Validate email 7", args: struct{ e string }{e: "adsf@asdf.ru"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmailValid(tt.args.e); got != tt.want {
				t.Errorf("IsEmailValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
