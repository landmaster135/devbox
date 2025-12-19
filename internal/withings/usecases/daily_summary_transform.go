package usecases

// FlattenedDailySummaryResponse は CLI やエクスポート用に整形したサマリレスポンスです。
type FlattenedDailySummaryResponse struct {
	Summaries []FlattenedDailySummary `json:"summaries"`
	Timezone  string                  `json:"timezone,omitempty"`
}

// FlattenedDailySummary は計測値と活動値をフラットに表現した 1 日分のサマリです。
type FlattenedDailySummary struct {
	MeasuredAt string `json:"measured_at,omitempty"`
	*FlattenedMeasures
	*FlattenedActivity
}

// FlattenedMeasures は計測値をプレフィックス付きフィールドとして保持します。
type FlattenedMeasures struct {
	MeasuresWeightKg                   *float64 `json:"measures_weight_kg,omitempty"`
	MeasuresHeightMeter                *float64 `json:"measures_height_meter,omitempty"`
	MeasuresFatFreeMassKg              *float64 `json:"measures_fat_free_mass_kg,omitempty"`
	MeasuresFatRatioPercent            *float64 `json:"measures_fat_ratio_percent,omitempty"`
	MeasuresFatMassKg                  *float64 `json:"measures_fat_mass_kg,omitempty"`
	MeasuresDiastolicBpMmhg            *float64 `json:"measures_diastolic_bp_mmhg,omitempty"`
	MeasuresSystolicBpMmhg             *float64 `json:"measures_systolic_bp_mmhg,omitempty"`
	MeasuresHeartPulseBpm              *float64 `json:"measures_heart_pulse_bpm,omitempty"`
	MeasuresTemperatureC               *float64 `json:"measures_temperature_c,omitempty"`
	MeasuresSpo2Percent                *float64 `json:"measures_spo2_percent,omitempty"`
	MeasuresBodyTemperatureC           *float64 `json:"measures_body_temperature_c,omitempty"`
	MeasuresSkinTemperatureC           *float64 `json:"measures_skin_temperature_c,omitempty"`
	MeasuresMuscleMassKg               *float64 `json:"measures_muscle_mass_kg,omitempty"`
	MeasuresHydrationKg                *float64 `json:"measures_hydration_kg,omitempty"`
	MeasuresBoneMassKg                 *float64 `json:"measures_bone_mass_kg,omitempty"`
	MeasuresPulseWaveVelocityMPerS     *float64 `json:"measures_pulse_wave_velocity_m_per_s,omitempty"`
	MeasuresVo2MaxMlPerMinPerKg        *float64 `json:"measures_vo2max_ml_per_min_per_kg,omitempty"`
	MeasuresAtrialFibrillationResult   *float64 `json:"measures_atrial_fibrillation_result,omitempty"`
	MeasuresQrsDurationMs              *float64 `json:"measures_qrs_duration_ms,omitempty"`
	MeasuresPrDurationMs               *float64 `json:"measures_pr_duration_ms,omitempty"`
	MeasuresQtDurationMs               *float64 `json:"measures_qt_duration_ms,omitempty"`
	MeasuresQtCorrectedDurationMs      *float64 `json:"measures_qt_corrected_duration_ms,omitempty"`
	MeasuresAtrialFibrillationPpg      *float64 `json:"measures_atrial_fibrillation_ppg,omitempty"`
	MeasuresVascularAgeYears           *float64 `json:"measures_vascular_age_years,omitempty"`
	MeasuresNerveHealthConductanceFeet *float64 `json:"measures_nerve_health_conductance_feet,omitempty"`
	MeasuresExtracellularWaterKg       *float64 `json:"measures_extracellular_water_kg,omitempty"`
	MeasuresIntracellularWaterKg       *float64 `json:"measures_intracellular_water_kg,omitempty"`
	MeasuresVisceralFatIndex           *float64 `json:"measures_visceral_fat_index,omitempty"`
	MeasuresSegmentFatFreeMassKg       *float64 `json:"measures_segment_fat_free_mass_kg,omitempty"`
	MeasuresSegmentFatMassKg           *float64 `json:"measures_segment_fat_mass_kg,omitempty"`
	MeasuresSegmentMuscleMassKg        *float64 `json:"measures_segment_muscle_mass_kg,omitempty"`
	MeasuresElectrodermalActivityFeet  *float64 `json:"measures_electrodermal_activity_feet,omitempty"`
	MeasuresBasalMetabolicRate         *float64 `json:"measures_basal_metabolic_rate,omitempty"`
	MeasuresMetabolicAgeYears          *float64 `json:"measures_metabolic_age_years,omitempty"`
	MeasuresElectrochemicalSkinConduct *float64 `json:"measures_electrochemical_skin_conductance,omitempty"`
}

// FlattenedActivity は活動情報をプレフィックス付きフィールドとして保持します。
type FlattenedActivity struct {
	ActivitySteps             *int     `json:"activity_steps,omitempty"`
	ActivityDistanceMeter     *float64 `json:"activity_distance_meter,omitempty"`
	ActivityElevationMeter    *float64 `json:"activity_elevation_meter,omitempty"`
	ActivityCaloriesKcal      *float64 `json:"activity_calories_kcal,omitempty"`
	ActivityTotalCaloriesKcal *float64 `json:"activity_total_calories_kcal,omitempty"`
	ActivitySoftSeconds       *int     `json:"activity_soft_seconds,omitempty"`
	ActivityModerateSeconds   *int     `json:"activity_moderate_seconds,omitempty"`
	ActivityIntenseSeconds    *int     `json:"activity_intense_seconds,omitempty"`
	ActivityActiveSeconds     *int     `json:"activity_active_seconds,omitempty"`
	ActivityHrAverageBPM      *float64 `json:"activity_hr_average_bpm,omitempty"`
	ActivityHrMinBPM          *float64 `json:"activity_hr_min_bpm,omitempty"`
	ActivityHrMaxBPM          *float64 `json:"activity_hr_max_bpm,omitempty"`
	ActivityDeviceBrand       *int     `json:"activity_device_brand,omitempty"`
	ActivityDeviceModelID     *int     `json:"activity_device_model_id,omitempty"`
	ActivityDeviceModelName   *string  `json:"activity_device_model_name,omitempty"`
	ActivityIsTracker         *bool    `json:"activity_is_tracker,omitempty"`
}

// FlattenDailySummaryResponse は DailySummaryResponse をフラットな構造に変換します。
func FlattenDailySummaryResponse(resp *DailySummaryResponse) FlattenedDailySummaryResponse {
	if resp == nil {
		return FlattenedDailySummaryResponse{}
	}

	summaries := make([]FlattenedDailySummary, len(resp.Summaries))
	for i, summary := range resp.Summaries {
		flattened := FlattenedDailySummary{
			MeasuredAt:        summary.Date,
			FlattenedMeasures: toFlattenedMeasures(summary.Measures),
			FlattenedActivity: toFlattenedActivity(summary.Activity),
		}
		summaries[i] = flattened
	}

	return FlattenedDailySummaryResponse{
		Summaries: summaries,
		Timezone:  resp.Timezone,
	}
}

func toFlattenedMeasures(measures *DailySummaryMeasures) *FlattenedMeasures {
	if measures == nil {
		return nil
	}

	return &FlattenedMeasures{
		MeasuresWeightKg:                   measures.WeightKg,
		MeasuresHeightMeter:                measures.HeightMeter,
		MeasuresFatFreeMassKg:              measures.FatFreeMassKg,
		MeasuresFatRatioPercent:            measures.FatRatioPercent,
		MeasuresFatMassKg:                  measures.FatMassKg,
		MeasuresDiastolicBpMmhg:            measures.DiastolicBpMmhg,
		MeasuresSystolicBpMmhg:             measures.SystolicBpMmhg,
		MeasuresHeartPulseBpm:              measures.HeartPulseBpm,
		MeasuresTemperatureC:               measures.TemperatureC,
		MeasuresSpo2Percent:                measures.Spo2Percent,
		MeasuresBodyTemperatureC:           measures.BodyTemperatureC,
		MeasuresSkinTemperatureC:           measures.SkinTemperatureC,
		MeasuresMuscleMassKg:               measures.MuscleMassKg,
		MeasuresHydrationKg:                measures.HydrationKg,
		MeasuresBoneMassKg:                 measures.BoneMassKg,
		MeasuresPulseWaveVelocityMPerS:     measures.PulseWaveVelocityMPerS,
		MeasuresVo2MaxMlPerMinPerKg:        measures.Vo2MaxMlPerMinPerKg,
		MeasuresAtrialFibrillationResult:   measures.AtrialFibrillationResult,
		MeasuresQrsDurationMs:              measures.QrsDurationMs,
		MeasuresPrDurationMs:               measures.PrDurationMs,
		MeasuresQtDurationMs:               measures.QtDurationMs,
		MeasuresQtCorrectedDurationMs:      measures.QtCorrectedDurationMs,
		MeasuresAtrialFibrillationPpg:      measures.AtrialFibrillationPpg,
		MeasuresVascularAgeYears:           measures.VascularAgeYears,
		MeasuresNerveHealthConductanceFeet: measures.NerveHealthConductanceFeet,
		MeasuresExtracellularWaterKg:       measures.ExtracellularWaterKg,
		MeasuresIntracellularWaterKg:       measures.IntracellularWaterKg,
		MeasuresVisceralFatIndex:           measures.VisceralFatIndex,
		MeasuresSegmentFatFreeMassKg:       measures.SegmentFatFreeMassKg,
		MeasuresSegmentFatMassKg:           measures.SegmentFatMassKg,
		MeasuresSegmentMuscleMassKg:        measures.SegmentMuscleMassKg,
		MeasuresElectrodermalActivityFeet:  measures.ElectrodermalActivityFeet,
		MeasuresBasalMetabolicRate:         measures.BasalMetabolicRate,
		MeasuresMetabolicAgeYears:          measures.MetabolicAgeYears,
		MeasuresElectrochemicalSkinConduct: measures.ElectrochemicalSkinConduct,
	}
}

func toFlattenedActivity(activity *ActivitySummary) *FlattenedActivity {
	if activity == nil {
		return nil
	}

	return &FlattenedActivity{
		ActivitySteps:             activity.Steps,
		ActivityDistanceMeter:     activity.DistanceMeter,
		ActivityElevationMeter:    activity.ElevationMeter,
		ActivityCaloriesKcal:      activity.CaloriesKcal,
		ActivityTotalCaloriesKcal: activity.TotalCaloriesKcal,
		ActivitySoftSeconds:       activity.SoftSeconds,
		ActivityModerateSeconds:   activity.ModerateSeconds,
		ActivityIntenseSeconds:    activity.IntenseSeconds,
		ActivityActiveSeconds:     activity.ActiveSeconds,
		ActivityHrAverageBPM:      activity.HrAverageBPM,
		ActivityHrMinBPM:          activity.HrMinBPM,
		ActivityHrMaxBPM:          activity.HrMaxBPM,
		ActivityDeviceBrand:       activity.DeviceBrand,
		ActivityDeviceModelID:     activity.DeviceModelID,
		ActivityDeviceModelName:   activity.DeviceModelName,
		ActivityIsTracker:         activity.IsTracker,
	}
}
