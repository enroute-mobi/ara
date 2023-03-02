Feature: Subscription Management

  Background:
    Given a Referential "test" is created

  Scenario: Delete subscriptions resources on different subscriptions after a Partner reload
    Given a SIRI server "fake_partner" on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
    And a Subscription exist with the following attributes:
      | Kind                  | StopMonitoringBroadcast                            |
      | ExternalId            | ExternalId                                         |
      | ReferenceArray[0]     | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | SubscribeResourcesNow | true                                               |
    And show me ara subscriptions for partner "test"
    And 30 seconds have passed
    And a Subscription exist with the following attributes:
      | Kind                  | StopMonitoringBroadcast                            |
      | ExternalId            | ExternalId                                         |
      | ReferenceArray[0]     | StopArea, "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | SubscribeResourcesNow | true                                               |
    And a minute has passed
    And show me ara subscriptions for partner "test"
    Then Subscriptions exist with the following attributes:
      | internal | NINOXE:StopPoint:SP:24:LOC |
      | internal | NINOXE:StopPoint:SP:25:LOC |
    When I edit the "fake_partner" SIRI server with new ServiceStartedTime "2017-01-01T14:00:20.000+02:00"
    And 30 seconds have passed
    Then No Subscriptions exist with the following resources:
      | internal | NINOXE:StopPoint:SP:24:LOC |
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:StopPoint:SP:25:LOC |
    When I edit the "fake_partner" SIRI server with new ServiceStartedTime "2017-01-01T14:01:00.000+02:00"
    And 30 seconds have passed
    Then no Subscription exists

  Scenario: Delete multiple resources on a single subscription after Partner reload
    Given a SIRI server "fake_partner" on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
    And a Subscription exist with the following attributes:
      | Kind                  | StopMonitoringBroadcast                            |
      | ExternalId            | ExternalId                                         |
      | ReferenceArray[0]     | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
      | ReferenceArray[1]     | StopArea, "internal": "NINOXE:StopPoint:SP:25:LOC" |
      | SubscribeResourcesNow | true                                               |
    And show me ara subscriptions for partner "test"
    And 30 seconds have passed
    Then Subscriptions exist with the following attributes:
      | internal | NINOXE:StopPoint:SP:24:LOC |
      | internal | NINOXE:StopPoint:SP:25:LOC |
    When I edit the "fake_partner" SIRI server with new ServiceStartedTime "2017-01-01T14:00:20.000+02:00"
    And 30 seconds have passed
    Then No Subscriptions exist with the following resources:
      | internal | NINOXE:StopPoint:SP:24:LOC |
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:StopPoint:SP:25:LOC |
    When I edit the "fake_partner" SIRI server with new ServiceStartedTime "2017-01-01T14:01:00.000+02:00"
    And 30 seconds have passed
    Then no Subscription exists
