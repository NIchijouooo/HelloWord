package setting

import "testing"

func Test_crc16Xmodem(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := crc16Xmodem(tt.args.buf); got != tt.want {
				t.Errorf("crc16Xmodem() = %v, want %v", got, tt.want)
			}
		})
	}
}
