package core

import (
	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/state"
	"errors"
	"slices"
)

type ServiceAlertsBroadcaster struct {
	state.Startable
	connector

	cache *cache.CachedItem
}

type ServiceAlertsBroadcasterFactory struct{}

func (factory *ServiceAlertsBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewServiceAlertsBroadcaster(partner)
}

func (factory *ServiceAlertsBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
}

func NewServiceAlertsBroadcaster(partner *Partner) *ServiceAlertsBroadcaster {
	connector := &ServiceAlertsBroadcaster{}
	connector.partner = partner

	return connector
}

func (connector *ServiceAlertsBroadcaster) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(GTFS_RT_SERVICE_ALERTS_BROADCASTER)
	connector.cache = cache.NewCachedItem("ServiceAlerts", connector.partner.CacheTimeout(GTFS_RT_SERVICE_ALERTS_BROADCASTER), nil, func(...interface{}) (interface{}, error) { return connector.handleGtfs() })
}

func (connector *ServiceAlertsBroadcaster) HandleGtfs(feed *gtfs.FeedMessage) {
	entities, _ := connector.cache.Value()
	feedEntities := entities.([]*gtfs.FeedEntity)

	feed.Entity = append(feed.Entity, feedEntities...)
}

func (connector *ServiceAlertsBroadcaster) handleGtfs() (entities []*gtfs.FeedEntity, err error) {
	situations := connector.partner.Model().Situations().FindAll()

	for _, situation := range situations {
		if !connector.canBroadcast(situation) {
			continue
		}
		var situationNumber string
		code, present := situation.Code(connector.remoteCodeSpace)
		if present {
			situationNumber = code.Value()
		} else {
			code, present = situation.Code(model.Default)
			if !present {
				logger.Log.Debugf("Unknown Code for Situation %s", situation.Id())
				return
			}
			situationNumber = connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "InfoMessage", Id: code.Value()})
		}

		alert := &gtfs.Alert{}

		// InformedEntities
		for _, affect := range situation.Affects {
			informedEntities, _, err := model.AffectToProto(affect, connector.remoteCodeSpace, connector.Partner().Model())
			if err != nil {
				logger.Log.Debugf("Error for affect: %v", err)
			}
			alert.InformedEntity = append(alert.InformedEntity, informedEntities...)
		}

		if len(alert.InformedEntity) == 0 {
			return nil, errors.New("no informed entities")
		}

		feedEntity := &gtfs.FeedEntity{
			Id: &situationNumber,
		}

		// Periods
		var ts []*gtfs.TimeRange
		for _, period := range situation.ValidityPeriods {
			t := &gtfs.TimeRange{}
			if err := period.ToProto(t); err != nil {
				logger.Log.Debugf("Error for period: %v", err)
			}

			ts = append(ts, t)
		}

		if len(ts) != 0 {
			alert.ActivePeriod = ts
		}

		// Effect
		// we choose the first one ...
		var e gtfs.Alert_Effect
		if len(situation.Consequences) != 0 {
			if err := situation.Consequences[0].Condition.ToProto(&e); err != nil {
				logger.Log.Debugf("Error for Condition: %v", err)
			}
		} else {
			var c model.SituationCondition
			c.ToProto(&e)
		}
		alert.Effect = &e

		// Cause
		var c gtfs.Alert_Cause
		if err := situation.AlertCause.ToProto(&c); err != nil {
			logger.Log.Debugf("Error for alert cause: %v", err)
		} else {
			alert.Cause = &c
		}

		// Severity
		var s gtfs.Alert_SeverityLevel
		if err := situation.Severity.ToProto(&s); err != nil {
			logger.Log.Debugf("Error for severity: %v", err)
		} else {
			alert.SeverityLevel = &s
		}

		// HeaderText
		var h gtfs.TranslatedString_Translation
		if err := situation.Summary.ToProto(&h); err != nil {
			logger.Log.Debugf("Error for summary: %v", err)
		} else {
			translatedString := gtfs.TranslatedString{}
			translatedString.Translation = append(translatedString.Translation, &h)
			alert.HeaderText = &translatedString
		}

		// DescriptionText
		var d gtfs.TranslatedString_Translation
		if err := situation.Description.ToProto(&d); err != nil {
			logger.Log.Debugf("Error for description: %v", err)
		} else {
			translatedString := gtfs.TranslatedString{}
			translatedString.Translation = append(translatedString.Translation, &d)
			alert.DescriptionText = &translatedString
		}

		feedEntity.Alert = alert
		entities = append(entities, feedEntity)
	}

	return
}

func (connector *ServiceAlertsBroadcaster) canBroadcast(situation model.Situation) bool {
	tagsToBroadcast := connector.partner.BroadcastSituationsInternalTags()
	if len(tagsToBroadcast) != 0 {
		for _, tag := range situation.InternalTags {
			if slices.Contains(tagsToBroadcast, tag) {
				return true
			}
		}
		return false
	}

	return true
}
