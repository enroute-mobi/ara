Feature: Support SIRI FacilityMonitoring by request

  Background:
    Given a Referential "test" is created

  @ARA-1731
  Scenario: Performs a SIRI FacilityMonitoring request to a Partner
    Given a SIRI server waits GetFacilityMonitoring request on "http://localhost:8090" to respond with
      """
     <?xml version='1.0' encoding='utf-8'?>
     <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
       <soap:Body>
         <sw:GetFacilityMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
           <ServiceDeliveryInfo>
             <siri:ResponseTimestamp>2030-01-01T12:01:10.000Z</siri:ResponseTimestamp>
           </ServiceDeliveryInfo>
           <Answer>
             <siri:FacilityMonitoringDelivery>
               <siri:ResponseTimestamp>2030-01-01T15:00:00.000Z</siri:ResponseTimestamp>
               <siri:FacilityCondition>
                 <siri:FacilityRef>NINOXE:Facility:ABC1:LOC</siri:FacilityRef>
                 <siri:FacilityStatus>
                   <siri:Status>available</siri:Status>
                 </siri:FacilityStatus>
               </siri:FacilityCondition>
             </siri:FacilityMonitoringDelivery>
           </Answer>
           <AnswerExtension/>
         </sw:GetFacilityMonitoringResponse>
       </soap:Body>
     </soap:Envelope>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-facility-monitoring-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | remote_code_space | internal              |
    And a minute has passed
    And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:ABC1:LOC |
    When a minute has passed
    And the SIRI server has received a GetFacilityMonitoring request
    Then one Facility has the following attributes:
      | Codes[internal] | NINOXE:Facility:ABC1:LOC |
      | Status          | available                |
    And an audit event should exist with these attributes:
      | Protocol  | siri                      |
      | Direction | sent                      |
      | Status    | OK                        |
      | Type      | FacilityMonitoringRequest |

  Scenario: Handle a SIRI FacilityMonitoring request
    Given a SIRI Partner "test" exists with connectors [siri-facility-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:ABC1:LOC |
      | Status          | available                |
    When I send this SIRI request
    """
    <?xml version='1.0' encoding='utf-8'?>
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetFacilityMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2016-09-22T07:54:52.977Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>FacilityMonitoring:Test:0</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request>
            <siri:RequestTimestamp>2016-09-22T07:54:52.977Z</siri:RequestTimestamp>
            <siri:FacilityRef>NINOXE:Facility:ABC1:LOC</siri:FacilityRef>
          </Request>
          <RequestExtension/>
        </sw:GetFacilityMonitoring>
      </soap:Body>
    </soap:Envelope>
    """
    Then I should receive this SIRI response
    """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetFacilityMonitoringResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>FacilityMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:FacilityMonitoringDelivery version='2.0:FR-IDF-2.4'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>FacilityMonitoring:Test:0</siri:RequestMessageRef>
                <siri:FacilityCondition>
                  <siri:FacilityRef>NINOXE:Facility:ABC1:LOC</siri:FacilityRef>
                  <siri:FacilityStatus>
                    <siri:Status>available</siri:Status>
                  </siri:FacilityStatus>
                </siri:FacilityCondition>
              </siri:FacilityMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetFacilityMonitoringResponse>
        </S:Body>
      </S:Envelope>
    """
    And an audit event should exist with these attributes:
      | Protocol        | siri                           |
      | Direction       | received                       |
      | Status          | OK                             |
      | Type            | FacilityMonitoringRequest      |

  Scenario: Handle a SIRI FacilityMonitoring request on an unknown Facility
    Given a SIRI Partner "test" exists with connectors [siri-facility-monitoring-request-broadcaster] and the following settings:
      | local_credential  | test     |
      | remote_code_space | internal |
    And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:ABC1:LOC |
      | Status          | available                |
    When I send this SIRI request
    """
    <?xml version='1.0' encoding='utf-8'?>
    <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
        <sw:GetFacilityMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <ServiceRequestInfo>
            <siri:RequestTimestamp>2016-09-22T07:54:52.977Z</siri:RequestTimestamp>
            <siri:RequestorRef>test</siri:RequestorRef>
            <siri:MessageIdentifier>FacilityMonitoring:Test:0</siri:MessageIdentifier>
          </ServiceRequestInfo>
          <Request>
            <siri:RequestTimestamp>2016-09-22T07:54:52.977Z</siri:RequestTimestamp>
            <siri:FacilityRef>UNKNOWN</siri:FacilityRef>
          </Request>
          <RequestExtension/>
        </sw:GetFacilityMonitoring>
      </soap:Body>
    </soap:Envelope>
    """
    Then I should receive this SIRI response
    """
      <?xml version='1.0' encoding='UTF-8'?>
      <S:Envelope xmlns:S='http://schemas.xmlsoap.org/soap/envelope/'>
        <S:Body>
          <sw:GetFacilityMonitoringResponse xmlns:sw='http://wsdl.siri.org.uk' xmlns:siri='http://www.siri.org.uk/siri'>
            <ServiceDeliveryInfo>
              <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
              <siri:ProducerRef>Ara</siri:ProducerRef>
              <siri:ResponseMessageIdentifier>RATPDev:ResponseMessage::6ba7b814-9dad-11d1-3-00c04fd430c8:LOC</siri:ResponseMessageIdentifier>
              <siri:RequestMessageRef>FacilityMonitoring:Test:0</siri:RequestMessageRef>
            </ServiceDeliveryInfo>
            <Answer>
              <siri:FacilityMonitoringDelivery version='2.0:FR-IDF-2.4'>
                <siri:ResponseTimestamp>2017-01-01T12:00:00.000Z</siri:ResponseTimestamp>
                <siri:RequestMessageRef>FacilityMonitoring:Test:0</siri:RequestMessageRef>
                <siri:ErrorCondition>
                  <siri:InvalidDataReferencesError>
                    <siri:ErrorText>Facility not found: 'UNKNOWN'</siri:ErrorText>
                  </siri:InvalidDataReferencesError>
                </siri:ErrorCondition>
              </siri:FacilityMonitoringDelivery>
            </Answer>
            <AnswerExtension/>
          </sw:GetFacilityMonitoringResponse>
        </S:Body>
      </S:Envelope>
    """
    And an audit event should exist with these attributes:
      | Protocol     | siri                                                      |
      | Direction    | received                                                  |
      | Type         | FacilityMonitoringRequest                                 |
      | Status       | Error                                                     |
      | ErrorDetails | InvalidDataReferencesError: Facility not found: 'UNKNOWN' |

  @ARA-1731
  Scenario: Performs a raw SIRI FacilityMonitoring request to a Partner
    Given a raw SIRI server waits FacilityMonitoring request on "http://localhost:8090" to respond with
      """
     <?xml version="1.0" encoding="UTF-8"?>
     <Siri xmlns="http://www.siri.org.uk/siri" xmlns:acsb="http://www.ifopt.org.uk/acsb" xmlns:ifopt="http://www.ifopt.org.uk/ifopt" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="2.0" xsi:schemaLocation="http://www.siri.org.uk/siri ../../xsd/siri.xsd">
         <ServiceDelivery>
             <ResponseTimestamp>2030-01-01T12:01:10.000Z</ResponseTimestamp>
             <FacilityMonitoringDelivery>
                 <ResponseTimestamp>2030-01-01T15:00:00.000Z</ResponseTimestamp>
                 <FacilityCondition>
                     <FacilityRef>NINOXE:Facility:ABC1:LOC</FacilityRef>
                     <FacilityStatus>
                         <Status>available</Status>
                     </FacilityStatus>
                 </FacilityCondition>
             </FacilityMonitoringDelivery>
         </ServiceDelivery>
     </Siri>
      """
    And a Partner "test" exists with connectors [siri-check-status-client, siri-facility-monitoring-request-collector] and the following settings:
      | remote_url        | http://localhost:8090 |
      | remote_credential | test                  |
      | remote_code_space | internal              |
      | siri.envelope     | raw                   |
    And a minute has passed
    And a Facility exists with the following attributes:
      | Codes[internal] | NINOXE:Facility:ABC1:LOC |
    When a minute has passed
    And the SIRI server has received a FacilityMonitoring request
    Then one Facility has the following attributes:
      | Codes[internal] | NINOXE:Facility:ABC1:LOC |
      | Status          | available                |
    And an audit event should exist with these attributes:
      | Protocol  | siri                      |
      | Direction | sent                      |
      | Status    | OK                        |
      | Type      | FacilityMonitoringRequest |
