Feature: Support SIRI SituationExchange by subscription

  Background:
      Given a Referential "test" is created

  @ARA-1450 @siri-valid
  Scenario: Manage a SX Subscription
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-situation-exchange-subscription-collector] and the following settings:
        | remote_url                       | http://localhost:8090 |
        | remote_credential                | test                  |
        | local_credential                 | NINOXE:default        |
        | remote_code_space                | internal              |
        | collect.include_lines            | NINOXE:Line::3:LOC    |
        | collect.situations.internal_tags | first,second          |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:BP:LOC" |
      | Name  | Ligne BP Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "internal": "NINOXE:StopPoint:SP:24:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test last stop                           |
      | Codes | "internal": "NINOXE:StopPoint:SP:25:LOC" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3534                            |
      | Codes | "internal": "STIF:StopPoint:Q:3534:" |
    And a StopArea exists with the following attributes:
      | Name  | Test 3533                            |
      | Codes | "internal": "STIF:StopPoint:Q:3533:" |
    And a Subscription exist with the following attributes:
      | Kind              | SituationExchangeCollect          |
      | ReferenceArray[0] | "SituationExchangeCollect": "all" |
    And a minute has passed
    When I send this SIRI request
      """
     <?xml version='1.0' encoding='utf-8'?>
     <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
       <S:Body>
          <sw:NotifySituationExchange xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Notification>
              <siri:SituationExchangeDelivery>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:SubscriptionRef>6ba7b814-9dad-11d1-9-00c04fd430c8</siri:SubscriptionRef>
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
                    <siri:Reality>technicalExercise</siri:Reality>
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
                    <siri:Summary xml:lang="FR">Nouveau pass Navigo</siri:Summary>
                    <siri:Summary xml:lang="EN">New pass Navigo</siri:Summary>
                    <siri:Description>La nouvelle carte d'abonnement est disponible</siri:Description>
                    <siri:Description xml:lang="EN">The new pass is available</siri:Description>
                    <siri:InfoLinks>
                      <siri:InfoLink>
                        <siri:Uri>http://example.com</siri:Uri>
                        <siri:Label>Label Sample</siri:Label>
                        <siri:Image>
                          <siri:ImageRef>http://www.example.com/image.png</siri:ImageRef>
                        </siri:Image>
                        <siri:LinkContent>relatedSite</siri:LinkContent>
                      </siri:InfoLink>
                    </siri:InfoLinks>
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
                        <siri:Condition>changeOfPlatform</siri:Condition>
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
                    <siri:PublishingActions>
                      <siri:PublishToWebAction>
                        <siri:ActionStatus>published</siri:ActionStatus>
                        <siri:Description xml:lang="NO">Toget vil bytte togmateriell på Dovre. Du må dessverre bytte tog på denne stasjonen.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>DataName</siri:Name>
                          <siri:Type>dummy</siri:Type>
                          <siri:Value>dummy1</siri:Value>
                          <siri:Prompt xml:lang="NO">Toget vil bytte togmateriell på Dovre.</siri:Prompt>
                          <siri:Prompt xml:lang="EN">You must change trains at Dovre. We apologize for the inconvenience.</siri:Prompt>
                          <siri:PublishAtScope>
                          <siri:ScopeType>line</siri:ScopeType>
                          <siri:Affects>
                            <siri:Networks>
                              <siri:AffectedNetwork>
                                <siri:AffectedLine>
                                  <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                </siri:AffectedLine>
                              </siri:AffectedNetwork>
                            </siri:Networks>
                          </siri:Affects>
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:Incidents>true</siri:Incidents>
                        <siri:HomePage>true</siri:HomePage>
                        <siri:Ticker>false</siri:Ticker>
                        <siri:SocialNetwork>facebook.com</siri:SocialNetwork>
                      </siri:PublishToWebAction>
                      <siri:PublishToMobileAction>
                        <siri:ActionStatus>open</siri:ActionStatus>
                        <siri:Description xml:lang="DE">Der Zug wird in Dovre das Zugmaterial wechseln. Leider muss man an diesem Bahnhof umsteigen.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>AnotherDataName</siri:Name>
                          <siri:Type>dummy2</siri:Type>
                          <siri:Value>dummy3</siri:Value>
                          <siri:Prompt xml:lang="DE">Der Zug wird in Dovre das Zugmaterial wechseln.</siri:Prompt>
                          <siri:Prompt xml:lang="HU">A vonat Dovre-ban módosítja a vonat anyagát. Sajnos ezen az állomáson át kell szállni.</siri:Prompt>
                          <siri:PublishAtScope>
                            <siri:ScopeType>stopPlace</siri:ScopeType>
                            <siri:Affects>
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
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:Incidents>false</siri:Incidents>
                        <siri:HomePage>true</siri:HomePage>
                      </siri:PublishToMobileAction>
                      <siri:PublishToDisplayAction>
                        <siri:ActionStatus>open</siri:ActionStatus>
                        <siri:Description xml:lang="IS">Lestin mun skipta um efni í Dovre. Því miður þarftu að flytja á þessari stöð.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>AnotherNewDataName</siri:Name>
                          <siri:Type>dummy6</siri:Type>
                          <siri:Value>dummy7</siri:Value>
                          <siri:Prompt xml:lang="DK">Toget skifter materialet på Dovre.</siri:Prompt>
                          <siri:Prompt xml:lang="FI">Juna vaihtaa junamateriaalia Dovressa. Valitettavasti sinun täytyy vaihtaa tällä asemalla.</siri:Prompt>
                          <siri:PublishAtScope>
                            <siri:ScopeType>network</siri:ScopeType>
                            <siri:Affects>
                              <siri:Networks>
                                <siri:AffectedNetwork>
                                  <siri:AffectedLine>
                                    <siri:LineRef>NINOXE:Line:3:LOC</siri:LineRef>
                                  </siri:AffectedLine>
                                </siri:AffectedNetwork>
                              </siri:Networks>
                            </siri:Affects>
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:OnPlace>false</siri:OnPlace>
                        <siri:OnBoard>true</siri:OnBoard>
                      </siri:PublishToDisplayAction>
                      <siri:PublishToDisplayAction>
                        <siri:ActionStatus>closed</siri:ActionStatus>
                        <siri:Description xml:lang="ES">El tren cambiará su material en Dovre. Lamentablemente tienes que hacer transbordo en esta estación.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>NewDataName</siri:Name>
                          <siri:Type>dummy4</siri:Type>
                          <siri:Value>dummy5</siri:Value>
                          <siri:Prompt xml:lang="PL">Pociąg zmieni materiał składowy w Dovre.</siri:Prompt>
                          <siri:Prompt xml:lang="BG">Влакът променя материала на влака в Dovre. За съжаление трябва да сменяте на тази гара.</siri:Prompt>
                          <siri:PublishAtScope>
                            <siri:ScopeType>general</siri:ScopeType>
                            <siri:Affects>
                              <siri:Networks>
                                <siri:AffectedNetwork>
                                  <siri:AffectedLine>
                                    <siri:LineRef>NINOXE:Line:BP:LOC</siri:LineRef>
                                  </siri:AffectedLine>
                                </siri:AffectedNetwork>
                              </siri:Networks>
                            </siri:Affects>
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:OnPlace>true</siri:OnPlace>
                        <siri:OnBoard>false</siri:OnBoard>
                      </siri:PublishToDisplayAction>
                    </siri:PublishingActions>
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Notification>
            <SiriExtension/>
          </sw:NotifySituationExchange>
        </S:Body>
      </S:Envelope>
      """
    Then one Situation has the following attributes:
      | Codes                                                                              | "internal" : "test"                           |
      | InternalTags                                                                       | ["first","second"]                            |
      | RecordedAt                                                                         | 2017-01-01T01:02:03+02:00                     |
      | Version                                                                            | 1                                             |
      | Keywords                                                                           | ["Commercial", "Test"]                        |
      | ReportType                                                                         | general                                       |
      | ParticipantRef                                                                     | "535"                                         |
      | VersionedAt                                                                        | 2017-01-01T01:02:03+02:00                     |
      | Progress                                                                           | published                                     |
      | Reality                                                                            | technicalExercise                             |
      | Severity                                                                           | slight                                        |
      | ValidityPeriods[0]#StartTime                                                       | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                     |
      | PublicationWindows[0]#StartTime                                                    | 2017-09-01T01:00:00+02:00                     |
      | PublicationWindows[0]#EndTime                                                      | 2017-09-25T01:00:00+02:00                     |
      | AlertCause                                                                         | maintenanceWork                               |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible |
      | Description[Translations]#EN                                                       | The new pass is available                     |
      | Summary[Translations]#FR                                                           | Nouveau pass Navigo                           |
      | Summary[Translations]#EN                                                           | New pass Navigo                               |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/FirstStop      | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/LastStop       | 6ba7b814-9dad-11d1-6-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedRoutes[0]/RouteRef         | Route:66:LOC                                  |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedRoutes[1]/RouteRef         | Route:77:LOC                                  |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[0]   | 6ba7b814-9dad-11d1-7-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedRoutes[0]/StopAreaIds[1]   | 6ba7b814-9dad-11d1-8-00c04fd430c8             |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-5-00c04fd430c8             |
    And the Situation "internal":"test" has an InfoLink with the following attributes:
      | Uri         | http://example.com               |
      | Label       | Label Sample                     |
      | ImageRef    | http://www.example.com/image.png |
      | LinkContent | relatedSite                      |
    Then one Situation has the following attributes:
      | Codes                        | "internal" : "test2"              |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00         |
      | Version                      | 5                                 |
      | Keywords                     | ["Commercial", "Test2"]           |
      | ReportType                   | general                           |
      | Severity                     | noImpact                          |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | Description[DefaultValue]    | carte d'abonnement                |
      | Affects[Line]                | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And the Situation "internal":"test2" has a Consequence with the following attributes:
      | Periods[0]#StartTime                                                          | 2023-09-18T05:30:59Z              |
      | Periods[0]#EndTime                                                            | 2023-09-18T08:00:54Z              |
      | Severity                                                                      | verySlight                        |
      | Condition                                                                     | changeOfPlatform                  |
      | Affects[Line]                                                                 | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/FirstStop | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Affects[Line=6ba7b814-9dad-11d1-3-00c04fd430c8]/AffectedSections[0]/LastStop  | 6ba7b814-9dad-11d1-6-00c04fd430c8 |
      | Affects[StopArea]                                                             | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Blocking[JourneyPlanner]                                                      | true                              |
      | Blocking[RealTime]                                                            | true                              |
      | Affects[StopArea=6ba7b814-9dad-11d1-2-00c04fd430c8]/LineIds[0]                | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-3-00c04fd430c8]/LineIds[1]                | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And the Situation "internal":"test2" has a PublishToWebAction with the following attributes:
      | ActionStatus                    | published                                                                            |
      | Incidents                       | true                                                                                 |
      | HomePage                        | true                                                                                 |
      | Ticker                          | false                                                                                |
      | SocialNetworks                  | ["facebook.com"]                                                                     |
      | Description[Translations]#NO    | Toget vil bytte togmateriell på Dovre. Du må dessverre bytte tog på denne stasjonen. |
      | PublicationWindows[0]#StartTime | 2017-12-01T09:00:00Z                                                                 |
      | PublicationWindows[0]#EndTime   | 2017-12-09T17:00:00Z                                                                 |
      | Name                            | DataName                                                                             |
      | ActionType                      | dummy                                                                                |
      | Value                           | dummy1                                                                               |
      | Prompt[Translations]#NO         | Toget vil bytte togmateriell på Dovre.                                               |
      | Prompt[Translations]#EN         | You must change trains at Dovre. We apologize for the inconvenience.                 |
      | ScopeType                       | line                                                                                 |
      | Affects[Line]                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                                                    |
    And the Situation "internal":"test2" has a PublishToMobileAction with the following attributes:
      | ActionStatus                    | open                                                                                         |
      | Incidents                       | false                                                                                        |
      | HomePage                        | true                                                                                         |
      | Description[Translations]#DE    | Der Zug wird in Dovre das Zugmaterial wechseln. Leider muss man an diesem Bahnhof umsteigen. |
      | PublicationWindows[0]#StartTime | 2017-12-01T09:00:00Z                                                                         |
      | PublicationWindows[0]#EndTime   | 2017-12-09T17:00:00Z                                                                         |
      | Name                    | AnotherDataName                                                                        |
      | ActionType              | dummy2                                                                                 |
      | Value                   | dummy3                                                                                 |
      | Prompt[Translations]#DE | Der Zug wird in Dovre das Zugmaterial wechseln.                                        |
      | Prompt[Translations]#HU | A vonat Dovre-ban módosítja a vonat anyagát. Sajnos ezen az állomáson át kell szállni. |
      | ScopeType                                                      | stopPlace                         |
      | Affects[StopArea]                                              | 6ba7b814-9dad-11d1-5-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-5-00c04fd430c8]/LineIds[0] | 6ba7b814-9dad-11d1-3-00c04fd430c8 |
      | Affects[StopArea=6ba7b814-9dad-11d1-5-00c04fd430c8]/LineIds[1] | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And the Situation "internal":"test2" has a PublishToDisplayAction with the following attributes:
      | ActionStatus                    | closed                                                                                               |
      | OnBoard                         | false                                                                                                |
      | OnPlace                         | true                                                                                                 |
      | Description[Translations]#ES    | El tren cambiará su material en Dovre. Lamentablemente tienes que hacer transbordo en esta estación. |
      | PublicationWindows[0]#StartTime | 2017-12-01T09:00:00Z                                                                                 |
      | PublicationWindows[0]#EndTime   | 2017-12-09T17:00:00Z                                                                                 |
      | Name                    | NewDataName                                                                             |
      | ActionType              | dummy4                                                                                  |
      | Value                   | dummy5                                                                                  |
      | Prompt[Translations]#PL | Pociąg zmieni materiał składowy w Dovre.                                                |
      | Prompt[Translations]#BG | Влакът променя материала на влака в Dovre. За съжаление трябва да сменяте на тази гара. |
      | ScopeType     | general                           |
      | Affects[Line] | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And the Situation "internal":"test2" has a PublishToDisplayAction with the following attributes:
      | ActionStatus                    | open                                                                                       |
      | OnBoard                         | true                                                                                       |
      | OnPlace                         | false                                                                                      |
      | Description[Translations]#IS    | Lestin mun skipta um efni í Dovre. Því miður þarftu að flytja á þessari stöð.              |
      | PublicationWindows[0]#StartTime | 2017-12-01T09:00:00Z                                                                       |
      | PublicationWindows[0]#EndTime   | 2017-12-09T17:00:00Z                                                                       |
      | Name                            | AnotherNewDataName                                                                         |
      | ActionType                      | dummy6                                                                                     |
      | Value                           | dummy7                                                                                     |
      | Prompt[Translations]#DK         | Toget skifter materialet på Dovre.                                                         |
      | Prompt[Translations]#FI         | Juna vaihtaa junamateriaalia Dovressa. Valitettavasti sinun täytyy vaihtaa tällä asemalla. |
      | ScopeType                       | network                                                                                    |
      | Affects[Line]                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                                                          |
    And an audit event should exist with these attributes:
      | Protocol  | siri                                                                                                             |
      | Direction | received                                                                                                         |
      | Status    | OK                                                                                                               |
      | Type      | NotifySituationExchange                                                                                          |
      | StopAreas | ["STIF:StopPoint:Q:3534:", "STIF:StopPoint:Q:3533:", "NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:25:LOC"] |
      | Lines     | ["NINOXE:Line:3:LOC", "NINOXE:Line:BP:LOC"]                                                                      |

  @siri-valid @ARA-1582
  Scenario: Manage a SX Subscription with All affected Lines
    Given a SIRI server on "http://localhost:8090"
    And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-situation-exchange-subscription-collector] and the following settings:
      | remote_url                       | http://localhost:8090 |
      | remote_credential                | test                  |
      | local_credential                 | NINOXE:default        |
      | remote_code_space                | external              |
      | collect.include_lines            | NINOXE:Line::3:LOC    |
      | collect.situations.internal_tags | first,second          |
    And 30 seconds have passed
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a Subscription exist with the following attributes:
      | Kind              | SituationExchangeCollect          |
      | ReferenceArray[0] | "SituationExchangeCollect": "all" |
    And a minute has passed
    When I send this SIRI request
      """
     <?xml version='1.0' encoding='utf-8'?>
     <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
       <S:Body>
          <sw:NotifySituationExchange xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Notification>
              <siri:SituationExchangeDelivery>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:SubscriptionRef>6ba7b814-9dad-11d1-5-00c04fd430c8</siri:SubscriptionRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:30:06.000+02:00</siri:CreationTime>
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
                    <siri:AlertCause>maintenanceWork</siri:AlertCause>
                    <siri:Summary xml:lang="FR">Nouveau pass Navigo</siri:Summary>
                    <siri:Description xml:lang="EN">The new pass is available</siri:Description>
                    <siri:Affects>
                      <siri:Networks>
                        <siri:AffectedNetwork>
                          <siri:AllLines/>
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
            </Notification>
            <SiriExtension/>
          </sw:NotifySituationExchange>
        </S:Body>
      </S:Envelope>
      """
    Then one Situation has the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T01:02:03+02:00         |
      | Version                      | 1                                 |
      | VersionedAt                  | 2017-01-01T01:02:03+02:00         |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | AlertCause                   | maintenanceWork                   |
      | Description[Translations]#EN | The new pass is available         |
      | Summary[Translations]#FR     | Nouveau pass Navigo               |
      | Affects[AllLines]            |                                   |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-4-00c04fd430c8 |

  @ARA-1450
  Scenario: Send DeleteSubscriptionRequests for wrong Subscription
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
         <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</ns5:RequestMessageRef>
       </SubscriptionAnswerInfo>
       <Answer
         xmlns:ns2="http://www.ifopt.org.uk/acsb"
         xmlns:ns3="http://www.ifopt.org.uk/ifopt"
         xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
         xmlns:ns5="http://www.siri.org.uk/siri"
         xmlns:ns6="http://wsdl.siri.org.uk/siri">
         <ns5:ResponseStatus>
             <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
             <ns5:RequestMessageRef>RATPDev:Message::6ba7b814-9dad-11d1-6-00c04fd430c8:LOC</ns5:RequestMessageRef>
             <ns5:SubscriberRef>NINOXE:default</ns5:SubscriberRef>
             <ns5:SubscriptionRef>6ba7b814-9dad-11d1-3-00c04fd430c8</ns5:SubscriptionRef>
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
        | remote_url                      | http://localhost:8090 |
        | remote_credential               | test                  |
        | local_credential                | NINOXE:default        |
        | remote_code_space               | internal              |
        | collect.filter_general_messages | true                  |
        | collect.include_lines           | NINOXE:Line::3:LOC    |
      And a Line exists with the following attributes:
        | Name              | Test                            |
        | Codes             | "internal":"NINOXE:Line::3:LOC" |
        | CollectSituations | true                            |
      And a Subscription exist with the following attributes:
        | Kind              | SituationExchangeCollect               |
        | ReferenceArray[0] | Line, "internal": "NINOXE:Line::3:LOC" |
      Then show me ara subscriptions for partner "test"
      When I send this SIRI request
        """
     <?xml version='1.0' encoding='utf-8'?>
     <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
       <S:Body>
          <sw:NotifySituationExchange xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>NINOXE:default</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-4-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>33170d7c-35e3-11ee-8a32-7f95f59ec38f</siri:RequestMessageRef>

            </ServiceDeliveryInfo>
            <Notification>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:SubscriptionRef>wrong</siri:SubscriptionRef>
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
                </siri:PtSituationElement>
                </siri:Situations>
              </siri:SituationExchangeDelivery>
            </Notification>
            <SiriExtension/>
          </sw:NotifySituationExchange>
        </S:Body>
      </S:Envelope>
    """
    Then the SIRI server should have received 1 DeleteSubscription request
    And an audit event should exist with these attributes:
      | Protocol                | siri                      |
      | Direction               | sent                      |
      | Type                    | DeleteSubscriptionRequest |
      | SubscriptionIdentifiers | ["wrong"]                 |

  @ARA-1450 @siri-valid
  Scenario: Send Subscriptions requests
   Given a SIRI server on "http://localhost:8090"
      And a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server,siri-situation-exchange-subscription-collector] and the following settings:
        | remote_url                 | http://localhost:8090      |
        | remote_credential          | test                       |
        | local_credential           | NINOXE:default             |
        | remote_code_space          | internal                   |
        | collect.filter_situations  | true                       |
        | collect.include_lines      | NINOXE:Line::3:LOC         |
        | collect.include_stop_areas | NINOXE:StopPoint:SP:24:LOC |
      And 30 seconds have passed
      And a Line exists with the following attributes:
        | Name              | Test                            |
        | Codes             | "internal":"NINOXE:Line::3:LOC" |
        | CollectSituations | true                            |
      And a Line exists with the following attributes:
        | Name              | Test                             |
        | Codes             | "internal":"NINOXE:Line::BP:LOC" |
        | CollectSituations | true                             |
      And a StopArea exists with the following attributes:
        | Name              | Test                                     |
        | Codes             | "internal": "NINOXE:StopPoint:SP:24:LOC" |
        | CollectSituations | true                                     |
      And a minute has passed
      And a minute has passed
      Then the SIRI server should have received a SituationExchangeSubscriptionRequest request with:
        | //siri:LineRef      | NINOXE:Line::3:LOC         |
        | //siri:StopPointRef | NINOXE:StopPoint:SP:24:LOC |
      And an audit event should exist with these attributes:
        | Protocol  | siri                                 |
        | Direction | sent                                 |
        | Type      | SituationExchangeSubscriptionRequest |
        | StopAreas | ["NINOXE:StopPoint:SP:24:LOC"]       |
        | Lines     | ["NINOXE:Line::3:LOC"]               |

  @ARA-1451 @siri-valid
  Scenario: Handle SituationExchange subscription request to all Situations
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-situation-exchange-subscription-broadcaster] and the following settings:
       | remote_url        | http://localhost:8090 |
       | remote_credential | test                  |
       | local_credential  | NINOXE:default        |
       | remote_code_space | internal              |
    And a Line exists with the following attributes:
      | Codes | "another": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                  |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:A:BUS" |
      | Name  | Ligne A Bus                     |
    And a minute has passed
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:SituationExchangeSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>6ba7b814-9dad-11d1--00c04fd430c8</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:SituationExchangeRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:LineRef>dummy</siri:LineRef>
              </siri:SituationExchangeRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:SituationExchangeSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
      """
    Then Subscriptions exist with the following resources:
      | SituationResource | Situation |

  @ARA-1451 @siri-valid
  Scenario: Brodcast a SituationExchange Notification after modification of a Situation
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-situation-exchange-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | internal              |
    And a Subscription exist with the following attributes:
      | Kind              | SituationExchangeBroadcast                  |
      | ExternalId        | externalId                                  |
      | SubscriberRef     | subscriber                                  |
      | ReferenceArray[0] | Situation, "SituationResource": "Situation" |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "internal":"1234" |
      | CollectSituations | true              |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:BP:LOC" |
      | Name  | Ligne BP Metro                   |
    And a Situation exists with the following attributes:
      | Codes                                                                              | "internal" : "NINOXE:GeneralMessage:27_1"     |
      | RecordedAt                                                                         | 2017-01-01T03:30:06+02:00                     |
      | Version                                                                            | 1                                             |
      | Keywords                                                                           | ["Perturbation"]                              |
      | Reality                                                                            | technicalExercise                             |
      | Description[DefaultValue]                                                          | La nouvelle carte d'abonnement est disponible |
      | Description[Translations]#EN                                                       | The new pass is available                     |
      | Summary[Translations]#FR                                                           | Nouveau pass Navigo                           |
      | Summary[Translations]#EN                                                           | New pass Navigo                               |
      | ValidityPeriods[0]#StartTime                                                       | 2017-01-01T01:30:06+02:00                     |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                     |
      | Description[DefaultValue]                                                          | a very very very long message                 |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-3-00c04fd430c8             |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-7-00c04fd430c8             |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-8-00c04fd430c8             |
    And a Situation exists with the following attributes:
      | Codes                        | "internal" : "NINOXE:SituationExchange:01_1" |
      | RecordedAt                   | 2017-01-01T03:30:06+02:00                    |
      | Version                      | 1                                            |
      | Keywords                     | ["test"]                                     |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00                    |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00                    |
      | Description[DefaultValue]    | An Another Very Long Message                 |
      | Affects[Line]                | 6ba7b814-9dad-11d1-3-00c04fd430c8            |
    When the Situation "internal":"NINOXE:SituationExchange:01_1" is edited with an InfoLink with the following attributes:
      | Uri         | http://example.com               |
      | Label       | Label Sample                     |
      | ImageRef    | http://www.example.com/image.png |
      | LinkContent | relatedSite                      |
    When the Situation "internal":"NINOXE:SituationExchange:01_1" is edited with an InfoLink with the following attributes:
      | Uri         | http://anotherexample.com               |
      | Label       | another Label Sambple                   |
      | ImageRef    | http://www.example.com/t2_line_info.pdf |
      | LinkContent | advice                                  |
    When the Situation "internal":"NINOXE:SituationExchange:01_1" is edited with a Consequence with the following attributes:
      | Periods[0]#StartTime                                                          | 2023-09-18T05:30:59Z              |
      | Periods[0]#EndTime                                                            | 2023-09-18T08:00:54Z              |
      | Severity                                                                      | verySlight                        |
      | Condition                                                                     | changeOfPlatform                  |
    And the Situation "internal":"NINOXE:SituationExchange:01_1" is edited with a PublishToWebAction with the following attributes:
      | ActionStatus                    | published                                                                            |
      | Incidents                       | true                                                                                 |
      | HomePage                        | true                                                                                 |
      | Ticker                          | false                                                                                |
      | SocialNetworks                  | ["facebook.com"]                                                                     |
      | Description[Translations]#NO    | Toget vil bytte togmateriell på Dovre. Du må dessverre bytte tog på denne stasjonen. |
      | PublicationWindows[0]#StartTime | 2017-12-01T09:00:00Z                                                                 |
      | PublicationWindows[0]#EndTime   | 2017-12-09T17:00:00Z                                                                 |
      | Name                            | DataName                                                                             |
      | ActionType                      | dummy                                                                                |
      | Value                           | dummy1                                                                               |
      | Prompt[Translations]#NO         | Toget vil bytte togmateriell på Dovre.                                               |
      | Prompt[Translations]#EN         | You must change trains at Dovre. We apologize for the inconvenience.                 |
      | ScopeType                       | line                                                                                 |
      | Affects[Line]                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                                                    |
    And the Situation "internal":"NINOXE:SituationExchange:01_1" is edited with a PublishToMobileAction with the following attributes:
      | ActionStatus                                                   | open                                                                                         |
      | Incidents                                                      | false                                                                                        |
      | HomePage                                                       | true                                                                                         |
      | Description[Translations]#DE                                   | Der Zug wird in Dovre das Zugmaterial wechseln. Leider muss man an diesem Bahnhof umsteigen. |
      | PublicationWindows[0]#StartTime                                | 2017-12-01T09:00:00Z                                                                         |
      | PublicationWindows[0]#EndTime                                  | 2017-12-09T17:00:00Z                                                                         |
      | Name                                                           | AnotherDataName                                                                              |
      | ActionType                                                     | dummy2                                                                                       |
      | Value                                                          | dummy3                                                                                       |
      | Prompt[Translations]#DE                                        | Der Zug wird in Dovre das Zugmaterial wechseln.                                              |
      | Prompt[Translations]#HU                                        | A vonat Dovre-ban módosítja a vonat anyagát. Sajnos ezen az állomáson át kell szállni.       |
      | ScopeType                                                      | stopPlace                                                                                    |
      | Affects[StopArea]                                              | 6ba7b814-9dad-11d1-7-00c04fd430c8                                                            |
      | Affects[StopArea=6ba7b814-9dad-11d1-4-00c04fd430c8]/LineIds[0] | 6ba7b814-9dad-11d1-3-00c04fd430c8                                                            |
      | Affects[StopArea=6ba7b814-9dad-11d1-4-00c04fd430c8]/LineIds[1] | 6ba7b814-9dad-11d1-4-00c04fd430c8                                                            |
    And the Situation "internal":"NINOXE:SituationExchange:01_1" is edited with a PublishToDisplayAction with the following attributes:
      | ActionStatus                    | closed                                                                                               |
      | OnBoard                         | false                                                                                                |
      | OnPlace                         | true                                                                                                 |
      | Description[Translations]#ES    | El tren cambiará su material en Dovre. Lamentablemente tienes que hacer transbordo en esta estación. |
      | PublicationWindows[0]#StartTime | 2017-12-01T09:00:00Z                                                                                 |
      | PublicationWindows[0]#EndTime   | 2017-12-09T17:00:00Z                                                                                 |
      | Name                            | NewDataName                                                                                          |
      | ActionType                      | dummy4                                                                                               |
      | Value                           | dummy5                                                                                               |
      | Prompt[Translations]#PL         | Pociąg zmieni materiał składowy w Dovre.                                                             |
      | Prompt[Translations]#BG         | Влакът променя материала на влака в Dovre. За съжаление трябва да сменяте на тази гара.              |
      | ScopeType                       | general                                                                                              |
      | Affects[Line]                   | 6ba7b814-9dad-11d1-4-00c04fd430c8                                                                    |
    And the Situation "internal":"NINOXE:SituationExchange:01_1" is edited with a PublishToDisplayAction with the following attributes:
      | ActionStatus                    | open                                                                                       |
      | OnBoard                         | true                                                                                       |
      | OnPlace                         | false                                                                                      |
      | Description[Translations]#IS    | Lestin mun skipta um efni í Dovre. Því miður þarftu að flytja á þessari stöð.              |
      | PublicationWindows[0]#StartTime | 2017-12-01T09:00:00Z                                                                       |
      | PublicationWindows[0]#EndTime   | 2017-12-09T17:00:00Z                                                                       |
      | Name                            | AnotherNewDataName                                                                         |
      | ActionType                      | dummy6                                                                                     |
      | Value                           | dummy7                                                                                     |
      | Prompt[Translations]#DK         | Toget skifter materialet på Dovre.                                                         |
      | Prompt[Translations]#FI         | Juna vaihtaa junamateriaalia Dovressa. Valitettavasti sinun täytyy vaihtaa tällä asemalla. |
      | ScopeType                       | network                                                                                    |
      | Affects[Line]                   | 6ba7b814-9dad-11d1-3-00c04fd430c8                                                          |
    And a StopArea exists with the following attributes:
      | Name              | Test                                    |
      | Codes             | "internal":"NINOXE:StopPoint:SP:24:LOC" |
      | CollectSituations | true                                    |
    And a StopArea exists with the following attributes:
      | Name              | Test                                    |
      | Codes             | "internal":"NINOXE:StopPoint:SP:12:LOC" |
      | CollectSituations | true                                    |
    And 10 seconds have passed
    When the Situation "6ba7b814-9dad-11d1-5-00c04fd430c8" is edited with the following attributes:
      | RecordedAt                   | 2017-01-01T03:50:06+02:00              |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00              |
      | ValidityPeriods[0]#EndTime   | 2017-10-24T20:30:06+02:00              |
      | Description[DefaultValue]    | an ANOTHER very very very long message |
      | Version                      | 2                                      |
    When the Situation "6ba7b814-9dad-11d1-6-00c04fd430c8" is edited with the following attributes:
      | RecordedAt                   | 2017-01-01T03:50:06+02:00                   |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00                   |
      | ValidityPeriods[0]#EndTime   | 2017-10-24T20:30:06+02:00                   |
      | Description[DefaultValue]    | a SUPER ANOTHER very very very long message |
      | Version                      | 3                                           |
    And 15 seconds have passed
    Then the SIRI server should receive this response
    """
     <?xml version='1.0' encoding='utf-8'?>
     <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
       <S:Body>
         <sw:NotifySituationExchange xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
           <ServiceDeliveryInfo>
             <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
             <siri:ProducerRef>test</siri:ProducerRef>
             <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-a-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
           </ServiceDeliveryInfo>
           <Notification>
             <siri:SituationExchangeDelivery version="2.0:FR-IDF-2.4" xmlns:stif="http://wsdl.siri.org.uk/siri">
               <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
               <siri:SubscriberRef>subscriber</siri:SubscriberRef>
               <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
               <siri:Status>true</siri:Status>
               <siri:Situations>
                 <siri:PtSituationElement>
                   <siri:CreationTime>2017-01-01T03:50:06.000+02:00</siri:CreationTime>
                   <siri:SituationNumber>NINOXE:GeneralMessage:27_1</siri:SituationNumber>
                   <siri:Version>2</siri:Version>
                   <siri:Source>
                     <siri:SourceType>directReport</siri:SourceType>
                   </siri:Source>
                   <siri:Reality>technicalExercise</siri:Reality>
                   <siri:ValidityPeriod>
                     <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                     <siri:EndTime>2017-10-24T20:30:06.000+02:00</siri:EndTime>
                   </siri:ValidityPeriod>
                   <siri:UndefinedReason />
                   <siri:Keywords>Perturbation</siri:Keywords>
                   <siri:Summary xml:lang='EN'>New pass Navigo</siri:Summary>
                   <siri:Summary xml:lang='FR'>Nouveau pass Navigo</siri:Summary>
                   <siri:Description>an ANOTHER very very very long message</siri:Description>
                   <siri:Description xml:lang='EN'>The new pass is available</siri:Description>
                   <siri:Affects>
                     <siri:Networks>
                       <siri:AffectedNetwork>
                         <siri:AffectedLine>
                           <siri:LineRef>1234</siri:LineRef>
                           <siri:Destinations>
                             <siri:StopPlaceRef>NINOXE:StopPoint:SP:12:LOC</siri:StopPlaceRef>
                           </siri:Destinations>
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
                   <siri:CreationTime>2017-01-01T03:50:06.000+02:00</siri:CreationTime>
                   <siri:SituationNumber>NINOXE:SituationExchange:01_1</siri:SituationNumber>
                   <siri:Version>3</siri:Version>
                   <siri:Source>
                     <siri:SourceType>directReport</siri:SourceType>
                   </siri:Source>
                   <siri:ValidityPeriod>
                     <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                     <siri:EndTime>2017-10-24T20:30:06.000+02:00</siri:EndTime>
                   </siri:ValidityPeriod>
                   <siri:UndefinedReason/>
                   <siri:Keywords>test</siri:Keywords>
                   <siri:Description>a SUPER ANOTHER very very very long message</siri:Description>
                    <siri:InfoLinks>
                     <siri:InfoLink>
                       <siri:Uri>http://example.com</siri:Uri>
                       <siri:Label>Label Sample</siri:Label>
                       <siri:Image>
                         <siri:ImageRef>http://www.example.com/image.png</siri:ImageRef>
                       </siri:Image>
                       <siri:LinkContent>relatedSite</siri:LinkContent>
                     </siri:InfoLink>
                     <siri:InfoLink>
                       <siri:Uri>http://anotherexample.com</siri:Uri>
                       <siri:Label>another Label Sambple</siri:Label>
                       <siri:Image>
                         <siri:ImageRef>http://www.example.com/t2_line_info.pdf</siri:ImageRef>
                       </siri:Image>
                       <siri:LinkContent>advice</siri:LinkContent>
                     </siri:InfoLink>
                   </siri:InfoLinks>
                   <siri:Affects>
                     <siri:Networks>
                       <siri:AffectedNetwork>
                         <siri:AffectedLine>
                           <siri:LineRef>1234</siri:LineRef>
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
                        <siri:Condition>changeOfPlatform</siri:Condition>
                        <siri:Severity>verySlight</siri:Severity>
                      </siri:Consequence>
                    </siri:Consequences>
                    <siri:PublishingActions>
                      <siri:PublishToWebAction>
                        <siri:ActionStatus>published</siri:ActionStatus>
                        <siri:Description xml:lang='NO'>Toget vil bytte togmateriell på Dovre. Du må dessverre bytte tog på denne stasjonen.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>DataName</siri:Name>
                          <siri:Type>dummy</siri:Type>
                          <siri:Value>dummy1</siri:Value>
                          <siri:Prompt xml:lang='EN'>You must change trains at Dovre. We apologize for the inconvenience.</siri:Prompt>
                          <siri:Prompt xml:lang='NO'>Toget vil bytte togmateriell på Dovre.</siri:Prompt>
                          <siri:PublishAtScope>
                            <siri:ScopeType>line</siri:ScopeType>
                            <siri:Affects>
                              <siri:Networks>
                                <siri:AffectedNetwork>
                                  <siri:AffectedLine>
                                    <siri:LineRef>1234</siri:LineRef>
                                  </siri:AffectedLine>
                                </siri:AffectedNetwork>
                              </siri:Networks>
                            </siri:Affects>
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:Incidents>true</siri:Incidents>
                        <siri:HomePage>true</siri:HomePage>
                        <siri:SocialNetwork>facebook.com</siri:SocialNetwork>
                      </siri:PublishToWebAction>
                      <siri:PublishToMobileAction>
                        <siri:ActionStatus>open</siri:ActionStatus>
                        <siri:Description xml:lang='DE'>Der Zug wird in Dovre das Zugmaterial wechseln. Leider muss man an diesem Bahnhof umsteigen.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>AnotherDataName</siri:Name>
                          <siri:Type>dummy2</siri:Type>
                          <siri:Value>dummy3</siri:Value>
                          <siri:Prompt xml:lang='DE'>Der Zug wird in Dovre das Zugmaterial wechseln.</siri:Prompt>
                          <siri:Prompt xml:lang='HU'>A vonat Dovre-ban módosítja a vonat anyagát. Sajnos ezen az állomáson át kell szállni.</siri:Prompt>
                          <siri:PublishAtScope>
                            <siri:ScopeType>stopPlace</siri:ScopeType>
                            <siri:Affects>
                              <siri:StopPoints>
                                <siri:AffectedStopPoint>
                                  <siri:StopPointRef>NINOXE:StopPoint:SP:24:LOC</siri:StopPointRef>
                                  <siri:Lines>
                                    <siri:AffectedLine>
                                      <siri:LineRef>1234</siri:LineRef>
                                    </siri:AffectedLine>
                                    <siri:AffectedLine>
                                      <siri:LineRef>NINOXE:Line:BP:LOC</siri:LineRef>
                                    </siri:AffectedLine>
                                  </siri:Lines>
                                </siri:AffectedStopPoint>
                              </siri:StopPoints>
                            </siri:Affects>
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:Incidents>false</siri:Incidents>
                        <siri:HomePage>true</siri:HomePage>
                      </siri:PublishToMobileAction>
                      <siri:PublishToDisplayAction>
                        <siri:ActionStatus>closed</siri:ActionStatus>
                        <siri:Description xml:lang='ES'>El tren cambiará su material en Dovre. Lamentablemente tienes que hacer transbordo en esta estación.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>NewDataName</siri:Name>
                          <siri:Type>dummy4</siri:Type>
                          <siri:Value>dummy5</siri:Value>
                          <siri:Prompt xml:lang='BG'>Влакът променя материала на влака в Dovre. За съжаление трябва да сменяте на тази гара.</siri:Prompt>
                          <siri:Prompt xml:lang='PL'>Pociąg zmieni materiał składowy w Dovre.</siri:Prompt>
                          <siri:PublishAtScope>
                            <siri:ScopeType>general</siri:ScopeType>
                            <siri:Affects>
                              <siri:Networks>
                                <siri:AffectedNetwork>
                                  <siri:AffectedLine>
                                    <siri:LineRef>NINOXE:Line:BP:LOC</siri:LineRef>
                                  </siri:AffectedLine>
                                </siri:AffectedNetwork>
                              </siri:Networks>
                            </siri:Affects>
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:OnPlace>true</siri:OnPlace>
                        <siri:OnBoard>false</siri:OnBoard>
                      </siri:PublishToDisplayAction>
                      <siri:PublishToDisplayAction>
                        <siri:ActionStatus>open</siri:ActionStatus>
                        <siri:Description xml:lang='IS'>Lestin mun skipta um efni í Dovre. Því miður þarftu að flytja á þessari stöð.</siri:Description>
                        <siri:ActionData>
                          <siri:Name>AnotherNewDataName</siri:Name>
                          <siri:Type>dummy6</siri:Type>
                          <siri:Value>dummy7</siri:Value>
                          <siri:Prompt xml:lang='DK'>Toget skifter materialet på Dovre.</siri:Prompt>
                          <siri:Prompt xml:lang='FI'>Juna vaihtaa junamateriaalia Dovressa. Valitettavasti sinun täytyy vaihtaa tällä asemalla.</siri:Prompt>
                          <siri:PublishAtScope>
                            <siri:ScopeType>network</siri:ScopeType>
                            <siri:Affects>
                              <siri:Networks>
                                <siri:AffectedNetwork>
                                  <siri:AffectedLine>
                                    <siri:LineRef>1234</siri:LineRef>
                                  </siri:AffectedLine>
                                </siri:AffectedNetwork>
                              </siri:Networks>
                            </siri:Affects>
                          </siri:PublishAtScope>
                        </siri:ActionData>
                        <siri:PublicationWindow>
                          <siri:StartTime>2017-12-01T09:00:00.000Z</siri:StartTime>
                          <siri:EndTime>2017-12-09T17:00:00.000Z</siri:EndTime>
                        </siri:PublicationWindow>
                        <siri:OnPlace>false</siri:OnPlace>
                        <siri:OnBoard>true</siri:OnBoard>
                      </siri:PublishToDisplayAction>
                    </siri:PublishingActions>
                 </siri:PtSituationElement>
               </siri:Situations>
             </siri:SituationExchangeDelivery>
           </Notification>
           <SiriExtension />
         </sw:NotifySituationExchange>
       </S:Body>
     </S:Envelope>
    """
    And an audit event should exist with these attributes:
      | Protocol                | siri                                                         |
      | Direction               | sent                                                         |
      | Status                  | OK                                                           |
      | Type                    | NotifySituationExchange                                      |
      | SubscriptionIdentifiers | ["externalId"]                                               |
      | StopAreas               | ["NINOXE:StopPoint:SP:24:LOC", "NINOXE:StopPoint:SP:12:LOC"] |
      | Lines                   | ["1234", "NINOXE:Line:BP:LOC"]                               |

  @ARA-1582 @siri-valid
  Scenario: Brodcast a SituationExchange Notification after modification of a Situation having All affected Lines
    Given a SIRI server on "http://localhost:8090"
    And a SIRI Partner "test" exists with connectors [siri-check-status-client, siri-situation-exchange-subscription-broadcaster] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | local_credential  | NINOXE:default        |
      | remote_code_space | external              |
    And a Subscription exist with the following attributes:
      | Kind              | SituationExchangeBroadcast                  |
      | ExternalId        | externalId                                  |
      | SubscriberRef     | subscriber                                  |
      | ReferenceArray[0] | Situation, "SituationResource": "Situation" |
    And a Line exists with the following attributes:
      | Codes | "internal": "NINOXE:Line:3:LOC" |
      | Name  | Ligne 3 Metro                   |
    And a StopArea exists with the following attributes:
      | Name  | Test                                     |
      | Codes | "external": "NINOXE:StopPoint:SP:24:LOC" |
    And a Situation exists with the following attributes:
      | Codes                        | "external" : "test"               |
      | RecordedAt                   | 2017-01-01T01:02:03+02:00         |
      | Version                      | 1                                 |
      | VersionedAt                  | 2017-01-01T01:02:03+02:00         |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00         |
      | ValidityPeriods[0]#EndTime   | 2017-01-01T20:30:06+02:00         |
      | AlertCause                   | maintenanceWork                   |
      | Description[Translations]#EN | The new pass is available         |
      | Summary[Translations]#FR     | Nouveau pass Navigo               |
      | Affects[AllLines]            |                                   |
      | Affects[StopArea]            | 6ba7b814-9dad-11d1-4-00c04fd430c8 |
    And 10 seconds have passed
    When the Situation "6ba7b814-9dad-11d1-5-00c04fd430c8" is edited with the following attributes:
      | RecordedAt                   | 2017-01-01T03:50:06+02:00              |
      | ValidityPeriods[0]#StartTime | 2017-01-01T01:30:06+02:00              |
      | ValidityPeriods[0]#EndTime   | 2017-10-24T20:30:06+02:00              |
      | Description[DefaultValue]    | an ANOTHER very very very long message |
      | Version                      | 2                                      |
    And 15 seconds have passed
    Then the SIRI server should receive this response
      """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:NotifySituationExchange xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>test</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-7-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
            </ServiceDeliveryInfo>
            <Notification>
              <siri:SituationExchangeDelivery version='2.0:FR-IDF-2.4' xmlns:stif='http://wsdl.siri.org.uk/siri'>
                <siri:ResponseTimestamp>2017-01-01T12:00:25.000Z</siri:ResponseTimestamp>
                <siri:SubscriberRef>subscriber</siri:SubscriberRef>
                <siri:SubscriptionRef>externalId</siri:SubscriptionRef>
                <siri:Status>true</siri:Status>
                <siri:Situations>
                  <siri:PtSituationElement>
                    <siri:CreationTime>2017-01-01T03:50:06.000+02:00</siri:CreationTime>
                    <siri:SituationNumber>test</siri:SituationNumber>
                    <siri:Version>2</siri:Version>
                    <siri:Source>
                      <siri:SourceType>directReport</siri:SourceType>
                    </siri:Source>
                    <siri:VersionedAtTime>2017-01-01T01:02:03.000+02:00</siri:VersionedAtTime>
                    <siri:ValidityPeriod>
                      <siri:StartTime>2017-01-01T01:30:06.000+02:00</siri:StartTime>
                      <siri:EndTime>2017-10-24T20:30:06.000+02:00</siri:EndTime>
                    </siri:ValidityPeriod>
                    <siri:AlertCause>maintenanceWork</siri:AlertCause>
                    <siri:Summary xml:lang='FR'>Nouveau pass Navigo</siri:Summary>
                    <siri:Description>an ANOTHER very very very long message</siri:Description>
                    <siri:Description xml:lang='EN'>The new pass is available</siri:Description>
                    <siri:Affects>
                      <siri:Networks>
                        <siri:AffectedNetwork>
                          <siri:AllLines/>
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
            </Notification>
            <SiriExtension/>
          </sw:NotifySituationExchange>
        </S:Body>
      </S:Envelope>
      """
    And an audit event should exist with these attributes:
      | Protocol                | siri                           |
      | Direction               | sent                           |
      | Status                  | OK                             |
      | Type                    | NotifySituationExchange        |
      | SubscriptionIdentifiers | ["externalId"]                 |
      | StopAreas               | ["NINOXE:StopPoint:SP:24:LOC"] |

  @ARA-1451 @siri-valid
  Scenario: Handle SituationExchange subscription request to an unknowm line
    Given a SIRI server on "http://localhost:8090"
    Given a Partner "test" exists with connectors [siri-check-status-client,siri-check-status-server ,siri-situation-exchange-subscription-broadcaster] and the following settings:
       | remote_url        | http://localhost:8090 |
       | remote_credential | test                  |
       | local_credential  | NINOXE:default        |
       | remote_code_space | internal              |
    And a Line exists with the following attributes:
      | Name              | Test              |
      | Codes             | "internal":"1234" |
      | CollectSituations | true              |
    And a Situation exists with the following attributes:
      | Codes                                                                              | "internal" : "NINOXE:GeneralMessage:27_1" |
      | RecordedAt                                                                         | 2017-01-01T03:30:06+02:00                 |
      | Version                                                                            | 1                                         |
      | Keywords                                                                           | ["Perturbation"]                          |
      | ValidityPeriods[0]#StartTime                                                       | 2017-01-01T01:30:06+02:00                 |
      | ValidityPeriods[0]#EndTime                                                         | 2017-01-01T20:30:06+02:00                 |
      | Description[DefaultValue]                                                          | a very very very long message             |
      | Affects[Line]                                                                      | 6ba7b814-9dad-11d1-2-00c04fd430c8         |
      | Affects[StopArea]                                                                  | 6ba7b814-9dad-11d1-5-00c04fd430c8         |
      | Affects[Line=6ba7b814-9dad-11d1-2-00c04fd430c8]/AffectedDestinations[0]/StopAreaId | 6ba7b814-9dad-11d1-6-00c04fd430c8         |
    And a StopArea exists with the following attributes:
      | Name              | Test                                    |
      | Codes             | "internal":"NINOXE:StopPoint:SP:24:LOC" |
      | CollectSituations | true                                    |
    And a StopArea exists with the following attributes:
      | Name              | Test                                    |
      | Codes             | "internal":"NINOXE:StopPoint:SP:12:LOC" |
      | CollectSituations | true                                    |
    And a minute has passed
    When I send this SIRI request
      """
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:SituationExchangeSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>6ba7b814-9dad-11d1--00c04fd430c8</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:SituationExchangeRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:LineRef>dummy</siri:LineRef>
              </siri:SituationExchangeRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:SituationExchangeSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>
    </S:Envelope>
      """
    And 30 seconds have passed
    And the SIRI server should not have received a NotifySituationExchange request
