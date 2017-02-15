Feature: Manage Stop Visits

  Background:
    Given a Referential "test" is created

  Scenario: Create a StopVisit
  When a StopVisit is created with the following attributes:
    | ObjectIDs                   | "internal": "1234"                |
    | Attribute[PassageOrder]     | 4                                 |
    | Reference[StopAreaId]       | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    | Reference[VehicleJourneyId] | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
  Then the StopVisit "6ba7b814-9dad-11d1-1-00c04fd430c8" has the following attributes:
    | ObjectIDs                 | "internal":"1234"                 |
    | Attribute[PassageOrder]   | 4                                 |
    | Attribute[StopArea]       | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    | Attribute[VehicleJourney] | 6ba7b814-9dad-11d1-4-00c04fd430c8 |

  @wip
  # See Issue 2603
  Scenario: Create a StopVisit
  When a StopVisit is created with the following attributes:
    | Attribute[ArrivalStatus]    | ontime                            |
    | Attribute[DepartureStatus]  | ontime                            |
    | ObjectIDs                   | "internal": "1234"                |
    | Attribute[PassageOrder      | 4                                 |
    | Attribut[RecordedAt]        | 2017-01-01T11:00:00.000Z          |
    | Schedule[aimed]#Arrival     | 2017-01-01T13:00:00.000Z          |
    | Schedule[aimed]#Departure   | 2017-01-01T13:02:00.000Z          |
    | Reference[StopAreaId]       | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    | Reference[VehicleJourneyId] | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
  Then the StopVisit "6ba7b814-9dad-11d1-1-00c04fd430c8" has the following attributes:
    | Attribute[ArrivalStatus]    | ontime                            |
    | Attribute[DepartureStatus]  | ontime                            |
    | ObjectIDs                   | "internal":"1234"                 |
    | Attribute[PassageOrder]     | 4                                 |
    | RecordedAt                  | 2017-01-01T11:00:00.000Z          |
    | Schedule[aimed]#Arrival     | 2017-01-01T13:00:00.000Z          |
    | Schedule[aimed]#Departure   | 2017-01-01T13:02:00.000Z          |
    | Reference[StopAreaId]       | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    | Reference[VehicleJourneyId] | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
