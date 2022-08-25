package twiliosendmail

import (
	"io/ioutil"
	"myapp/domain/sentmail"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    sentmail.EMailer
		wantErr bool
	}{
		{
			name:    "正常系_オブジェクト生成できること",
			want:    &SendGridMail{apiKey: "abc123"},
			wantErr: false,
		},
	}
	os.Setenv("API_KEY", "abc123")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSendGridConfig()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGridEMailWrap_Send(t *testing.T) {
	os.Setenv("API_KEY", "abc123")

	// 添付テスト用
	attch, err := ioutil.ReadFile("testdata/test.zip")
	if err != nil {
		t.Fatalf("添付用ファイルオープンエラー: %v", err)
	}

	type args struct {
		tos          []sentmail.EmailAddress
		ccs          []sentmail.EmailAddress
		bccs         []sentmail.EmailAddress
		from         sentmail.EmailAddress
		subject      string
		body         string
		attachment   []sentmail.Attachment
		replacements map[string]string
	}
	tests := []struct {
		name    string
		m       *SendGridMail
		args    args
		wantErr bool
	}{
		{
			name: "正常系_添付無し送信成功",
			m:    NewSendGridConfig().(*SendGrid),
			args: args{
				tos: []sentmail.EmailAddress{{
					Name:    "宛先１",
					Address: "gounittest1@sink.sendgrid.net",
				}},
				ccs:  []sentmail.EmailAddress{},
				bccs: []sentmail.EmailAddress{},
				from: sentmail.EmailAddress{
					Name:    "送り主",
					Address: "k_hirano@wadass.com",
				},
				subject:      "正常系_添付無し送信成功",
				body:         "本文:正常系_添付無し送信成功",
				attachment:   []sentmail.Attachment{},
				replacements: map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "正常系_添付無し複数宛先送信成功",
			m:    NewSendGridConfig().(*SendGrid),
			args: args{
				tos: []sentmail.EmailAddress{
					{
						Name:    "宛先２",
						Address: "gounittest2@sink.sendgrid.net",
					},
					{
						Name:    "宛先３",
						Address: "gounittest3@sink.sendgrid.net",
					},
					{
						Name:    "宛先４",
						Address: "gounittest4@sink.sendgrid.net",
					},
				},
				ccs: []sentmail.EmailAddress{
					{Name: "宛先CC２", Address: "gounittestcc2@sink.sendgrid.net"},
					{Name: "宛先CC３", Address: "gounittestcc3@sink.sendgrid.net"},
					{Name: "宛先CC４", Address: "gounittestcc4@sink.sendgrid.net"},
					{Name: "宛先CC5", Address: "gounittestcc5@sink.sendgrid.net"},
					{Name: "宛先CC6", Address: "gounittestcc6@sink.sendgrid.net"},
					{Name: "宛先CC7", Address: "gounittestcc7@sink.sendgrid.net"},
					{Name: "宛先CC8", Address: "gounittestcc8@sink.sendgrid.net"},
					{Name: "宛先CC9", Address: "gounittestcc9@sink.sendgrid.net"},
					{Name: "宛先CC10", Address: "gounittestcc10@sink.sendgrid.net"},
					{Name: "宛先CC11", Address: "gounittestcc11@sink.sendgrid.net"},
					{Name: "宛先CC12", Address: "gounittestcc12@sink.sendgrid.net"},
					{Name: "宛先CC13", Address: "gounittestcc13@sink.sendgrid.net"},
					{Name: "宛先CC14", Address: "gounittestcc14@sink.sendgrid.net"},
					{Name: "宛先CC15", Address: "gounittestcc15@sink.sendgrid.net"},
					{Name: "宛先CC16", Address: "gounittestcc16@sink.sendgrid.net"},
					{Name: "宛先CC17", Address: "gounittestcc17@sink.sendgrid.net"},
					{Name: "宛先CC18", Address: "gounittestcc18@sink.sendgrid.net"},
					{Name: "宛先CC19", Address: "gounittestcc19@sink.sendgrid.net"},
					{Name: "宛先CC20", Address: "gounittestcc20@sink.sendgrid.net"},
					{Name: "宛先CC21", Address: "gounittestcc21@sink.sendgrid.net"},
					{Name: "宛先CC22", Address: "gounittestcc22@sink.sendgrid.net"},
					{Name: "宛先CC23", Address: "gounittestcc23@sink.sendgrid.net"},
					{Name: "宛先CC24", Address: "gounittestcc24@sink.sendgrid.net"},
					{Name: "宛先CC25", Address: "gounittestcc25@sink.sendgrid.net"},
				},
				bccs: []sentmail.EmailAddress{
					{
						Name:    "宛先BCC２",
						Address: "gounittestbcc2@sink.sendgrid.net",
					},
					{
						Name:    "宛先BCC３",
						Address: "gounittestbcc3@sink.sendgrid.net",
					},
					{
						Name:    "宛先BCC４",
						Address: "gounittestbcc4@sink.sendgrid.net",
					},
				},
				from: sentmail.EmailAddress{
					Name:    "送り主",
					Address: "k_hirano@wadass.com",
				},
				subject:      "正常系_添付無し複数宛先送信",
				body:         "本文:正常系_添付無し複数宛先送信",
				attachment:   []sentmail.Attachment{},
				replacements: map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "正常系_添付無し複数宛先送信成功",
			m:    NewSendGridConfig().(*SendGrid),
			args: args{
				tos: []sentmail.EmailAddress{
					{Name: "宛先２", Address: "gounittest2@sink.sendgrid.net"},
					{Name: "宛先３", Address: "gounittest3@sink.sendgrid.net"},
				},
				ccs: []sentmail.EmailAddress{
					{Name: "宛先CC２", Address: "gounittestcc2@sink.sendgrid.net"},
				},
				bccs: []sentmail.EmailAddress{
					{Name: "宛先BCC２", Address: "gounittestbcc2@sink.sendgrid.net"},
				},
				from: sentmail.EmailAddress{
					Name:    "送り主",
					Address: "k_hirano@wadass.com",
				},
				subject: "正常系_添付無し複数宛先送信",
				body:    "本文:正常系_添付無し複数宛先送信",
				attachment: []sentmail.Attachment{{
					Data:     attch,
					FileType: "",
					FileName: "",
				}},
				replacements: map[string]string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.m.Send(
				tt.args.tos,
				tt.args.ccs,
				tt.args.bccs,
				tt.args.from,
				tt.args.subject,
				tt.args.body,
				tt.args.attachment,
				tt.args.replacements,
			); (err != nil) != tt.wantErr {
				t.Errorf("gridEMailWrap.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
