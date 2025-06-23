Feature: Support SIRI subscription

  Background:
    Given a Referential "test" is created

  Scenario: 4377 - Change status of subscription with termination request
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_code_space | internal              |
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
            <DeleteSubscriptionInfo
             xmlns:ns2="http://www.ifopt.org.uk/acsb"
             xmlns:ns3="http://www.ifopt.org.uk/ifopt"
             xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
             xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:RequestTimestamp>2016-09-22T08:01:20.227+02:00</ns5:RequestTimestamp>
              <ns5:RequestorRef>NINOXE:default</ns5:RequestorRef>
              <ns5:MessageIdentifier>TermSubReq:Test:0</ns5:MessageIdentifier>
            </DeleteSubscriptionInfo>
            <Request>
              <ns5:SubscriptionRef>ExternalId</ns5:SubscriptionRef>
            </Request>
            <RequestExtension/>
          </ns1:DeleteSubscription>
        </S:Body>
      </S:Envelope>
      """
    Then no Subscription exists

  @ARA-1066
  Scenario: Handle DeleteSubscription on an unknown subscription using a SOAP envelope
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ara                   |
      | local_credential     | NINOXE:default        |
      | remote_code_space | internal              |
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
            <DeleteSubscriptionInfo
             xmlns:ns2="http://www.ifopt.org.uk/acsb"
             xmlns:ns3="http://www.ifopt.org.uk/ifopt"
             xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
             xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:RequestTimestamp>2016-09-22T08:01:20.227+02:00</ns5:RequestTimestamp>
              <ns5:RequestorRef>NINOXE:default</ns5:RequestorRef>
              <ns5:MessageIdentifier>TermSubReq:Test:0</ns5:MessageIdentifier>
            </DeleteSubscriptionInfo>
            <Request>
              <ns5:SubscriptionRef>UnknownExternalId</ns5:SubscriptionRef>
            </Request>
            <RequestExtension/>
          </ns1:DeleteSubscription>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
    """
<?xml version='1.0' encoding='UTF-8'?>
<S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
  <S:Body>
    <sw:DeleteSubscriptionResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
      <DeleteSubscriptionAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>ara</siri:ResponderRef>
        <siri:RequestMessageRef>TermSubReq:Test:0</siri:RequestMessageRef>
      </DeleteSubscriptionAnswerInfo>
      <Answer>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>ara</siri:ResponderRef>
        <siri:RequestMessageRef>TermSubReq:Test:0</siri:RequestMessageRef>
        <siri:TerminationResponseStatus>
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:SubscriberRef>NINOXE:default</siri:SubscriberRef>
          <siri:SubscriptionRef>UnknownExternalId</siri:SubscriptionRef>
          <siri:Status>false</siri:Status>
          <siri:ErrorCondition>
            <siri:UnknownSubscriptionError>
              <siri:ErrorText>Subscription not found: 'UnknownExternalId'</siri:ErrorText>
            </siri:UnknownSubscriptionError>
          </siri:ErrorCondition>
        </siri:TerminationResponseStatus>
      </Answer>
      <AnswerExtension/>
    </sw:DeleteSubscriptionResponse>
  </S:Body>
</S:Envelope>
    """
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:StopPoint:SP:24:LOC |

  @ARA-1267
  Scenario: Ignore DeleteSubscription if setting broadcast.siri.ignore_terminate_subscription_requests is set to true
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-stop-monitoring-subscription-broadcaster] and the following settings:
      | remote_url                                            | http://localhost:8090 |
      | remote_credential                                     | ara                   |
      | local_credential                                      | NINOXE:default        |
      | remote_code_space                                  | internal              |
      | broadcast.siri.ignore_terminate_subscription_requests | true                  |
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
            <DeleteSubscriptionInfo
             xmlns:ns2="http://www.ifopt.org.uk/acsb"
             xmlns:ns3="http://www.ifopt.org.uk/ifopt"
             xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
             xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:RequestTimestamp>2016-09-22T08:01:20.227+02:00</ns5:RequestTimestamp>
              <ns5:RequestorRef>NINOXE:default</ns5:RequestorRef>
              <ns5:MessageIdentifier>TermSubReq:Test:0</ns5:MessageIdentifier>
            </DeleteSubscriptionInfo>
            <Request>
              <ns5:SubscriptionRef>ExternalId</ns5:SubscriptionRef>
            </Request>
            <RequestExtension/>
          </ns1:DeleteSubscription>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
    """
<?xml version='1.0' encoding='UTF-8'?>
<S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
  <S:Body>
    <sw:DeleteSubscriptionResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
      <DeleteSubscriptionAnswerInfo>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>ara</siri:ResponderRef>
        <siri:RequestMessageRef>TermSubReq:Test:0</siri:RequestMessageRef>
      </DeleteSubscriptionAnswerInfo>
      <Answer>
        <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
        <siri:ResponderRef>ara</siri:ResponderRef>
        <siri:RequestMessageRef>TermSubReq:Test:0</siri:RequestMessageRef>
        <siri:TerminationResponseStatus>
          <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
          <siri:SubscriberRef>NINOXE:default</siri:SubscriberRef>
          <siri:SubscriptionRef>ExternalId</siri:SubscriptionRef>
          <siri:Status>false</siri:Status>
          <siri:ErrorCondition>
            <siri:CapabilityNotSupportedError>
              <siri:ErrorText>Subscription Termination is disabled for this Subscriber</siri:ErrorText>
            </siri:CapabilityNotSupportedError>
          </siri:ErrorCondition>
        </siri:TerminationResponseStatus>
      </Answer>
      <AnswerExtension/>
    </sw:DeleteSubscriptionResponse>
  </S:Body>
</S:Envelope>
    """
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:StopPoint:SP:24:LOC |
    Then an audit event should exist with these attributes:
      | Type                    | DeleteSubscriptionRequest                                |
      | Protocol                | siri                                                     |
      | Direction               | received                                                 |
      | Status                  | Error                                                    |
      | Partner                 | test                                                     |
      | ErrorDetails            | Subscription Termination is disabled for this Subscriber |
      | SubscriptionIdentifiers | ["ExternalId"]                                           |

  @ARA-1267
  Scenario: Ignore DeleteSubscription if setting broadcast.siri.ignore_terminate_subscription_requests is set to true for a TerminateSubscriptionRequest for All subscriptions
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url                                            | http://localhost:8090 |
      | remote_credential                                     | ara                   |
      | local_credential                                      | NINOXE:default        |
      | remote_code_space                                  | internal              |
      | siri.envelope                                         | raw                   |
      | broadcast.siri.ignore_terminate_subscription_requests | true                  |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | AnotherExternalId                     |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:A:BUS" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast            |
      | ExternalId        | SpecialExternalId                      |
      | SubscriberRef     | subscriber                             |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:C:Tram" |
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
      <SubscriptionRef>AnotherExternalId</SubscriptionRef>
      <Status>false</Status>
      <ErrorCondition>
        <CapabilityNotSupportedError>
          <ErrorText>Subscription Termination is disabled for this Subscriber</ErrorText>
        </CapabilityNotSupportedError>
      </ErrorCondition>
    </TerminationResponseStatus>
    <TerminationResponseStatus>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>false</Status>
      <ErrorCondition>
        <CapabilityNotSupportedError>
          <ErrorText>Subscription Termination is disabled for this Subscriber</ErrorText>
        </CapabilityNotSupportedError>
      </ErrorCondition>
    </TerminationResponseStatus>
    <TerminationResponseStatus>
      <SubscriptionRef>SpecialExternalId</SubscriptionRef>
      <Status>false</Status>
      <ErrorCondition>
        <CapabilityNotSupportedError>
          <ErrorText>Subscription Termination is disabled for this Subscriber</ErrorText>
        </CapabilityNotSupportedError>
      </ErrorCondition>
    </TerminationResponseStatus>
  </TerminateSubscriptionResponse>
</Siri>
      """
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:Line:3:LOC  |
      | internal | NINOXE:Line:A:BUS  |
      | internal | NINOXE:Line:C:Tram |
    Then an audit event should exist with these attributes:
      | Type         | DeleteSubscriptionRequest                                |
      | Protocol     | siri                                                     |
      | Direction    | received                                                 |
      | Status       | Error                                                    |
      | Partner      | test                                                     |
      | ErrorDetails | Subscription Termination is disabled for this Subscriber |


  @ARA-1066
  Scenario: Remove all subscriptions at once with termination request using raw envelope
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ara                   |
      | local_credential     | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope        | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | AnotherExternalId                     |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:A:BUS" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast            |
      | ExternalId        | SpecialExternalId                      |
      | SubscriberRef     | subscriber                             |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:C:Tram" |
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
      <SubscriptionRef>AnotherExternalId</SubscriptionRef>
      <Status>true</Status>
    </TerminationResponseStatus>
    <TerminationResponseStatus>
      <SubscriptionRef>externalId</SubscriptionRef>
      <Status>true</Status>
    </TerminationResponseStatus>
    <TerminationResponseStatus>
      <SubscriptionRef>SpecialExternalId</SubscriptionRef>
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
      | remote_code_space | internal              |
      | siri.envelope        | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
      | ExternalId        | externalId                            |
      | SubscriberRef     | subscriber                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
     And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
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
    Then No Subscriptions exist with the following resources:
      | internal | NINOXE:Line:3:LOC |
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:Line:A:BUS |

  @ARA-1066
  Scenario: Handle TerminateSubscriptionRequest on an unnown subscription using raw envelope
    Given a raw SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-estimated-timetable-subscription-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ara                   |
      | local_credential     | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope        | raw                   |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast            |
      | ExternalId        | externalId                             |
      | SubscriberRef     | subscriber                             |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC"  |
      | ReferenceArray[1] | Line, "internal": "NINOXE:Line:C:Tram" |
     And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableBroadcast           |
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
        <UnknownSubscriptionError>
          <ErrorText>Subscription not found: 'UnknownExternalId'</ErrorText>
        </UnknownSubscriptionError>
      </ErrorCondition>
    </TerminationResponseStatus>
  </TerminateSubscriptionResponse>
</Siri>
     """
    Then Subscriptions exist with the following resources:
      | internal | NINOXE:Line:A:BUS  |
      | internal | NINOXE:Line:3:LOC  |
      | internal | NINOXE:Line:C:Tram |

  @ARA-1432
  Scenario: Accept response to an EstimatedTimetable subscripτion with only one Subscription with missingg RequestMessageRef in ResponseStatus
    Given a raw SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
<SubscriptionResponse>
        <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
        <ResponderRef>NINOXE:default</ResponderRef>
        <ResponseStatus>
            <ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ResponseTimestamp>
            <SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</SubscriptionRef>
            <Status>true</Status>
            <ValidUntil>2016-09-22T08:01:20.227+02:00</ValidUntil>
        </ResponseStatus>
        <ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ServiceStartedTime>
</SubscriptionResponse>
</Siri>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | local_credential      | NINOXE:default        |
      | remote_code_space     | internal              |
      | siri.envelope         | raw                   |
      | collect.include_lines | NINOXE:Line:3:LOC     |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name   | Test                            |
      | Codes  | "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | EstimatedTimetableCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z    |

  @ARA-1476
  Scenario: Accept response to a StopMonitoring subscription with only one Subscription with missing RequestMessageRef in ResponseStatus
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <SubscriptionAnswerInfo
        xmlns:ns2="http://www.ifopt.org.uk/acsb"
        xmlns:ns3="http://www.ifopt.org.uk/ifopt"
        xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns5="http://www.siri.org.uk/siri"
        xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
        <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
        <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
        <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">Subscription:Test:0</ns5:RequestMessageRef>
      </SubscriptionAnswerInfo>
      <Answer
        xmlns:ns2="http://www.ifopt.org.uk/acsb"
        xmlns:ns3="http://www.ifopt.org.uk/ifopt"
        xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns5="http://www.siri.org.uk/siri"
        xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
            <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</ns5:SubscriptionRef>
            <ns5:Status>true</ns5:Status>
            <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
        </ns5:ResponseStatus>
        <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
      </Answer>
      <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
      </ns1:SubscribeResponse>
      </S:Body>
      </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url                 | http://localhost:8090 |
      | remote_credential          | test                  |
      | local_credential           | NINOXE:default        |
      | remote_code_space          | internal              |
      | collect.include_stop_areas | NINOXE:StopArea:A:LOC |
    And 30 seconds have passed
    And a StopArea exists with the following attributes:
      | Name  | Test                                 |
      | Codes | "internal": "NINOXE:StopArea:A:LOC"  |
    And a Subscription exist with the following attributes:
      | Kind              | StopMonitoringCollect                         |
      | SubscriberRef     | subscriber                                    |
      | ExternalId        | externalId                                    |
      | ReferenceArray[0] | StopArea, "internal": "NINOXE:StopArea:A:LOC" |
    And a minute has passed
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | StopMonitoringCollect  |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z |

  @ARA-1476
  Scenario: Accept response to a SituationExchange subscription with only one Subscription with missing RequestMessageRef in ResponseStatus
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
      <ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <SubscriptionAnswerInfo
        xmlns:ns2="http://www.ifopt.org.uk/acsb"
        xmlns:ns3="http://www.ifopt.org.uk/ifopt"
        xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns5="http://www.siri.org.uk/siri"
        xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
        <ns5:Address>http://appli.chouette.mobi/siri_france/siri</ns5:Address>
        <ns5:ResponderRef>NINOXE:default</ns5:ResponderRef>
        <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">Subscription:Test:0</ns5:RequestMessageRef>
      </SubscriptionAnswerInfo>
      <Answer
        xmlns:ns2="http://www.ifopt.org.uk/acsb"
        xmlns:ns3="http://www.ifopt.org.uk/ifopt"
        xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
        xmlns:ns5="http://www.siri.org.uk/siri"
        xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
            <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
            <ns5:SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</ns5:SubscriptionRef>
            <ns5:Status>true</ns5:Status>
            <ns5:ValidUntil>2016-09-22T08:01:20.227+02:00</ns5:ValidUntil>
        </ns5:ResponseStatus>
        <ns5:ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ns5:ServiceStartedTime>
      </Answer>
      <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
      </ns1:SubscribeResponse>
      </S:Body>
      </S:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-situation-exchange-subscription-collector] and the following settings:
      | remote_url                 | http://localhost:8090 |
      | remote_credential          | test                  |
      | local_credential           | NINOXE:default        |
      | remote_code_space          | internal              |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name  | Test                            |
      | Codes | "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
        | Kind              | SituationExchangeCollect               |
        | ReferenceArray[0] | Line, "internal": "NINOXE:Line::3:LOC" |
    And a minute has passed
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | SituationExchangeCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z   |

  @ARA-1454
  Scenario: Handle TerminateSubscriptionRequest when receiving an notification with unknown subscription using raw envelope
    Given a raw SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
<SubscriptionResponse>
        <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
        <ResponderRef>NINOXE:default</ResponderRef>
        <ResponseStatus>
            <ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ResponseTimestamp>
            <RequestMessageRef>{LastRequestMessageRef}</RequestMessageRef>
            <SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</SubscriptionRef>
            <Status>true</Status>
            <ValidUntil>2016-09-22T08:01:20.227+02:00</ValidUntil>
        </ResponseStatus>
        <ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ServiceStartedTime>
</SubscriptionResponse>
</Siri>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-estimated-timetable-subscription-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name  | Test                            |
      | Codes | "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | EstimatedTimetableCollect             |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | EstimatedTimetableCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z    |
    When I send this SIRI request
        """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns='http://www.siri.org.uk/siri'>
 <ServiceDelivery>
    <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
    <ProducerRef>NINOXE:default</ProducerRef>
    <ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ResponseMessageIdentifier>
    <EstimatedTimetableDelivery>
      <ResponseTimestamp>2017-01-01T12:00:20.000Z</ResponseTimestamp>
      <SubscriberRef>subscriber</SubscriberRef>
      <SubscriptionRef>DUMMY</SubscriptionRef>
      <Status>true</Status>
      <EstimatedJourneyVersionFrame>
        <RecordedAtTime>2017-01-01T12:00:20.000Z</RecordedAtTime>
        <EstimatedVehicleJourney>
          <LineRef>NINOXE:Line:3:LOC</LineRef>
          <DirectionRef>Aller</DirectionRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</DataFrameRef>
            <DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <DestinationRef>ThisIsTheEnd</DestinationRef>
          <OperatorRef>CdF:Company::410:LOC</OperatorRef>
          <EstimatedCalls>
            <EstimatedCall>
              <StopPointRef>NINOXE:StopPoint:SP:24:LOC</StopPointRef>
              <Order>4</Order>
              <StopPointName>Test</StopPointName>
              <ExpectedArrivalTime>2017-01-01T15:01:01.000Z</ExpectedArrivalTime>
              <ArrivalStatus>delayed</ArrivalStatus>
            </EstimatedCall>
          </EstimatedCalls>
        </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
    </EstimatedTimetableDelivery>
 </ServiceDelivery>
</Siri>
      """
    Then the SIRI server should receive this response
     """
<?xml version='1.0' encoding='utf-8'?>
<Siri xmlns='http://www.siri.org.uk/siri' version='2.0'>
<TerminateSubscriptionRequest>
	<RequestTimestamp>2017-01-01T12:01:30.000Z</RequestTimestamp>
	<RequestorRef>test</RequestorRef>
	<MessageIdentifier>6ba7b814-9dad-11d1-8-00c04fd430c8</MessageIdentifier>
	<SubscriptionRef>DUMMY</SubscriptionRef>
</TerminateSubscriptionRequest>
</Siri>
     """
    Then an audit event should exist with these attributes:
      | Type                    | DeleteSubscriptionRequest |
      | Protocol                | siri                      |
      | Direction               | sent                      |
      | Partner                 | test                      |
      | SubscriptionIdentifiers | ["DUMMY"]                 |

  @ARA-1458
  Scenario: Accept response to a VehicleMonitoring subscripτion with only one Subscription with missing RequestMessageRef in ResponseStatus
    Given a raw SIRI server waits Subscribe request on "http://localhost:8090" to respond with
      """
<?xml version="1.0" encoding="utf-8"?>
<Siri xmlns="http://www.siri.org.uk/siri" version="2.0">
<SubscriptionResponse>
        <ResponseTimestamp>2017-01-01T12:01:00.000Z</ResponseTimestamp>
        <ResponderRef>NINOXE:default</ResponderRef>
        <ResponseStatus>
            <ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ResponseTimestamp>
            <SubscriptionRef>6ba7b814-9dad-11d1-4-00c04fd430c8</SubscriptionRef>
            <Status>true</Status>
            <ValidUntil>2016-09-22T08:01:20.227+02:00</ValidUntil>
        </ResponseStatus>
        <ServiceStartedTime>2016-09-22T08:01:20.227+02:00</ServiceStartedTime>
</SubscriptionResponse>
</Siri>
      """
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-vehicle-monitoring-subscription-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | local_credential      | NINOXE:default        |
      | remote_code_space     | internal              |
      | siri.envelope         | raw                   |
      | collect.include_lines | NINOXE:Line:3:LOC     |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Name   | Test                            |
      | Codes  | "internal": "NINOXE:Line:3:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | VehicleMonitoringCollect              |
      | SubscriberRef     | subscriber                            |
      | ExternalId        | externalId                            |
      | ReferenceArray[0] | Line, "internal": "NINOXE:Line:3:LOC" |
    And a minute has passed
    And a minute has passed
    Then one Subscription exists with the following attributes:
      | Kind                      | VehicleMonitoringCollect |
      | Resources[0]/SubscribedAt | > 2017-01-01T12:01:00Z   |
