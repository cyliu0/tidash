package pd

import (
	"net/http"
	"testing"
)

func TestPDClient_GetTrend(t *testing.T) {
	type fields struct {
		ApiAddr string
		Client  *http.Client
	}
	tests := []struct {
		name      string
		fields    fields
		wantTrend *Trend
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "testGetTrend",
			fields: fields{
				ApiAddr: "http://127.0.0.1:32770/pd/api/v1",
				Client:  &http.Client{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdc := &PDClient{
				ApiAddr: tt.fields.ApiAddr,
				Client:  tt.fields.Client,
			}
			_, err := pdc.GetTrend()
			if (err != nil) != tt.wantErr {
				t.Errorf("PDClient.GetTrend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/*
				if !reflect.DeepEqual(gotTrend, tt.wantTrend) {
					t.Errorf("PDClient.GetTrend() = %v, want %v", gotTrend, tt.wantTrend)
				}
			*/
		})
	}
}
