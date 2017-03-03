Feature: Support SIRI StopMonitoring

  Background:
    Given a Referential "test" is created

  @wip
  Scenario: Handle a SIRI StopDiscovery request
    Given a Partner "test" exists with connectors [siri-stop-discovery-request-broadcaster] and the following settings:
      | remote_url           | http://localhost:8090 |
      | remote_credential    | test                  |
      | remote_objectid_kind | internal              |
    And a StopArea exists with the following attributes:
      | Name                  | Test                                     |
      | ObjectIDs             | "internal": "NINOXE:StopPoint:BP:6:LOC"   |
      | Attribute[Monitored]  | true   |
      | Attribute[StopName]   | Hotel des compagnons de la musique des alpages (R)   |
      | Attribute[LineRef]    | NINOXE:Line:1:LOC   |
      | Attribute[Longitude]  | 5.4175472259521480   |
      | Attribute[Latitude]   | 46.2270561080357000   |
    And a StopArea exists with the following attributes:
      | Name                  | Test 2                                    |
      | ObjectIDs             | "internal": "NINOXE:StopPoint:SP:16:LOC"   |
      | Attribute[Monitored]  | true   |
      | Attribute[StopName]   | Musée de la céramique orientale   |
      | Attribute[LineRef]    | NINOXE:Line:1:LOC   |
      | Attribute[Longitude]  | 5.4896450042724610   |
      | Attribute[Latitude]   | 46.3326469609286900   |
    And a StopArea exists with the following attributes:
      | Name                  | Test 3                                    |
      | ObjectIDs             | "internal": "NINOXE:StopPoint:BP:7:LOC"   |
      | Attribute[Monitored]  | true   |
      | Attribute[StopName]   | L'abreuvoir des neuf chèvres (A)   |
      | Attribute[LineRef]    | NINOXE:Line:1:LOC   |
      | Attribute[Longitude]  | 5.4488754272460940   |
      | Attribute[Latitude]   | 46.3530890624813000   |
    When I send this SIRI request
      """
      <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"
            xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
        <S:Body>
          <ns7:StopPointsDiscovery xmlns:ns2="http://www.siri.org.uk/siri"
                           xmlns:ns3="http://www.ifopt.org.uk/acsb"
                           xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                           xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                           xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
            <Request>
              <ns2:RequestTimestamp>2017-03-03T11:28:00.359Z</ns2:RequestTimestamp>
              <ns2:RequestorRef>STIF</ns2:RequestorRef>
            </Request>
            <RequestExtension />
          </ns7:StopPointsDiscovery>
        </S:Body>
      </S:Envelope>
        """
    Then I should receive this SIRI response
      """
      <?xml version="1.0" encoding="UTF-8"?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <ns8:StopPointsDiscoveryResponse xmlns:ns8="http://wsdl.siri.org.uk" xmlns:ns3="http://www.siri.org.uk/siri" xmlns:ns4="http://www.ifopt.org.uk/acsb" xmlns:ns5="http://www.ifopt.org.uk/ifopt" xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns7="http://scma/siri" xmlns:ns9="http://wsdl.siri.org.uk/siri">
      <Answer version="2.0">
        <ns3:ResponseTimestamp>2017-03-03T13:54:18.888+01:00</ns3:ResponseTimestamp>
        <ns3:Status>true</ns3:Status>
        <ns3:AnnotatedStopPointRef>
          <ns3:StopPointRef>NINOXE:StopPoint:BP:6:LOC</ns3:StopPointRef>
          <ns3:Monitored>true</ns3:Monitored>
          <ns3:StopName>Hotel des compagnons de la musique des alpages (R)</ns3:StopName>
          <ns3:Lines>
            <ns3:LineRef>NINOXE:Line:1:LOC</ns3:LineRef>
          </ns3:Lines>
          <ns3:Location srsName="EPSG:4328">
            <ns3:Longitude>5.4175472259521480</ns3:Longitude>
            <ns3:Latitude>46.2270561080357000</ns3:Latitude>
          </ns3:Location>
        </ns3:AnnotatedStopPointRef>
        <ns3:AnnotatedStopPointRef>
          <ns3:StopPointRef>NINOXE:StopPoint:SP:16:LOC</ns3:StopPointRef>
          <ns3:Monitored>true</ns3:Monitored>
          <ns3:StopName>Musée de la céramique orientale</ns3:StopName>
          <ns3:Lines>
            <ns3:LineRef>NINOXE:Line:1:LOC</ns3:LineRef>
          </ns3:Lines>
          <ns3:Location srsName="EPSG:4328">
            <ns3:Longitude>5.4896450042724610</ns3:Longitude>
            <ns3:Latitude>46.3326469609286900</ns3:Latitude>
          </ns3:Location>
        </ns3:AnnotatedStopPointRef>
        <ns3:AnnotatedStopPointRef>
          <ns3:StopPointRef>NINOXE:StopPoint:BP:7:LOC</ns3:StopPointRef>
          <ns3:Monitored>true</ns3:Monitored>
          <ns3:StopName>L'abreuvoir des neuf chèvres (A)</ns3:StopName>
          <ns3:Lines>
            <ns3:LineRef>NINOXE:Line:1:LOC</ns3:LineRef>
          </ns3:Lines>
          <ns3:Location srsName="EPSG:4328">
            <ns3:Longitude>5.4488754272460940</ns3:Longitude>
            <ns3:Latitude>46.3530890624813000</ns3:Latitude>
          </ns3:Location>
        </ns3:AnnotatedStopPointRef>
      </Answer>
      <AnswerExtension />
    </ns8:StopPointsDiscoveryResponse>
  </S:Body>
</S:Envelope>
      """
