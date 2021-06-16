package ftp_test

import (
	"reflect"
	"testing"
	"time"

	base "github.com/jlaffaye/ftp"
	ftp "github.com/molotovtv/go-ftp"
	"github.com/molotovtv/go-ftp/mocks"
	"github.com/stretchr/testify/mock"
)

func TestFTP_GetExtensionFile(t *testing.T) {
	type fields struct {
		Addr     string
		Password string
		Timeout  time.Duration
		Username string
	}
	type args struct {
		oFile *base.Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test JSON",
			// fields: &fields{},
			args: args{
				oFile: &base.Entry{
					Name: "test.json",
					Type: base.EntryTypeFile,
					Size: 1000,
					Time: time.Now(),
				},
			},
			want: "json",
		},
		{
			name: "Test JSON.DONE",
			// fields: &fields{},
			args: args{
				oFile: &base.Entry{
					Name: "test.json.done",
					Type: base.EntryTypeFile,
					Size: 1000,
					Time: time.Now(),
				},
			},
			want: "done",
		},
		{
			name: "Test XML",
			// fields: &fields{},
			args: args{
				oFile: &base.Entry{
					Name: "test.XML",
					Type: base.EntryTypeFile,
					Size: 1000,
					Time: time.Now(),
				},
			},
			want: "xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ftp.FTP{
				Addr:     tt.fields.Addr,
				Password: tt.fields.Password,
				Timeout:  tt.fields.Timeout,
				Username: tt.fields.Username,
			}
			if got := f.GetExtensionFile(tt.args.oFile); got != tt.want {
				t.Errorf("base.getExtensionFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFTP_GetFileNameWithoutExtension(t *testing.T) {
	type fields struct {
		Addr     string
		Password string
		Timeout  time.Duration
		Username string
	}
	type args struct {
		sFileName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test JSON",
			// fields: &fields{},
			args: args{
				sFileName: "test.json",
			},
			want: "test",
		},
		{
			name: "Test JSON.DONE",
			// fields: &fields{},
			args: args{
				sFileName: "test.json.done",
			},
			want: "test.json",
		},
		{
			name: "Test jépétay-capu.prout",
			// fields: &fields{},
			args: args{
				sFileName: "jépétay-capu.prout",
			},
			want: "jépétay-capu",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ftp.FTP{
				Addr:     tt.fields.Addr,
				Password: tt.fields.Password,
				Timeout:  tt.fields.Timeout,
				Username: tt.fields.Username,
			}
			if got := f.GetFileNameWithoutExtension(tt.args.sFileName); got != tt.want {
				t.Errorf("base.GetFileNameWithoutExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFTP_List(t *testing.T) {
	type fields struct {
		Addr     string
		Password string
		Timeout  time.Duration
		Username string
		Dialer   ftp.Dialer
	}
	type args struct {
		sFolder            string
		aExtensionsAllowed []string
		sPattern           string
	}

	aFiles := getListOfFiles()
	oFtp := NewFtp(getMockOfServerConnexion(aFiles))

	aExpectedSimple := aFiles[:7]
	aExpectedByExtension := aFiles[1:4]
	aExpectedByPattern := []*base.Entry{}
	aExpectedByPattern = append(aExpectedByPattern, aFiles[1:5]...)
	aExpectedByPattern = append(aExpectedByPattern, aFiles[6])

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*base.Entry
	}{
		{
			name: "List simple",
			args: args{
				sFolder:            "",
				aExtensionsAllowed: []string{},
				sPattern:           "",
			},
			want: aExpectedSimple,
		},
		{
			name: "List filre by extension",
			args: args{
				sFolder:            "",
				aExtensionsAllowed: []string{"json", "DONE"},
				sPattern:           "",
			},
			want: aExpectedByExtension,
		},
		{
			name: "List filtre by pattern",
			args: args{
				sFolder:            "",
				aExtensionsAllowed: []string{},
				sPattern:           "test",
			},
			want: aExpectedByPattern,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oFtp.List(tt.args.sFolder, tt.args.aExtensionsAllowed, tt.args.sPattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("base.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func prepareTestFTP_list() {

}

func NewFtp(oConnexion ftp.ServerConnexion) *ftp.FTP {
	oDialer := &mocks.Dialer{}
	oDialer.On("Dial", mock.Anything).Return(oConnexion, nil)
	return ftp.New(ftp.Configuration{}, oDialer)
}

func getMockOfServerConnexion(aFiles []*base.Entry) ftp.ServerConnexion {
	oConnexion := &mocks.ServerConnexion{}
	oConnexion.On("Login", mock.Anything, mock.Anything).Return(nil)
	oConnexion.On("Quit").Return(nil)
	oConnexion.On("List", mock.Anything).Return(aFiles, nil)
	return oConnexion
}

func getListOfFiles() []*base.Entry {
	var aFiles []*base.Entry
	return append(aFiles, &base.Entry{
		Name: "te-st.XML",
		Type: base.EntryTypeFile,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: "testicule.XML.done",
		Type: base.EntryTypeFile,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: "ab-test.json",
		Type: base.EntryTypeFile,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: "test-amant.JSON.done",
		Type: base.EntryTypeFile,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: "ab-test-amant.mp4",
		Type: base.EntryTypeFile,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: "folder",
		Type: base.EntryTypeFolder,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: "folder-test",
		Type: base.EntryTypeFolder,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: ".",
		Type: base.EntryTypeFolder,
		Size: 1000,
		Time: time.Now(),
	}, &base.Entry{
		Name: "..",
		Type: base.EntryTypeFolder,
		Size: 1000,
		Time: time.Now(),
	})
}

// func TestFTP_FileSize(t *testing.T) {
// 	type fields struct {
// 		Addr     string
// 		Password string
// 		Timeout  time.Duration
// 		Username string
// 	}
// 	type args struct {
// 		src string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantS   int64
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &FTP{
// 				Addr:     tt.fields.Addr,
// 				Password: tt.fields.Password,
// 				Timeout:  tt.fields.Timeout,
// 				Username: tt.fields.Username,
// 			}
// 			gotS, err := f.FileSize(tt.args.src)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("base.FileSize() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotS != tt.wantS {
// 				t.Errorf("base.FileSize() = %v, want %v", gotS, tt.wantS)
// 			}
// 		})
// 	}
// }

func TestFTP_Exists(t *testing.T) {
	type fields struct {
		Addr     string
		Password string
		Timeout  time.Duration
		Username string
	}
	type args struct {
		sFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantB   bool
		wantErr bool
	}{
		{
			name: "Test file exists",
			// fields: fields{},
			args: args{
				sFilePath: "test.txt",
			},
			wantB:   true,
			wantErr: false,
		},
		{
			name: "Test file doesn't exists",
			// fields: fields{},
			args: args{
				sFilePath: "test2.txt",
			},
			wantB:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := ftp.New(ftp.Configuration{
				Addr:     tt.fields.Addr,
				Password: tt.fields.Password,
				Timeout:  tt.fields.Timeout,
				Username: tt.fields.Username,
			}, ftp.NewDefaultDialer())
			f.Connect()

			gotB, err := f.Exists(tt.args.sFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FTP.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotB != tt.wantB {
				t.Errorf("FTP.Exists() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
