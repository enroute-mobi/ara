Feature: Support SIRI Situation Exchange by request
  Background:
      Given a Referential "test" is created

  @siri-valid @ARA-1342
  Scenario: Handle a SIRI SituationExchange request
    Given a Situation exists with the following attributes:
      | Codes                                                                              | "external" : "test"                           |
      | RecordedAt                                                                         | 2017-01-01T03:30:06+02:00                     |
      | Version                                                                            | 1                                             |
      | Keywords                                                                           | ["Commercial", "Test"]                        |
      | ReportType                                                                         | general                                       |
      | ParticipantRef                                                                     | "535"                                         |
      | VersionedAt                                                                        | 2017-01-01T01:02:03+02:00                     |
      | Progress                                                                           | published                                     |
      | ValidityPeriods[0]#StartTime                                                       | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                     |
      | PublicationWindows[0]#StartTime                                                    | 2017-09-01T01:00:00+02:00                     |
      | PublicationWindows[0]#EndTime                                                      | 2017-09-25T01:00:00+02:00                     |
      | AlertCause                                                                         | maintenanceWork                               |
      | Severity                                                                           | normal                                        |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-2-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStopId     | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStopId    | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/RouteRef         | Route:66:LOC                                  |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[0]   | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[1]   | 6ba7b814-9dad-11d1-6-00c04fd430c8             |
    When the Situation "external":"test" is edited with a Consequence with the following attributes:
      | Periods[0]#StartTime                                                          | 2023-09-18T05:30:59Z              |
      | Periods[0]#EndTime                                                            | 2023-09-18T08:00:54Z              |
      | Severity                                                                      | verySlight                        |
      | Affects[Line]                                                                 | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStop | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStop  | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Affects[StopArea]                                                             | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Blocking[JourneyPlanner]                                                      | true                              |
      | Blocking[RealTime]                                                            | true                              |
      | Affects[StopArea=6ba7b814-9dad-11d1-4-00c04fd430c8]/LineIds[0]                | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-4-00c04fd430c8]/LineIds[1]                | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test last stop                           |
      | Codes | "external": "NINOXE:StopPoint:SP:25:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3534                            |
      | Codes | "external": "STIF:StopPoint:Q:3534:" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3533                            |
      | Codes | "external": "STIF:StopPoint:Q:3533:" |
    And a SIRI Partner "test" exists with connectors [siri-situation-exchange-request-broadcaster] and the following settings:
      | local_credential  | NINOXE:default |
      | remote_code_space | external       |
    When I send this SIRI request
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <sw:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </soap:Body>
  </soap:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-8-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
                    <siri:ParticipantRef>535</siri:ParticipantRef>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:VersionedAtTime>2017-01-01T01:02:03.000+02:00</siri:VersionedAtTime>
                    <siri:Progress>published</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:PublicationWindow>
                      <siri:StartTime>2017-09-01T01:00:00.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-09-25T01:00:00.000+02:00</siri:EndTime>
                    </siri:PublicationWindow>
                    <siri:AlertCause>maintenanceWork</siri:AlertCause>
                    <siri:Severity>normal</siri:Severity>
                    <siri:ReportType>general</siri:ReportType>
                    <siri:Keywords>Commercial Test</siri:Keywords>
                    <siri:Description>La nouvelle carte d'abonnement est disponible</siri:Description>
                    <siri:Affects>
                      <siri:Networks>
                        <siri:AffectedNetwork>
                          <siri:AffectedLine>
                            <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                            <siri:Destinations>
                              <siri:StopPlaceRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPlaceRef>
                            </siri:Destinations>
                            <siri:Routes>
                              <siri:AffectedRoute>
                                <siri:RouteRef>Route:66:LOC</siri:RouteRef>
                                <siri:StopPoints>
                                  <siri:AffectedStopPoint>
                                    <siri:StopPointRef>STIF:StopPoint:Q:3534:</siri:StopPointRef>
                                  </siri:AffectedStopPoint>
                                  <siri:AffectedStopPoint>
                                    <siri:StopPointRef>STIF:StopPoint:Q:3533:</siri:StopPointRef>
                                  </siri:AffectedStopPoint>
                                </siri:StopPoints>
                              </siri:AffectedRoute>
                            </siri:Routes>
                            <siri:Sections>
                              <siri:AffectedSection>
                                <siri:IndirectSectionRef>
                                  <siri:FirstStopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:FirstStopPointRef>
                                  <siri:LastStopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:LastStopPointRef>
                                </siri:IndirectSectionRef>
                              </siri:AffectedSection>
                            </siri:Sections>
                          </siri:AffectedLine>
                        </siri:AffectedNetwork>
                      </siri:Networks>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                    <siri:Consequences>
                      <siri:Consequence>
                        <siri:Period>
                          <siri:StartTime>2023-09-18T05:30:59.000Z</siri:StartTime>
                          <siri:EndTime>2023-09-18T08:00:54.000Z</siri:EndTime>
                        </siri:Period>
                        <siri:Severity>verySlight</siri:Severity>
                        <siri:Affects>
                          <siri:Networks>
                            <siri:AffectedNetwork>
                              <siri:AffectedLine>
                                <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                <siri:Sections>
                                  <siri:AffectedSection>
                                    <siri:IndirectSectionRef>
                                      <siri:FirstStopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:FirstStopPointRef>
                                      <siri:LastStopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:LastStopPointRef>
                                    </siri:IndirectSectionRef>
                                  </siri:AffectedSection>
                                </siri:Sections>
                              </siri:AffectedLine>
                            </siri:AffectedNetwork>
                          </siri:Networks>
                          <siri:StopPoints>
                            <siri:AffectedStopPoint>
                              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                              <siri:Lines>
                                <siri:AffectedLine>
                                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                </siri:AffectedLine>
                              </siri:Lines>
                            </siri:AffectedStopPoint>
                          </siri:StopPoints>
                        </siri:Affects>
                        <siri:Blocking>
                          <siri:JourneyPlanner>true</siri:JourneyPlanner>
                          <siri:RealTime>true</siri:RealTime>
                        </siri:Blocking>
                      </siri:Consequence>
                    </siri:Consequences>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                         |
      | Direction | received                                                     |
      | Status    | OK                                                           |
      | Type      | SituationExchangeRequest                                     |
      | StopAreas | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | Lines     | ["NINOXE:Line:3:LOC"]                                        |

  @siri-valid @ARA-1342
  Scenario: Handle a SIRI SituationExchange request without any situation
    And a SIRI Partner "test" exists with connectors [siri-situation-exchange-request-broadcaster] and the following settings:
      | local_credential  | NINOXE:default |
      | remote_code_space | external       |
    When I send this SIRI request
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <sw:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </soap:Body>
  </soap:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol  | siri                     |
      | Direction | received                 |
      | Status    | OK                       |
      | Type      | SituationExchangeRequest |
      | StopAreas | []                       |
      | Lines     | []                       |

  @siri-valid @ARA-1397
  Scenario: Handle a SX response (ServiceDelivery)
    Given a SIRI server waits SituationExchangeRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
                    <siri:ParticipantRef>535</siri:ParticipantRef>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:VersionedAtTime>2017-01-01T01:02:03.000+02:00</siri:VersionedAtTime>
                    <siri:Progress>published</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:PublicationWindow>
                      <siri:StartTime>2017-09-01T01:00:00.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-09-25T01:00:00.000+02:00</siri:EndTime>
                    </siri:PublicationWindow>
                    <siri:AlertCause>maintenanceWork</siri:AlertCause>
                    <siri:Severity>slight</siri:Severity>
                    <siri:ReportType>general</siri:ReportType>
                    <siri:Keywords>Commercial Test</siri:Keywords>
                    <siri:Description>La nouvelle carte d'abonnement est disponible</siri:Description>
                    <siri:Affects>
                      <siri:Networks>
                        <siri:AffectedNetwork>
                          <siri:AffectedLine>
                            <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                            <siri:Destinations>
                              <siri:StopPlaceRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPlaceRef>
                            </siri:Destinations>
                            <siri:Routes>
                              <siri:AffectedRoute>
                                <siri:RouteRef>Route:66:LOC</siri:RouteRef>
                                <siri:StopPoints>
                                   <siri:AffectedStopPoint>
                                       <siri:StopPointRef>STIF:StopPoint:Q:3534:</siri:StopPointRef>
                                   </siri:AffectedStopPoint>
                                   <siri:AffectedStopPoint>
                                       <siri:StopPointRef>STIF:StopPoint:Q:3533:</siri:StopPointRef>
                                   </siri:AffectedStopPoint>
                                 </siri:StopPoints>
                              </siri:AffectedRoute>
                              <siri:AffectedRoute>
                                <siri:RouteRef>Route:77:LOC</siri:RouteRef>
                              </siri:AffectedRoute>
                            </siri:Routes>
                            <siri:Sections>
                              <siri:AffectedSection>
                                <siri:IndirectSectionRef>
                                  <siri:FirstStopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:FirstStopPointRef>
                                  <siri:LastStopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:LastStopPointRef>
                                </siri:IndirectSectionRef>
                              </siri:AffectedSection>
                            </siri:Sections>
                          </siri:AffectedLine>
                        </siri:AffectedNetwork>
                      </siri:Networks>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                </siri:PtSituationElement>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
                    <siri:SituationNumber>test2</siri:SituationNumber>
                    <siri:Version>5</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
                    <siri:Severity>noImpact</siri:Severity>
                    <siri:ReportType>general</siri:ReportType>
                    <siri:Keywords>Commercial Test2</siri:Keywords>
                    <siri:Description>carte d'abonnement</siri:Description>
                    <siri:Affects>
                      <siri:Networks>
                        <siri:AffectedNetwork>
                          <siri:AffectedLine>
                            <siri:LineRef>NINOXE:Line:BP:LOC</siri:LineRef>
                          </siri:AffectedLine>
                        </siri:AffectedNetwork>
                      </siri:Networks>
                    </siri:Affects>
                    <siri:Consequences>
                      <siri:Consequence>
                        <siri:Period>
                          <siri:StartTime>2023-09-18T05:30:59.000Z</siri:StartTime>
                          <siri:EndTime>2023-09-18T08:00:54.000Z</siri:EndTime>
                        </siri:Period>
                        <siri:Severity>verySlight</siri:Severity>
                        <siri:Affects>
                          <siri:Networks>
                            <siri:AffectedNetwork>
                              <siri:AffectedLine>
                                <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                <siri:Sections>
                                  <siri:AffectedSection>
                                    <siri:IndirectSectionRef>
                                      <siri:FirstStopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:FirstStopPointRef>
                                      <siri:LastStopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:LastStopPointRef>
                                    </siri:IndirectSectionRef>
                                  </siri:AffectedSection>
                                </siri:Sections>
                              </siri:AffectedLine>
                            </siri:AffectedNetwork>
                          </siri:Networks>
                          <siri:StopPoints>
                            <siri:AffectedStopPoint>
                              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                              <siri:Lines>
                                  <siri:AffectedLine>
                                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                  </siri:AffectedLine>
                                  <siri:AffectedLine>
                                    <siri:LineRef>NINOXE:Line:BP:LOC</siri:LineRef>
                                  </siri:AffectedLine>
                              </siri:Lines>
                            </siri:AffectedStopPoint>
                          </siri:StopPoints>
                        </siri:Affects>
                        <siri:Blocking>
                          <siri:JourneyPlanner>true</siri:JourneyPlanner>
                          <siri:RealTime>true</siri:RealTime>
                        </siri:Blocking>
                      </siri:Consequence>
                    </siri:Consequences>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-situation-exchange-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | ineo                  |
      | remote_code_space | external              |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:BP:LOC" |
      | Name  | Ligne BP Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test last stop                           |
      | Codes | "external": "NINOXE:StopPoint:SP:25:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3534                            |
      | Codes | "external": "STIF:StopPoint:Q:3534:" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3533                            |
      | Codes | "external": "STIF:StopPoint:Q:3533:" |
    And a minute has passed
    When a minute has passed
    Then one Situation has the following attributes:
      | Codes                                                                              | "external" : "test"                           |
      | RecordedAt                                                                         | 2017-01-01T01:02:03+02:00                     |
      | Version                                                                            | 1                                             |
      | Keywords                                                                           | ["Commercial", "Test"]                        |
      | ReportType                                                                         | general                                       |
      | ParticipantRef                                                                     | "535"                                         |
      | VersionedAt                                                                        | 2017-01-01T01:02:03+02:00                     |
      | Progress                                                                           | published                                     |
      | Severity                                                                           | slight                                        |
      | ValidityPeriods[0]#StartTime                                                       | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                     |
      | PublicationWindows[0]#StartTime                                                    | 2017-09-01T01:00:00+02:00                     |
      | PublicationWindows[0]#EndTime                                                      | 2017-09-25T01:00:00+02:00                     |
      | AlertCause                                                                         | maintenanceWork                               |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-2-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStop      | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStop       | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/RouteRef         | Route:66:LOC                                  |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[1]/RouteRef         | Route:77:LOC                                  |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[0]   | 6ba7b814-9dad-11d1-6-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[1]   | 6ba7b814-9dad-11d1-7-00c04fd430c8             |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
    Then one Situation has the following attributes:
      | Codes                        | "external" : "test2"              |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00         |
      | Version                      | 5                                 |
      | Keywords                     | ["Commercial", "Test2"]           |
      | ReportType                   | general                           |
      | Severity                     | noImpact                          |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Description[DefaultValue]    | carte d'abonnement                |
      | Affects[Line]                | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And the Situation "external":"test2" has a Consequence with the following attributes:
      | Periods[0]#StartTime                                                          | 2023-09-18T05:30:59Z              |
      | Periods[0]#EndTime                                                            | 2023-09-18T08:00:54Z              |
      | Severity                                                                      | verySlight                        |
      | Affects[Line]                                                                 | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStop | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStop  | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Affects[StopArea]                                                             | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Blocking[JourneyPlanner]                                                      | true                              |
      | Blocking[RealTime]                                                            | true                              |
      | Affects[StopArea=6ba7b814-9dad-11d1-2-00c04fd430c8]/LineIds[0]                | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-3-00c04fd430c8]/LineIds[1]                | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                                                                             |
      | Direction | sent                                                                                                             |
      | Status    | OK                                                                                                               |
      | Type      | SituationExchangeRequest                                                                                         |
      | StopAreas | ["STIF:StopPoint:Q:3534:", "STIF:StopPoint:Q:3533:", "NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | Lines     | ["NINOXE:Line:3:LOC", "NINOXE:Line:BP:LOC"]                                                                      |

  @ARA-1397 @siri-valid
  Scenario: SituationExchange collect should send GetSituationExchange request to partner
   Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-situation-exchange-request-collector] and the following settings:
      | remote_url            | http://localhost:8090 |
      | remote_credential     | test                  |
      | remote_code_space     | internal              |
      | collect.include_lines | RLA_Bus:Line::05:LOC  |
      | local_credential      | ara                   |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name  | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   And 20 seconds have passed
   Then the SIRI server should have received 1 GetSituationExchange request
   And an audit event should exist with these attributes:
      | Protocol  | siri                     |
      | Direction | sent                     |
      | Type      | SituationExchangeRequest |
      | Lines     | nil                      |

  @ARA-1461 @siri-valid
  Scenario: Manage a Request with a Line filter
    Given a SIRI server on "http://localhost:8090"
And a Partner "test" exists with connectors [siri-check-status-client,siri-situation-exchange-request-collector] and the following settings:
      | remote_url                | http://localhost:8090 |
      | remote_credential         | test                  |
      | remote_code_space         | internal              |
      | collect.filter_situations | true                  |
    And a minute has passed
    And a Line exists with the following attributes:
      | Name              | Test               |
      | Codes             | "internal": "1234" |
      | CollectSituations | true               |
      And 10 seconds have passed
      And the SIRI server has received a GetGeneralMessage request
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetSituationExchange xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:01:10.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>6ba7b814-9dad-11d1-4-00c04fd430c8</siri:MessageIdentifier>
      </ServiceRequestInfo>
      <Request version="2.0:FR-IDF-2.4">
        <siri:RequestTimestamp>2017-01-01T12:01:10.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>6ba7b814-9dad-11d1-4-00c04fd430c8</siri:MessageIdentifier>
        <siri:PreviewInterval>PT12H</siri:PreviewInterval>
        <siri:StartTime>2017-01-01T00:01:10.000Z</siri:StartTime>
        <siri:LineRef>1234</siri:LineRef>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </S:Body>
  </S:Envelope>
      """

  @ARA-1461 @siri-valid
  Scenario: Manage a Request with a StopArea filter
    Given a SIRI server on "http://localhost:8090"
      And a Partner "test" exists with connectors [siri-check-status-client,siri-situation-exchange-request-collector] and the following settings:
      | remote_url                | http://localhost:8090 |
      | remote_credential         | test                  |
      | remote_code_space         | internal              |
      | collect.filter_situations | true                  |
    And a minute has passed
    And a StopArea exists with the following attributes:
        | Name              | Test              |
        | Codes             | "internal":"1234" |
        | CollectSituations | true              |
      And 10 seconds have passed
      And the SIRI server has received a GetGeneralMessage request
    Then the SIRI server should receive this response
      """
<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <sw:GetSituationExchange xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:01:10.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test</siri:RequestorRef>
        <siri:MessageIdentifier>6ba7b814-9dad-11d1-4-00c04fd430c8</siri:MessageIdentifier>
      </ServiceRequestInfo>
      <Request version="2.0:FR-IDF-2.4">
        <siri:RequestTimestamp>2017-01-01T12:01:10.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>6ba7b814-9dad-11d1-4-00c04fd430c8</siri:MessageIdentifier>
        <siri:PreviewInterval>PT12H</siri:PreviewInterval>
        <siri:StartTime>2017-01-01T00:01:10.000Z</siri:StartTime>
        <siri:StopPointRef>1234</siri:StopPointRef>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </S:Body>
  </S:Envelope>
      """

  @siri-valid @ARA-1471
  Scenario: Broadcast via a SIRI SX request a Situation without EndTime
    Given a Situation exists with the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00         |
      | Version                      | 1                                 |
      | ReportType                   | general                           |
      | Progress                     | published                         |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      # | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Description[DefaultValue]    | Description Sample                |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Name  | Stop Area Sample                     |
      | Codes | "external": "sample" |
    And a SIRI Partner "test" exists with connectors [siri-situation-exchange-request-broadcaster] and the following settings:
      | local_credential  | NINOXE:default |
      | remote_code_space | external       |
    When I send this SIRI request
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <sw:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </soap:Body>
</soap:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:Progress>published</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
                    <siri:ReportType>general</siri:ReportType>
                    <siri:Description>Description Sample</siri:Description>
                    <siri:Affects>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>sample</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """


  @siri-valid @ARA-1472
  Scenario: Ensure a Situation is updated with a VersionedAtTime changed and Version remains the same
    Given a SIRI server waits SituationExchangeRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:00+02:00</siri:CreationTime>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:VersionedAtTime>2017-01-01T12:00:00.000+02:00</siri:VersionedAtTime>
                    <siri:Progress>closed</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
                    <siri:Description>La nouvelle carte d'abonnement est disponible</siri:Description>
                    <siri:Affects>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-situation-exchange-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | ineo                  |
      | remote_code_space | external              |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a Situation exists with the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T03:30:00+02:00         |
      | VersionedAt                  | 2017-01-01T03:30:00+02:00         |
      | Version                      | 1                                 |
      | ReportType                   | general                           |
      | Progress                     | open                              |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    When a minute has passed
    And a minute has passed
    Then one Situation has the following attributes:
      | Codes                        | "external" : "test"                           |
      | RecordedAt                   | 2017-01-01T12:00:00+02:00                     |
      | VersionedAt                  | 2017-01-01T12:00:00+02:00                     |
      | Version                      | 1                                             |
      | Progress                     | closed                                        |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime   | 0001-01-01T00:00:00Z                          |
      | Description[DefaultValue]    | La nouvelle carte d'abonnement est disponible |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8             |

  @siri-valid @ARA-1472
  Scenario: Ensure a Situation is updated with a Version changed and no VersionAtTime
    Given a SIRI server waits SituationExchangeRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:00+02:00</siri:CreationTime>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>2</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:Progress>closed</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
                    <siri:Description>La nouvelle carte d'abonnement est disponible</siri:Description>
                    <siri:Affects>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-situation-exchange-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | ineo                  |
      | remote_code_space | external              |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a Situation exists with the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T03:30:00+02:00         |
      | Version                      | 1                                 |
      | ReportType                   | general                           |
      | Progress                     | open                              |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    When a minute has passed
    And a minute has passed
    Then one Situation has the following attributes:
      | Codes                        | "external" : "test"                           |
      | RecordedAt                   | 2017-01-01T03:30:00+02:00                     |
      | Version                      | 2                                             |
      | Progress                     | closed                                        |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime   | 0001-01-01T00:00:00Z                          |
      | Description[DefaultValue]    | La nouvelle carte d'abonnement est disponible |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8             |


  @siri-valid @ARA-1443
  Scenario: Collect SituationExchange with internal tags
    Given a SIRI server waits SituationExchangeRequest request on "http://localhost:8090" to respond with
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:00+02:00</siri:CreationTime>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>2</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:Progress>closed</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
                    <siri:Description>La nouvelle carte d'abonnement est disponible</siri:Description>
                    <siri:Affects>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    And a Partner "ineo" exists with connectors [siri-check-status-client, siri-situation-exchange-request-collector] and the following settings:
      | remote_url                       | http://localhost:8090 |
      | remote_credential                | ineo                  |
      | remote_code_space                | external              |
      | collect.situations.internal_tags | first,second          |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    When a minute has passed
    And a minute has passed
    Then one Situation has the following attributes:
      | Codes        | "external" : "test" |
      | InternalTags | ["first","second"]  |

  @siri-valid @ARA-1444
  Scenario: Broadcast via a SIRI SX request a Situation with matching internal tags
    Given a Situation exists with the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00         |
      | Version                      | 1                                 |
      | ReportType                   | general                           |
      | Progress                     | published                         |
      | InternalTags                 | ["first","second"]                |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Description[DefaultValue]    | Description Sample                |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Name  | Stop Area Sample     |
      | Codes | "external": "sample" |
    And a SIRI Partner "test" exists with connectors [siri-situation-exchange-request-broadcaster] and the following settings:
      | local_credential                   | NINOXE:default |
      | remote_code_space                  | external       |
      | broadcast.situations.internal_tags | first,another  |
    When I send this SIRI request
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <sw:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </soap:Body>
</soap:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:Progress>published</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
                    <siri:ReportType>general</siri:ReportType>
                    <siri:Description>Description Sample</siri:Description>
                    <siri:Affects>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>sample</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """

  @siri-valid @ARA-1444
  Scenario: Do not broadcast a Situation with no matching internal tags
    Given a Situation exists with the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00         |
      | Version                      | 1                                 |
      | ReportType                   | general                           |
      | Progress                     | published                         |
      | InternalTags                 | ["first","second"]                |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Description[DefaultValue]    | Description Sample                |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Name  | Stop Area Sample     |
      | Codes | "external": "sample" |
    And a SIRI Partner "test" exists with connectors [siri-situation-exchange-request-broadcaster] and the following settings:
      | local_credential                   | NINOXE:default |
      | remote_code_space                  | external       |
      | broadcast.situations.internal_tags | another        |
    When I send this SIRI request
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <sw:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </soap:Body>
</soap:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """

  @siri-valid @ARA-1493
  Scenario: Handle referent lines in a SIRI SituationExchange request
    Given a Situation exists with the following attributes:
      | Codes                                                                              | "external" : "test"                           |
      | RecordedAt                                                                         | 2017-01-01T03:30:06+02:00                     |
      | Version                                                                            | 1                                             |
      | Keywords                                                                           | ["Commercial", "Test"]                        |
      | ReportType                                                                         | general                                       |
      | ParticipantRef                                                                     | "535"                                         |
      | VersionedAt                                                                        | 2017-01-01T01:02:03+02:00                     |
      | Progress                                                                           | published                                     |
      | ValidityPeriods[0]#StartTime                                                       | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                     |
      | PublicationWindows[0]#StartTime                                                    | 2017-09-01T01:00:00+02:00                     |
      | PublicationWindows[0]#EndTime                                                      | 2017-09-25T01:00:00+02:00                     |
      | AlertCause                                                                         | maintenanceWork                               |
      | Severity                                                                           | normal                                        |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/LastStopId     | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/FirstStopId    | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedRoutes[0]/RouteRef         | Route:66:LOC                                  |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[0]   | 6ba7b814-9dad-11d1-6-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[1]   | 6ba7b814-9dad-11d1-7-00c04fd430c8             |
    When the Situation "external":"test" is edited with a Consequence with the following attributes:
      | Periods[0]#StartTime                                                          | 2023-09-18T05:30:59Z              |
      | Periods[0]#EndTime                                                            | 2023-09-18T08:00:54Z              |
      | Severity                                                                      | verySlight                        |
      | Affects[Line]                                                                 | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/FirstStop | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/LastStop  | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Affects[StopArea]                                                             | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
      | Blocking[JourneyPlanner]                                                      | true                              |
      | Blocking[RealTime]                                                            | true                              |
      | Affects[StopArea=6ba7b814-9dad-11d1-5-00c04fd430c8]/LineIds[0]                | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-5-00c04fd430c8]/LineIds[1]                | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:4:LOC" |
      | Name  | Ligne 3 Metro                   |
      | ReferentId | 6ba7b814-9dad-11d1-2-00c04fd430c8 |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test last stop                           |
      | Codes | "external": "NINOXE:StopPoint:SP:25:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3534                            |
      | Codes | "external": "STIF:StopPoint:Q:3534:" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3533                            |
      | Codes | "external": "STIF:StopPoint:Q:3533:" |
    And a SIRI Partner "test" exists with connectors [siri-situation-exchange-request-broadcaster] and the following settings:
      | local_credential  | NINOXE:default |
      | remote_code_space | external       |
    When I send this SIRI request
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <sw:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </sw:GetSituationExchange>
  </soap:Body>
  </soap:Envelope>
      """
    Then I should receive this SIRI response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetSituationExchangeResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-9-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
                    <siri:ParticipantRef>535</siri:ParticipantRef>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:VersionedAtTime>2017-01-01T01:02:03.000+02:00</siri:VersionedAtTime>
                    <siri:Progress>published</siri:Progress>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:PublicationWindow>
                      <siri:StartTime>2017-09-01T01:00:00.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-09-25T01:00:00.000+02:00</siri:EndTime>
                    </siri:PublicationWindow>
                    <siri:AlertCause>maintenanceWork</siri:AlertCause>
                    <siri:Severity>normal</siri:Severity>
                    <siri:ReportType>general</siri:ReportType>
                    <siri:Keywords>Commercial Test</siri:Keywords>
                    <siri:Description>La nouvelle carte d'abonnement est disponible</siri:Description>
                    <siri:Affects>
                      <siri:Networks>
                        <siri:AffectedNetwork>
                          <siri:AffectedLine>
                            <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                            <siri:Destinations>
                              <siri:StopPlaceRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPlaceRef>
                            </siri:Destinations>
                            <siri:Routes>
                              <siri:AffectedRoute>
                                <siri:RouteRef>Route:66:LOC</siri:RouteRef>
                                <siri:StopPoints>
                                  <siri:AffectedStopPoint>
                                    <siri:StopPointRef>STIF:StopPoint:Q:3534:</siri:StopPointRef>
                                  </siri:AffectedStopPoint>
                                  <siri:AffectedStopPoint>
                                    <siri:StopPointRef>STIF:StopPoint:Q:3533:</siri:StopPointRef>
                                  </siri:AffectedStopPoint>
                                </siri:StopPoints>
                              </siri:AffectedRoute>
                            </siri:Routes>
                            <siri:Sections>
                              <siri:AffectedSection>
                                <siri:IndirectSectionRef>
                                  <siri:FirstStopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:FirstStopPointRef>
                                  <siri:LastStopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:LastStopPointRef>
                                </siri:IndirectSectionRef>
                              </siri:AffectedSection>
                            </siri:Sections>
                          </siri:AffectedLine>
                        </siri:AffectedNetwork>
                      </siri:Networks>
                      <siri:StopPoints>
                        <siri:AffectedStopPoint>
                          <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                        </siri:AffectedStopPoint>
                      </siri:StopPoints>
                    </siri:Affects>
                    <siri:Consequences>
                      <siri:Consequence>
                        <siri:Period>
                          <siri:StartTime>2023-09-18T05:30:59.000Z</siri:StartTime>
                          <siri:EndTime>2023-09-18T08:00:54.000Z</siri:EndTime>
                        </siri:Period>
                        <siri:Severity>verySlight</siri:Severity>
                        <siri:Affects>
                          <siri:Networks>
                            <siri:AffectedNetwork>
                              <siri:AffectedLine>
                                <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                <siri:Sections>
                                  <siri:AffectedSection>
                                    <siri:IndirectSectionRef>
                                      <siri:FirstStopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:FirstStopPointRef>
                                      <siri:LastStopPointRef>NINOXE:StopPoint:SP:25:LOC</siri:LastStopPointRef>
                                    </siri:IndirectSectionRef>
                                  </siri:AffectedSection>
                                </siri:Sections>
                              </siri:AffectedLine>
                            </siri:AffectedNetwork>
                          </siri:Networks>
                          <siri:StopPoints>
                            <siri:AffectedStopPoint>
                              <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                              <siri:Lines>
                                <siri:AffectedLine>
                                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                </siri:AffectedLine>
                              </siri:Lines>
                            </siri:AffectedStopPoint>
                          </siri:StopPoints>
                        </siri:Affects>
                        <siri:Blocking>
                          <siri:JourneyPlanner>true</siri:JourneyPlanner>
                          <siri:RealTime>true</siri:RealTime>
                        </siri:Blocking>
                      </siri:Consequence>
                    </siri:Consequences>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetSituationExchangeResponse>
        </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                         |
      | Direction | received                                                     |
      | Status    | OK                                                           |
      | Type      | SituationExchangeRequest                                     |
      | StopAreas | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | Lines     | ["NINOXE:Line:3:LOC"]                                        |