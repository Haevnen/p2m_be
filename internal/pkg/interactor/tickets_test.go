package interactor

import (
	"testing"

	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

func TestTicketManagement_parseFolderPathToGetInternalLink(t1 *testing.T) {
	type fields struct {
		userManagement   interactorinterface.UserManagementInterface
		clientManagement interactorinterface.ClientManagementInterface
		txManager        interactorinterface.TxManager
	}
	type args struct {
		root   string
		folder string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test get path",
			fields: fields{
				userManagement:   nil,
				clientManagement: nil,
				txManager:        nil,
			},
			args: args{
				root:   "/volume5",
				folder: "/volume5/CLIENTS/SAW/UPLOAD/2024/9/14/TestAuto",
			},
			want:    "/CLIENTS/SAW/UPLOAD/2024/9/14/TestAuto",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TicketManagement{
				userManagement:   tt.fields.userManagement,
				clientManagement: tt.fields.clientManagement,
				txManager:        tt.fields.txManager,
			}
			_, _, got, err := t.parseFolderPathToGetTicketMetadata(tt.args.root, tt.args.folder)
			if (err != nil) != tt.wantErr {
				t1.Errorf("parseFolderPathToGetTicketMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("parseFolderPathToGetTicketMetadata() got = %v, want %v", got, tt.want)
			}
		})
	}
}
