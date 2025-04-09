package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RecordLog(t *testing.T) {
	type args struct {
		scriptName   string
		functionName string
		isRecording  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RecordLog(tt.args.scriptName, tt.args.functionName, tt.args.isRecording))
		})
	}
}
