package telemetry

import (
	"fmt"
	"os"
	"time"
)

type GameData struct {
	Keys    []string
	Data    map[string]float32
	RawData []byte
}

type ConverterInterface interface {
	ChannelInit(now time.Time, channel chan GameData, port int)
	Convert(now time.Time, data GameData, port int)
}

type TelemetryInterface interface {
	InitAndRun(port int) error
}

type TelemetryHandler struct {
	Telemetries map[string]TelemetryData
	Keys        []string
	Adapters    []ConverterInterface
}

type TelemetryData struct {
	Position    int
	Name        string
	DataType    string
	StartOffset int
	EndOffset   int
}

// DisplayLog Check if flag was passed
func DisplayLog(flagName string, logText any) {
	if os.Getenv("DEBUG_MODE") == flagName {
		fmt.Println(logText)
	}
}

//nolint:lll
func Telemetries() (map[string]TelemetryData, []string) {
	return map[string]TelemetryData{
			"IsRaceOn":                             {Position: 0, Name: "IsRaceOn", DataType: "S32", StartOffset: 0, EndOffset: 4},
			"TimestampMS":                          {Position: 1, Name: "TimestampMS", DataType: "U32", StartOffset: 4, EndOffset: 8},
			"EngineMaxRpm":                         {Position: 2, Name: "EngineMaxRpm", DataType: "F32", StartOffset: 8, EndOffset: 12},
			"EngineIdleRpm":                        {Position: 3, Name: "EngineIdleRpm", DataType: "F32", StartOffset: 12, EndOffset: 16},
			"CurrentEngineRpm":                     {Position: 4, Name: "CurrentEngineRpm", DataType: "F32", StartOffset: 16, EndOffset: 20},
			"AccelerationX":                        {Position: 5, Name: "AccelerationX", DataType: "F32", StartOffset: 20, EndOffset: 24},
			"AccelerationY":                        {Position: 6, Name: "AccelerationY", DataType: "F32", StartOffset: 24, EndOffset: 28},
			"AccelerationZ":                        {Position: 7, Name: "AccelerationZ", DataType: "F32", StartOffset: 28, EndOffset: 32},
			"VelocityX":                            {Position: 8, Name: "VelocityX", DataType: "F32", StartOffset: 32, EndOffset: 36},
			"VelocityY":                            {Position: 9, Name: "VelocityY", DataType: "F32", StartOffset: 36, EndOffset: 40},
			"VelocityZ":                            {Position: 10, Name: "VelocityZ", DataType: "F32", StartOffset: 40, EndOffset: 44},
			"AngularVelocityX":                     {Position: 11, Name: "AngularVelocityX", DataType: "F32", StartOffset: 44, EndOffset: 48},
			"AngularVelocityY":                     {Position: 12, Name: "AngularVelocityY", DataType: "F32", StartOffset: 48, EndOffset: 52},
			"AngularVelocityZ":                     {Position: 13, Name: "AngularVelocityZ", DataType: "F32", StartOffset: 52, EndOffset: 56},
			"Yaw":                                  {Position: 14, Name: "Yaw", DataType: "F32", StartOffset: 56, EndOffset: 60},
			"Pitch":                                {Position: 15, Name: "Pitch", DataType: "F32", StartOffset: 60, EndOffset: 64},
			"Roll":                                 {Position: 16, Name: "Roll", DataType: "F32", StartOffset: 64, EndOffset: 68},
			"NormalizedSuspensionTravelFrontLeft":  {Position: 17, Name: "NormalizedSuspensionTravelFrontLeft", DataType: "F32", StartOffset: 68, EndOffset: 72},
			"NormalizedSuspensionTravelFrontRight": {Position: 18, Name: "NormalizedSuspensionTravelFrontRight", DataType: "F32", StartOffset: 72, EndOffset: 76},
			"NormalizedSuspensionTravelRearLeft":   {Position: 19, Name: "NormalizedSuspensionTravelRearLeft", DataType: "F32", StartOffset: 76, EndOffset: 80},
			"NormalizedSuspensionTravelRearRight":  {Position: 20, Name: "NormalizedSuspensionTravelRearRight", DataType: "F32", StartOffset: 80, EndOffset: 84},
			"TireSlipRatioFrontLeft":               {Position: 21, Name: "TireSlipRatioFrontLeft", DataType: "F32", StartOffset: 84, EndOffset: 88},
			"TireSlipRatioFrontRight":              {Position: 22, Name: "TireSlipRatioFrontRight", DataType: "F32", StartOffset: 88, EndOffset: 92},
			"TireSlipRatioRearLeft":                {Position: 23, Name: "TireSlipRatioRearLeft", DataType: "F32", StartOffset: 92, EndOffset: 96},
			"TireSlipRatioRearRight":               {Position: 24, Name: "TireSlipRatioRearRight", DataType: "F32", StartOffset: 96, EndOffset: 100},
			"WheelRotationSpeedFrontLeft":          {Position: 25, Name: "WheelRotationSpeedFrontLeft", DataType: "F32", StartOffset: 100, EndOffset: 104},
			"WheelRotationSpeedFrontRight":         {Position: 26, Name: "WheelRotationSpeedFrontRight", DataType: "F32", StartOffset: 104, EndOffset: 108},
			"WheelRotationSpeedRearLeft":           {Position: 27, Name: "WheelRotationSpeedRearLeft", DataType: "F32", StartOffset: 108, EndOffset: 112},
			"WheelRotationSpeedRearRight":          {Position: 28, Name: "WheelRotationSpeedRearRight", DataType: "F32", StartOffset: 112, EndOffset: 116},
			"WheelOnRumbleStripFrontLeft":          {Position: 29, Name: "WheelOnRumbleStripFrontLeft", DataType: "S32", StartOffset: 116, EndOffset: 120},
			"WheelOnRumbleStripFrontRight":         {Position: 30, Name: "WheelOnRumbleStripFrontRight", DataType: "S32", StartOffset: 120, EndOffset: 124},
			"WheelOnRumbleStripRearLeft":           {Position: 31, Name: "WheelOnRumbleStripRearLeft", DataType: "S32", StartOffset: 124, EndOffset: 128},
			"WheelOnRumbleStripRearRight":          {Position: 32, Name: "WheelOnRumbleStripRearRight", DataType: "S32", StartOffset: 128, EndOffset: 132},
			"WheelInPuddleDepthFrontLeft":          {Position: 33, Name: "WheelInPuddleDepthFrontLeft", DataType: "F32", StartOffset: 132, EndOffset: 136},
			"WheelInPuddleDepthFrontRight":         {Position: 34, Name: "WheelInPuddleDepthFrontRight", DataType: "F32", StartOffset: 136, EndOffset: 140},
			"WheelInPuddleDepthRearLeft":           {Position: 35, Name: "WheelInPuddleDepthRearLeft", DataType: "F32", StartOffset: 140, EndOffset: 144},
			"WheelInPuddleDepthRearRight":          {Position: 36, Name: "WheelInPuddleDepthRearRight", DataType: "F32", StartOffset: 144, EndOffset: 148},
			"SurfaceRumbleFrontLeft":               {Position: 37, Name: "SurfaceRumbleFrontLeft", DataType: "F32", StartOffset: 148, EndOffset: 152},
			"SurfaceRumbleFrontRight":              {Position: 38, Name: "SurfaceRumbleFrontRight", DataType: "F32", StartOffset: 152, EndOffset: 156},
			"SurfaceRumbleRearLeft":                {Position: 39, Name: "SurfaceRumbleRearLeft", DataType: "F32", StartOffset: 156, EndOffset: 160},
			"SurfaceRumbleRearRight":               {Position: 40, Name: "SurfaceRumbleRearRight", DataType: "F32", StartOffset: 160, EndOffset: 164},
			"TireSlipAngleFrontLeft":               {Position: 41, Name: "TireSlipAngleFrontLeft", DataType: "F32", StartOffset: 164, EndOffset: 168},
			"TireSlipAngleFrontRight":              {Position: 42, Name: "TireSlipAngleFrontRight", DataType: "F32", StartOffset: 168, EndOffset: 172},
			"TireSlipAngleRearLeft":                {Position: 43, Name: "TireSlipAngleRearLeft", DataType: "F32", StartOffset: 172, EndOffset: 176},
			"TireSlipAngleRearRight":               {Position: 44, Name: "TireSlipAngleRearRight", DataType: "F32", StartOffset: 176, EndOffset: 180},
			"TireCombinedSlipFrontLeft":            {Position: 45, Name: "TireCombinedSlipFrontLeft", DataType: "F32", StartOffset: 180, EndOffset: 184},
			"TireCombinedSlipFrontRight":           {Position: 46, Name: "TireCombinedSlipFrontRight", DataType: "F32", StartOffset: 184, EndOffset: 188},
			"TireCombinedSlipRearLeft":             {Position: 47, Name: "TireCombinedSlipRearLeft", DataType: "F32", StartOffset: 188, EndOffset: 192},
			"TireCombinedSlipRearRight":            {Position: 48, Name: "TireCombinedSlipRearRight", DataType: "F32", StartOffset: 192, EndOffset: 196},
			"SuspensionTravelMetersFrontLeft":      {Position: 49, Name: "SuspensionTravelMetersFrontLeft", DataType: "F32", StartOffset: 196, EndOffset: 200},
			"SuspensionTravelMetersFrontRight":     {Position: 50, Name: "SuspensionTravelMetersFrontRight", DataType: "F32", StartOffset: 200, EndOffset: 204},
			"SuspensionTravelMetersRearLeft":       {Position: 51, Name: "SuspensionTravelMetersRearLeft", DataType: "F32", StartOffset: 204, EndOffset: 208},
			"SuspensionTravelMetersRearRight":      {Position: 52, Name: "SuspensionTravelMetersRearRight", DataType: "F32", StartOffset: 208, EndOffset: 212},
			"CarOrdinal":                           {Position: 53, Name: "CarOrdinal", DataType: "S32", StartOffset: 212, EndOffset: 216},
			"CarClass":                             {Position: 54, Name: "CarClass", DataType: "S32", StartOffset: 216, EndOffset: 220},
			"CarPerformanceIndex":                  {Position: 55, Name: "CarPerformanceIndex", DataType: "S32", StartOffset: 220, EndOffset: 224},
			"DrivetrainType":                       {Position: 56, Name: "DrivetrainType", DataType: "S32", StartOffset: 224, EndOffset: 228},
			"NumCylinders":                         {Position: 57, Name: "NumCylinders", DataType: "S32", StartOffset: 228, EndOffset: 232},
			"PositionX":                            {Position: 58, Name: "PositionX", DataType: "F32", StartOffset: 232, EndOffset: 236},
			"PositionY":                            {Position: 59, Name: "PositionY", DataType: "F32", StartOffset: 236, EndOffset: 240},
			"PositionZ":                            {Position: 60, Name: "PositionZ", DataType: "F32", StartOffset: 240, EndOffset: 244},
			"Speed":                                {Position: 61, Name: "Speed", DataType: "F32", StartOffset: 244, EndOffset: 248},
			"Power":                                {Position: 62, Name: "Power", DataType: "F32", StartOffset: 248, EndOffset: 252},
			"Torque":                               {Position: 63, Name: "Torque", DataType: "F32", StartOffset: 252, EndOffset: 256},
			"TireTempFrontLeft":                    {Position: 64, Name: "TireTempFrontLeft", DataType: "F32", StartOffset: 256, EndOffset: 260},
			"TireTempFrontRight":                   {Position: 65, Name: "TireTempFrontRight", DataType: "F32", StartOffset: 260, EndOffset: 264},
			"TireTempRearLeft":                     {Position: 66, Name: "TireTempRearLeft", DataType: "F32", StartOffset: 264, EndOffset: 268},
			"TireTempRearRight":                    {Position: 67, Name: "TireTempRearRight", DataType: "F32", StartOffset: 268, EndOffset: 272},
			"Boost":                                {Position: 68, Name: "Boost", DataType: "F32", StartOffset: 272, EndOffset: 276},
			"Fuel":                                 {Position: 69, Name: "Fuel", DataType: "F32", StartOffset: 276, EndOffset: 280},
			"DistanceTraveled":                     {Position: 70, Name: "DistanceTraveled", DataType: "F32", StartOffset: 280, EndOffset: 284},
			"BestLap":                              {Position: 71, Name: "BestLap", DataType: "F32", StartOffset: 284, EndOffset: 288},
			"LastLap":                              {Position: 72, Name: "LastLap", DataType: "F32", StartOffset: 288, EndOffset: 292},
			"CurrentLap":                           {Position: 73, Name: "CurrentLap", DataType: "F32", StartOffset: 292, EndOffset: 296},
			"CurrentRaceTime":                      {Position: 74, Name: "CurrentRaceTime", DataType: "F32", StartOffset: 296, EndOffset: 300},
			"LapNumber":                            {Position: 75, Name: "LapNumber", DataType: "U16", StartOffset: 300, EndOffset: 302},
			"RacePosition":                         {Position: 76, Name: "RacePosition", DataType: "U8", StartOffset: 302, EndOffset: 303},
			"Accel":                                {Position: 77, Name: "Accel", DataType: "U8", StartOffset: 303, EndOffset: 304},
			"Brake":                                {Position: 78, Name: "Brake", DataType: "U8", StartOffset: 304, EndOffset: 305},
			"Clutch":                               {Position: 79, Name: "Clutch", DataType: "U8", StartOffset: 305, EndOffset: 306},
			"HandBrake":                            {Position: 80, Name: "HandBrake", DataType: "U8", StartOffset: 306, EndOffset: 307},
			"Gear":                                 {Position: 81, Name: "Gear", DataType: "U8", StartOffset: 307, EndOffset: 308},
			"Steer":                                {Position: 82, Name: "Steer", DataType: "S8", StartOffset: 308, EndOffset: 309},
			"NormalizedDrivingLine":                {Position: 83, Name: "NormalizedDrivingLine", DataType: "S8", StartOffset: 309, EndOffset: 310},
			"NormalizedAIBrakeDifference":          {Position: 84, Name: "NormalizedAIBrakeDifference", DataType: "S8", StartOffset: 310, EndOffset: 311},
			"TireWearFrontLeft":                    {Position: 85, Name: "TireWearFrontLeft", DataType: "F32", StartOffset: 311, EndOffset: 315},
			"TireWearFrontRight":                   {Position: 86, Name: "TireWearFrontRight", DataType: "F32", StartOffset: 315, EndOffset: 319},
			"TireWearRearLeft":                     {Position: 87, Name: "TireWearRearLeft", DataType: "F32", StartOffset: 319, EndOffset: 323},
			"TireWearRearRight":                    {Position: 88, Name: "TireWearRearRight", DataType: "F32", StartOffset: 323, EndOffset: 327},
			"TrackOrdinal":                         {Position: 89, Name: "TrackOrdinal", DataType: "S32", StartOffset: 327, EndOffset: 331},
		}, []string{
			"IsRaceOn",
			"TimestampMS",
			"EngineMaxRpm",
			"EngineIdleRpm",
			"CurrentEngineRpm",
			"AccelerationX",
			"AccelerationY",
			"AccelerationZ",
			"VelocityX",
			"VelocityY",
			"VelocityZ",
			"AngularVelocityX",
			"AngularVelocityY",
			"AngularVelocityZ",
			"Yaw",
			"Pitch",
			"Roll",
			"NormalizedSuspensionTravelFrontLeft",
			"NormalizedSuspensionTravelFrontRight",
			"NormalizedSuspensionTravelRearLeft",
			"NormalizedSuspensionTravelRearRight",
			"TireSlipRatioFrontLeft",
			"TireSlipRatioFrontRight",
			"TireSlipRatioRearLeft",
			"TireSlipRatioRearRight",
			"WheelRotationSpeedFrontLeft",
			"WheelRotationSpeedFrontRight",
			"WheelRotationSpeedRearLeft",
			"WheelRotationSpeedRearRight",
			"WheelOnRumbleStripFrontLeft",
			"WheelOnRumbleStripFrontRight",
			"WheelOnRumbleStripRearLeft",
			"WheelOnRumbleStripRearRight",
			"WheelInPuddleDepthFrontLeft",
			"WheelInPuddleDepthFrontRight",
			"WheelInPuddleDepthRearLeft",
			"WheelInPuddleDepthRearRight",
			"SurfaceRumbleFrontLeft",
			"SurfaceRumbleFrontRight",
			"SurfaceRumbleRearLeft",
			"SurfaceRumbleRearRight",
			"TireSlipAngleFrontLeft",
			"TireSlipAngleFrontRight",
			"TireSlipAngleRearLeft",
			"TireSlipAngleRearRight",
			"TireCombinedSlipFrontLeft",
			"TireCombinedSlipFrontRight",
			"TireCombinedSlipRearLeft",
			"TireCombinedSlipRearRight",
			"SuspensionTravelMetersFrontLeft",
			"SuspensionTravelMetersFrontRight",
			"SuspensionTravelMetersRearLeft",
			"SuspensionTravelMetersRearRight",
			"CarOrdinal",
			"CarClass",
			"CarPerformanceIndex",
			"DrivetrainType",
			"NumCylinders",
			"PositionX",
			"PositionY",
			"PositionZ",
			"Speed",
			"Power",
			"Torque",
			"TireTempFrontLeft",
			"TireTempFrontRight",
			"TireTempRearLeft",
			"TireTempRearRight",
			"Boost",
			"Fuel",
			"DistanceTraveled",
			"BestLap",
			"LastLap",
			"CurrentLap",
			"CurrentRaceTime",
			"LapNumber",
			"RacePosition",
			"Accel",
			"Brake",
			"Clutch",
			"HandBrake",
			"Gear",
			"Steer",
			"NormalizedDrivingLine",
			"NormalizedAIBrakeDifference",
			"TireWearFrontLeft",
			"TireWearFrontRight",
			"TireWearRearLeft",
			"TireWearRearRight",
			"TrackOrdinal",
		}
}
