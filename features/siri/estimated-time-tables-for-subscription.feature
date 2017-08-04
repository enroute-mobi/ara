Feature: Support SIRI EstimatedTimeTable by subscription

  Background:
    Given a Referential "test" is created

  @wip
    Scenario: 4233 - Manage a ETT Subscription

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
             <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">Subscription:Test:0</ns5:RequestMessageRef>
           </SubscriptionAnswerInfo>
           <Answer
             xmlns:ns2="http://www.ifopt.org.uk/acsb"
             xmlns:ns3="http://www.ifopt.org.uk/ifopt"
             xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
             xmlns:ns5="http://www.siri.org.uk/siri"
             xmlns:ns6="http://wsdl.siri.org.uk/siri">
             <ns5:ResponseStatus>
                 <ns5:ResponseTimestamp>2016-09-22T08:01:20.227+02:00</ns5:ResponseTimestamp>
                 <ns5:RequestMessageRef>Subscription:Test:0</ns5:RequestMessageRef>
                 <ns5:SubscriberRef>SubscriberRef</ns5:SubscriberRef>
                 <ns5:SubscriptionRef>SubscriptionIdentifier</ns5:SubscriptionRef>
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
