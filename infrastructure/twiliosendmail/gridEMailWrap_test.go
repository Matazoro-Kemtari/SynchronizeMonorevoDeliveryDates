package twiliosendmail_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/report"
	"SynchronizeMonorevoDeliveryDates/infrastructure/twiliosendmail"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func TestSendGridMail_Send(t *testing.T) {
	err_read := godotenv.Load(`../../.env`)
	if err_read != nil {
		os.Exit(1)
	}
	logger, _ := zap.NewDevelopment()
	appcnf := appsetting_obtain_case.TestAppSettingDtoCreate(
		appsetting_obtain_case.OptSandboxMode(
			*appsetting_obtain_case.TestSandboxModeDtoCreate(
				appsetting_obtain_case.OptSandboxModeSendGrid(false),
			),
		),
	)
	sendgridConfig := twiliosendmail.NewSendGridConfig()
	type args struct {
		tos                []report.EmailAddress
		ccs                []report.EmailAddress
		bccs               []report.EmailAddress
		from               report.EmailAddress
		replyTo            report.EmailAddress
		subject            string
		editedPropositions []report.EditedProposition
		prefixReport       string
		suffixReport       string
	}
	tests := []struct {
		name    string
		m       *twiliosendmail.SendGridMail
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正常系_添付無し送信成功",
			m: twiliosendmail.NewSendGridMail(
				logger.Sugar(),
				appcnf,
				sendgridConfig,
			),
			args: args{
				tos:     []report.EmailAddress{{Name: "宛先１", Address: "gounittest1@sink.sendgrid.net"}},
				ccs:     nil,
				bccs:    nil,
				from:    report.EmailAddress{Name: "送り主", Address: "no-reply@exsample.com"},
				replyTo: report.EmailAddress{Name: "返信先", Address: "reply@exsample.com"},
				subject: "正常系_添付無し送信成功",
				editedPropositions: []report.EditedProposition{
					{
						WorkedNumber:        "99A-0001",
						DET:                 "1",
						Successful:          true,
						DeliveryDate:        time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 1, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0002",
						DET:                 "2",
						Successful:          false,
						DeliveryDate:        time.Date(2099, 2, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 2, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0003",
						DET:                 "3",
						Successful:          true,
						DeliveryDate:        time.Date(2099, 3, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 3, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0004",
						DET:                 "4",
						Successful:          false,
						DeliveryDate:        time.Date(2099, 4, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 4, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0005",
						DET:                 "5",
						Successful:          true,
						DeliveryDate:        time.Date(2099, 5, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 5, 20, 0, 0, 0, 0, time.UTC),
					},
				},
				prefixReport: "次の納期を変更した",
				suffixReport: "以上",
			},
			want:    time.Now().Format("Mon, 02 Jan 2006"),
			wantErr: false,
		},
		{
			name: "正常系_添付無し複数宛先送信成功",
			m: twiliosendmail.NewSendGridMail(
				logger.Sugar(),
				appcnf,
				sendgridConfig,
			),
			args: args{
				tos: []report.EmailAddress{
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
				ccs: []report.EmailAddress{
					{Name: "宛先CC2", Address: "gounittestcc2@sink.sendgrid.net"},
					{Name: "宛先CC3", Address: "gounittestcc3@sink.sendgrid.net"},
					{Name: "宛先CC4", Address: "gounittestcc4@sink.sendgrid.net"},
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
				bccs: []report.EmailAddress{
					{
						Name:    "宛先BCC2",
						Address: "gounittestbcc2@sink.sendgrid.net",
					},
					{
						Name:    "宛先BCC3",
						Address: "gounittestbcc3@sink.sendgrid.net",
					},
					{
						Name:    "宛先BCC4",
						Address: "gounittestbcc4@sink.sendgrid.net",
					},
				},
				from: report.EmailAddress{
					Name:    "送り主",
					Address: "no-reply@exsample.com",
				},
				replyTo: report.EmailAddress{Name: "返信先", Address: "reply@exsample.com"},
				subject: "正常系_添付無し複数宛先送信",
				editedPropositions: []report.EditedProposition{
					{
						WorkedNumber:        "99A-0001",
						DET:                 "1",
						Successful:          true,
						DeliveryDate:        time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 1, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0002",
						DET:                 "2",
						Successful:          false,
						DeliveryDate:        time.Date(2099, 2, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 2, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0003",
						DET:                 "3",
						Successful:          true,
						DeliveryDate:        time.Date(2099, 3, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 3, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0004",
						DET:                 "4",
						Successful:          false,
						DeliveryDate:        time.Date(2099, 4, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 4, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber:        "99A-0005",
						DET:                 "5",
						Successful:          true,
						DeliveryDate:        time.Date(2099, 5, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2099, 5, 20, 0, 0, 0, 0, time.UTC),
					},
				},
				prefixReport: "次の納期を変更した",
				suffixReport: "以上",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Send(tt.args.tos, tt.args.ccs, tt.args.bccs, tt.args.from, tt.args.replyTo, tt.args.subject, tt.args.editedPropositions, tt.args.prefixReport, tt.args.suffixReport)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendGridMail.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("SendGridMail.Send() = %v, want %v", got, tt.want)
			}
		})
	}
}
