Feature: Support SIRI StopMonitoring by request

  Background:
      Given a Referential "test" is created

  Scenario: 3754a - Handle a SIRI StopMonitoring request with filter PreviewInterval
    # si StartTime absent, alors ça démarre à l'heure courante (=celle de la request)
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T14:00:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T14:30:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:00:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:30:00.000+02:00                                        |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <ServiceRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
            </ServiceRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
              <ns2:PreviewInterval>PT1H</ns2:PreviewInterval>
              <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
              <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
            </Request>
            <RequestExtension />
          </ns7:GetStopMonitoring>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>Edwig</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                    <ns3:VehicleJourneyName>Passage 32</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Test</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T14:30:00.000+02:00</ns3:AimedArrivalTime>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                    <ns3:VehicleJourneyName>Passage 32</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Test</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T15:00:00.000+02:00</ns3:AimedArrivalTime>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns3:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3754b - Handle a SIRI StopMonitoring request with filter PreviewInterval and StartTime
    # si StartTime absent, alors ça démarre à l'heure courante (=celle de la request)
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:00:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:30:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T14:00:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T14:30:00.000+02:00                                        |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <ServiceRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
            </ServiceRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
              <ns2:PreviewInterval>PT1H</ns2:PreviewInterval>
              <ns2:StartTime>2017-01-01T13:35:00.000+02:00</ns2:StartTime>
              <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
              <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
            </Request>
            <RequestExtension />
          </ns7:GetStopMonitoring>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>Edwig</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                    <ns3:VehicleJourneyName>Passage 32</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Test</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T14:30:00.000+02:00</ns3:AimedArrivalTime>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns3:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3754c - Handle a SIRI StopMonitoring request with filter LineRef
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:4:LOC" |
      | Name      | Ligne 4 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Passage 32                              |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
      | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
    And a VehicleJourney exists with the following attributes:
      | Name      | Le 15.                                  |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:202" |
      | LineId    | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:00:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-5-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T13:30:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-6-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T14:02:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-6-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T14:30:00.000+02:00                                        |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <ServiceRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
            </ServiceRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
              <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
              <ns2:LineRef>NINOXE:Line:4:LOC</ns2:LineRef>
              <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
            </Request>
            <RequestExtension />
          </ns7:GetStopMonitoring>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>Edwig</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-b-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:4:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 4 Metro</ns3:PublishedLineName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                    <ns3:VehicleJourneyName>Le 15.</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Test</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T14:02:00.000+02:00</ns3:AimedArrivalTime>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:4:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 4 Metro</ns3:PublishedLineName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                    <ns3:VehicleJourneyName>Le 15.</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Test</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedArrivalTime>2017-01-01T14:30:00.000+02:00</ns3:AimedArrivalTime>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns3:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

@wip
  Scenario: 3754d - Handle a SIRI StopMonitoring request with filter StopVisitTypes "departure"
    Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | test     |
      | remote_objectid_kind | internal |
    And a StopArea exists with the following attributes:
      | Name      | Test                                     |
      | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
      | Name      | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | ObjectIDs | "internal": "NINOXE:Line:4:LOC" |
      | Name      | Ligne 4 Metro                   |
    And a VehicleJourney exists with the following attributes:
      | Name      | Le 15.                                  |
      | ObjectIDs | "internal": "NINOXE:VehicleJourney:202" |
      | LineId    | 6ba7b814-9dad-11d1-4-00c04fd430c8       |
    And a StopVisit exists with the following attributes:
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-6-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:01:00.000+02:00                                        |
      | Schedule[aimed]#Departure       | 2017-01-01T15:02:00.000+02:00                                        |
    And a StopVisit exists with the following attributes:
    # StopVisit correspondant au terminus, donc un arrival mais pas de departure
      | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4" |
      | PassageOrder                    | 4                                                                    |
      | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
      | VehicleJourneyId                | 6ba7b814-9dad-11d1-6-00c04fd430c8                                    |
      | VehicleAtStop                   | true                                                                 |
      | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
      | Schedule[aimed]#Arrival         | 2017-01-01T15:30:00.000+02:00                                        |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:ns3="http://www.ifopt.org.uk/acsb"
                                 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                                 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                                 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <ServiceRequestInfo>
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:RequestorRef>test</ns2:RequestorRef>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
            </ServiceRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
              <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
              <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
              <ns2:StopVisitTypes>departures</ns2:StopVisitTypes>
            </Request>
            <RequestExtension />
          </ns7:GetStopMonitoring>
        </S:Body>
      </S:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='utf-8'?>
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
          xmlns:ns4="http://www.ifopt.org.uk/acsb"
          xmlns:ns5="http://www.ifopt.org.uk/ifopt"
          xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
          xmlns:ns7="http://scma/siri"
          xmlns:ns8="http://wsdl.siri.org.uk"
          xmlns:ns9="http://wsdl.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
              <ns3:ProducerRef>Edwig</ns3:ProducerRef>
              <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
              <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
                <ns3:Status>true</ns3:Status>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:4:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 4 Metro</ns3:PublishedLineName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                    <ns3:VehicleJourneyName>Le 15.</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Test</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedDepartureTime>2017-01-01T15:02:00.000+02:00</ns3:AimedDepartureTime>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
                <ns3:MonitoredStopVisit>
                  <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                  <ns3:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4</ns3:ItemIdentifier>
                  <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                  <ns3:MonitoredVehicleJourney>
                    <ns3:LineRef>NINOXE:Line:4:LOC</ns3:LineRef>
                    <ns3:FramedVehicleJourneyRef>
                      <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                      <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</ns3:DatedVehicleJourneyRef>
                    </ns3:FramedVehicleJourneyRef>
                    <ns3:PublishedLineName>Ligne 4 Metro</ns3:PublishedLineName>
                    <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                    <ns3:VehicleJourneyName>Le 15.</ns3:VehicleJourneyName>
                    <ns3:MonitoredCall>
                      <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                      <ns3:Order>4</ns3:Order>
                      <ns3:StopPointName>Test</ns3:StopPointName>
                      <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                      <ns3:AimedDepartureTime>2017-01-01T14:30:00.000+02:00</ns3:AimedDepartureTime>
                    </ns3:MonitoredCall>
                  </ns3:MonitoredVehicleJourney>
                </ns3:MonitoredStopVisit>
              </ns3:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension />
          </ns8:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

      Scenario: 3754f - Handle a SIRI StopMonitoring request with filter MaximumStopVisits
        # si StartTime absent, alors ça démarre à l'heure courante (=celle de la request)
        Given a Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
          | local_credential     | test     |
          | remote_objectid_kind | internal |
        And a StopArea exists with the following attributes:
          | Name      | Test                                     |
          | ObjectIDs | "internal": "NINOXE:StopPoint:SP:24:LOC" |
        And a Line exists with the following attributes:
          | ObjectIDs | "internal": "NINOXE:Line:3:LOC" |
          | Name      | Ligne 3 Metro                   |
        And a VehicleJourney exists with the following attributes:
          | Name      | Passage 32                              |
          | ObjectIDs | "internal": "NINOXE:VehicleJourney:201" |
          | LineId    | 6ba7b814-9dad-11d1-3-00c04fd430c8       |
        And a StopVisit exists with the following attributes:
          | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
          | PassageOrder                    | 4                                                                    |
          | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
          | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
          | VehicleAtStop                   | true                                                                 |
          | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
          | Schedule[aimed]#Arrival         | 2017-01-01T16:00:00.000+02:00                                        |
        And a StopVisit exists with the following attributes:
          | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
          | PassageOrder                    | 4                                                                    |
          | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
          | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
          | VehicleAtStop                   | true                                                                 |
          | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
          | Schedule[aimed]#Arrival         | 2017-01-01T16:30:00.000+02:00                                        |
        And a StopVisit exists with the following attributes:
          | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
          | PassageOrder                    | 4                                                                    |
          | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
          | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
          | VehicleAtStop                   | true                                                                 |
          | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
          | Schedule[aimed]#Arrival         | 2017-01-01T17:00:00.000+02:00                                        |
        And a StopVisit exists with the following attributes:
          | ObjectIDs                       | "internal": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3" |
          | PassageOrder                    | 4                                                                    |
          | StopAreaId                      | 6ba7b814-9dad-11d1-2-00c04fd430c8                                    |
          | VehicleJourneyId                | 6ba7b814-9dad-11d1-4-00c04fd430c8                                    |
          | VehicleAtStop                   | true                                                                 |
          | Reference[OperatorRef]#ObjectID | "internal": "CdF:Company::410:LOC"                                   |
          | Schedule[aimed]#Arrival         | 2017-01-01T18:30:00.000+02:00                                        |
        When I send this SIRI request
          """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
      <SOAP-ENV:Header />
      <S:Body>
        <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                               xmlns:ns3="http://www.ifopt.org.uk/acsb"
                               xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                               xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                               xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
          <ServiceRequestInfo>
            <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
            <ns2:RequestorRef>test</ns2:RequestorRef>
            <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
          </ServiceRequestInfo>

          <Request version="2.0:FR-IDF-2.4">
            <ns2:RequestTimestamp>2017-01-01T12:26:10.116+02:00</ns2:RequestTimestamp>
            <ns2:MessageIdentifier>StopMonitoring:Test:0</ns2:MessageIdentifier>
            <ns2:MaximumStopVisits>2</ns2:MaximumStopVisits>
            <ns2:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns2:MonitoringRef>
            <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
          </Request>
          <RequestExtension />
        </ns7:GetStopMonitoring>
      </S:Body>
      </S:Envelope>
          """
        Then I should receive this SIRI response
        """
        <?xml version='1.0' encoding='utf-8'?>
        <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
          <S:Body>
            <ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
            xmlns:ns4="http://www.ifopt.org.uk/acsb"
            xmlns:ns5="http://www.ifopt.org.uk/ifopt"
            xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
            xmlns:ns7="http://scma/siri"
            xmlns:ns8="http://wsdl.siri.org.uk"
            xmlns:ns9="http://wsdl.siri.org.uk/siri">
              <ServiceDeliveryInfo>
                <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                <ns3:ProducerRef>Edwig</ns3:ProducerRef>
                <ns3:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</ns3:ResponseMessageIdentifier>
                <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
              </ServiceDeliveryInfo>
              <Answer>
                <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                  <ns3:ResponseTimestamp>2017-01-01T12:00:00.000Z</ns3:ResponseTimestamp>
                  <ns3:RequestMessageRef>StopMonitoring:Test:0</ns3:RequestMessageRef>
                  <ns3:Status>true</ns3:Status>
                  <ns3:MonitoredStopVisit>
                    <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                    <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                    <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                    <ns3:MonitoredVehicleJourney>
                      <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                      <ns3:FramedVehicleJourneyRef>
                        <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                        <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                      </ns3:FramedVehicleJourneyRef>
                      <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                      <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                      <ns3:VehicleJourneyName>Passage 32</ns3:VehicleJourneyName>
                      <ns3:MonitoredCall>
                        <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                        <ns3:Order>4</ns3:Order>
                        <ns3:StopPointName>Test</ns3:StopPointName>
                        <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                        <ns3:AimedArrivalTime>2017-01-01T16:00:00.000+02:00</ns3:AimedArrivalTime>
                      </ns3:MonitoredCall>
                    </ns3:MonitoredVehicleJourney>
                  </ns3:MonitoredStopVisit>
                  <ns3:MonitoredStopVisit>
                    <ns3:RecordedAtTime>0001-01-01T00:00:00.000Z</ns3:RecordedAtTime>
                    <ns3:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</ns3:ItemIdentifier>
                    <ns3:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</ns3:MonitoringRef>
                    <ns3:MonitoredVehicleJourney>
                      <ns3:LineRef>NINOXE:Line:3:LOC</ns3:LineRef>
                      <ns3:FramedVehicleJourneyRef>
                        <ns3:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</ns3:DataFrameRef>
                        <ns3:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</ns3:DatedVehicleJourneyRef>
                      </ns3:FramedVehicleJourneyRef>
                      <ns3:PublishedLineName>Ligne 3 Metro</ns3:PublishedLineName>
                      <ns3:OperatorRef>RATPDev:Operator::9901377d84631ed7c2c09bbb32d70effaee59cc0:</ns3:OperatorRef>
                      <ns3:VehicleJourneyName>Passage 32</ns3:VehicleJourneyName>
                      <ns3:MonitoredCall>
                        <ns3:StopPointRef>NINOXE:StopPoint:SP:24:LOC</ns3:StopPointRef>
                        <ns3:Order>4</ns3:Order>
                        <ns3:StopPointName>Test</ns3:StopPointName>
                        <ns3:VehicleAtStop>true</ns3:VehicleAtStop>
                        <ns3:AimedArrivalTime>2017-01-01T16:30:00.000+02:00</ns3:AimedArrivalTime>
                      </ns3:MonitoredCall>
                    </ns3:MonitoredVehicleJourney>
                  </ns3:MonitoredStopVisit>
                </ns3:StopMonitoringDelivery>
              </Answer>
              <AnswerExtension />
            </ns8:GetStopMonitoringResponse>
          </S:Body>
        </S:Envelope>
        """
