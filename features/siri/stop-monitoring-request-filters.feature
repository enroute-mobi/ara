Feature: Support SIRI StopMonitoring by request

  Background:
      Given a Referential "test" is created

  Scenario: 3754a - Handle a SIRI StopMonitoring request with filter PreviewInterval
    # si StartTime absent, alors ça démarre à l'heure courante (=celle de la request)
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                        |
      | Codes[internal] | NINOXE:VehicleJourney:201         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                              |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T14:00:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T14:30:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:00:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:30:00.000+02:00                          |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:siri="http://www.ifopt.org.uk/acsb"
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
          <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2017-01-01T14:30:00.000+02:00</siri:AimedArrivalTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2017-01-01T15:00:00.000+02:00</siri:AimedArrivalTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3754b - Handle a SIRI StopMonitoring request with filter PreviewInterval and StartTime
    # si StartTime absent, alors ça démarre à l'heure courante (=celle de la request)
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                        |
      | Codes[internal] | NINOXE:VehicleJourney:201         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                              |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T13:00:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T13:30:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T14:00:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T14:30:00.000+02:00                          |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:siri="http://www.ifopt.org.uk/acsb"
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
          <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2017-01-01T14:30:00.000+02:00</siri:AimedArrivalTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3754c - Handle a SIRI StopMonitoring request with filter LineRef
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:4:LOC |
      | Name            | Ligne 4 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                        |
      | Codes[internal] | NINOXE:VehicleJourney:201         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                              |
    And a VehicleJourney exists with the following attributes:
      | Name            | Le 15.                            |
      | Codes[internal] | NINOXE:VehicleJourney:202         |
      | LineId          | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored       | true                              |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-1 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T13:00:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-2 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T13:30:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T14:02:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-6-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T14:30:00.000+02:00                          |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:siri="http://www.ifopt.org.uk/acsb"
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
          <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-b-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 4 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Le 15.</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2017-01-01T14:02:00.000+02:00</siri:AimedArrivalTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 4 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Le 15.</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2017-01-01T14:30:00.000+02:00</siri:AimedArrivalTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3754d - Handle a SIRI StopMonitoring request with filter StopVisitTypes "departure"
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:4:LOC |
      | Name            | Ligne 4 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name            | Le 15.                            |
      | Codes[internal] | NINOXE:VehicleJourney:202         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                              |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:01:00.000+02:00                          |
      | Schedule[aimed]#Departure   | 2017-01-01T15:02:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
    # StopVisit correspondant au terminus, donc un arrival = departure
      | Codes[internal]             | NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:30:00.000+02:00                          |
      | Schedule[aimed]#Departure   | 2017-01-01T15:30:00.000+02:00                          |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:siri="http://www.ifopt.org.uk/acsb"
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
          <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 4 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Le 15.</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedDepartureTime>2017-01-01T15:02:00.000+02:00</siri:AimedDepartureTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 4 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Le 15.</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedDepartureTime>2017-01-01T15:30:00.000+02:00</siri:AimedDepartureTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """


  Scenario: 3754e - Handle a SIRI StopMonitoring request with filter StopVisitTypes "arrivals"
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:4:LOC |
      | Name            | Ligne 4 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name            | Le 15.                            |
      | Codes[internal] | NINOXE:VehicleJourney:202         |
      | LineId          | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Monitored       | true                              |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:01:00.000+02:00                          |
      | Schedule[aimed]#Departure   | 2017-01-01T15:02:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
    # StopVisit correspondant au point de départ de la course, donc un arrival = departure
      | Codes[internal]             | NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-5-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T15:30:00.000+02:00                          |
      | Schedule[aimed]#Departure   | 2017-01-01T15:30:00.000+02:00                          |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header />
        <S:Body>
          <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                                 xmlns:siri="http://www.ifopt.org.uk/acsb"
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
              <ns2:StopVisitTypes>arrivals</ns2:StopVisitTypes>
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
          <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 4 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Le 15.</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2017-01-01T15:01:00.000+02:00</siri:AimedArrivalTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                  <siri:ItemIdentifier>NINOXE:VehicleJourney:202-NINOXE:StopPoint:SP:24:LOC-4</siri:ItemIdentifier>
                  <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>NINOXE:Line:4:LOC</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:202</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:PublishedLineName>Ligne 4 Metro</siri:PublishedLineName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:VehicleJourneyName>Le 15.</siri:VehicleJourneyName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                      <siri:Order>4</siri:Order>
                      <siri:StopPointName>Test</siri:StopPointName>
                      <siri:VehicleAtStop>true</siri:VehicleAtStop>
                      <siri:AimedArrivalTime>2017-01-01T15:30:00.000+02:00</siri:AimedArrivalTime>
                    </siri:MonitoredCall>
                  </siri:MonitoredVehicleJourney>
                </siri:MonitoredStopVisit>
              </siri:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetStopMonitoringResponse>
        </S:Body>
      </S:Envelope>
      """

  Scenario: 3754f - Handle a SIRI StopMonitoring request with filter MaximumStopVisits
    # si StartTime absent, alors ça démarre à l'heure courante (=celle de la request)
    Given a SIRI Partner "test" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a StopArea exists with the following attributes:
      | Name            | Test                       |
      | Codes[internal] | NINOXE:StopPoint:SP:24:LOC |
      | Monitored       | true                       |
    And a Line exists with the following attributes:
      | Codes[internal] | NINOXE:Line:3:LOC |
      | Name            | Ligne 3 Metro     |
    And a VehicleJourney exists with the following attributes:
      | Name            | Passage 32                        |
      | Codes[internal] | NINOXE:VehicleJourney:201         |
      | LineId          | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Monitored       | true                              |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T16:00:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T16:30:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-5 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T17:00:00.000+02:00                          |
    And a StopVisit exists with the following attributes:
      | Codes[internal]             | NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-6 |
      | PassageOrder                | 4                                                      |
      | StopAreaId                  | 6ba7b814-9dad-11d1-2-00c04fd430c8                      |
      | VehicleJourneyId            | 6ba7b814-9dad-11d1-4-00c04fd430c8                      |
      | VehicleAtStop               | true                                                   |
      | Reference[OperatorRef]#Code | "internal": "CdF:Company::410:LOC"                     |
      | Schedule[aimed]#Arrival     | 2017-01-01T18:30:00.000+02:00                          |
    When I send this SIRI request
      """
  <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header />
  <S:Body>
    <ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:siri="http://www.ifopt.org.uk/acsb"
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
        <ns2:StopVisitTypes>all</ns2:StopVisitTypes>
        <ns2:MaximumStopVisits>2</ns2:MaximumStopVisits>
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
        <sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceDeliveryInfo>
            <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
            <siri:ProducerRef>Ara</siri:ProducerRef>
            <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
          </ServiceDeliveryInfo>
          <Answer>
            <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:RequestMessageRef>StopMonitoring:Test:0</siri:RequestMessageRef>
              <siri:Status>true</siri:Status>
              <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
              <siri:MonitoredStopVisit>
                <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3</siri:ItemIdentifier>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                  <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:MonitoredCall>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <siri:Order>4</siri:Order>
                    <siri:StopPointName>Test</siri:StopPointName>
                    <siri:VehicleAtStop>true</siri:VehicleAtStop>
                    <siri:AimedArrivalTime>2017-01-01T16:00:00.000+02:00</siri:AimedArrivalTime>
                  </siri:MonitoredCall>
                </siri:MonitoredVehicleJourney>
              </siri:MonitoredStopVisit>
              <siri:MonitoredStopVisit>
                <siri:RecordedAtTime>0001-01-01T00:00:00.000Z</siri:RecordedAtTime>
                <siri:ItemIdentifier>NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-4</siri:ItemIdentifier>
                <siri:MonitoringRef>NINOXE:StopPoint:SP:24:LOC</siri:MonitoringRef>
                <siri:MonitoredVehicleJourney>
                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                  <siri:FramedVehicleJourneyRef>
                    <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                    <siri:DatedVehicleJourneyRef>NINOXE:VehicleJourney:201</siri:DatedVehicleJourneyRef>
                  </siri:FramedVehicleJourneyRef>
                  <siri:PublishedLineName>Ligne 3 Metro</siri:PublishedLineName>
                  <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                  <siri:VehicleJourneyName>Passage 32</siri:VehicleJourneyName>
                  <siri:Monitored>true</siri:Monitored>
                  <siri:MonitoredCall>
                    <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                    <siri:Order>4</siri:Order>
                    <siri:StopPointName>Test</siri:StopPointName>
                    <siri:VehicleAtStop>true</siri:VehicleAtStop>
                    <siri:AimedArrivalTime>2017-01-01T16:30:00.000+02:00</siri:AimedArrivalTime>
                  </siri:MonitoredCall>
                </siri:MonitoredVehicleJourney>
              </siri:MonitoredStopVisit>
            </siri:StopMonitoringDelivery>
          </Answer>
          <AnswerExtension/>
        </sw:GetStopMonitoringResponse>
      </S:Body>
    </S:Envelope>
    """
