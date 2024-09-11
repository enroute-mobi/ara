package model

import (
	"errors"
	"fmt"

	"bitbucket.org/enroute-mobi/ara/gtfs"
)

type SituationSeverity string

const (
	SituationSeverityNoImpact   SituationSeverity = "noImpact"
	SituationSeverityVerySlight SituationSeverity = "verySlight"
	SituationSeveritySlight     SituationSeverity = "slight"
	SituationSeverityNormal     SituationSeverity = "normal"
	SituationSeveritySevere     SituationSeverity = "severe"
	SituationSeverityUndefinded SituationSeverity = "undefined"
	SituationSeverityUnknown    SituationSeverity = "unknown"
	SituationSeverityVerySevere SituationSeverity = "verySevere"
)

func (s *SituationSeverity) FromProto(value interface{}) error {
	switch v := value.(type) {
	case gtfs.Alert_SeverityLevel:
		switch v {
		case gtfs.Alert_UNKNOWN_SEVERITY:
			*s = SituationSeverityUnknown
		case gtfs.Alert_INFO:
			*s = SituationSeverityNoImpact
		case gtfs.Alert_WARNING:
			*s = SituationSeverityNormal
		case gtfs.Alert_SEVERE:
			*s = SituationSeveritySevere
		default:
			return errors.New("invalid value")
		}
	default:
		return fmt.Errorf("unsupported value %T", value)
	}
	return nil
}

func (s SituationSeverity) ToProto(dest interface{}) error {
	switch v := dest.(type) {
	case *gtfs.Alert_SeverityLevel:
		switch s {
		case SituationSeverityUnknown:
			*v = gtfs.Alert_UNKNOWN_SEVERITY
		case SituationSeverityNoImpact:
			*v = gtfs.Alert_INFO
		case SituationSeverityNormal:
			*v = gtfs.Alert_WARNING
		case SituationSeveritySevere:
			*v = gtfs.Alert_SEVERE
		default:
			*v = gtfs.Default_Alert_SeverityLevel
		}
	default:
		return fmt.Errorf("unsupported destination %T", dest)
	}
	return nil
}

func (severity *SituationSeverity) FromString(s string) error {
	switch SituationSeverity(s) {
	case SituationSeverityNoImpact:
		fallthrough
	case SituationSeverityVerySlight:
		fallthrough
	case SituationSeveritySlight:
		fallthrough
	case SituationSeverityNormal:
		fallthrough
	case SituationSeveritySevere:
		fallthrough
	case SituationSeverityUndefinded:
		fallthrough
	case SituationSeverityUnknown:
		fallthrough
	case SituationSeverityVerySevere:
		*severity = SituationSeverity(s)
		return nil
	}
	return fmt.Errorf("invalid severity %s", s)
}

type SituationProgress string

const (
	SituationProgressDraft           SituationProgress = "draft"
	SituationProgressPendingApproval SituationProgress = "pendingApproval"
	SituationProgressApprovedDraft   SituationProgress = "approvedDraft"
	SituationProgressOpend           SituationProgress = "open"
	SituationProgressPublished       SituationProgress = "published"
	SituationProgressClosing         SituationProgress = "closing"
	SituationProgressClosed          SituationProgress = "closed"
)

func (progress *SituationProgress) FromString(s string) error {
	switch SituationProgress(s) {
	case SituationProgressDraft:
		fallthrough
	case SituationProgressPendingApproval:
		fallthrough
	case SituationProgressApprovedDraft:
		fallthrough
	case SituationProgressOpend:
		fallthrough
	case SituationProgressPublished:
		fallthrough
	case SituationProgressClosing:
		fallthrough
	case SituationProgressClosed:
		*progress = SituationProgress(s)
		return nil
	}
	return fmt.Errorf("invalid progress %s", s)
}

type SituationAlertCause string

const (
	SituationAlertCauseAccident                          SituationAlertCause = "accident"
	SituationAlertCauseAirraid                           SituationAlertCause = "airRaid"
	SituationAlertCauseAltercation                       SituationAlertCause = "altercation"
	SituationAlertCauseAnimalontheline                   SituationAlertCause = "animalOnTheLine"
	SituationAlertCauseAsphalting                        SituationAlertCause = "asphalting"
	SituationAlertCauseAssault                           SituationAlertCause = "assault"
	SituationAlertCauseAttack                            SituationAlertCause = "attack"
	SituationAlertCauseAvalanches                        SituationAlertCause = "avalanches"
	SituationAlertCauseAwaitingapproach                  SituationAlertCause = "awaitingApproach"
	SituationAlertCauseAwaitingoncomingvehicle           SituationAlertCause = "awaitingOncomingVehicle"
	SituationAlertCauseAwaitingshuttle                   SituationAlertCause = "awaitingShuttle"
	SituationAlertCauseBlizzardconditions                SituationAlertCause = "blizzardConditions"
	SituationAlertCauseBoardingdelay                     SituationAlertCause = "boardingDelay"
	SituationAlertCauseBombalert                         SituationAlertCause = "bombAlert"
	SituationAlertCauseBombdisposal                      SituationAlertCause = "bombDisposal"
	SituationAlertCauseBombexplosion                     SituationAlertCause = "bombExplosion"
	SituationAlertCauseBordercontrol                     SituationAlertCause = "borderControl"
	SituationAlertCauseBreakdown                         SituationAlertCause = "breakDown"
	SituationAlertCauseBridgedamage                      SituationAlertCause = "bridgeDamage"
	SituationAlertCauseBridgestrike                      SituationAlertCause = "bridgeStrike"
	SituationAlertCauseBrokenrail                        SituationAlertCause = "brokenRail"
	SituationAlertCauseCablefire                         SituationAlertCause = "cableFire"
	SituationAlertCauseCabletheft                        SituationAlertCause = "cableTheft"
	SituationAlertCauseChangeincarriages                 SituationAlertCause = "changeInCarriages"
	SituationAlertCauseCivilemergency                    SituationAlertCause = "civilEmergency"
	SituationAlertCauseClosedformaintenance              SituationAlertCause = "closedForMaintenance"
	SituationAlertCauseCollision                         SituationAlertCause = "collision"
	SituationAlertCauseCongestion                        SituationAlertCause = "congestion"
	SituationAlertCauseConstructionwork                  SituationAlertCause = "constructionWork"
	SituationAlertCauseContractorstaffinjury             SituationAlertCause = "contractorStaffInjury"
	SituationAlertCauseDefectivecctv                     SituationAlertCause = "defectiveCctv"
	SituationAlertCauseDefectivefirealarmequipment       SituationAlertCause = "defectiveFireAlarmEquipment"
	SituationAlertCauseDefectiveplatformedgedoors        SituationAlertCause = "defectivePlatformEdgeDoors"
	SituationAlertCauseDefectivepublicannouncementsystem SituationAlertCause = "defectivePublicAnnouncementSystem"
	SituationAlertCauseDefectivesecuritysystem           SituationAlertCause = "defectiveSecuritySystem"
	SituationAlertCauseDefectivetrain                    SituationAlertCause = "defectiveTrain"
	SituationAlertCauseDefectivevehicle                  SituationAlertCause = "defectiveVehicle"
	SituationAlertCauseDeicingwork                       SituationAlertCause = "deicingWork"
	SituationAlertCauseDemonstration                     SituationAlertCause = "demonstration"
	SituationAlertCauseDerailment                        SituationAlertCause = "derailment"
	SituationAlertCauseDoorfailure                       SituationAlertCause = "doorFailure"
	SituationAlertCauseDriftingsnow                      SituationAlertCause = "driftingSnow"
	SituationAlertCauseEarthquakedamage                  SituationAlertCause = "earthquakeDamage"
	SituationAlertCauseEmergencybrake                    SituationAlertCause = "emergencyBrake"
	SituationAlertCauseEmergencyengineeringwork          SituationAlertCause = "emergencyEngineeringWork"
	SituationAlertCauseEmergencymedicalservices          SituationAlertCause = "emergencyMedicalServices"
	SituationAlertCauseEmergencyservices                 SituationAlertCause = "emergencyServices"
	SituationAlertCauseEmergencyservicescall             SituationAlertCause = "emergencyServicesCall"
	SituationAlertCauseEnginefailure                     SituationAlertCause = "engineFailure"
	SituationAlertCauseEscalatorfailure                  SituationAlertCause = "escalatorFailure"
	SituationAlertCauseEvacuation                        SituationAlertCause = "evacuation"
	SituationAlertCauseExplosion                         SituationAlertCause = "explosion"
	SituationAlertCauseExplosionhazard                   SituationAlertCause = "explosionHazard"
	SituationAlertCauseFallenleaves                      SituationAlertCause = "fallenLeaves"
	SituationAlertCauseFallentree                        SituationAlertCause = "fallenTree"
	SituationAlertCauseFallentreeontheline               SituationAlertCause = "fallenTreeOnTheLine"
	SituationAlertCauseFatality                          SituationAlertCause = "fatality"
	SituationAlertCauseFilterblockade                    SituationAlertCause = "filterBlockade"
	SituationAlertCauseFire                              SituationAlertCause = "fire"
	SituationAlertCauseFireatstation                     SituationAlertCause = "fireAtStation"
	SituationAlertCauseFireatthestation                  SituationAlertCause = "fireAtTheStation"
	SituationAlertCauseFirebrigadeorder                  SituationAlertCause = "fireBrigadeOrder"
	SituationAlertCauseFirebrigadesafetychecks           SituationAlertCause = "fireBrigadeSafetyChecks"
	SituationAlertCauseFirerun                           SituationAlertCause = "fireRun"
	SituationAlertCauseFlashfloods                       SituationAlertCause = "flashFloods"
	SituationAlertCauseFlooding                          SituationAlertCause = "flooding"
	SituationAlertCauseFog                               SituationAlertCause = "fog"
	SituationAlertCauseForeigndisturbances               SituationAlertCause = "foreignDisturbances"
	SituationAlertCauseFrozen                            SituationAlertCause = "frozen"
	SituationAlertCauseFuelproblem                       SituationAlertCause = "fuelProblem"
	SituationAlertCauseFuelshortage                      SituationAlertCause = "fuelShortage"
	SituationAlertCauseGangwayproblem                    SituationAlertCause = "gangwayProblem"
	SituationAlertCauseGlazedfrost                       SituationAlertCause = "glazedFrost"
	SituationAlertCauseGrassfire                         SituationAlertCause = "grassFire"
	SituationAlertCauseGunfireonroadway                  SituationAlertCause = "gunfireOnRoadway"
	SituationAlertCauseHail                              SituationAlertCause = "hail"
	SituationAlertCauseHeavyrain                         SituationAlertCause = "heavyRain"
	SituationAlertCauseHeavysnowfall                     SituationAlertCause = "heavySnowfall"
	SituationAlertCauseHeavytraffic                      SituationAlertCause = "heavyTraffic"
	SituationAlertCauseHightemperatures                  SituationAlertCause = "highTemperatures"
	SituationAlertCauseHightide                          SituationAlertCause = "highTide"
	SituationAlertCauseHighwaterlevel                    SituationAlertCause = "highWaterLevel"
	SituationAlertCauseHoliday                           SituationAlertCause = "holiday"
	SituationAlertCauseIce                               SituationAlertCause = "ice"
	SituationAlertCauseIcedrift                          SituationAlertCause = "iceDrift"
	SituationAlertCauseIceoncarriages                    SituationAlertCause = "iceOnCarriages"
	SituationAlertCauseIceonrailway                      SituationAlertCause = "iceOnRailway"
	SituationAlertCauseIllvehicleoccupants               SituationAlertCause = "illVehicleOccupants"
	SituationAlertCauseIncident                          SituationAlertCause = "incident"
	SituationAlertCauseIndustrialaction                  SituationAlertCause = "industrialAction"
	SituationAlertCauseInsufficientdemand                SituationAlertCause = "insufficientDemand"
	SituationAlertCauseLackofoperationalstock            SituationAlertCause = "lackOfOperationalStock"
	SituationAlertCauseLandslide                         SituationAlertCause = "landslide"
	SituationAlertCauseLatefinishtoengineeringwork       SituationAlertCause = "lateFinishToEngineeringWork"
	SituationAlertCauseLeaderboardfailure                SituationAlertCause = "leaderBoardFailure"
	SituationAlertCauseLevelcrossingaccident             SituationAlertCause = "levelCrossingAccident"
	SituationAlertCauseLevelcrossingblocked              SituationAlertCause = "levelCrossingBlocked"
	SituationAlertCauseLevelcrossingfailure              SituationAlertCause = "levelCrossingFailure"
	SituationAlertCauseLevelcrossingincident             SituationAlertCause = "levelCrossingIncident"
	SituationAlertCauseLiftfailure                       SituationAlertCause = "liftFailure"
	SituationAlertCauseLightingfailure                   SituationAlertCause = "lightingFailure"
	SituationAlertCauseLightningstrike                   SituationAlertCause = "lightningStrike"
	SituationAlertCauseLinesidefire                      SituationAlertCause = "linesideFire"
	SituationAlertCauseLogisticproblems                  SituationAlertCause = "logisticProblems"
	SituationAlertCauseLowtide                           SituationAlertCause = "lowTide"
	SituationAlertCauseLowwaterlevel                     SituationAlertCause = "lowWaterLevel"
	SituationAlertCauseLuggagecarouselproblem            SituationAlertCause = "luggageCarouselProblem"
	SituationAlertCauseMaintenancework                   SituationAlertCause = "maintenanceWork"
	SituationAlertCauseMarch                             SituationAlertCause = "march"
	SituationAlertCauseMiscellaneous                     SituationAlertCause = "miscellaneous"
	SituationAlertCauseMudslide                          SituationAlertCause = "mudslide"
	SituationAlertCauseNearmiss                          SituationAlertCause = "nearMiss"
	SituationAlertCauseObjectontheline                   SituationAlertCause = "objectOnTheLine"
	SituationAlertCauseOperatorceasedtrading             SituationAlertCause = "operatorCeasedTrading"
	SituationAlertCauseOperatorsuspended                 SituationAlertCause = "operatorSuspended"
	SituationAlertCauseOvercrowded                       SituationAlertCause = "overcrowded"
	SituationAlertCauseOverheadobstruction               SituationAlertCause = "overheadObstruction"
	SituationAlertCauseOverheadwirefailure               SituationAlertCause = "overheadWireFailure"
	SituationAlertCauseOvertaking                        SituationAlertCause = "overtaking"
	SituationAlertCausePassengeraction                   SituationAlertCause = "passengerAction"
	SituationAlertCausePassengersblockingdoors           SituationAlertCause = "passengersBlockingDoors"
	SituationAlertCausePaving                            SituationAlertCause = "paving"
	SituationAlertCausePersonhitbytrain                  SituationAlertCause = "personHitByTrain"
	SituationAlertCausePersonhitbyvehicle                SituationAlertCause = "personHitByVehicle"
	SituationAlertCausePersonillonvehicle                SituationAlertCause = "personIllOnVehicle"
	SituationAlertCausePersonontheline                   SituationAlertCause = "personOnTheLine"
	SituationAlertCausePersonundertrain                  SituationAlertCause = "personUnderTrain"
	SituationAlertCausePointsfailure                     SituationAlertCause = "pointsFailure"
	SituationAlertCausePointsproblem                     SituationAlertCause = "pointsProblem"
	SituationAlertCausePoliceactivity                    SituationAlertCause = "policeActivity"
	SituationAlertCausePoliceorder                       SituationAlertCause = "policeOrder"
	SituationAlertCausePoorrailconditions                SituationAlertCause = "poorRailConditions"
	SituationAlertCausePoorweather                       SituationAlertCause = "poorWeather"
	SituationAlertCausePowerproblem                      SituationAlertCause = "powerProblem"
	SituationAlertCausePrecedingvehicle                  SituationAlertCause = "precedingVehicle"
	SituationAlertCausePreviousdisturbances              SituationAlertCause = "previousDisturbances"
	SituationAlertCauseProblemsatborderpost              SituationAlertCause = "problemsAtBorderPost"
	SituationAlertCauseProblemsatcustomspost             SituationAlertCause = "problemsAtCustomsPost"
	SituationAlertCauseProblemsonlocalroad               SituationAlertCause = "problemsOnLocalRoad"
	SituationAlertCauseProcession                        SituationAlertCause = "procession"
	SituationAlertCauseProvisiondelay                    SituationAlertCause = "provisionDelay"
	SituationAlertCausePublicdisturbance                 SituationAlertCause = "publicDisturbance"
	SituationAlertCauseRailwaycrime                      SituationAlertCause = "railwayCrime"
	SituationAlertCauseRepairwork                        SituationAlertCause = "repairWork"
	SituationAlertCauseRiskofavalanches                  SituationAlertCause = "riskOfAvalanches"
	SituationAlertCauseRiskofflooding                    SituationAlertCause = "riskOfFlooding"
	SituationAlertCauseRiskoflandslide                   SituationAlertCause = "riskOfLandslide"
	SituationAlertCauseRoadclosed                        SituationAlertCause = "roadClosed"
	SituationAlertCauseRoadmaintenance                   SituationAlertCause = "roadMaintenance"
	SituationAlertCauseRoadwaydamage                     SituationAlertCause = "roadwayDamage"
	SituationAlertCauseRoadworks                         SituationAlertCause = "roadworks"
	SituationAlertCauseRockfalls                         SituationAlertCause = "rockfalls"
	SituationAlertCauseRoughsea                          SituationAlertCause = "roughSea"
	SituationAlertCauseRouteblockage                     SituationAlertCause = "routeBlockage"
	SituationAlertCauseRoutediversion                    SituationAlertCause = "routeDiversion"
	SituationAlertCauseSabotage                          SituationAlertCause = "sabotage"
	SituationAlertCauseSafetyviolation                   SituationAlertCause = "safetyViolation"
	SituationAlertCauseSecurityalert                     SituationAlertCause = "securityAlert"
	SituationAlertCauseSecurityincident                  SituationAlertCause = "securityIncident"
	SituationAlertCauseServicedisruption                 SituationAlertCause = "serviceDisruption"
	SituationAlertCauseServicefailure                    SituationAlertCause = "serviceFailure"
	SituationAlertCauseServiceindicatorfailure           SituationAlertCause = "serviceIndicatorFailure"
	SituationAlertCauseSeweragemaintenance               SituationAlertCause = "sewerageMaintenance"
	SituationAlertCauseSeweroverflow                     SituationAlertCause = "sewerOverflow"
	SituationAlertCauseSightseersobstructingaccess       SituationAlertCause = "sightseersObstructingAccess"
	SituationAlertCauseSignalandswitchfailure            SituationAlertCause = "signalAndSwitchFailure"
	SituationAlertCauseSignalfailure                     SituationAlertCause = "signalFailure"
	SituationAlertCauseSignalpassedatdanger              SituationAlertCause = "signalPassedAtDanger"
	SituationAlertCauseSignalproblem                     SituationAlertCause = "signalProblem"
	SituationAlertCauseSleet                             SituationAlertCause = "sleet"
	SituationAlertCauseSlipperiness                      SituationAlertCause = "slipperiness"
	SituationAlertCauseSlipperytrack                     SituationAlertCause = "slipperyTrack"
	SituationAlertCauseSmokedetectedonvehicle            SituationAlertCause = "smokeDetectedOnVehicle"
	SituationAlertCauseSpecialevent                      SituationAlertCause = "specialEvent"
	SituationAlertCauseSpeedrestrictions                 SituationAlertCause = "speedRestrictions"
	SituationAlertCauseStaffabsence                      SituationAlertCause = "staffAbsence"
	SituationAlertCauseStaffassault                      SituationAlertCause = "staffAssault"
	SituationAlertCauseStaffinjury                       SituationAlertCause = "staffInjury"
	SituationAlertCauseStaffinwrongplace                 SituationAlertCause = "staffInWrongPlace"
	SituationAlertCauseStaffshortage                     SituationAlertCause = "staffShortage"
	SituationAlertCauseStaffsickness                     SituationAlertCause = "staffSickness"
	SituationAlertCauseStationoverrun                    SituationAlertCause = "stationOverrun"
	SituationAlertCauseStormconditions                   SituationAlertCause = "stormConditions"
	SituationAlertCauseStormdamage                       SituationAlertCause = "stormDamage"
	SituationAlertCauseStrongwinds                       SituationAlertCause = "strongWinds"
	SituationAlertCauseSubsidence                        SituationAlertCause = "subsidence"
	SituationAlertCauseSuspectvehicle                    SituationAlertCause = "suspectVehicle"
	SituationAlertCauseSwingbridgefailure                SituationAlertCause = "swingBridgeFailure"
	SituationAlertCauseTechnicalproblem                  SituationAlertCause = "technicalProblem"
	SituationAlertCauseTelephonedthreat                  SituationAlertCause = "telephonedThreat"
	SituationAlertCauseTerroristincident                 SituationAlertCause = "terroristIncident"
	SituationAlertCauseTheft                             SituationAlertCause = "theft"
	SituationAlertCauseTicketingsystemnotavailable       SituationAlertCause = "ticketingSystemNotAvailable"
	SituationAlertCauseTidalrestrictions                 SituationAlertCause = "tidalRestrictions"
	SituationAlertCauseTrackcircuitproblem               SituationAlertCause = "trackCircuitProblem"
	SituationAlertCauseTractionfailure                   SituationAlertCause = "tractionFailure"
	SituationAlertCauseTrafficmanagementsystemfailure    SituationAlertCause = "trafficManagementSystemFailure"
	SituationAlertCauseTraincoupling                     SituationAlertCause = "trainCoupling"
	SituationAlertCauseTraindoor                         SituationAlertCause = "trainDoor"
	SituationAlertCauseTrainstruckanimal                 SituationAlertCause = "trainStruckAnimal"
	SituationAlertCauseTrainstruckobject                 SituationAlertCause = "trainStruckObject"
	SituationAlertCauseTrainwarningsystemproblem         SituationAlertCause = "trainWarningSystemProblem"
	SituationAlertCauseUnattendedbag                     SituationAlertCause = "unattendedBag"
	SituationAlertCauseUndefinedalertcause               SituationAlertCause = "undefinedAlertCause"
	SituationAlertCauseUndefinedenvironmentalproblem     SituationAlertCause = "undefinedEnvironmentalProblem"
	SituationAlertCauseUndefinedequipmentproblem         SituationAlertCause = "undefinedEquipmentProblem"
	SituationAlertCauseUndefinedpersonnelproblem         SituationAlertCause = "undefinedPersonnelProblem"
	SituationAlertCauseUndefinedproblem                  SituationAlertCause = "undefinedProblem"
	SituationAlertCauseUnknown                           SituationAlertCause = "unknown"
	SituationAlertCauseUnofficialindustrialaction        SituationAlertCause = "unofficialIndustrialAction"
	SituationAlertCauseUnscheduledconstructionwork       SituationAlertCause = "unscheduledConstructionWork"
	SituationAlertCauseVandalism                         SituationAlertCause = "vandalism"
	SituationAlertCauseVegetation                        SituationAlertCause = "vegetation"
	SituationAlertCauseVehicleblockingtrack              SituationAlertCause = "vehicleBlockingTrack"
	SituationAlertCauseVehiclefailure                    SituationAlertCause = "vehicleFailure"
	SituationAlertCauseVehicleontheline                  SituationAlertCause = "vehicleOnTheLine"
	SituationAlertCauseVehiclestruckanimal               SituationAlertCause = "vehicleStruckAnimal"
	SituationAlertCauseVehiclestruckobject               SituationAlertCause = "vehicleStruckObject"
	SituationAlertCauseViaductfailure                    SituationAlertCause = "viaductFailure"
	SituationAlertCauseWaitingfortransferpassengers      SituationAlertCause = "waitingForTransferPassengers"
	SituationAlertCauseWaterlogged                       SituationAlertCause = "waterlogged"
	SituationAlertCauseWheelimpactload                   SituationAlertCause = "wheelImpactLoad"
	SituationAlertCauseWheelproblem                      SituationAlertCause = "wheelProblem"
	SituationAlertCauseWildlandfire                      SituationAlertCause = "wildlandFire"
	SituationAlertCauseWorktorule                        SituationAlertCause = "workToRule"
)

func (ac *SituationAlertCause) FromProto(value interface{}) error {
	switch v := value.(type) {
	case gtfs.Alert_Cause:
		switch v {
		case gtfs.Alert_UNKNOWN_CAUSE:
			*ac = SituationAlertCauseUnknown
		case gtfs.Alert_OTHER_CAUSE:
			*ac = SituationAlertCauseMiscellaneous
		case gtfs.Alert_TECHNICAL_PROBLEM:
			*ac = SituationAlertCauseTechnicalproblem
		case gtfs.Alert_STRIKE:
			*ac = SituationAlertCauseIndustrialaction
		case gtfs.Alert_DEMONSTRATION:
			*ac = SituationAlertCauseDemonstration
		case gtfs.Alert_ACCIDENT:
			*ac = SituationAlertCauseAccident
		case gtfs.Alert_HOLIDAY:
			*ac = SituationAlertCauseHoliday
		case gtfs.Alert_WEATHER:
			*ac = SituationAlertCausePoorweather
		case gtfs.Alert_MAINTENANCE:
			*ac = SituationAlertCauseMaintenancework
		case gtfs.Alert_CONSTRUCTION:
			*ac = SituationAlertCauseConstructionwork
		case gtfs.Alert_POLICE_ACTIVITY:
			*ac = SituationAlertCausePoliceactivity
		case gtfs.Alert_MEDICAL_EMERGENCY:
			*ac = SituationAlertCauseEmergencymedicalservices
		default:
			return errors.New("invalid value")
		}
	default:
		return fmt.Errorf("unsupported value %T", value)

	}

	return nil
}

func (ac SituationAlertCause) ToProto(dest interface{}) error {
	switch v := dest.(type) {
	case *gtfs.Alert_Cause:
		switch ac {
		case SituationAlertCauseUnknown:
			*v = gtfs.Alert_UNKNOWN_CAUSE
		case SituationAlertCauseMiscellaneous:
			*v = gtfs.Alert_OTHER_CAUSE
		case SituationAlertCauseTechnicalproblem:
			*v = gtfs.Alert_TECHNICAL_PROBLEM
		case SituationAlertCauseIndustrialaction:
			*v = gtfs.Alert_STRIKE
		case SituationAlertCauseDemonstration:
			*v = gtfs.Alert_DEMONSTRATION
		case SituationAlertCauseAccident:
			*v = gtfs.Alert_ACCIDENT
		case SituationAlertCauseHoliday:
			*v = gtfs.Alert_HOLIDAY
		case SituationAlertCausePoorweather:
			*v = gtfs.Alert_WEATHER
		case SituationAlertCauseMaintenancework:
			*v = gtfs.Alert_MAINTENANCE
		case SituationAlertCauseConstructionwork:
			*v = gtfs.Alert_CONSTRUCTION
		case SituationAlertCausePoliceactivity:
			*v = gtfs.Alert_POLICE_ACTIVITY
		case SituationAlertCauseEmergencymedicalservices:
			*v = gtfs.Alert_MEDICAL_EMERGENCY
		default:
			*v = gtfs.Default_Alert_Cause
		}
	default:
		return fmt.Errorf("unsupported destination %T", dest)
	}
	return nil
}

func (alertCause *SituationAlertCause) FromString(s string) error {
	switch SituationAlertCause(s) {
	case SituationAlertCauseAccident:
		fallthrough
	case SituationAlertCauseAirraid:
		fallthrough
	case SituationAlertCauseAltercation:
		fallthrough
	case SituationAlertCauseAnimalontheline:
		fallthrough
	case SituationAlertCauseAsphalting:
		fallthrough
	case SituationAlertCauseAssault:
		fallthrough
	case SituationAlertCauseAttack:
		fallthrough
	case SituationAlertCauseAvalanches:
		fallthrough
	case SituationAlertCauseAwaitingapproach:
		fallthrough
	case SituationAlertCauseAwaitingoncomingvehicle:
		fallthrough
	case SituationAlertCauseAwaitingshuttle:
		fallthrough
	case SituationAlertCauseBlizzardconditions:
		fallthrough
	case SituationAlertCauseBoardingdelay:
		fallthrough
	case SituationAlertCauseBombalert:
		fallthrough
	case SituationAlertCauseBombdisposal:
		fallthrough
	case SituationAlertCauseBombexplosion:
		fallthrough
	case SituationAlertCauseBordercontrol:
		fallthrough
	case SituationAlertCauseBreakdown:
		fallthrough
	case SituationAlertCauseBridgedamage:
		fallthrough
	case SituationAlertCauseBridgestrike:
		fallthrough
	case SituationAlertCauseBrokenrail:
		fallthrough
	case SituationAlertCauseCablefire:
		fallthrough
	case SituationAlertCauseCabletheft:
		fallthrough
	case SituationAlertCauseChangeincarriages:
		fallthrough
	case SituationAlertCauseCivilemergency:
		fallthrough
	case SituationAlertCauseClosedformaintenance:
		fallthrough
	case SituationAlertCauseCollision:
		fallthrough
	case SituationAlertCauseCongestion:
		fallthrough
	case SituationAlertCauseConstructionwork:
		fallthrough
	case SituationAlertCauseContractorstaffinjury:
		fallthrough
	case SituationAlertCauseDefectivecctv:
		fallthrough
	case SituationAlertCauseDefectivefirealarmequipment:
		fallthrough
	case SituationAlertCauseDefectiveplatformedgedoors:
		fallthrough
	case SituationAlertCauseDefectivepublicannouncementsystem:
		fallthrough
	case SituationAlertCauseDefectivesecuritysystem:
		fallthrough
	case SituationAlertCauseDefectivetrain:
		fallthrough
	case SituationAlertCauseDefectivevehicle:
		fallthrough
	case SituationAlertCauseDeicingwork:
		fallthrough
	case SituationAlertCauseDemonstration:
		fallthrough
	case SituationAlertCauseDerailment:
		fallthrough
	case SituationAlertCauseDoorfailure:
		fallthrough
	case SituationAlertCauseDriftingsnow:
		fallthrough
	case SituationAlertCauseEarthquakedamage:
		fallthrough
	case SituationAlertCauseEmergencybrake:
		fallthrough
	case SituationAlertCauseEmergencyengineeringwork:
		fallthrough
	case SituationAlertCauseEmergencymedicalservices:
		fallthrough
	case SituationAlertCauseEmergencyservices:
		fallthrough
	case SituationAlertCauseEmergencyservicescall:
		fallthrough
	case SituationAlertCauseEnginefailure:
		fallthrough
	case SituationAlertCauseEscalatorfailure:
		fallthrough
	case SituationAlertCauseEvacuation:
		fallthrough
	case SituationAlertCauseExplosion:
		fallthrough
	case SituationAlertCauseExplosionhazard:
		fallthrough
	case SituationAlertCauseFallenleaves:
		fallthrough
	case SituationAlertCauseFallentree:
		fallthrough
	case SituationAlertCauseFallentreeontheline:
		fallthrough
	case SituationAlertCauseFatality:
		fallthrough
	case SituationAlertCauseFilterblockade:
		fallthrough
	case SituationAlertCauseFire:
		fallthrough
	case SituationAlertCauseFireatstation:
		fallthrough
	case SituationAlertCauseFireatthestation:
		fallthrough
	case SituationAlertCauseFirebrigadeorder:
		fallthrough
	case SituationAlertCauseFirebrigadesafetychecks:
		fallthrough
	case SituationAlertCauseFirerun:
		fallthrough
	case SituationAlertCauseFlashfloods:
		fallthrough
	case SituationAlertCauseFlooding:
		fallthrough
	case SituationAlertCauseFog:
		fallthrough
	case SituationAlertCauseForeigndisturbances:
		fallthrough
	case SituationAlertCauseFrozen:
		fallthrough
	case SituationAlertCauseFuelproblem:
		fallthrough
	case SituationAlertCauseFuelshortage:
		fallthrough
	case SituationAlertCauseGangwayproblem:
		fallthrough
	case SituationAlertCauseGlazedfrost:
		fallthrough
	case SituationAlertCauseGrassfire:
		fallthrough
	case SituationAlertCauseGunfireonroadway:
		fallthrough
	case SituationAlertCauseHail:
		fallthrough
	case SituationAlertCauseHeavyrain:
		fallthrough
	case SituationAlertCauseHeavysnowfall:
		fallthrough
	case SituationAlertCauseHeavytraffic:
		fallthrough
	case SituationAlertCauseHightemperatures:
		fallthrough
	case SituationAlertCauseHightide:
		fallthrough
	case SituationAlertCauseHighwaterlevel:
		fallthrough
	case SituationAlertCauseHoliday:
		fallthrough
	case SituationAlertCauseIce:
		fallthrough
	case SituationAlertCauseIcedrift:
		fallthrough
	case SituationAlertCauseIceoncarriages:
		fallthrough
	case SituationAlertCauseIceonrailway:
		fallthrough
	case SituationAlertCauseIllvehicleoccupants:
		fallthrough
	case SituationAlertCauseIncident:
		fallthrough
	case SituationAlertCauseIndustrialaction:
		fallthrough
	case SituationAlertCauseInsufficientdemand:
		fallthrough
	case SituationAlertCauseLackofoperationalstock:
		fallthrough
	case SituationAlertCauseLandslide:
		fallthrough
	case SituationAlertCauseLatefinishtoengineeringwork:
		fallthrough
	case SituationAlertCauseLeaderboardfailure:
		fallthrough
	case SituationAlertCauseLevelcrossingaccident:
		fallthrough
	case SituationAlertCauseLevelcrossingblocked:
		fallthrough
	case SituationAlertCauseLevelcrossingfailure:
		fallthrough
	case SituationAlertCauseLevelcrossingincident:
		fallthrough
	case SituationAlertCauseLiftfailure:
		fallthrough
	case SituationAlertCauseLightingfailure:
		fallthrough
	case SituationAlertCauseLightningstrike:
		fallthrough
	case SituationAlertCauseLinesidefire:
		fallthrough
	case SituationAlertCauseLogisticproblems:
		fallthrough
	case SituationAlertCauseLowtide:
		fallthrough
	case SituationAlertCauseLowwaterlevel:
		fallthrough
	case SituationAlertCauseLuggagecarouselproblem:
		fallthrough
	case SituationAlertCauseMaintenancework:
		fallthrough
	case SituationAlertCauseMarch:
		fallthrough
	case SituationAlertCauseMiscellaneous:
		fallthrough
	case SituationAlertCauseMudslide:
		fallthrough
	case SituationAlertCauseNearmiss:
		fallthrough
	case SituationAlertCauseObjectontheline:
		fallthrough
	case SituationAlertCauseOperatorceasedtrading:
		fallthrough
	case SituationAlertCauseOperatorsuspended:
		fallthrough
	case SituationAlertCauseOvercrowded:
		fallthrough
	case SituationAlertCauseOverheadobstruction:
		fallthrough
	case SituationAlertCauseOverheadwirefailure:
		fallthrough
	case SituationAlertCauseOvertaking:
		fallthrough
	case SituationAlertCausePassengeraction:
		fallthrough
	case SituationAlertCausePassengersblockingdoors:
		fallthrough
	case SituationAlertCausePaving:
		fallthrough
	case SituationAlertCausePersonhitbytrain:
		fallthrough
	case SituationAlertCausePersonhitbyvehicle:
		fallthrough
	case SituationAlertCausePersonillonvehicle:
		fallthrough
	case SituationAlertCausePersonontheline:
		fallthrough
	case SituationAlertCausePersonundertrain:
		fallthrough
	case SituationAlertCausePointsfailure:
		fallthrough
	case SituationAlertCausePointsproblem:
		fallthrough
	case SituationAlertCausePoliceactivity:
		fallthrough
	case SituationAlertCausePoliceorder:
		fallthrough
	case SituationAlertCausePoorrailconditions:
		fallthrough
	case SituationAlertCausePoorweather:
		fallthrough
	case SituationAlertCausePowerproblem:
		fallthrough
	case SituationAlertCausePrecedingvehicle:
		fallthrough
	case SituationAlertCausePreviousdisturbances:
		fallthrough
	case SituationAlertCauseProblemsatborderpost:
		fallthrough
	case SituationAlertCauseProblemsatcustomspost:
		fallthrough
	case SituationAlertCauseProblemsonlocalroad:
		fallthrough
	case SituationAlertCauseProcession:
		fallthrough
	case SituationAlertCauseProvisiondelay:
		fallthrough
	case SituationAlertCausePublicdisturbance:
		fallthrough
	case SituationAlertCauseRailwaycrime:
		fallthrough
	case SituationAlertCauseRepairwork:
		fallthrough
	case SituationAlertCauseRiskofavalanches:
		fallthrough
	case SituationAlertCauseRiskofflooding:
		fallthrough
	case SituationAlertCauseRiskoflandslide:
		fallthrough
	case SituationAlertCauseRoadclosed:
		fallthrough
	case SituationAlertCauseRoadmaintenance:
		fallthrough
	case SituationAlertCauseRoadwaydamage:
		fallthrough
	case SituationAlertCauseRoadworks:
		fallthrough
	case SituationAlertCauseRockfalls:
		fallthrough
	case SituationAlertCauseRoughsea:
		fallthrough
	case SituationAlertCauseRouteblockage:
		fallthrough
	case SituationAlertCauseRoutediversion:
		fallthrough
	case SituationAlertCauseSabotage:
		fallthrough
	case SituationAlertCauseSafetyviolation:
		fallthrough
	case SituationAlertCauseSecurityalert:
		fallthrough
	case SituationAlertCauseSecurityincident:
		fallthrough
	case SituationAlertCauseServicedisruption:
		fallthrough
	case SituationAlertCauseServicefailure:
		fallthrough
	case SituationAlertCauseServiceindicatorfailure:
		fallthrough
	case SituationAlertCauseSeweragemaintenance:
		fallthrough
	case SituationAlertCauseSeweroverflow:
		fallthrough
	case SituationAlertCauseSightseersobstructingaccess:
		fallthrough
	case SituationAlertCauseSignalandswitchfailure:
		fallthrough
	case SituationAlertCauseSignalfailure:
		fallthrough
	case SituationAlertCauseSignalpassedatdanger:
		fallthrough
	case SituationAlertCauseSignalproblem:
		fallthrough
	case SituationAlertCauseSleet:
		fallthrough
	case SituationAlertCauseSlipperiness:
		fallthrough
	case SituationAlertCauseSlipperytrack:
		fallthrough
	case SituationAlertCauseSmokedetectedonvehicle:
		fallthrough
	case SituationAlertCauseSpecialevent:
		fallthrough
	case SituationAlertCauseSpeedrestrictions:
		fallthrough
	case SituationAlertCauseStaffabsence:
		fallthrough
	case SituationAlertCauseStaffassault:
		fallthrough
	case SituationAlertCauseStaffinjury:
		fallthrough
	case SituationAlertCauseStaffinwrongplace:
		fallthrough
	case SituationAlertCauseStaffshortage:
		fallthrough
	case SituationAlertCauseStaffsickness:
		fallthrough
	case SituationAlertCauseStationoverrun:
		fallthrough
	case SituationAlertCauseStormconditions:
		fallthrough
	case SituationAlertCauseStormdamage:
		fallthrough
	case SituationAlertCauseStrongwinds:
		fallthrough
	case SituationAlertCauseSubsidence:
		fallthrough
	case SituationAlertCauseSuspectvehicle:
		fallthrough
	case SituationAlertCauseSwingbridgefailure:
		fallthrough
	case SituationAlertCauseTechnicalproblem:
		fallthrough
	case SituationAlertCauseTelephonedthreat:
		fallthrough
	case SituationAlertCauseTerroristincident:
		fallthrough
	case SituationAlertCauseTheft:
		fallthrough
	case SituationAlertCauseTicketingsystemnotavailable:
		fallthrough
	case SituationAlertCauseTidalrestrictions:
		fallthrough
	case SituationAlertCauseTrackcircuitproblem:
		fallthrough
	case SituationAlertCauseTractionfailure:
		fallthrough
	case SituationAlertCauseTrafficmanagementsystemfailure:
		fallthrough
	case SituationAlertCauseTraincoupling:
		fallthrough
	case SituationAlertCauseTraindoor:
		fallthrough
	case SituationAlertCauseTrainstruckanimal:
		fallthrough
	case SituationAlertCauseTrainstruckobject:
		fallthrough
	case SituationAlertCauseTrainwarningsystemproblem:
		fallthrough
	case SituationAlertCauseUnattendedbag:
		fallthrough
	case SituationAlertCauseUndefinedalertcause:
		fallthrough
	case SituationAlertCauseUndefinedenvironmentalproblem:
		fallthrough
	case SituationAlertCauseUndefinedequipmentproblem:
		fallthrough
	case SituationAlertCauseUndefinedpersonnelproblem:
		fallthrough
	case SituationAlertCauseUndefinedproblem:
		fallthrough
	case SituationAlertCauseUnknown:
		fallthrough
	case SituationAlertCauseUnofficialindustrialaction:
		fallthrough
	case SituationAlertCauseUnscheduledconstructionwork:
		fallthrough
	case SituationAlertCauseVandalism:
		fallthrough
	case SituationAlertCauseVegetation:
		fallthrough
	case SituationAlertCauseVehicleblockingtrack:
		fallthrough
	case SituationAlertCauseVehiclefailure:
		fallthrough
	case SituationAlertCauseVehicleontheline:
		fallthrough
	case SituationAlertCauseVehiclestruckanimal:
		fallthrough
	case SituationAlertCauseVehiclestruckobject:
		fallthrough
	case SituationAlertCauseViaductfailure:
		fallthrough
	case SituationAlertCauseWaitingfortransferpassengers:
		fallthrough
	case SituationAlertCauseWaterlogged:
		fallthrough
	case SituationAlertCauseWheelimpactload:
		fallthrough
	case SituationAlertCauseWheelproblem:
		fallthrough
	case SituationAlertCauseWildlandfire:
		fallthrough
	case SituationAlertCauseWorktorule:
		*alertCause = SituationAlertCause(s)
		return nil
	default:
		return fmt.Errorf("invalid alert cause %s", s)
	}
}

type SituationReality string

const (
	SituationRealityReal              SituationReality = "real"
	SituationRealitySecurityExercise  SituationReality = "securityExercise"
	SituationRealityTechnicalExercise SituationReality = "technicalExercise"
	SituationRealityTest              SituationReality = "test"
)

func (reality *SituationReality) FromString(s string) error {
	switch SituationReality(s) {
	case SituationRealityReal:
		fallthrough
	case SituationRealitySecurityExercise:
		fallthrough
	case SituationRealityTechnicalExercise:
		fallthrough
	case SituationRealityTest:
		*reality = SituationReality(s)
		return nil
	default:
		return fmt.Errorf("invalid reality %s", s)
	}
}

type SituationCondition string

const (
	SituationConditionAdditionalRide              SituationCondition = "additionalRide"
	SituationConditionAdditionalService           SituationCondition = "additionalService"
	SituationContitionAdditionalStop              SituationCondition = "additionalStop"
	SituationConditionAltered                     SituationCondition = "altered"
	SituationConditionArrivesEarly                SituationCondition = "arrivesEarly"
	SituationConditionBoarding                    SituationCondition = "boarding"
	SituationConditionCancelled                   SituationCondition = "cancelled"
	SituationConditionChangeOfPlatform            SituationCondition = "changeOfPlatform"
	SituationConditionDelay                       SituationCondition = "delay"
	SituationConditionDelayed                     SituationCondition = "delayed"
	SituationConditionDiscontinuedOperation       SituationCondition = "discontinuedOperation"
	SituationConditionDisrupted                   SituationCondition = "disrupted"
	SituationConditionDisruption                  SituationCondition = "disruption"
	SituationConditionDisturbanceRectified        SituationCondition = "disturbanceRectified"
	SituationConditionDiverted                    SituationCondition = "diverted"
	SituationConditionExtendedService             SituationCondition = "extendedService"
	SituationConditionFullLengthService           SituationCondition = "fullLengthService"
	SituationConditionGoToGate                    SituationCondition = "goToGate"
	SituationConditionIntermittentService         SituationCondition = "intermittentService"
	SituationConditionIrregularTraffic            SituationCondition = "irregularTraffic"
	SituationConditionLimitedOperation            SituationCondition = "limitedOperation"
	SituationConditionLineCancellation            SituationCondition = "lineCancellation"
	SituationConditionMajorDelays                 SituationCondition = "majorDelays"
	SituationConditionMinorDelays                 SituationCondition = "minorDelays"
	SituationConditionNormalService               SituationCondition = "normalService"
	SituationConditionNoService                   SituationCondition = "noService"
	SituationConditionOnTime                      SituationCondition = "onTime"
	SituationConditionOperationTimeExtension      SituationCondition = "operationTimeExtension"
	SituationConditionReplacementRide             SituationCondition = "replacementRide"
	SituationConditionReplacementService          SituationCondition = "replacementService"
	SituationConditionReplacementTransport        SituationCondition = "replacementTransport"
	SituationConditionShortFormedService          SituationCondition = "shortFormedService"
	SituationConditionShuttleService              SituationCondition = "shuttleService"
	SituationConditionSpecialService              SituationCondition = "specialService"
	SituationConditionSplittingTrain              SituationCondition = "splittingTrain"
	SituationConditionStopCancelled               SituationCondition = "stopCancelled"
	SituationConditionStopMove                    SituationCondition = "stopMoved"
	SituationConditionStopOnDemand                SituationCondition = "stopOnDemand"
	SituationConditionSubstitutedStop             SituationCondition = "substitutedStop"
	SituationConditionTemporarilyNonStopping      SituationCondition = "temporarilyNonStopping"
	SituationConditionTemporaryStopplace          SituationCondition = "temporaryStopplace"
	SituationConditionTrainShortened              SituationCondition = "trainShortened"
	SituationConditionTripCancellation            SituationCondition = "tripCancellation"
	SituationConditionUndefinedServiceInformation SituationCondition = "undefinedServiceInformation"
	SituationConditionUndefinedStatus             SituationCondition = "undefinedStatus"
	SituationConditionUnknown                     SituationCondition = "unknown"
	SituationConditionWagonOrderChanged           SituationCondition = "wagonOrderChanged"
)

func (condition *SituationCondition) FromString(c string) error {
	switch SituationCondition(c) {
	case SituationConditionAdditionalRide:
		fallthrough
	case SituationConditionAdditionalService:
		fallthrough
	case SituationContitionAdditionalStop:
		fallthrough
	case SituationConditionAltered:
		fallthrough
	case SituationConditionArrivesEarly:
		fallthrough
	case SituationConditionBoarding:
		fallthrough
	case SituationConditionCancelled:
		fallthrough
	case SituationConditionChangeOfPlatform:
		fallthrough
	case SituationConditionDelay:
		fallthrough
	case SituationConditionDelayed:
		fallthrough
	case SituationConditionDiscontinuedOperation:
		fallthrough
	case SituationConditionDisrupted:
		fallthrough
	case SituationConditionDisruption:
		fallthrough
	case SituationConditionDisturbanceRectified:
		fallthrough
	case SituationConditionDiverted:
		fallthrough
	case SituationConditionExtendedService:
		fallthrough
	case SituationConditionFullLengthService:
		fallthrough
	case SituationConditionGoToGate:
		fallthrough
	case SituationConditionIntermittentService:
		fallthrough
	case SituationConditionIrregularTraffic:
		fallthrough
	case SituationConditionLimitedOperation:
		fallthrough
	case SituationConditionLineCancellation:
		fallthrough
	case SituationConditionMajorDelays:
		fallthrough
	case SituationConditionMinorDelays:
		fallthrough
	case SituationConditionNormalService:
		fallthrough
	case SituationConditionNoService:
		fallthrough
	case SituationConditionOnTime:
		fallthrough
	case SituationConditionOperationTimeExtension:
		fallthrough
	case SituationConditionReplacementRide:
		fallthrough
	case SituationConditionReplacementService:
		fallthrough
	case SituationConditionReplacementTransport:
		fallthrough
	case SituationConditionShortFormedService:
		fallthrough
	case SituationConditionShuttleService:
		fallthrough
	case SituationConditionSpecialService:
		fallthrough
	case SituationConditionSplittingTrain:
		fallthrough
	case SituationConditionStopCancelled:
		fallthrough
	case SituationConditionStopMove:
		fallthrough
	case SituationConditionStopOnDemand:
		fallthrough
	case SituationConditionSubstitutedStop:
		fallthrough
	case SituationConditionTemporarilyNonStopping:
		fallthrough
	case SituationConditionTemporaryStopplace:
		fallthrough
	case SituationConditionTrainShortened:
		fallthrough
	case SituationConditionTripCancellation:
		fallthrough
	case SituationConditionUndefinedServiceInformation:
		fallthrough
	case SituationConditionUndefinedStatus:
		fallthrough
	case SituationConditionUnknown:
		fallthrough
	case SituationConditionWagonOrderChanged:
		*condition = SituationCondition(c)
		return nil
	default:
		return fmt.Errorf("invalid condition %s", c)
	}
}

func (c *SituationCondition) FromProto(value interface{}) error {
	switch v := value.(type) {
	case gtfs.Alert_Effect:
		switch v {
		case gtfs.Alert_NO_SERVICE:
			*c = SituationConditionNoService
		case gtfs.Alert_REDUCED_SERVICE:
			*c = SituationConditionAltered
		case gtfs.Alert_SIGNIFICANT_DELAYS:
			*c = SituationConditionDisrupted
		case gtfs.Alert_DETOUR:
			*c = SituationConditionDiverted
		case gtfs.Alert_ADDITIONAL_SERVICE:
			*c = SituationConditionAdditionalService
		case gtfs.Alert_STOP_MOVED:
			*c = SituationConditionDiverted
		case gtfs.Alert_NO_EFFECT:
			*c = SituationConditionNormalService
		default:
			return errors.New("invalid value")
		}
	default:
		return fmt.Errorf("unsupported value %T", value)
	}
	return nil
}

func (c SituationCondition) ToProto(dest interface{}) error {
	switch v := dest.(type) {
	case *gtfs.Alert_Effect:
		switch c {
		case SituationConditionNoService:
			*v = gtfs.Alert_NO_SERVICE
		case SituationConditionAltered:
			*v = gtfs.Alert_REDUCED_SERVICE
		case SituationConditionDisrupted:
			*v = gtfs.Alert_SIGNIFICANT_DELAYS
		case SituationConditionDiverted:
			*v = gtfs.Alert_DETOUR
		case SituationConditionAdditionalService:
			*v = gtfs.Alert_ADDITIONAL_SERVICE
		case SituationConditionUnknown:
			*v = gtfs.Alert_UNKNOWN_EFFECT
		case SituationConditionNormalService:
			*v = gtfs.Alert_NO_EFFECT
		default:
			*v = gtfs.Default_Alert_Effect
		}
	default:
		return fmt.Errorf("unsupported destination %T", dest)
	}
	return nil
}

type SituationActionStatus string

const (
	SituationActionStatusOpen		SituationActionStatus = "open"
	SituationActionStatusPublished 	SituationActionStatus = "published"
	SituationActionStatusClosed    	SituationActionStatus = "closed"
)

func (as *SituationActionStatus) FromString(s string) error {
	switch SituationActionStatus(s) {
	case SituationActionStatusOpen:
		fallthrough
	case SituationActionStatusPublished:
		fallthrough
	case SituationActionStatusClosed:
		*as = SituationActionStatus(s)
		return nil
	}
	return fmt.Errorf("invalid action status %s", s)
}

type SituationScopeType string

const (
	SituationScopeTypeGeneral        SituationScopeType = "general"
	SituationScopeTypeOperator       SituationScopeType = "operator"
	SituationScopeTypeNetwork        SituationScopeType = "network"
	SituationScopeTypeRoute          SituationScopeType = "route"
	SituationScopeTypeLine           SituationScopeType = "line"
	SituationScopeTypePlace          SituationScopeType = "place"
	SituationScopeTypeStopPlace      SituationScopeType = "stopPlace"
	SituationScopeTypeStopPoint      SituationScopeType = "stopPoint"
	SituationScopeTypeVehicleJourney SituationScopeType = "VehicleJourney"
)

func (sc *SituationScopeType) FromString(s string) error {
	switch SituationScopeType(s) {
	case SituationScopeTypeGeneral:
		fallthrough
	case SituationScopeTypeOperator:
		fallthrough
	case SituationScopeTypeNetwork:
		fallthrough
	case SituationScopeTypeRoute:
		fallthrough
	case SituationScopeTypeLine:
		fallthrough
	case SituationScopeTypePlace:
		fallthrough
	case SituationScopeTypeStopPlace:
		fallthrough
	case SituationScopeTypeStopPoint:
		fallthrough
	case SituationScopeTypeVehicleJourney:
		*sc = SituationScopeType(s)
		return nil
	}
	return fmt.Errorf("invalid scope type %s", s)
}
