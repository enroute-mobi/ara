Feature: Support SIRI Situation Exchange by request
  Background:
      Given a Referential "test" is created

  @siri-valid @ARA-1342
  Scenario: Handle a SIRI SituationExchange request
    Given a Situation exists with the following attributes:
      | Codes                                                                               | "external" : "test"                           |
      | RecordedAt                                                                          | 2017-01-01T03:30:06+02:00                     |
      | Version                                                                             | 1                                             |
      | Keywords                                                                            | ["Commercial", "Test"]                        |
      | ReportType                                                                          | general                                       |
      | ValidityPeriods[0]#StartTime                                                        | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                          | 2017-01-01T20:30:06+02:00                     |
      | Description[DefaultValue]                                                           | La nouvelle carte d'abonnement est disponible |
      | Affects[StopArea]                                                                   | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line]                                                                       | 6ba7b814-9dad-11d1-2-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId] | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStopId      | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStopId     | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/RouteRef          | Route:66:LOC                                  |
    And a Line exists with the following attributes:
      | Codes | "external": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test last stop                           |
      | Codes | "external": "NINOXE:StopPoint:SP:25:LOC" |
    And a SIRI Partner "test" exists with connectors [siri-situation-exchange-request-broadcaster] and the following settings:
      | local_credential  | NINOXE:default |
      | remote_code_space | external       |
    When I send this SIRI request
      """
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ws:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:ws="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </ws:GetSituationExchange>
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
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
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
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
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
    <ws:GetSituationExchange xmlns:siri="http://www.siri.org.uk/siri" xmlns:ws="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
      </ServiceRequestInfo>
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:MessageIdentifier>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:MessageIdentifier>
      </Request>
      <RequestExtension />
    </ws:GetSituationExchange>
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
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>1</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-01-01T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:UndefinedReason/>
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
    And a minute has passed
    When a minute has passed
    Then one Situation has the following attributes:
      | Codes                                                                              | "external" : "test"                           |
      | RecordedAt                                                                         | 2017-01-01T03:30:06+02:00                     |
      | Version                                                                            | 1                                             |
      | Keywords                                                                           | ["Commercial", "Test"]                        |
      | ReportType                                                                         | general                                       |
      | ValidityPeriods[0]#StartTime                                                       | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                     |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-2-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/FirstStop      | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedSections[0]/LastStop       | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedRoutes[0]/RouteRef         | Route:66:LOC                                  |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-4-00c04fd430c8             |
    Then one Situation has the following attributes:
      | Codes                        | "external" : "test2"              |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00         |
      | Version                      | 5                                 |
      | Keywords                     | ["Commercial", "Test2"]           |
      | ReportType                   | general                           |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Description[DefaultValue]    | carte d'abonnement                |
      | Affects[Line]                | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                         |
      | Direction | sent                                                         |
      | Status    | OK                                                           |
      | Type      | SituationExchangeRequest                                     |
      | StopAreas | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | Lines     | ["NINOXE:Line:3:LOC", "NINOXE:Line:BP:LOC"]                  |

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
      | Name      | Test 1                             |
      | Codes | "internal": "RLA_Bus:Line::05:LOC" |
   And a minute has passed
   And 20 seconds have passed
   Then the SIRI server should have received 1 GetSituationExchange request
   And an audit event should exist with these attributes:
      | Protocol  | siri                     |
      | Direction | sent                     |
      | Type      | SituationExchangeRequest |
      | Lines     | nil                      |
