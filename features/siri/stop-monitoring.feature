Feature: Support SIRI StopMonitoring

  Background:
      Given a Referential "test" is created

  Scenario: Handle a SIRI StopMonitoring response after SM Request to a SIRI server
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
        """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
  <soap:Body>
    <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:StopMonitoringDelivery version="1.3">
          <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
          <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
          <ns5:Status>true</ns5:Status>
          <ns5:MonitoredStopVisit>
            <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
            <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
              <ns5:PublishedLineName>415</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
              <ns5:DestinationRef>boabonn</ns5:DestinationRef>
              <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
              <ns5:Monitored>true</ns5:Monitored>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                <ns5:Order>44</ns5:Order>
                <ns5:StopPointName>Arletty</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2017-01-01T13:43:05.000+01:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2017-01-01T13:43:05.000+01:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2017-01-01T13:43:05.000+01:00</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2017-01-01T13:43:05.000+01:00</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney>
          </ns5:MonitoredStopVisit>
        </ns5:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
    </ns1:GetStopMonitoringResponse>
  </soap:Body>
</soap:Envelope>
        """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ratpdev               |
      | remote_objectid_kind | internal              |
    And a Partner "stif" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | STIF                                           |
      | remote_objectid_kind | external                                       |
      | remote_credential    | RATPDev                                        |
      | local_url            | https://api.concerto.ratpdev.com/concerto/siri |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Ligne 415                                                         |
      | ObjectIDs | "internal": "CdF:Line::415:LOC", "external": "STIF:Line::C00001:" |
    And a StopArea exists with the following attributes:
      | Name      | Arletty                                                                |
      | ObjectIDs | "internal": "boaarle", "external": "STIF:StopPoint:Q:eeft52df543d:" |
    And a StopArea exists with the following attributes:
      | Name            | Test 2                                                                  |
      | ObjectIDs       | "internal": "boabonn", "external": "STIF:StopPoint:Q:875fdetgyh765:" |
      | CollectedAlways | false                                                                   |
    And a minute has passed
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
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:RequestorRef>STIF</ns2:RequestorRef>
        <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
      </ServiceRequestInfo>

      <Request version="2.0:FR-IDF-2.4">
        <ns2:RequestTimestamp>2017-01-01T12:00:00.000Z</ns2:RequestTimestamp>
        <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
        <ns2:MonitoringRef>STIF:StopPoint:Q:eeft52df543d:</ns2:MonitoringRef>
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
        <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
        <siri:ProducerRef>RATPDev</siri:ProducerRef>
        <siri:Address>https://api.concerto.ratpdev.com/concerto/siri</siri:Address>
        <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-11-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
        <siri:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer>
        <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
          <siri:ResponseTimestamp>2017-01-01T12:02:00.000Z</siri:ResponseTimestamp>
          <siri:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
          <siri:MonitoringRef>STIF:StopPoint:Q:eeft52df543d:</siri:MonitoringRef>
          <siri:Status>true</siri:Status>
          <siri:MonitoredStopVisit>
            <siri:RecordedAtTime>2017-01-01T11:47:15.600+01:00</siri:RecordedAtTime>
            <siri:ItemIdentifier>RATPDev:Item::63a1ebcfe85a0e9a548c91a611cfb572f4a545af:LOC</siri:ItemIdentifier>
            <siri:MonitoringRef>STIF:StopPoint:Q:eeft52df543d:</siri:MonitoringRef>
            <siri:MonitoredVehicleJourney>
              <siri:LineRef>STIF:Line::C00001:</siri:LineRef>
              <siri:FramedVehicleJourneyRef>
                <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                <siri:DatedVehicleJourneyRef>RATPDev:VehicleJourney::760857fec357472a48012c5ca776d940f8c49681:LOC</siri:DatedVehicleJourneyRef>
              </siri:FramedVehicleJourneyRef>
              <siri:JourneyPatternRef>RATPDev:JourneyPattern::983a5c43233dc44a0ed956117ee55d257fea06eb:LOC</siri:JourneyPatternRef>
              <siri:PublishedLineName>Ligne 415</siri:PublishedLineName>
              <siri:DirectionName>Aller</siri:DirectionName>
              <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
              <siri:DestinationRef>STIF:StopPoint:Q:875fdetgyh765:</siri:DestinationRef>
              <siri:DestinationName>Méliès - Croix Bonnet</siri:DestinationName>
              <siri:Monitored>true</siri:Monitored>
              <siri:MonitoredCall>
                <siri:StopPointRef>STIF:StopPoint:Q:eeft52df543d:</siri:StopPointRef>
                <siri:Order>44</siri:Order>
                <siri:StopPointName>Arletty</siri:StopPointName>
                <siri:VehicleAtStop>false</siri:VehicleAtStop>
                <siri:DestinationDisplay>Méliès - Croix Bonnet</siri:DestinationDisplay>
                <siri:AimedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:AimedArrivalTime>
                <siri:ExpectedArrivalTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedArrivalTime>
                <siri:ArrivalStatus>onTime</siri:ArrivalStatus>
                <siri:AimedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:AimedDepartureTime>
                <siri:ExpectedDepartureTime>2017-01-01T13:43:05.000+01:00</siri:ExpectedDepartureTime>
                <siri:DepartureStatus>onTime</siri:DepartureStatus>
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

  Scenario: Handles invalid GetStopMonitoring response
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
        <html><title>Error</title></body>Error 500</body></html>
      """
    And a Partner "invalid" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | ObjectIDs | "internal": "dummy" |
    When a minute has passed
    And the SIRI server has received a GetStopMonitoring request
    Then a StopArea exists with the following attributes:
      | ObjectIDs   | "internal": "dummy" |
      | CollectedAt | -                   |

  Scenario: Handle a SIRI StopMonitoring response after SM cancellation from a SIRI server
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
 """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
  <soap:Body>
    <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
      <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
        <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
        <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
      </ServiceDeliveryInfo>
      <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:StopMonitoringDelivery version="1.3">
          <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
          <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
          <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
          <ns5:Status>true</ns5:Status>
          <ns5:MonitoredStopVisit>
            <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
            <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
            <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
            <ns5:MonitoredVehicleJourney>
              <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
              <ns5:FramedVehicleJourneyRef>
                <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
              </ns5:FramedVehicleJourneyRef>
              <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
              <ns5:PublishedLineName>415</ns5:PublishedLineName>
              <ns5:DirectionName>Aller</ns5:DirectionName>
              <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
              <ns5:DestinationRef>boabonn</ns5:DestinationRef>
              <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
              <ns5:Monitored>true</ns5:Monitored>
              <ns5:MonitoredCall>
                <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                <ns5:Order>44</ns5:Order>
                <ns5:StopPointName>Arletty</ns5:StopPointName>
                <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                <ns5:AimedArrivalTime>2017-01-01T12:43:05.000+00:00</ns5:AimedArrivalTime>
                <ns5:ExpectedArrivalTime>2017-01-01T12:43:05.000+00:00</ns5:ExpectedArrivalTime>
                <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                <ns5:AimedDepartureTime>2017-01-01T12:43:05.000</ns5:AimedDepartureTime>
                <ns5:ExpectedDepartureTime>2017-01-01T12:43:05.000</ns5:ExpectedDepartureTime>
                <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
              </ns5:MonitoredCall>
            </ns5:MonitoredVehicleJourney>
          </ns5:MonitoredStopVisit>
        </ns5:StopMonitoringDelivery>
      </Answer>
      <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
    </ns1:GetStopMonitoringResponse>
  </soap:Body>
</soap:Envelope>
        """
     And a Partner "ineo" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | ratpdev               |
      | remote_objectid_kind | internal              |
    And a Partner "stif" exists with connectors [siri-stop-monitoring-request-broadcaster] and the following settings:
      | local_credential     | STIF     |
      | remote_objectid_kind | external |
      | remote_credential    | RATPDev  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name      | Ligne 415                                                         |
      | ObjectIDs | "internal": "CdF:Line::415:LOC", "external": "STIF:Line::C00001:" |
    And a StopArea exists with the following attributes:
      | Name      | Arletty                                                                |
      | ObjectIDs | "internal": "boaarle", "external": "RATPDev:StopPoint:Q:eeft52df543d:" |
    And a StopArea exists with the following attributes:
      | Name            | Test 2                                                                  |
      | ObjectIDs       | "internal": "boabonn", "external": "RATPDev:StopPoint:Q:875fdetgyh765:" |
      | CollectedAlways | false                                                                   |
    And a minute has passed
    And the SIRI server waits GetStopMonitoring request to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:ArrivalStatus>cancelled</ns5:ArrivalStatus>
                      <ns5:DepartureStatus>cancelled</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And 2 minutes have passed
    When the SIRI server has received 2 GetStopMonitoring requests
    And I send this SIRI request
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
              <ns2:RequestTimestamp>2017-01-01T12:03:00.000Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>STIF</ns2:RequestorRef>
              <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
            </ServiceRequestInfo>

            <Request version="2.0:FR-IDF-2.4">
              <ns2:RequestTimestamp>2017-01-01T12:03:00.000Z</ns2:RequestTimestamp>
              <ns2:MessageIdentifier>STIF:Message::2345Fsdfrg35df:LOC</ns2:MessageIdentifier>
              <ns2:MonitoringRef>RATPDev:StopPoint:Q:eeft52df543d:</ns2:MonitoringRef>
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
              <siri:ResponseTimestamp>2017-01-01T12:04:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>RATPDev</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-17-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
                <siri:ResponseTimestamp>2017-01-01T12:04:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>STIF:Message::2345Fsdfrg35df:LOC</siri:RequestMessageRef>
                <siri:MonitoringRef>RATPDev:StopPoint:Q:eeft52df543d:</siri:MonitoringRef>
                <siri:Status>true</siri:Status>
                <siri:MonitoredStopVisit>
                  <siri:RecordedAtTime>2017-01-01T11:47:15.600+01:00</siri:RecordedAtTime>
                  <siri:ItemIdentifier>RATPDev:Item::63a1ebcfe85a0e9a548c91a611cfb572f4a545af:LOC</siri:ItemIdentifier>
                  <siri:MonitoringRef>RATPDev:StopPoint:Q:eeft52df543d:</siri:MonitoringRef>
                  <siri:MonitoredVehicleJourney>
                    <siri:LineRef>STIF:Line::C00001:</siri:LineRef>
                    <siri:FramedVehicleJourneyRef>
                      <siri:DataFrameRef>RATPDev:DataFrame::2017-01-01:LOC</siri:DataFrameRef>
                      <siri:DatedVehicleJourneyRef>RATPDev:VehicleJourney::760857fec357472a48012c5ca776d940f8c49681:LOC</siri:DatedVehicleJourneyRef>
                    </siri:FramedVehicleJourneyRef>
                    <siri:JourneyPatternRef>RATPDev:JourneyPattern::983a5c43233dc44a0ed956117ee55d257fea06eb:LOC</siri:JourneyPatternRef>
                    <siri:PublishedLineName>Ligne 415</siri:PublishedLineName>
                    <siri:DirectionName>Aller</siri:DirectionName>
                    <siri:OperatorRef>CdF:Company::410:LOC</siri:OperatorRef>
                    <siri:DestinationRef>RATPDev:StopPoint:Q:875fdetgyh765:</siri:DestinationRef>
                    <siri:DestinationName>Méliès - Croix Bonnet</siri:DestinationName>
                    <siri:Monitored>true</siri:Monitored>
                    <siri:MonitoredCall>
                      <siri:StopPointRef>RATPDev:StopPoint:Q:eeft52df543d:</siri:StopPointRef>
                      <siri:Order>44</siri:Order>
                      <siri:StopPointName>Arletty</siri:StopPointName>
                      <siri:VehicleAtStop>false</siri:VehicleAtStop>
                      <siri:DestinationDisplay>Méliès - Croix Bonnet</siri:DestinationDisplay>
                      <siri:ArrivalStatus>cancelled</siri:ArrivalStatus>
                      <siri:DepartureStatus>cancelled</siri:DepartureStatus>
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

  Scenario: Manage a passed StopVisit
    Given a SIRI server waits GetStopMonitoring request on "http://localhost:8090" to respond with
    # include a MonitoredStopVisit/ItemIdentifier A at 13:00
    # include a MonitoredStopVisit/ItemIdentifier B arrival 12:02:30 / departure 12:03
    # include a MonitoredStopVisit/ItemIdentifier C at 15:00
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:00:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:A</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:B</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T12:02:30.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T12:02:30.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T12:03:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T12:03:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:C</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | Test                  |
      | remote_objectid_kind | internal              |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | Name      | Arletty               |
      | ObjectIDs | "internal": "boaarle" |
    And a minute has passed
    And the SIRI server waits GetStopMonitoring request to respond with
      # include a MonitoredStopVisit/ItemIdentifier A at 14:00
      # no MonitoredStopVisit/ItemIdentifier B
      # include a MonitoredStopVisit/ItemIdentifier C at 15:00
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>SQYBUS</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>NAVINEO:SM:RQ:107</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:A</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T13:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T13:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>StopVisit:C</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>SQYBUS:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And 90 seconds have passed
    When the SIRI server has received 2 GetStopMonitoring requests
    Then the StopVisit "6ba7b814-9dad-11d1-f-00c04fd430c8" has the following attributes:
      # "internal": "A"
      | DepartureStatus   | onTime          |
      | ArrivalStatus     | onTime          |
    And the StopVisit "6ba7b814-9dad-11d1-10-00c04fd430c8" has the following attributes:
      # "internal": "B"
      | Collected   | false                |
      | CollectedAt | 2017-01-01T12:02:00Z |
    And the StopVisit "6ba7b814-9dad-11d1-11-00c04fd430c8" has the following attributes:
      # "internal": "C"
      | DepartureStatus   | onTime          |
      | ArrivalStatus     | onTime          |
    And 10 seconds have passed
    And the StopVisit "6ba7b814-9dad-11d1-10-00c04fd430c8" has the following attributes:
      # "internal": "B"
      | Collected       | false                |
      | CollectedAt     | 2017-01-01T12:02:00Z |
      | DepartureStatus | departed             |
      | ArrivalStatus   | cancelled            |

  Scenario: 2939 - Partner Setting collect.include_stop_areas is used to select the best Partner
    Given a SIRI server "first" waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>first</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>first:ResponseMessage::6ba:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>first:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a SIRI server "second" waits GetStopMonitoring request on "http://localhost:8091" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>second</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>second:ResponseMessage::tf7:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:1</ns5:RequestMessageRef>
                <siri:MonitoringRef>boaboon</siri:MonitoringRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193250</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaboon</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>second:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaboon</ns5:StopPointRef>
                      <ns5:Order>45</ns5:Order>
                      <ns5:StopPointName>Charles</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:04:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:04:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:05.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:05:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a Partner "first" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8090 |
      | collect.include_stop_areas | first                 |
      | remote_objectid_kind       | external              |
      | remote_credential          | dummy                 |
    And a Partner "second" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url                 | http://localhost:8091 |
      | collect.include_stop_areas | second                |
      | remote_objectid_kind       | external              |
      | remote_credential          | dummy                 |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | ObjectIDs       | "external": "first" |
    And a StopArea exists with the following attributes:
      | ObjectIDs       | "external": "second" |
    When a minute has passed
    Then the "first" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | first |
    And the "second" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | second |

  Scenario: 2939 - Partner Setting collect.priority is used to select the best Partner
    Given a SIRI server "first" waits GetStopMonitoring request on "http://localhost:8090" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>first</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>first:ResponseMessage::6ba:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>first:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a SIRI server "second" waits GetStopMonitoring request on "http://localhost:8091" to respond with
      """
      <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
        <SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"/>
        <soap:Body>
          <ns1:GetStopMonitoringResponse xmlns:ns1="http://wsdl.siri.org.uk">
            <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
              <ns5:ProducerRef>first</ns5:ProducerRef>
              <ns5:ResponseMessageIdentifier>first:ResponseMessage::6ba:LOC</ns5:ResponseMessageIdentifier>
              <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
              <ns5:StopMonitoringDelivery version="1.3">
                <ns5:ResponseTimestamp>2017-01-01T12:02:00.000+01:00</ns5:ResponseTimestamp>
                <ns5:RequestMessageRef>StopMonitoring:Test:0</ns5:RequestMessageRef>
                <siri:MonitoringRef>boaarle</siri:MonitoringRef>
                <ns5:Status>true</ns5:Status>
                <ns5:MonitoredStopVisit>
                  <ns5:RecordedAtTime>2017-01-01T11:47:15.600+01:00</ns5:RecordedAtTime>
                  <ns5:ItemIdentifier>SIRI:33193249</ns5:ItemIdentifier>
                  <ns5:MonitoringRef>boaarle</ns5:MonitoringRef>
                  <ns5:MonitoredVehicleJourney>
                    <ns5:LineRef>CdF:Line::415:LOC</ns5:LineRef>
                    <ns5:FramedVehicleJourneyRef>
                      <ns5:DataFrameRef>first:Version:1.0:LOC</ns5:DataFrameRef>
                      <ns5:DatedVehicleJourneyRef>1STD721687165983</ns5:DatedVehicleJourneyRef>
                    </ns5:FramedVehicleJourneyRef>
                    <ns5:JourneyPatternRef>CdF:JourneyPattern::L415P289:LOC</ns5:JourneyPatternRef>
                    <ns5:PublishedLineName>415</ns5:PublishedLineName>
                    <ns5:DirectionName>Aller</ns5:DirectionName>
                    <ns5:OperatorRef>CdF:Company::410:LOC</ns5:OperatorRef>
                    <ns5:DestinationRef>boabonn</ns5:DestinationRef>
                    <ns5:DestinationName>Méliès - Croix Bonnet</ns5:DestinationName>
                    <ns5:Monitored>true</ns5:Monitored>
                    <ns5:MonitoredCall>
                      <ns5:StopPointRef>boaarle</ns5:StopPointRef>
                      <ns5:Order>44</ns5:Order>
                      <ns5:StopPointName>Arletty</ns5:StopPointName>
                      <ns5:VehicleAtStop>false</ns5:VehicleAtStop>
                      <ns5:DestinationDisplay>Méliès - Croix Bonnet</ns5:DestinationDisplay>
                      <ns5:AimedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:AimedArrivalTime>
                      <ns5:ExpectedArrivalTime>2017-01-01T15:00:00.000+01:00</ns5:ExpectedArrivalTime>
                      <ns5:ArrivalStatus>onTime</ns5:ArrivalStatus>
                      <ns5:AimedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:AimedDepartureTime>
                      <ns5:ExpectedDepartureTime>2017-01-01T15:01:00.000+01:00</ns5:ExpectedDepartureTime>
                      <ns5:DepartureStatus>onTime</ns5:DepartureStatus>
                    </ns5:MonitoredCall>
                  </ns5:MonitoredVehicleJourney>
                </ns5:MonitoredStopVisit>
              </ns5:StopMonitoringDelivery>
            </Answer>
            <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:siri="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
          </ns1:GetStopMonitoringResponse>
        </soap:Body>
      </soap:Envelope>
        """
    And a Partner "first" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8090 |
      | collect.priority     | 1                     |
      | remote_objectid_kind | external              |
      | remote_credential    | dummy                 |
    And a Partner "second" exists with connectors [siri-check-status-client, siri-stop-monitoring-request-collector] and the following settings:
      | remote_url           | http://localhost:8091 |
      | collect.priority     | 2                     |
      | remote_objectid_kind | external              |
      | remote_credential    | dummy                 |
    And a minute has passed
    And a StopArea exists with the following attributes:
      | ObjectIDs       | "external": "single"     |
    When a minute has passed
    Then the "first" SIRI server should not have received a GetStopMonitoring request
    Then the "second" SIRI server should have received a GetStopMonitoring request with:
      | //siri:MonitoringRef | single |
