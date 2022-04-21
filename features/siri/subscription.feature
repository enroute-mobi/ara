Feature: Support SIRI subscription

  Background:
    Given a Referential "test" is created

  Scenario: 4377 - Change status of subscription with termination request
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringBroadcast                            |
      | ExternalId        | ExternalId                                         |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopPoint:SP:24:LOC" |
    When I send this SIRI request
        """
        <?xml version='1.0' encoding='utf-8'?>
        <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns1:DeleteSubscription xmlns:ns1="http://wsdl.siri.org.uk" xmlns:ns5="http://www.siri.org.uk/siri">
            <ServiceRequestInfo
             xmlns:ns2="http://www.ifopt.org.uk/acsb"
             xmlns:ns3="http://www.ifopt.org.uk/ifopt"
             xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
             xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
              <ns5:RequestorRef>NINOXE:default</ns5:RequestorRef>
              <ns2:MessageIdentifier>TermSubReq:Test:0</ns2:MessageIdentifier>
            </ServiceRequestInfo>
            <Request version="2.0:FR-IDF-2.4">
              <ns5:SubscriptionRef>ExternalId</ns5:SubscriptionRef>
            </Request>
            <RequestExtension/>
          </ns1:DeleteSubscription>
        </S:Body>
      </S:Envelope>
      """
    Then no Subscription exists

  @ARA-1066
  Scenario: Remove all subscriptions at once with termination request using raw envelope
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ara                   |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
      | siri.envelope        | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimeTableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimeTableBroadcast           |
      | ExternalId        | AnotherExternalId                     |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:A:BUS" |
    When I send this SIRI request
        """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
   <TerminateSubscriptionRequest>
      <RequestTimestamp>2016-09-22T08:01:20.227+02:00</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <MessageIdentifier />
      <All />
   </TerminateSubscriptionRequest>
</Siri>
      """
    Then I should receive this SIRI response
      """
<?xml version='1.0' encoding='UTF-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
  <TerminateSubscriptionResponse>
    <ResponseTimestamp>2017-01-01T12:00:00.000Z</ResponseTimestamp>
    <ResponderRef>ara</ResponderRef>
    <TerminationResponseStatus>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
    </TerminationResponseStatus>
    <TerminationResponseStatus>
      <SubscriptionRef>AnotherExternalId</SubscriptionRef>
      <Status>true</Status>
    </TerminationResponseStatus>
  </TerminateSubscriptionResponse>
</Siri>
      """
    Then no Subscription exists

  @ARA-1066
  Scenario: Remove one subscription with termination request using raw envelope
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ara                   |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
      | siri.envelope        | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimeTableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
     And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimeTableBroadcast           |
      | ExternalId        | AnotherExternalId                     |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:A:BUS" |
   When I send this SIRI request
        """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
   <TerminateSubscriptionRequest>
      <RequestTimestamp>2016-09-22T08:01:20.227+02:00</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <SubscriptionRef>externalId</SubscriptionRef>
   </TerminateSubscriptionRequest>
</Siri>
      """
    Then I should receive this SIRI response
     """
<?xml version='1.0' encoding='UTF-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
  <TerminateSubscriptionResponse>
    <ResponseTimestamp>2017-01-01T12:00:00.000Z</ResponseTimestamp>
    <ResponderRef>ara</ResponderRef>
    <TerminationResponseStatus>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
    </TerminationResponseStatus>
  </TerminateSubscriptionResponse>
</Siri>
     """
    Then No Subscriptions exist with the following attributes:
      | internal | NINOXE:Line:3:LOC |
    Then Subscriptions exist with the following attributes:
      | internal | NINOXE:Line:A:BUS |

  @ARA-1066
  Scenario: Handle TerminateSubscriptionRequest on an unnown subscription using raw envelope
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ara                   |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
      | siri.envelope        | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimeTableBroadcast            |
      | ExternalId        | externalId                             |
      | SubscriberRef     | subscriber                             |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC"  |
      | ReferenceArray[1] | Line, "internal": "NINOXE:Line:C:Tram" |
     And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimeTableBroadcast           |
      | ExternalId        | AnotherExternalId                     |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:A:BUS" |
   When I send this SIRI request
        """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xmlns="http://www.siri.org.uk/siri">
   <TerminateSubscriptionRequest>
      <RequestTimestamp>2016-09-22T08:01:20.227+02:00</RequestTimestamp>
      <RequestorRef>NINOXE:default</RequestorRef>
      <SubscriptionRef>UnknownExternalId</SubscriptionRef>
   </TerminateSubscriptionRequest>
</Siri>
      """
    Then I should receive this SIRI response
     """
<?xml version='1.0' encoding='UTF-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
  <TerminateSubscriptionResponse>
    <ResponseTimestamp>2017-01-01T12:00:00.000Z</ResponseTimestamp>
    <ResponderRef>ara</ResponderRef>
    <TerminationResponseStatus>
      <SubscriptionRef>UnknownExternalId</SubscriptionRef>
      <Status>false</Status>
      <ErrorCondition>
        <InvalidDataReferencesError>
          <ErrorText>Subscription not found: 'UnknownExternalId'</ErrorText>
        </InvalidDataReferencesError>
      </ErrorCondition>
    </TerminationResponseStatus>
  </TerminateSubscriptionResponse>
</Siri>
     """
    Then Subscriptions exist with the following attributes:
      | internal | NINOXE:Line:A:BUS  |
      | internal | NINOXE:Line:3:LOC  |
      | internal | NINOXE:Line:C:Tram |
