Feature: Test Database load and save

  @nostart @database
  Scenario: Load Referentials and partners
    Given the table "referentials" has the following data:
    | referential_id                         | slug     | settings          | tokens          |
    | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'first'  | '{"key":"value"}' | '["testtoken"]' |
    | '6ba7b814-9dad-11d1-0001-00c04fd430c8' | 'second' | '{}'              | '["testtoken"]' |
    And the table "partners" has the following data:
    | id                                     | referential_id                         | slug             | name     | settings                                                                                                | connector_types                                                          |
    | '6ba7b814-9dad-11d1-0002-00c04fd430c8' | '6ba7b814-9dad-11d1-0000-00c04fd430c8' | 'first_partner'  | 'first'  | '{"remote_url": "http://localhost", "remote_objectid_kind": "Reflex", "remote_credential": "ara_cred"}' | '["siri-stop-monitoring-request-collector", "siri-check-status-client"]' |
    | '6ba7b814-9dad-11d1-0003-00c04fd430c8' | '6ba7b814-9dad-11d1-0001-00c04fd430c8' | 'second_partner' | 'second' | '{}'                                                                                                    | '[]'                                                                     |
    When I start Ara
    Then one Referential has the following attributes:
    | Id       | 6ba7b814-9dad-11d1-0000-00c04fd430c8 |
    | Slug     | first                                |
    | Settings | {"key":"value"}                      |
    And one Referential has the following attributes:
    | Id   | 6ba7b814-9dad-11d1-0001-00c04fd430c8 |
    | Slug | second                               |
    And one Partner in Referential "first" has the following attributes:
    | Id             | 6ba7b814-9dad-11d1-0002-00c04fd430c8                                                                  |
    | Slug           | first_partner                                                                                         |
    | Name           | first                                                                                                 |
    | Settings       | {"remote_url": "http://localhost", "remote_objectid_kind": "Reflex", "remote_credential": "ara_cred"} |
    | ConnectorTypes | ["siri-check-status-client", "siri-stop-monitoring-request-collector"]                                |
    And one Partner in Referential "second" has the following attributes:
    | Id   | 6ba7b814-9dad-11d1-0003-00c04fd430c8 |
    | Slug | second_partner                       |

  @database
  Scenario: Save referentials
    Given a Referential "first" exists with the following settings:
      | model.reload_at | 01:00 |
    And a Referential "second" exists with the following settings:
      | model.reload_at | 02:00 |
    When I save all referentials
    Then the table "referentials" has rows with the following values:
      | slug   | settings                    |
      | first  | {"model.reload_at":"01:00"} |
      | second | {"model.reload_at":"02:00"} |

  @database
  Scenario: Remove a deleted referential
    Given a Referential "first" exists
    And a Referential "second" exists
    And I save all referentials
    When the Referential "second" is destroyed
    And I save all referentials
    Then the table "referentials" has a row with the following values:
      | slug    | first   |
    And the table "referentials" has no row with the following values:
      | slug    | second  |
