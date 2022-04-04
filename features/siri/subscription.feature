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
