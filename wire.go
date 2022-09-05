//go:build wireinject
// +build wireinject

package main

import (
	"SynchronizeMonorevoDeliveryDates/domain/compare"
	"SynchronizeMonorevoDeliveryDates/infrastructure/jobbook"
	"SynchronizeMonorevoDeliveryDates/infrastructure/proposition"
	"SynchronizeMonorevoDeliveryDates/infrastructure/twiliosendmail"
	"SynchronizeMonorevoDeliveryDates/presentation"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"SynchronizeMonorevoDeliveryDates/usecase/difference_extract_case"
	"SynchronizeMonorevoDeliveryDates/usecase/jobbook_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_post_case"
	"SynchronizeMonorevoDeliveryDates/usecase/report_send_case"
	"SynchronizeMonorevoDeliveryDates/usecase/reportsetting_obtain_case"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func InitializeSynchronize(
	log *zap.SugaredLogger,
	ap *appsetting_obtain_case.AppSettingDto,
	rep *reportsetting_obtain_case.ReportSettingDto,
) *presentation.SynchronizingDeliveryDate {
	wire.Build(
		presentation.Set,
		difference_extract_case.Set,
		jobbook_fetch_case.Set,
		proposition_fetch_case.Set,
		proposition_post_case.Set,
		report_send_case.Set,
		jobbook.Set,
		proposition.Set,
		twiliosendmail.Set,
		compare.Set,
	)
	return nil
}
