package parse

import "testing"

func TestGin(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"case-1", args{"./testdata"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Gin(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Gin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
