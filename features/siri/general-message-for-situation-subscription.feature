Feature: Support SIRI GeneralMessage by subscription

  Background:
      Given a Referential "test" is created

@wip
   Scenario: 3863 - Manage a GM Subscription
    Given a SIRI server waits Subscribe request on "http://localhost:8090" to respond with
        """
    ...
        """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-subscription-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | local_credential     | NINOXE:default        |
      | remote_objectid_kind | internal              |
    And a minute has passed
    When I send this SIRI request
        """
        <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7: GeneralMessageSubscriptionRequest xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <SubscriptionRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:00:00.000+02:00</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>GMSubsccription:Test:0</ns2:MessageIdentifier>
            </SubscriptionRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:00:00.000+02:00</ns2:RequestTimestamp>
              <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
              <ns5:SubscriptionRef>NINOXE:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns5:SubscriptionRef>
              <ns5:InitialTerminationTime>2017-02-01T12:00:00.000+02:00</ns5:InitialTerminationTime>
            </Request>
            <RequestExtension />
          </ns7:GetStopMonitoring>
        </S:Body>
      </S:Envelope>
        """
    Then Then I should receive this SIRI response
    """
    ...
    """

