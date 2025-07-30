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
