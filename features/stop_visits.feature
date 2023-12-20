Feature: Manage Stop Visits

  Background:
    Given a Referential "test" is created

  Scenario: Create a StopVisit
    When a StopVisit is created with the following attributes:
      | Codes            | "internal": "1234"                |
      | PassageOrder     | 4                                 |
      | StopAreaId       | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    Then the StopVisit "6ba7b814-9dad-11d1-1-00c04fd430c8" has the following attributes:
      | Codes            | "internal":"1234"                 |
      | PassageOrder     | 4                                 |
      | StopAreaId       | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId | 6ba7b814-9dad-11d1-4-00c04fd430c8 |


      # See Issue 2603
  Scenario: Create a StopVisit
    When a StopVisit is created with the following attributes:
      | ArrivalStatus                 | ontime                            |
      | DepartureStatus               | ontime                            |
      | Codes                         | "internal": "1234"                |
      | PassageOrder                  | 4                                 |
      | RecordedAt                    | 2017-01-01T11:00:00Z              |
      | Schedule[aimed]#ArrivalTime   | 2017-01-01T13:00:00Z              |
      | Schedule[aimed]#DepartureTime | 2017-01-01T13:02:00Z              |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    Then the StopVisit "6ba7b814-9dad-11d1-1-00c04fd430c8" has the following attributes:
      | ArrivalStatus                 | ontime                            |
      | Collected                     | false                             |
      | DepartureStatus               | ontime                            |
      | Id                            | 6ba7b814-9dad-11d1-1-00c04fd430c8 |
      | Codes                         | "internal":"1234"                 |
      | PassageOrder                  | 4                                 |
      | RecordedAt                    | 2017-01-01T11:00:00Z              |
      | Schedule[aimed]#ArrivalTime   | 2017-01-01T13:00:00Z              |
      | Schedule[aimed]#DepartureTime | 2017-01-01T13:02:00Z              |
      | StopAreaId                    | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | VehicleAtStop                 | false                             |
      | VehicleJourneyId              | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
